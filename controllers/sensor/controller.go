/*
Copyright 2020 BlackRock, Inc.

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

package sensor

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sort"
	"strings"

	"github.com/argoproj/argo-events/common"
	apicommon "github.com/argoproj/argo-events/pkg/apis/common"
	eventsourcev1alpha1 "github.com/argoproj/argo-events/pkg/apis/eventsource/v1alpha1"
	sensorv1alpha1 "github.com/argoproj/argo-events/pkg/apis/sensor/v1alpha1"
)

const (
	// ControllerName is name of the controller
	ControllerName = "sensor-controller"

	finalizerName = ControllerName
)

type eventSourceEvent struct {
	eventName       string
	eventSourceName string
}

func newEventSourceEvent(eventName, eventSourceName string) *eventSourceEvent {
	return &eventSourceEvent{eventName: eventName, eventSourceName: eventSourceName}
}

type reconciler struct {
	client client.Client
	scheme *runtime.Scheme

	sensorImage string
	logger      *zap.SugaredLogger
}

// NewReconciler returns a new reconciler
func NewReconciler(client client.Client, scheme *runtime.Scheme, sensorImage string, logger *zap.SugaredLogger) reconcile.Reconciler {
	return &reconciler{client: client, scheme: scheme, sensorImage: sensorImage, logger: logger}
}

func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	sensor := &sensorv1alpha1.Sensor{}
	if err := r.client.Get(ctx, req.NamespacedName, sensor); err != nil {
		if apierrors.IsNotFound(err) {
			r.logger.Warnw("WARNING: sensor not found", "request", req)
			return reconcile.Result{}, nil
		}
		r.logger.Errorw("unable to get sensor ctl", zap.Any("request", req), zap.Error(err))
		return ctrl.Result{}, err
	}
	log := r.logger.With("namespace", sensor.Namespace).With("sensor", sensor.Name)
	sensorCopy := sensor.DeepCopy()
	reconcileErr := r.reconcile(ctx, sensorCopy)
	if reconcileErr != nil {
		log.Errorw("reconcile error", zap.Error(reconcileErr))
	}
	if r.needsUpdate(sensor, sensorCopy) {
		if err := r.client.Update(ctx, sensorCopy); err != nil {
			return reconcile.Result{}, err
		}
	}
	if err := r.client.Status().Update(ctx, sensorCopy); err != nil {
		return reconcile.Result{}, err
	}
	return ctrl.Result{}, reconcileErr
}

// reconcile does the real logic
func (r *reconciler) reconcile(ctx context.Context, sensor *sensorv1alpha1.Sensor) error {
	log := r.logger.With("namespace", sensor.Namespace).With("sensor", sensor.Name)
	if !sensor.DeletionTimestamp.IsZero() {
		log.Info("deleting sensor")
		if controllerutil.ContainsFinalizer(sensor, finalizerName) {
			// Finalizer logic should be added here.
			controllerutil.RemoveFinalizer(sensor, finalizerName)
		}
		return nil
	}
	controllerutil.AddFinalizer(sensor, finalizerName)

	sensor.Status.InitConditions()
	if err := ValidateSensor(sensor); err != nil {
		log.Errorw("validation error", "error", err)
		return err
	}
	sensorCopy := sensor.DeepCopy()
	err := r.recreateDependencies(ctx, sensor)
	if err != nil {
		log.Errorw("failed to map dependencies", "error", err)
		return err
	}
	if r.needsValidation(sensorCopy, sensor) {
		if err := ValidateSensor(sensor); err != nil {
			log.Errorw("validation error", "error", err)
			return err
		}
	}
	args := &AdaptorArgs{
		Image:  r.sensorImage,
		Sensor: sensor,
		Labels: map[string]string{
			"controller":           "sensor-controller",
			common.LabelSensorName: sensor.Name,
			common.LabelOwnerName:  sensor.Name,
		},
	}
	return Reconcile(r.client, args, log)
}

func (r *reconciler) needsUpdate(old, new *sensorv1alpha1.Sensor) bool {
	if old == nil {
		return true
	}
	return !equality.Semantic.DeepEqual(old.Finalizers, new.Finalizers)
}

func (r *reconciler) needsValidation(old, new *sensorv1alpha1.Sensor) bool {
	if old == nil {
		return true
	}
	return !equality.Semantic.DeepEqual(old.Spec.Dependencies, new.Spec.Dependencies)
}

func (r *reconciler) recreateDependencies(ctx context.Context, sensor *sensorv1alpha1.Sensor) error {
	currDeps := sensor.Spec.Dependencies
	newDeps := make([]sensorv1alpha1.EventDependency, 0, len(currDeps))
	var filterDeps []sensorv1alpha1.EventDependency
	for _, dep := range currDeps {
		eventSourceType := dep.EventSourceFilter
		if len(eventSourceType) != 0 {
			filterDeps = append(filterDeps, dep)
		} else {
			newDeps = append(newDeps, dep)
		}
	}

	mappedDeps, err := r.mapFilterDependenciesToRegularDependencies(ctx, sensor.Namespace, filterDeps)
	if err != nil {
		return err
	}

	newDeps = append(newDeps, mappedDeps...)
	sensor.Spec.Dependencies = newDeps

	return nil
}

func (r *reconciler) mapFilterDependenciesToRegularDependencies(ctx context.Context, namespace string, filterDeps []sensorv1alpha1.EventDependency) ([]sensorv1alpha1.EventDependency, error) {
	esList := &eventsourcev1alpha1.EventSourceList{}
	err := r.client.List(ctx, esList, &client.ListOptions{
		Namespace: namespace,
	})
	if err != nil {
		return nil, err
	}

	esEventsByType := make(map[apicommon.EventSourceType][]*eventSourceEvent)
	for _, es := range esList.Items {
		r.applyEventSourceEventsGroupedByTypes(&es, esEventsByType)
	}

	r.sortDependenciesByName(filterDeps)
	resultDeps := make([]sensorv1alpha1.EventDependency, 0, len(esList.Items))
	for _, fd := range filterDeps {
		esType := fd.EventSourceFilter
		events := esEventsByType[esType]
		r.sortEventSourceEvents(events)
		for i, e := range events {
			resultDeps = append(resultDeps, sensorv1alpha1.EventDependency{
				Name:            fmt.Sprintf("%s-%d", fd.Name, i),
				EventName:       e.eventName,
				EventSourceName: e.eventSourceName,
			})
		}
	}

	return resultDeps, nil
}

func (r *reconciler) applyEventSourceEventsGroupedByTypes(eventSource *eventsourcev1alpha1.EventSource, eventNamesMap map[apicommon.EventSourceType][]*eventSourceEvent) {
	esName := eventSource.Name
	if len(eventSource.Spec.AMQP) != 0 {
		for eventName := range eventSource.Spec.AMQP {
			eventNamesMap[apicommon.AMQPEvent] = append(eventNamesMap[apicommon.AMQPEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.AzureEventsHub) != 0 {
		for eventName := range eventSource.Spec.AzureEventsHub {
			eventNamesMap[apicommon.AzureEventsHub] = append(eventNamesMap[apicommon.AzureEventsHub], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Calendar) != 0 {
		for eventName := range eventSource.Spec.Calendar {
			eventNamesMap[apicommon.CalendarEvent] = append(eventNamesMap[apicommon.CalendarEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Emitter) != 0 {
		for eventName := range eventSource.Spec.Emitter {
			eventNamesMap[apicommon.EmitterEvent] = append(eventNamesMap[apicommon.EmitterEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.File) != 0 {
		for eventName := range eventSource.Spec.File {
			eventNamesMap[apicommon.FileEvent] = append(eventNamesMap[apicommon.FileEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Github) != 0 {
		for eventName := range eventSource.Spec.Github {
			eventNamesMap[apicommon.GithubEvent] = append(eventNamesMap[apicommon.GithubEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Gitlab) != 0 {
		for eventName := range eventSource.Spec.Gitlab {
			eventNamesMap[apicommon.GitlabEvent] = append(eventNamesMap[apicommon.GitlabEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.HDFS) != 0 {
		for eventName := range eventSource.Spec.HDFS {
			eventNamesMap[apicommon.HDFSEvent] = append(eventNamesMap[apicommon.HDFSEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Kafka) != 0 {
		for eventName := range eventSource.Spec.Kafka {
			eventNamesMap[apicommon.KafkaEvent] = append(eventNamesMap[apicommon.KafkaEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.MQTT) != 0 {
		for eventName := range eventSource.Spec.MQTT {
			eventNamesMap[apicommon.MQTTEvent] = append(eventNamesMap[apicommon.MQTTEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Minio) != 0 {
		for eventName := range eventSource.Spec.Minio {
			eventNamesMap[apicommon.MinioEvent] = append(eventNamesMap[apicommon.MinioEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.NATS) != 0 {
		for eventName := range eventSource.Spec.NATS {
			eventNamesMap[apicommon.NATSEvent] = append(eventNamesMap[apicommon.NATSEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.NSQ) != 0 {
		for eventName := range eventSource.Spec.NSQ {
			eventNamesMap[apicommon.NSQEvent] = append(eventNamesMap[apicommon.NSQEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.PubSub) != 0 {
		for eventName := range eventSource.Spec.PubSub {
			eventNamesMap[apicommon.PubSubEvent] = append(eventNamesMap[apicommon.PubSubEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Redis) != 0 {
		for eventName := range eventSource.Spec.Redis {
			eventNamesMap[apicommon.RedisEvent] = append(eventNamesMap[apicommon.RedisEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.SNS) != 0 {
		for eventName := range eventSource.Spec.SNS {
			eventNamesMap[apicommon.SNSEvent] = append(eventNamesMap[apicommon.SNSEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.SQS) != 0 {
		for eventName := range eventSource.Spec.SQS {
			eventNamesMap[apicommon.SQSEvent] = append(eventNamesMap[apicommon.SQSEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Slack) != 0 {
		for eventName := range eventSource.Spec.Slack {
			eventNamesMap[apicommon.SlackEvent] = append(eventNamesMap[apicommon.SlackEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.StorageGrid) != 0 {
		for eventName := range eventSource.Spec.StorageGrid {
			eventNamesMap[apicommon.StorageGridEvent] = append(eventNamesMap[apicommon.StorageGridEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Stripe) != 0 {
		for eventName := range eventSource.Spec.Stripe {
			eventNamesMap[apicommon.StripeEvent] = append(eventNamesMap[apicommon.StripeEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Webhook) != 0 {
		for eventName := range eventSource.Spec.Webhook {
			eventNamesMap[apicommon.WebhookEvent] = append(eventNamesMap[apicommon.WebhookEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Resource) != 0 {
		for eventName := range eventSource.Spec.Resource {
			eventNamesMap[apicommon.ResourceEvent] = append(eventNamesMap[apicommon.ResourceEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Pulsar) != 0 {
		for eventName := range eventSource.Spec.Pulsar {
			eventNamesMap[apicommon.PulsarEvent] = append(eventNamesMap[apicommon.PulsarEvent], newEventSourceEvent(eventName, esName))
		}
	}
	if len(eventSource.Spec.Generic) != 0 {
		for eventName := range eventSource.Spec.Generic {
			eventNamesMap[apicommon.GenericEvent] = append(eventNamesMap[apicommon.GenericEvent], newEventSourceEvent(eventName, esName))
		}
	}
}

func (r *reconciler) sortDependenciesByName(deps []sensorv1alpha1.EventDependency) {
	sort.Slice(deps, func(i, j int) bool {
		return strings.Compare(deps[i].Name, deps[j].Name) < 0
	})
}

func (r *reconciler) sortEventSourceEvents(events []*eventSourceEvent) {
	sort.Slice(events, func(i, j int) bool {
		comboKey1 := fmt.Sprintf("%s-%s", events[i].eventSourceName, events[i].eventName)
		comboKey2 := fmt.Sprintf("%s-%s", events[j].eventSourceName, events[j].eventName)
		return strings.Compare(comboKey1, comboKey2) < 0
	})
}