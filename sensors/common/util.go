package common

import (
	"encoding/hex"

	"github.com/pkg/errors"

	"github.com/argoproj/argo-events/pkg/apis/sensor/v1alpha1"
)

func ApplyEventLabels(labels map[string]string, events map[string]*v1alpha1.Event) error {
	for key, val := range events {
		decodedID, err := hex.DecodeString(val.Context.ID)
		if err != nil {
			return errors.Wrap(err, "failed to decode event ID")
		}

		labelName := "events.argoproj.io/event-" + key
		labels[labelName] = string(decodedID)
	}

	return nil
}

