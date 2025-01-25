package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/gbh007/hgraber-next-agent-example/application/agent"
	"github.com/gbh007/hgraber-next-agent-example/config"
	"github.com/gbh007/hgraber-next-agent-example/dataprovider/loader"
	"github.com/gbh007/hgraber-next-agent-example/domain/hgraber"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	agent.Serve(ctx, func(ctx context.Context, logger *slog.Logger, cfg config.Config[config.Parsers]) ([]hgraber.Parser, error) {
		if cfg.Parsers == nil {
			logger.DebugContext(ctx, "nil parser config, skipping")

			return []hgraber.Parser{}, nil
		}

		return loader.NewDefaultParsers(
			logger,
			cfg.Parsers.HG4Token,
			cfg.Application.ClientTimeout,
			cfg.Parsers.Enabled,
		), nil
	})
}
