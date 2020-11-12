package kisslogger

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc             string
		level            string
		expectedPriority uint
	}{
		{
			desc:             "should work for debug level",
			level:            "DEBUG",
			expectedPriority: 0,
		},
		{
			desc:             "should work for info level",
			level:            "INFO",
			expectedPriority: 1,
		},
		{
			desc:             "should work for warn level",
			level:            "WARN",
			expectedPriority: 2,
		},
		{
			desc:             "should work for error level",
			level:            "ERROR",
			expectedPriority: 3,
		},
		{
			desc:             "should default to info when input is unexpected",
			level:            "unexpected input",
			expectedPriority: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			instance := New(test.level)
			assert.Equal(t, test.expectedPriority, instance.priorityLevel)
		})
	}
}

func TestBuildJSONString(t *testing.T) {
	tests := []struct {
		desc           string
		level          string
		title          string
		body           Body
		expectedOutput map[string]interface{}
	}{
		{
			desc:  "should work with empty bodies",
			level: "fake-level",
			title: "fake-title",
			body:  Body{},
			expectedOutput: map[string]interface{}{
				"level": "fake-level",
				"title": "fake-title",
				// "timestamp": "can't compare timestamps",
			},
		},
		{
			desc:  "should work with non empty bodies",
			level: "fake-level",
			title: "fake-title",
			body: Body{
				"fake-key": "fake-timestamp",
			},
			expectedOutput: map[string]interface{}{
				"level": "fake-level",
				"title": "fake-title",
				// "timestamp": "can't compare timestamps",
				"fake-key": "fake-timestamp",
			},
		},
		{
			desc:  "should ignore reserved fields on body",
			level: "fake-level",
			title: "fake-title",
			body: Body{
				"level":     "fake-level2",
				"title":     "fake-title2",
				"timestamp": "fake-timestamp2",
			},
			expectedOutput: map[string]interface{}{
				"level": "fake-level",
				"title": "fake-title",
				// "timestamp": "can't compare timestamps",
			},
		},
		{
			desc:  "should output an error log when unable to marshal the body",
			level: "fake-level",
			title: "fake-title",
			body: Body{
				"value": CannotBeMarshaled{},
			},
			expectedOutput: map[string]interface{}{
				"level": "ERROR",
				"title": "could-not-marshal-log-body",
				"body":  fmt.Sprintf("%#v", Body{"value": CannotBeMarshaled{}}),
				// "timestamp": "can't compare timestamps",
			},
		},
	}

	for _, test := range tests {
		jsonString := buildJSONString(test.level, test.title, test.body)
		fmt.Println("String:", jsonString)

		var output map[string]interface{}
		err := json.Unmarshal([]byte(jsonString), &output)
		assert.Nil(t, err)

		timestring, ok := output["timestamp"].(string)
		assert.True(t, ok)

		_, err = time.Parse(time.RFC3339, timestring)
		assert.Nil(t, err)

		delete(output, "timestamp")
		assert.Equal(t, test.expectedOutput, output)
	}
}

type CannotBeMarshaled struct{}

func (c CannotBeMarshaled) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("fake-error-message")
}
