package klog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// Client represents our logger instance.
type Client struct {
	priorityLevel uint
	PrintlnFn     func(...interface{})

	ctxParsers []ContextParser
}

// ContextParser is used for reading a log Body from the
// context allowing the user of the library to customize
// how to get this information.
type ContextParser func(ctx context.Context) Body

// New builds a logger Client on the appropriate log level
func New(level string, parsers ...ContextParser) Client {
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
		PrintlnFn: func(args ...interface{}) {
			fmt.Println(args...)
		},
		ctxParsers: parsers,
	}
}

// Debug logs an entry on level "DEBUG" with the received title
// along with all the values collected from the input valueMaps and the context.
func (c Client) Debug(ctx context.Context, title string, valueMaps ...Body) {
	if c.priorityLevel > 0 {
		return
	}

	c.log(ctx, "DEBUG", title, valueMaps)
}

// Info logs an entry on level "INFO" with the received title
// along with all the values collected from the input valueMaps and the context.
func (c Client) Info(ctx context.Context, title string, valueMaps ...Body) {
	if c.priorityLevel > 1 {
		return
	}

	c.log(ctx, "INFO", title, valueMaps)
}

// Warn logs an entry on level "WARN" with the received title
// along with all the values collected from the input valueMaps and the context.
func (c Client) Warn(ctx context.Context, title string, valueMaps ...Body) {
	if c.priorityLevel > 2 {
		return
	}

	c.log(ctx, "WARN", title, valueMaps)
}

// Error logs an entry on level "ERROR" with the received title
// along with all the values collected from the input valueMaps and the context.
func (c Client) Error(ctx context.Context, title string, valueMaps ...Body) {
	if c.priorityLevel > 3 {
		return
	}

	c.log(ctx, "ERROR", title, valueMaps)
}

// Fatal logs an entry on level "ERROR" with the received title
// along with all the values collected from the input valueMaps and the context.
//
// After that it proceeds to exit the program with code 1.
func (c Client) Fatal(ctx context.Context, title string, valueMaps ...Body) {
	if c.priorityLevel > 3 {
		return
	}

	c.log(ctx, "ERROR", title, valueMaps)
	os.Exit(1)
}

func (c Client) log(ctx context.Context, level string, title string, valueMaps []Body) {
	body := Body{}
	for _, parser := range c.ctxParsers {
		MergeMaps(&body, parser(ctx))
	}
	MergeMaps(&body, valueMaps...)

	c.PrintlnFn(buildJSONString(level, title, body))
}

func buildJSONString(level string, title string, body Body) string {
	timestamp := time.Now().Format(time.RFC3339)

	// Remove reserved keys from the input map:
	delete(body, "timestamp")
	delete(body, "level")
	delete(body, "title")

	values := []string{
		`"timestamp":"` + timestamp + `"`,
		`"level":"` + level + `"`,
		`"title":` + escapeAsJSON(title),
	}

	for k, v := range body {
		values = append(values, escapeAsJSON(k)+`:`+escapeAsJSON(v))
	}
	sort.Strings(values[3:])

	return fmt.Sprint("{" + strings.Join(values, ",") + "}")
}

func escapeAsJSON(obj interface{}) string {
	rawJSON, err := json.Marshal(obj)
	if err != nil {
		return escapeAsJSON(fmt.Sprintf("%#v", obj))
	}
	return string(rawJSON)
}
