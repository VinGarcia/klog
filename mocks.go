package klog

import "context"

// Mock mocks the Provider interface
type Mock struct {
	DebugFn func(ctx context.Context, title string, valueMaps ...Body)
	InfoFn  func(ctx context.Context, title string, valueMaps ...Body)
	WarnFn  func(ctx context.Context, title string, valueMaps ...Body)
	ErrorFn func(ctx context.Context, title string, valueMaps ...Body)
	FatalFn func(ctx context.Context, title string, valueMaps ...Body)
}

// Debug implements the klog.Provider interface
func (m Mock) Debug(ctx context.Context, title string, valueMaps ...Body) {
	if m.DebugFn != nil {
		m.DebugFn(ctx, title, valueMaps...)
	}
}

// Info implements the klog.Provider interface
func (m Mock) Info(ctx context.Context, title string, valueMaps ...Body) {
	if m.InfoFn != nil {
		m.InfoFn(ctx, title, valueMaps...)
	}
}

// Warn implements the klog.Provider interface
func (m Mock) Warn(ctx context.Context, title string, valueMaps ...Body) {
	if m.WarnFn != nil {
		m.WarnFn(ctx, title, valueMaps...)
	}
}

// Error implements the klog.Provider interface
func (m Mock) Error(ctx context.Context, title string, valueMaps ...Body) {
	if m.ErrorFn != nil {
		m.ErrorFn(ctx, title, valueMaps...)
	}
}

// Fatal implements the klog.Provider interface
func (m Mock) Fatal(ctx context.Context, title string, valueMaps ...Body) {
	if m.FatalFn != nil {
		m.FatalFn(ctx, title, valueMaps...)
	}
}
