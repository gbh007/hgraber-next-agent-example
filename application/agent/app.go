package agent

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gbh007/hgraber-next-agent-example/config"
	"github.com/gbh007/hgraber-next-agent-example/controller/api"
	"github.com/gbh007/hgraber-next-agent-example/controller/async"
	"github.com/gbh007/hgraber-next-agent-example/dataprovider/files"
	"github.com/gbh007/hgraber-next-agent-example/dataprovider/loader"
	"github.com/gbh007/hgraber-next-agent-example/domain/hgraber"
	"github.com/gbh007/hgraber-next-agent-example/pkg"
	agentUC "github.com/gbh007/hgraber-next-agent-example/usecase/agent"
	"github.com/gbh007/hgraber-next-agent-example/usecase/exporter"
)

type ParserInit[T any] func(ctx context.Context, logger *slog.Logger, cfg config.Config[T]) ([]hgraber.Parser, error)

func Serve[T any](ctx context.Context, parserInit ParserInit[T]) {
	cfg, err := parseConfig[T]()
	if err != nil {
		// Поскольку на этот момент нет ни логгера ни вообще ничего то выкидываем панику.
		panic(err)
	}

	logger := initLogger(cfg)
	logger.InfoContext(ctx, "initializing system")

	if cfg.Application.TraceEndpoint != "" {
		err := initTrace(ctx, cfg.Application.TraceEndpoint)
		if err != nil {
			logger.ErrorContext(
				ctx, "fail init otel",
				slog.Any("error", err),
			)

			os.Exit(1)
		}
	}

	parsers, err := parserInit(ctx, logger, cfg)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		return
	}

	async := async.New(logger)
	loader := loader.New(
		logger,
		cfg.Application.ClientTimeout,
		parsers,
	)

	agentUseCases := agentUC.New(logger, loader)

	var exportUseCases api.ExportUseCases

	if cfg.Application.ExportPath != "" {
		fileStorage, err := files.New(cfg.Application.ExportPath, logger)
		if err != nil {
			logger.ErrorContext(ctx, err.Error())

			return
		}

		exportUseCases = exporter.New(logger, fileStorage)
	}

	parserNames := pkg.Map(parsers, func(parser hgraber.Parser) string {
		return parser.Name()
	})

	apiController, err := api.New(
		time.Now(),
		logger,
		agentUseCases,
		exportUseCases,
		cfg.API.Addr,
		cfg.Application.Debug,
		cfg.API.Token,
		parserNames,
	)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		return
	}

	async.RegisterRunner(ctx, apiController)

	logger.InfoContext(ctx, "service start")

	err = async.Serve(ctx)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		return
	}

	logger.InfoContext(ctx, "service stop")
}
