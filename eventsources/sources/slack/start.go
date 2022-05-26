/*
Copyright 2018 BlackRock, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"

	"github.com/argoproj/argo-events/common"
	"github.com/argoproj/argo-events/common/logging"
	eventsourcecommon "github.com/argoproj/argo-events/eventsources/common"
	"github.com/argoproj/argo-events/eventsources/common/webhook"
	"github.com/argoproj/argo-events/eventsources/sources"
	metrics "github.com/argoproj/argo-events/metrics"
	apicommon "github.com/argoproj/argo-events/pkg/apis/common"
	"github.com/argoproj/argo-events/pkg/apis/eventsource/v1alpha1"
)

// controller controls the webhook operations
var (
	controller = webhook.NewController()
)

// set up the activation and inactivation channels to control the state of routes.
func init() {
	go webhook.ProcessRouteStatus(controller)
}

// EventListener implements Eventing for slack event source
type EventListener struct {
	EventSourceName  string
	EventName        string
	SlackEventSource v1alpha1.SlackEventSource
	Metrics          *metrics.Metrics
}

// GetEventSourceName returns name of event source
func (el *EventListener) GetEventSourceName() string {
	return el.EventSourceName
}

// GetEventName returns name of event
func (el *EventListener) GetEventName() string {
	return el.EventName
}

// GetEventSourceType return type of event server
func (el *EventListener) GetEventSourceType() apicommon.EventSourceType {
	return apicommon.SlackEvent
}

// Router contains information about a REST endpoint
type Router struct {
	// route holds information to process an incoming request
	route *webhook.Route
	// slackEventSource is the event source which refers to configuration required to consume events from slack
	slackEventSource *v1alpha1.SlackEventSource
	// token is the slack token
	token string
	// refer to https://api.slack.com/docs/verifying-requests-from-slack
	signingSecret string
}

// Implement Router
// 1. GetRoute
// 2. HandleRoute
// 3. PostActivate
// 4. PostDeactivate

// GetRoute returns the route
func (rc *Router) GetRoute() *webhook.Route {
	return rc.route
}

// HandleRoute handles incoming requests on the route
func (rc *Router) HandleRoute(writer http.ResponseWriter, request *http.Request) {
	route := rc.route

	logger := route.Logger.With(
		logging.LabelEndpoint, route.Context.Endpoint,
		logging.LabelPort, route.Context.Port,
		logging.LabelHTTPMethod, route.Context.Method,
	)

	logger.Info("request a received, processing it...")

	if !route.Active {
		logger.Warn("endpoint is not active, won't process it")
		common.SendErrorResponse(writer, "endpoint is inactive")
		return
	}

	defer func(start time.Time) {
		route.Metrics.EventProcessingDuration(route.EventSourceName, route.EventName, float64(time.Since(start)/time.Millisecond))
	}(time.Now())

	logger.Info("verifying the request...")
	err := rc.verifyRequest(request)
	if err != nil {
		logger.Errorw("failed to validate the request", zap.Error(err))
		common.SendResponse(writer, http.StatusUnauthorized, err.Error())
		route.Metrics.EventProcessingFailed(route.EventSourceName, route.EventName)
		return
	}

	var data []byte

	// Interactive element actions are always
	// sent as application/x-www-form-urlencoded
	// If request was generated by an interactive element or a slash command, it will be a POST form
	if len(request.Header["Content-Type"]) > 0 && request.Header["Content-Type"][0] == "application/x-www-form-urlencoded" {
		if err := request.ParseForm(); err != nil {
			logger.Errorw("failed to parse form data", zap.Error(err))
			common.SendInternalErrorResponse(writer, err.Error())
			route.Metrics.EventProcessingFailed(route.EventSourceName, route.EventName)
			return
		}

		switch {
		case request.PostForm.Get("payload") != "":
			data, err = rc.handleInteraction(request)
			if err != nil {
				logger.Errorw("failed to process the interaction", zap.Error(err))
				common.SendInternalErrorResponse(writer, err.Error())
				route.Metrics.EventProcessingFailed(route.EventSourceName, route.EventName)
				return
			}

		case request.PostForm.Get("command") != "":
			data, err = rc.handleSlashCommand(request)
			if err != nil {
				logger.Errorw("failed to process the slash command", zap.Error(err))
				common.SendInternalErrorResponse(writer, err.Error())
				route.Metrics.EventProcessingFailed(route.EventSourceName, route.EventName)
				return
			}

		default:
			err = errors.New("could not determine slack type from form parameters")
			logger.Errorw("failed to determine type of slack post", zap.Error(err))
			common.SendInternalErrorResponse(writer, err.Error())
			route.Metrics.EventProcessingFailed(route.EventSourceName, route.EventName)
			return
		}
	} else {
		// If there's no payload in the post body, this is likely an
		// Event API request. Parse and process if valid.
		logger.Info("handling slack event...")
		var response []byte
		data, response, err = rc.handleEvent(request)
		if err != nil {
			logger.Errorw("failed to handle the event", zap.Error(err))
			common.SendInternalErrorResponse(writer, err.Error())
			route.Metrics.EventProcessingFailed(route.EventSourceName, route.EventName)
			return
		}
		if response != nil {
			writer.Header().Set("Content-Type", "text")
			if _, err := writer.Write(response); err != nil {
				logger.Errorw("failed to write the response for url verification", zap.Error(err))
				// don't return, we want to keep this running to give user chance to retry
			}
		}
	}

	if data != nil {
		logger.Info("dispatching event on route's data channel...")
		route.DataCh <- data
	}

	logger.Debug("request successfully processed")
	common.SendSuccessResponse(writer, "success")
}

// PostActivate performs operations once the route is activated and ready to consume requests
func (rc *Router) PostActivate() error {
	return nil
}

// PostInactivate performs operations after the route is inactivated
func (rc *Router) PostInactivate() error {
	return nil
}

// handleEvent parse the slack notification and validates the event type
func (rc *Router) handleEvent(request *http.Request) ([]byte, []byte, error) {
	var err error
	var response []byte
	var data []byte
	body, err := getRequestBody(request)
	if err != nil {
		return data, response, errors.Wrap(err, "failed to fetch request body")
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: rc.token}))
	if err != nil {
		return data, response, errors.Wrap(err, "failed to extract event")
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err = json.Unmarshal(body, &r)
		if err != nil {
			return data, response, errors.Wrap(err, "failed to verify the challenge")
		}
		response = []byte(r.Challenge)
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		data, err = json.Marshal(&eventsAPIEvent.InnerEvent)
		if err != nil {
			return data, response, errors.Wrap(err, "failed to marshal event data, rejecting the event...")
		}
	}

	return data, response, nil
}

func (rc *Router) handleInteraction(request *http.Request) ([]byte, error) {
	payload := request.PostForm.Get("payload")
	ie := &slack.InteractionCallback{}
	err := json.Unmarshal([]byte(payload), ie)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse interaction event")
	}

	data, err := json.Marshal(ie)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal action data")
	}

	return data, nil
}

func (rc *Router) handleSlashCommand(request *http.Request) ([]byte, error) {
	command, err := slack.SlashCommandParse(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse command")
	}

	data, err := json.Marshal(command)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal command data")
	}

	return data, nil
}

func getRequestBody(request *http.Request) ([]byte, error) {
	// Read request payload
	body, err := ioutil.ReadAll(request.Body)
	// Reset request.Body ReadCloser to prevent side-effect if re-read
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request body")
	}
	return body, nil
}

// If a signing secret is provided, validate the request against the
// X-Slack-Signature header value.
// The signature is a hash generated as per Slack documentation at:
// https://api.slack.com/docs/verifying-requests-from-slack
func (rc *Router) verifyRequest(request *http.Request) error {
	signingSecret := rc.signingSecret
	if len(signingSecret) > 0 {
		sv, err := slack.NewSecretsVerifier(request.Header, signingSecret)
		if err != nil {
			return errors.Wrap(err, "cannot create secrets verifier")
		}

		// Read the request body
		body, err := getRequestBody(request)
		if err != nil {
			return err
		}

		_, err = sv.Write(body)
		if err != nil {
			return errors.Wrap(err, "error writing body: cannot verify signature")
		}

		err = sv.Ensure()
		if err != nil {
			return errors.Wrap(err, "signature validation failed")
		}
	}
	return nil
}

// StartListening starts an event source
func (el *EventListener) StartListening(ctx context.Context, dispatch func([]byte, ...eventsourcecommon.Options) error) error {
	log := logging.FromContext(ctx).
		With(logging.LabelEventSourceType, el.GetEventSourceType(), logging.LabelEventName, el.GetEventName())

	log.Info("started processing the Slack event source...")
	defer sources.Recover(el.GetEventName())

	slackEventSource := &el.SlackEventSource
	log.Info("retrieving the slack token...")
	token, err := common.GetSecretFromVolume(slackEventSource.Token)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the token")
	}

	log.Info("retrieving the signing secret...")
	signingSecret, err := common.GetSecretFromVolume(slackEventSource.SigningSecret)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the signing secret")
	}

	route := webhook.NewRoute(slackEventSource.Webhook, log, el.GetEventSourceName(), el.GetEventName(), el.Metrics)

	return webhook.ManageRoute(ctx, &Router{
		route:            route,
		token:            token,
		signingSecret:    signingSecret,
		slackEventSource: slackEventSource,
	}, controller, dispatch)
}
