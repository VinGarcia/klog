package klog

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
				"level": "fake-level",
				"title": "fake-title",
				"value": fmt.Sprintf("%#v", CannotBeMarshaled{}),
				// "timestamp": "can't compare timestamps",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			jsonString := buildJSONString(test.level, test.title, test.body)

			var output map[string]interface{}
			err := json.Unmarshal([]byte(jsonString), &output)
			assert.Nil(t, err)

			timestring, ok := output["timestamp"].(string)
			assert.True(t, ok)

			_, err = time.Parse(time.RFC3339, timestring)
			assert.Nil(t, err)

			delete(output, "timestamp")
			assert.Equal(t, test.expectedOutput, output)
		})
	}
}

func TestLogFuncs(t *testing.T) {
	t.Run("debug logs should produce logs correctly", func(t *testing.T) {
		var output string
		ctx := context.TODO()
		client := Client{
			priorityLevel: 0,
			OutputHandler: func(level string, title string, body Body) {
				output = buildJSONString(level, title, body)
			},
			ctxParsers: []ContextParser{getCtxValues},
		}

		ctx = ctxWithValues(ctx, Body{
			"ctx_value1": "overwritten",
			"ctx_value2": "overwritten",
			"ctx_value3": "not-overwritten",
		})

		client.Debug(
			ctx,
			"fake-log-title",
			Body{
				"ctx_value1":   "overwrites",
				"body1_value1": "overwritten",
				"body1_value2": "not-overwritten",
			},
			Body{
				"ctx_value2":   "overwrites",
				"body1_value1": "overwrites",
				"body2_value1": "not-overwritten",
			},
		)

		// Default log values:
		assert.True(t, strings.Contains(output, `"level":"DEBUG"`))
		assert.True(t, strings.Contains(output, `"title":"fake-log-title"`))
		assert.True(t, strings.Contains(output, `"timestamp"`))

		// Final contextual values:
		assert.True(t, strings.Contains(output, `"ctx_value1":"overwrites"`))
		assert.True(t, strings.Contains(output, `"ctx_value2":"overwrites"`))
		assert.True(t, strings.Contains(output, `"body1_value1":"overwrites"`))
		assert.True(t, strings.Contains(output, `"ctx_value3":"not-overwritten"`))
		assert.True(t, strings.Contains(output, `"body1_value2":"not-overwritten"`))
		assert.True(t, strings.Contains(output, `"body2_value1":"not-overwritten"`))

		// No overwritten field should be present:
		assert.True(t, !strings.Contains(output, `"overwritten"`))
	})

	t.Run("debug logs should be ignored if priorityLevel > 0", func(t *testing.T) {
		var output string
		ctx := context.TODO()
		client := Client{
			priorityLevel: 1,
			OutputHandler: func(level string, title string, body Body) {
				output = buildJSONString(level, title, body)
			},
			ctxParsers: []ContextParser{getCtxValues},
		}

		client.Debug(ctx, "fake-log-title")

		assert.Equal(t, "", output)
	})

	t.Run("info logs should produce logs correctly", func(t *testing.T) {
		var output string
		ctx := context.TODO()
		client := Client{
			priorityLevel: 0,
			OutputHandler: func(level string, title string, body Body) {
				output = buildJSONString(level, title, body)
			},
			ctxParsers: []ContextParser{getCtxValues},
		}

		ctx = ctxWithValues(ctx, Body{
			"ctx_value1": "overwritten",
			"ctx_value2": "overwritten",
			"ctx_value3": "not-overwritten",
		})

		client.Info(
			ctx,
			"fake-log-title",
			Body{
				"ctx_value1":   "overwrites",
				"body1_value1": "overwritten",
				"body1_value2": "not-overwritten",
			},
			Body{
				"ctx_value2":   "overwrites",
				"body1_value1": "overwrites",
				"body2_value1": "not-overwritten",
			},
		)

		// Default log values:
		assert.True(t, strings.Contains(output, `"level":"INFO"`))
		assert.True(t, strings.Contains(output, `"title":"fake-log-title"`))
		assert.True(t, strings.Contains(output, `"timestamp"`))

		// Final contextual values:
		assert.True(t, strings.Contains(output, `"ctx_value1":"overwrites"`))
		assert.True(t, strings.Contains(output, `"ctx_value2":"overwrites"`))
		assert.True(t, strings.Contains(output, `"body1_value1":"overwrites"`))
		assert.True(t, strings.Contains(output, `"ctx_value3":"not-overwritten"`))
		assert.True(t, strings.Contains(output, `"body1_value2":"not-overwritten"`))
		assert.True(t, strings.Contains(output, `"body2_value1":"not-overwritten"`))

		// No overwritten field should be present:
		assert.True(t, !strings.Contains(output, `"overwritten"`))
	})

	t.Run("info logs should be ignored if priorityLevel > 1", func(t *testing.T) {
		var output string
		ctx := context.TODO()
		client := Client{
			priorityLevel: 2,
			OutputHandler: func(level string, title string, body Body) {
				output = buildJSONString(level, title, body)
			},
			ctxParsers: []ContextParser{getCtxValues},
		}

		client.Info(ctx, "fake-log-title")

		assert.Equal(t, "", output)
	})

	t.Run("warn logs should produce logs correctly", func(t *testing.T) {
		var output string
		ctx := context.TODO()
		client := Client{
			priorityLevel: 0,
			OutputHandler: func(level string, title string, body Body) {
				output = buildJSONString(level, title, body)
			},
			ctxParsers: []ContextParser{getCtxValues},
		}

		ctx = ctxWithValues(ctx, Body{
			"ctx_value1": "overwritten",
			"ctx_value2": "overwritten",
			"ctx_value3": "not-overwritten",
		})

		client.Warn(
			ctx,
			"fake-log-title",
			Body{
				"ctx_value1":   "overwrites",
				"body1_value1": "overwritten",
				"body1_value2": "not-overwritten",
			},
			Body{
				"ctx_value2":   "overwrites",
				"body1_value1": "overwrites",
				"body2_value1": "not-overwritten",
			},
		)

		// Default log values:
		assert.True(t, strings.Contains(output, `"level":"WARN"`))
		assert.True(t, strings.Contains(output, `"title":"fake-log-title"`))
		assert.True(t, strings.Contains(output, `"timestamp"`))

		// Final contextual values:
		assert.True(t, strings.Contains(output, `"ctx_value1":"overwrites"`))
		assert.True(t, strings.Contains(output, `"ctx_value2":"overwrites"`))
		assert.True(t, strings.Contains(output, `"body1_value1":"overwrites"`))
		assert.True(t, strings.Contains(output, `"ctx_value3":"not-overwritten"`))
		assert.True(t, strings.Contains(output, `"body1_value2":"not-overwritten"`))
		assert.True(t, strings.Contains(output, `"body2_value1":"not-overwritten"`))

		// No overwritten field should be present:
		assert.True(t, !strings.Contains(output, `"overwritten"`))
	})

	t.Run("warn logs should be ignored if priorityLevel > 2", func(t *testing.T) {
		var output string
		ctx := context.TODO()
		client := Client{
			priorityLevel: 3,
			OutputHandler: func(level string, title string, body Body) {
				output = buildJSONString(level, title, body)
			},
			ctxParsers: []ContextParser{getCtxValues},
		}

		client.Warn(ctx, "fake-log-title")

		assert.Equal(t, "", output)
	})

	t.Run("error logs should produce logs correctly", func(t *testing.T) {
		var output string
		ctx := context.TODO()
		client := Client{
			priorityLevel: 0,
			OutputHandler: func(level string, title string, body Body) {
				output = buildJSONString(level, title, body)
			},
			ctxParsers: []ContextParser{getCtxValues},
		}

		ctx = ctxWithValues(ctx, Body{
			"ctx_value1": "overwritten",
			"ctx_value2": "overwritten",
			"ctx_value3": "not-overwritten",
		})

		client.Error(
			ctx,
			"fake-log-title",
			Body{
				"ctx_value1":   "overwrites",
				"body1_value1": "overwritten",
				"body1_value2": "not-overwritten",
			},
			Body{
				"ctx_value2":   "overwrites",
				"body1_value1": "overwrites",
				"body2_value1": "not-overwritten",
			},
		)

		// Default log values:
		assert.True(t, strings.Contains(output, `"level":"ERROR"`))
		assert.True(t, strings.Contains(output, `"title":"fake-log-title"`))
		assert.True(t, strings.Contains(output, `"timestamp"`))

		// Final contextual values:
		assert.True(t, strings.Contains(output, `"ctx_value1":"overwrites"`))
		assert.True(t, strings.Contains(output, `"ctx_value2":"overwrites"`))
		assert.True(t, strings.Contains(output, `"body1_value1":"overwrites"`))
		assert.True(t, strings.Contains(output, `"ctx_value3":"not-overwritten"`))
		assert.True(t, strings.Contains(output, `"body1_value2":"not-overwritten"`))
		assert.True(t, strings.Contains(output, `"body2_value1":"not-overwritten"`))

		// No overwritten field should be present:
		assert.True(t, !strings.Contains(output, `"overwritten"`))
	})

	t.Run("error logs should be ignored if priorityLevel > 3", func(t *testing.T) {
		var output string
		ctx := context.TODO()
		client := Client{
			priorityLevel: 4,
			OutputHandler: func(level string, title string, body Body) {
				output = buildJSONString(level, title, body)
			},
			ctxParsers: []ContextParser{getCtxValues},
		}

		client.Error(ctx, "fake-log-title")

		assert.Equal(t, "", output)
	})

	t.Run("should run the middlewares when a log function is called", func(t *testing.T) {
		ctx := context.TODO()

		client := Client{
			priorityLevel: 1, // Should not run the debug log
			OutputHandler: func(level string, title string, body Body) {},
		}

		var args []map[string]any
		client.AddMiddleware(func(level string, title string, body Body) {
			args = append(args, map[string]any{
				"level": level,
				"title": title,
				"body":  body,
			})
		})

		client.Debug(ctx, "debug-log-title", Body{"debugLog": "won't appear"})
		client.Info(ctx, "info-log-title", Body{"infoLog": "value1"})
		client.Warn(ctx, "warn-log-title", Body{"warnLog": "value2"})
		client.Error(ctx, "error-log-title",
			Body{"errorLog": "overriden"},
			Body{"errorLog": "value3", "other": "value4"},
		)

		assert.Equal(t, args, []map[string]any{
			map[string]any{
				"level": "INFO",
				"title": "info-log-title",
				"body": Body{
					"infoLog": "value1",
				},
			},
			map[string]any{
				"level": "WARN",
				"title": "warn-log-title",
				"body": Body{
					"warnLog": "value2",
				},
			},
			map[string]any{
				"level": "ERROR",
				"title": "error-log-title",
				"body": Body{
					"errorLog": "value3",
					"other":    "value4",
				},
			},
		})
	})

	t.Run("should run the middlewares in the correct order", func(t *testing.T) {
		ctx := context.TODO()

		client := Client{
			priorityLevel: 0,
			OutputHandler: func(level string, title string, body Body) {},
		}

		var args []string
		client.AddMiddleware(func(level string, title string, body Body) {
			args = append(args, "m1")
		})
		client.AddMiddleware(func(level string, title string, body Body) {
			args = append(args, "m2")
		})
		client.AddMiddleware(func(level string, title string, body Body) {
			args = append(args, "m3")
		})

		client.Debug(ctx, "log-title")

		assert.Equal(t, args, []string{"m1", "m2", "m3"})
	})
}

type CannotBeMarshaled struct{}

func (c CannotBeMarshaled) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("fake-error-message")
}

type MyLogKey struct{}

func ctxWithValues(ctx context.Context, values Body) context.Context {
	m, _ := ctx.Value(MyLogKey{}).(Body)
	MergeMaps(&m, values)
	return context.WithValue(ctx, MyLogKey{}, m)
}

func getCtxValues(ctx context.Context) Body {
	m, _ := ctx.Value(MyLogKey{}).(Body)
	if m == nil {
		return Body{}
	}
	return m
}
