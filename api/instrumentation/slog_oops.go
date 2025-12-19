package instrumentation

import (
	"context"
	"log/slog"

	"github.com/samber/lo"
	"github.com/samber/oops"
)

func NewOopsFormatter(ignoreStacktrace bool, secondaryFormatter slog.Handler) *OopsFormatter {
	if secondaryFormatter == nil {
		secondaryFormatter = slog.Default().Handler()
	}

	return &OopsFormatter{
		formatter:        secondaryFormatter,
		ignoreStacktrace: ignoreStacktrace,
		ignoreSource:     true,
	}
}

type OopsFormatter struct {
	formatter        slog.Handler
	ignoreStacktrace bool
	ignoreSource     bool
}

func (f *OopsFormatter) clone() *OopsFormatter {
	return &OopsFormatter{
		formatter:        f.formatter,
		ignoreStacktrace: f.ignoreStacktrace,
		ignoreSource:     f.ignoreSource,
	}
}

func (f *OopsFormatter) Enabled(ctx context.Context, level slog.Level) bool {
	return f.formatter.Enabled(ctx, level)
}

func (f *OopsFormatter) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return f
	}

	f2 := f.clone()
	f2.formatter = f.formatter.WithAttrs(attrs)

	return f2
}

func (f *OopsFormatter) WithGroup(name string) slog.Handler {
	f2 := f.clone()
	f2.formatter = f.formatter.WithGroup(name)

	return f2
}

func (f *OopsFormatter) Handle(ctx context.Context, entry slog.Record) error {
	entry.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "error" {
			switch attr.Value.Kind() {
			case slog.KindAny:
				if oopsError, ok := lo.ErrorsAs[oops.OopsError](attr.Value.Any().(error)); ok {
					oopsErrorToEntryData(&oopsError, &entry)

					if !f.ignoreStacktrace {
						stacktrace := oopsError.Stacktrace()
						if len(stacktrace) > 0 {
							entry.Add(slog.String("stacktrace", stacktrace))
						}
					}

					if !f.ignoreSource {
						sources := oopsError.Sources()
						if len(sources) > 0 {
							entry.Add(slog.String("sources", sources))
						}
					}
				}
			default:
			}

			return false
		}

		return true
	})

	err := f.formatter.Handle(ctx, entry)
	if err != nil {
		return err
	}

	return nil
}

func oopsErrorToEntryData(err *oops.OopsError, entry *slog.Record) {
	entry.Time = err.Time()

	payload := err.ToMap()

	delete(payload, "stacktrace")
	delete(payload, "sources")

	for k, v := range payload {
		entry.Add(slog.Any(k, v))
	}
}
