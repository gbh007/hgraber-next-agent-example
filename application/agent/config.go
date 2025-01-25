package agent

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/gbh007/hgraber-next-agent-core/config"

	"go.opentelemetry.io/otel/trace"
)

func parseConfig[T any]() (config.Config[T], bool, error) {
	configPath := flag.String("config", "config.yaml", "path to config")
	generateConfig := flag.String("generate-config", "", "generate example config")
	scan := flag.Bool("scan", false, "scan zip file to register in db")
	flag.Parse()

	defaultParsers := func() *T { return nil } // TODO: вынести установку этой функции в обвязку

	if *generateConfig != "" {
		err := config.ExportToFile(config.DefaultConfig(defaultParsers), *generateConfig)
		if err != nil {
			panic(err)
		}

		os.Exit(0)
	}

	c, err := config.ImportConfig(*configPath, defaultParsers, true)

	return c, *scan, err
}

func initLogger[T any](cfg config.Config[T]) *slog.Logger {
	slogOpt := &slog.HandlerOptions{
		AddSource: cfg.Application.Debug,
		Level:     slog.LevelInfo,
	}

	if cfg.Application.Debug {
		slogOpt.Level = slog.LevelDebug
	}

	return slog.New(
		logHandler{
			Handler: slog.NewJSONHandler(
				os.Stderr,
				slogOpt,
			),
		},
	)
}

// TODO: в случае использования групп реализовать более безопасно.
type logHandler struct {
	slog.Handler
}

func (lh logHandler) Handle(ctx context.Context, r slog.Record) error {
	snapContext := trace.SpanContextFromContext(ctx)
	if snapContext.HasTraceID() {
		r.AddAttrs(slog.String("trace_id", snapContext.TraceID().String()))
	}

	return lh.Handler.Handle(ctx, r)
}

func (lh logHandler) WithGroup(name string) slog.Handler {
	return logHandler{
		Handler: lh.Handler.WithGroup(name),
	}
}

func (lh logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return logHandler{
		Handler: lh.Handler.WithAttrs(attrs),
	}
}
