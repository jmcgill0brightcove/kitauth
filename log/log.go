package log

import (
	levlog "github.com/go-kit/kit/log/levels"
	"golang.org/x/net/context"
)

const (
	LogTimestamp      = "ts"
	ContextBaseLogger = "kitauth/context/logger"
	ContextLogger     = "kitauth/context/levels"
)

func Crit(ctx context.Context, keyvals ...interface{}) {
	if l, ok := ctx.Value(ContextLogger).(levlog.Levels); ok {
		l.Crit().Log(keyvals...)
	}
}

func Debug(ctx context.Context, keyvals ...interface{}) {
	if l, ok := ctx.Value(ContextLogger).(levlog.Levels); ok {
		l.Debug().Log(keyvals...)
	}
}

func Error(ctx context.Context, keyvals ...interface{}) {
	if l, ok := ctx.Value(ContextLogger).(levlog.Levels); ok {
		l.Error().Log(keyvals...)
	}
}

func Info(ctx context.Context, keyvals ...interface{}) {
	if l, ok := ctx.Value(ContextLogger).(levlog.Levels); ok {
		l.Info().Log(keyvals...)
	}
}

func Warn(ctx context.Context, keyvals ...interface{}) {
	if l, ok := ctx.Value(ContextLogger).(levlog.Levels); ok {
		l.Warn().Log(keyvals...)
	}
}
