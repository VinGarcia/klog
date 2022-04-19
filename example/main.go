package main

import (
	"context"
	"klog"
)

type myLogKey struct{}

func main() {
	ctx := context.TODO()
	logger := klog.New("INFO", func(ctx context.Context) klog.Body {
		body, _ := ctx.Value(myLogKey{}).(klog.Body)
		return body
	})

	logger.Debug(ctx, "testing-debug-wont-show-up")

	logger.Info(ctx, "testing-log")

	logger.Warn(ctx, "testing-log-with-values", klog.Body{
		"msg": "it worked!",
	})

	ctx = context.WithValue(ctx, myLogKey{}, klog.Body{
		"user_id": 41,
	})
	logger.Error(ctx, "testing-log-with-context")

	ctx = context.WithValue(ctx, myLogKey{}, klog.Body{
		"user_id":    42,
		"company_id": 22,
	})
	logger.Info(ctx, "testing-log-with-merged-context-values")

	logger.Info(ctx, "testing-log-precedence", klog.Body{"user_id": 43, "company_id": 23}, klog.Body{
		"user_id": 44,
	})
}
