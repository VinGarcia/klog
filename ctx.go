package klog

import (
	"context"
)

// Declaring a unique private type for the ctx key
// guarantees that no key colision will ever happen:
type logCtxKeyType uint8

var logCtxKey logCtxKeyType

func CtxWithValues(ctx context.Context, values Body) context.Context {
	m, _ := ctx.Value(logCtxKey).(Body)
	return context.WithValue(ctx, logCtxKey, mergeMaps(m, values))
}

func GetCtxValues(ctx context.Context) Body {
	m, _ := ctx.Value(logCtxKey).(Body)
	if m == nil {
		return Body{}
	}
	return m
}

func mergeMaps(maps ...Body) Body {
	return mergeMapsUnsafe(Body{}, maps...)
}

func mergeMapsUnsafe(baseMap Body, maps ...Body) Body {
	for _, m := range maps {
		for k, v := range m {
			baseMap[k] = v
		}
	}

	return baseMap
}
