package slogerr

import (
	"context"
	"log/slog"
)

const (
	// ErrorKey is the key used to identify errors in slog records.
	ErrorKey = "error"
)

type ErrorNilDropper struct {
	slog.Handler
}

func NewErrorNilDropperHandler(h slog.Handler) *ErrorNilDropper {
	return &ErrorNilDropper{Handler: h}
}

// Handle processes the slog record and drops it if the error is nil.
// The error is identified by the ErrorKey attribute from this package.
func (h *ErrorNilDropper) Handle(ctx context.Context, r slog.Record) error {
	if r.Level != slog.LevelError {
		return h.Handler.Handle(ctx, r)
	}

	isNilError := false
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == ErrorKey &&
			a.Value.Kind() == slog.KindAny &&
			a.Value.Any() == nil {
			isNilError = true
			return false
		}
		return true
	})
	if isNilError {
		return nil
	}

	return h.Handler.Handle(ctx, r)
}

func (h *ErrorNilDropper) WithAttrs(as []slog.Attr) slog.Handler {
	return &ErrorNilDropper{Handler: h.Handler.WithAttrs(as)}
}

func (h *ErrorNilDropper) WithGroup(name string) slog.Handler {
	return &ErrorNilDropper{Handler: h.Handler.WithGroup(name)}
}
