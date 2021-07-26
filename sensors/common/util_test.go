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

package common

import (
	"fmt"
	"strings"
	"testing"

	"github.com/argoproj/argo-events/pkg/apis/sensor/v1alpha1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createFakeEvent(eventID string) *v1alpha1.Event {
	return &v1alpha1.Event{
		Context: &v1alpha1.EventContext{
			ID: eventID,
		},
	}
}

func TestApplyEventLabels(t *testing.T) {
	const maxK8sLabelLen int = 63
	const labelNamePrefix string = "events.argoproj.io/event-"

	t.Run("test event label name and value of a single event", func(t *testing.T) {
		uid := uuid.New()
		labelsMap := make(map[string]string)
		fakeEventsMap := map[string]*v1alpha1.Event{
			"test-dep": createFakeEvent(fmt.Sprintf("%x", uid)),
		}

		labelName := fmt.Sprintf("%s0", labelNamePrefix)
		err := ApplyEventLabels(labelsMap, fakeEventsMap)
		assert.Nil(t, err)
		assert.Equal(t, uid.String(), labelsMap[labelName])
		assert.True(t, len(labelsMap[labelName]) <= maxK8sLabelLen)
	})

	t.Run("test failure in case event ID isn't valid hex string", func(t *testing.T) {
		labelsMap := make(map[string]string)
		fakeEventsMap := map[string]*v1alpha1.Event{
			"test-dep": createFakeEvent("this is not hex"),
		}

		err := ApplyEventLabels(labelsMap, fakeEventsMap)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "failed to decode event ID"))
	})

	t.Run("test event labels prefix and value of multiple events", func(t *testing.T) {
		labelsMap := make(map[string]string)
		fakeEventsMap := map[string]*v1alpha1.Event{
			"test-dep1": createFakeEvent(fmt.Sprintf("%x", uuid.New())),
			"test-dep2": createFakeEvent(fmt.Sprintf("%x", uuid.New())),
		}

		err := ApplyEventLabels(labelsMap, fakeEventsMap)
		assert.Nil(t, err)
		assert.Equal(t, len(fakeEventsMap), len(labelsMap))

		for key, val := range labelsMap {
			assert.True(t, strings.HasPrefix(key, labelNamePrefix))
			_, err := uuid.Parse(val)
			assert.Nil(t, err)
			assert.True(t, len(val) <= maxK8sLabelLen)
		}
	})
}