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
