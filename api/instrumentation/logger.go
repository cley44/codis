package instrumentation

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"codis/config"

	"github.com/lmittmann/tint"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/samber/oops"
	slogmulti "github.com/samber/slog-multi"
	"golang.org/x/term"
)

func parseSlogLvl(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelWarn
	}
}

type LoggerService struct {
	injector do.Injector

	config *config.ConfigService
}

func NewLoggerService(injector do.Injector) (*LoggerService, error) {
	svc := LoggerService{
		injector: injector,

		config: do.MustInvoke[*config.ConfigService](injector),
	}

	svc.init()

	return &svc, nil
}

func (svc *LoggerService) init() {
	handlers := []slog.Handler{}

	// console logger
	handlers = append(
		handlers,
		NewOopsFormatter(
			true,
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      parseSlogLvl(svc.config.Instrumentation.LogLevel),
				TimeFormat: time.Kitchen,
				AddSource:  false,
				NoColor:    false,
				ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
					if attr.Key == "error" {
						return slog.String("error", attr.Value.String())
					}

					return attr
				},
			}),
		),
	)
	if term.IsTerminal(int(os.Stdout.Fd())) {
		handlers = append(
			handlers,
			slogmulti.NewHandleInlineHandler(
				func(ctx context.Context, groups []string, attrs []slog.Attr, record slog.Record) error {
					record.AddAttrs(attrs...)
					record.Attrs(func(attr slog.Attr) bool {
						if attr.Key == "error" {
							if eo, ok := lo.ErrorsAs[oops.OopsError](attr.Value.Any().(error)); ok {
								stacktrace := eo.Stacktrace()
								if len(stacktrace) > 0 {
									println(stacktrace)
								}
							}
							return false
						}

						return true
					})
					return nil
				},
			),
		)
	}

	recovery := slogmulti.RecoverHandlerError(
		func(ctx context.Context, record slog.Record, err error) {
			log.Println("recovered slog error:", err)
		},
	)

	slog.SetDefault(slog.New(recovery(slogmulti.Fanout(handlers...))))
}

func (svc *LoggerService) Shutdown() error {
	return nil
}
