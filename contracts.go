package klog

import "context"

// Provider describes the ways you can log information,
// depending on the log level some functions might just
// return without logging anything.
//
// The Fatal function will call `os.Exit(1)` after finishing
// writing the log.
type Provider interface {
	Debug(ctx context.Context, title string, valueMaps ...Body)
	Info(ctx context.Context, title string, valueMaps ...Body)
	Warn(ctx context.Context, title string, valueMaps ...Body)
	Error(ctx context.Context, title string, valueMaps ...Body)
	Fatal(ctx context.Context, title string, valueMaps ...Body)
}

// Body is the log body containing the keys and values
// used to build the structured logs
type Body = map[string]interface{}

// MiddlewareProvider describes the behavior of accepting
// middlewares that will be executed everytime a log function is called.
type MiddlewareProvider interface {
	AddBeforeEach(m Middleware)
	AddAfterEach(m Middleware)
}

// Middleware is a function that will be executed before or after
// an event is logged depending on whether it was registered
// with the AddBeforeEach or AddAfterEach method.
type Middleware func(*LogData)

// LogData represents all the data within a single log message
// and is used by the Middleware and OutputHandler functions.
type LogData struct {
	Level string
	Title string
	Body  Body
}
