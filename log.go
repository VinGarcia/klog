package kisslogger

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Body is the log body containing the keys and values
// used to build the structured logs
type Body map[string]interface{}

type Client struct {
	priorityLevel uint
	PrintlnFn     func(...interface{}) (int, error)
}

func New(level string) Client {
	var priority uint
	switch strings.ToUpper(level) {
	case "DEBUG":
		priority = 0
	case "INFO":
		priority = 1
	case "WARN":
		priority = 2
	case "ERROR":
		priority = 3
	default:
		priority = 1
	}

	return Client{
		priorityLevel: priority,
		PrintlnFn:     fmt.Println,
	}
}

func (c Client) Debug(ctx context.Context, title string, valueMaps ...Body) {
	if c.priorityLevel > 0 {
		return
	}

	logValues := mergeMapsUnsafe(Body{}, GetCtxValues(ctx))
	logValues = mergeMapsUnsafe(logValues, valueMaps...)

	c.PrintlnFn(buildJSONString("DEBUG", title, logValues))
}

func (c Client) Info(ctx context.Context, title string, valueMaps ...Body) {
	if c.priorityLevel > 1 {
		return
	}

	logValues := mergeMapsUnsafe(Body{}, GetCtxValues(ctx))
	logValues = mergeMapsUnsafe(logValues, valueMaps...)

	c.PrintlnFn(buildJSONString("INFO", title, logValues))
}

func (c Client) Warn(ctx context.Context, title string, valueMaps ...Body) {
	if c.priorityLevel > 2 {
		return
	}

	logValues := mergeMapsUnsafe(Body{}, GetCtxValues(ctx))
	logValues = mergeMapsUnsafe(logValues, valueMaps...)

	c.PrintlnFn(buildJSONString("WARN", title, logValues))
}

func (c Client) Error(ctx context.Context, title string, valueMaps ...Body) {
	if c.priorityLevel > 3 {
		return
	}

	logValues := mergeMapsUnsafe(Body{}, GetCtxValues(ctx))
	logValues = mergeMapsUnsafe(logValues, valueMaps...)

	c.PrintlnFn(buildJSONString("ERROR", title, logValues))
}

func buildJSONString(level string, title string, body Body) string {
	timestamp := time.Now().Format(time.RFC3339)

	// Remove reserved keys from the input map:
	delete(body, "level")
	delete(body, "title")
	delete(body, "timestamp")

	var prefix string
	var bodyJSON []byte = []byte("{}")
	var err error
	if len(body) > 0 {
		prefix = ","

		bodyJSON, err = json.Marshal(body)
		if err != nil {
			return fmt.Sprintf(
				`{"level":"ERROR","title":"could-not-marshal-log-body","body":"%#v","timestamp":"%s"}%s`,
				body,
				time.Now().Format(time.RFC3339),
				"\n",
			)
		}
	}

	titleJSON, _ := json.Marshal(title)

	return fmt.Sprintf(
		`{"level":"%s","title":%s,"timestamp":"%s"%s%s`,
		level,
		string(titleJSON),
		timestamp,
		prefix,
		string(bodyJSON[1:]),
	)
}
