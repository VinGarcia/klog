# Welcome to KLog

KLog stands for Keep it simple Logging Package.

It's a very simple and powerful package making logging
easy to write and easy to read.

Don't worry about building error messages that contain all
important information, write only an "error title" and keep
the important information as data, which is simpler and more
expressive.

KLog also supports saving log data in the context.

This is very useful for making sure that some important
information that should show up on every log such as
the `request_id` or an username will indeed show up without
the programmer having to explicity pass this value on every
call.

## The interface

KLog aims to be simple and so its interface is very simple and
we don't plan on adding new functions to it if we can avoid it:

```golang
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
```

## Example usage

```golang
package main

import (
	"context"
	"klog"
)

func main() {
	ctx := context.TODO()
	logger := klog.New("INFO")

	logger.Debug(ctx, "testing-debug-wont-show-up")

	logger.Info(ctx, "testing-log")

	logger.Warn(ctx, "testing-log-with-values", klog.Body{
		"msg": "it worked!",
	})

	ctx = klog.CtxWithValues(ctx, klog.Body{
		"user_id": 41,
	})
	logger.Error(ctx, "testing-log-with-context")

	ctx = klog.CtxWithValues(ctx, klog.Body{
		"user_id":    42,
		"company_id": 22,
	})
	logger.Info(ctx, "testing-log-with-merged-context-values")

	logger.Info(ctx, "testing-log-precedence", klog.Body{"user_id": 43, "company_id": 23}, klog.Body{
		"user_id": 44,
	})
}
```
