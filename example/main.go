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
