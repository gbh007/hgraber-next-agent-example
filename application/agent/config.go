package agent

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/gbh007/hgraber-next-agent-example/config"

	"github.com/kelseyhightower/envconfig"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/yaml.v3"
)

func parseConfig[T any]() (config.Config[T], error) {
	configPath := flag.String("config", "config.yaml", "path to config")
	flag.Parse()

	c := config.DefaultConfig[T](func() *T { return nil }) // TODO: вынести установку этой функции в обвязку

	f, err := os.Open(*configPath)
	if err != nil {
		return config.Config[T]{}, fmt.Errorf("open config file: %w", err)
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		return config.Config[T]{}, fmt.Errorf("decode yaml: %w", err)
	}

	err = envconfig.Process("APP", &c)
	if err != nil {
		return config.Config[T]{}, fmt.Errorf("decode env: %w", err)
	}

	return c, nil
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
