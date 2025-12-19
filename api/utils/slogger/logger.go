package slogger

import (
	"fmt"
	"log"
	"log/slog"
)

func formatText(msg string, args ...any) (string, []any) {
	// Count the number of format specifiers in msg
	numSpecifiers := 0
	for i := 0; i < len(msg); i++ {
		if msg[i] == '%' && i+1 < len(msg) && msg[i+1] != '%' {
			numSpecifiers++
		}
	}

	msg = fmt.Sprintf(msg, args[:numSpecifiers]...)

	// Remove the necessary parameters from the start of args
	if len(args) >= numSpecifiers {
		args = args[numSpecifiers:]
	}

	return msg, args
}

func Info(msg string, args ...any) {
	msg, args = formatText(msg, args...)

	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	msg, args = formatText(msg, args...)

	slog.Warn(msg, args...)
}

func Error(err error, args ...any) {
	args = append([]any{WithError(err)}, args...)
	slog.Error(err.Error(), args...)
}

func Errorf(msg string, args ...any) {
	msg, args = formatText(msg, args...)

	slog.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	msg, args = formatText(msg, args...)

	slog.Debug(msg, args...)
}

func Fatal(err error, args ...any) {
	log.Fatal(err)
}

func WithError(err error) slog.Attr {
	return slog.Any("error", err)
}
