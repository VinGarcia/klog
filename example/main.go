package main

import (
	"context"
	log "klog"
)

func main() {
	ctx := context.TODO()
	logger := log.New("INFO")

	logger.Debug(ctx, "testing-debug-wont-show-up")

	logger.Info(ctx, "testing-log")

	logger.Warn(ctx, "testing-log-with-values", log.Body{
		"msg": "it worked!",
	})

	ctx = log.CtxWithValues(ctx, log.Body{
		"user_id": 41,
	})
	logger.Error(ctx, "testing-log-with-context")

	ctx = log.CtxWithValues(ctx, log.Body{
		"user_id":    42,
		"company_id": 22,
	})
	logger.Info(ctx, "testing-log-with-merged-context-values")

	logger.Info(ctx, "testing-log-precedence", log.Body{"user_id": 43, "company_id": 23}, log.Body{
		"user_id": 44,
	})
}
