package agent

import (
	"app/internal/controller/api"
	"app/internal/controller/async"
	"app/internal/dataprovider/files"
	"app/internal/dataprovider/loader"
	agentUC "app/internal/usecase/agent"
	"app/internal/usecase/exporter"
	"context"
	"time"
)

func Serve(ctx context.Context) {
	cfg := parseFlag()

	logger := initLogger(cfg)
	logger.InfoContext(ctx, "initializing system")

	async := async.New(logger)
	loader := loader.New(logger, cfg.HG4Token)

	agentUseCases := agentUC.New(logger, loader)

	var exportUseCases api.ExportUseCases

	if cfg.ExportPath != "" {
		fileStorage, err := files.New(cfg.ExportPath, logger)
		if err != nil {
			logger.ErrorContext(ctx, err.Error())

			return
		}

		exportUseCases = exporter.New(logger, fileStorage)
	}

	apiController, err := api.New(
		time.Now(),
		logger,
		agentUseCases,
		exportUseCases,
		cfg.Addr,
		cfg.Debug,
		cfg.Token,
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
