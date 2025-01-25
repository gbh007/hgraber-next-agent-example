package agent

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gbh007/hgraber-next-agent-core/config"
	"github.com/gbh007/hgraber-next-agent-core/controller/api"
	"github.com/gbh007/hgraber-next-agent-core/controller/async"
	"github.com/gbh007/hgraber-next-agent-core/dataprovider/dataFS"
	"github.com/gbh007/hgraber-next-agent-core/dataprovider/exportFS"
	"github.com/gbh007/hgraber-next-agent-core/dataprovider/loader"
	"github.com/gbh007/hgraber-next-agent-core/dataprovider/masterAPI"
	"github.com/gbh007/hgraber-next-agent-core/dataprovider/storage"
	"github.com/gbh007/hgraber-next-agent-core/domain/hgraber"
	"github.com/gbh007/hgraber-next-agent-core/pkg"
	agentUC "github.com/gbh007/hgraber-next-agent-core/usecase/agent"
	"github.com/gbh007/hgraber-next-agent-core/usecase/exportAPI"
	"github.com/gbh007/hgraber-next-agent-core/usecase/exportDeduplicator"
	"go.opentelemetry.io/otel"
)

type ParserInit[T any] func(ctx context.Context, logger *slog.Logger, cfg config.Config[T]) ([]hgraber.Parser, error)

func Serve[T any](ctx context.Context, parserInit ParserInit[T]) {
	cfg, needScan, err := parseConfig[T]()
	if err != nil {
		// Поскольку на этот момент нет ни логгера ни вообще ничего то выкидываем панику.
		panic(err)
	}

	logger := initLogger(cfg)
	logger.InfoContext(ctx, "initializing system")

	if cfg.Application.TraceEndpoint != "" {
		err := initTrace(
			ctx,
			cfg.Application.TraceEndpoint,
			cfg.Application.ServiceName,
		)
		if err != nil {
			logger.ErrorContext(
				ctx, "fail init otel",
				slog.Any("error", err),
			)

			os.Exit(1)
		}
	}

	tracer := otel.GetTracerProvider().Tracer("hgraber-next-agent")

	async := async.New(logger)

	var (
		exportStorage api.ExportUseCases
		fileStorage   api.FileUseCases
		agentUseCases api.ParsingUseCases

		exportStorageRaw *exportFS.Storage
		dbRaw            *storage.Storage
		mAPI             *masterAPI.Client
	)

	parsers, err := parserInit(ctx, logger, cfg)
	if err != nil {
		logger.ErrorContext(
			ctx, "fail init parsers",
			slog.Any("error", err),
		)

		os.Exit(1)
	}

	if len(parsers) > 0 {
		loader := loader.New(
			logger,
			cfg.Application.ClientTimeout,
			parsers,
		)

		agentUseCases = agentUC.New(logger, loader)

		logger.DebugContext(
			ctx, "use parsing",
		)
	}

	if cfg.FSBase.ExportPath != "" {
		exportStorageRaw, err = exportFS.New(cfg.FSBase.ExportPath, logger, cfg.FSBase.ExportLimitOnFolder, cfg.Application.UseUnsafeCloser)
		if err != nil {
			logger.ErrorContext(
				ctx, "fail init export fs",
				slog.Any("error", err),
			)

			os.Exit(1)
		}

		exportStorage = exportStorageRaw

		logger.DebugContext(
			ctx, "use local export storage",
			slog.String("path", cfg.FSBase.ExportPath),
		)
	}

	if cfg.FSBase.FilePath != "" {
		fileStorage, err = dataFS.New(cfg.FSBase.FilePath, logger)
		if err != nil {
			logger.ErrorContext(
				ctx, "fail init data fs",
				slog.Any("error", err),
			)

			os.Exit(1)
		}

		logger.DebugContext(
			ctx, "use local file storage",
			slog.String("path", cfg.FSBase.FilePath),
		)
	}

	if cfg.Sqlite.FilePath != "" {
		dbRaw, err = storage.New(ctx, logger, cfg.Sqlite.FilePath)
		if err != nil {
			logger.ErrorContext(
				ctx, "fail init db",
				slog.Any("error", err),
			)

			os.Exit(1)
		}
	}

	if cfg.ZipScanner.MasterAddr != "" {
		mAPI, err = masterAPI.New(cfg.ZipScanner.MasterAddr, cfg.ZipScanner.MasterToken)
		if err != nil {
			logger.ErrorContext(
				ctx, "fail init master api",
				slog.Any("error", err),
			)

			os.Exit(1)
		}
	}

	if needScan {
		if dbRaw == nil || exportStorageRaw == nil || mAPI == nil {
			logger.ErrorContext(ctx, "invalid scan dependencies")

			os.Exit(1)
		}

		err = exportDeduplicator.New(logger, exportStorageRaw, dbRaw, mAPI).ScanZips(ctx)
		if err != nil {
			logger.ErrorContext(
				ctx, "fail scan zips",
				slog.Any("error", err),
			)

			os.Exit(1)
		}

		return
	}

	if cfg.FSBase.EnableDeduplication && dbRaw != nil && exportStorageRaw != nil {
		exportStorage = exportAPI.New(logger, dbRaw, exportStorageRaw)

		logger.DebugContext(ctx, "use export deduplication")
	}

	parserNames := pkg.Map(parsers, func(parser hgraber.Parser) string {
		return parser.Name()
	})

	apiController, err := api.New(
		time.Now(),
		logger,
		tracer,
		agentUseCases,
		exportStorage,
		fileStorage,
		cfg.API.Addr,
		cfg.Application.Debug,
		cfg.API.Token,
		parserNames,
	)
	if err != nil {
		logger.ErrorContext(
			ctx, "fail init api controller",
			slog.Any("error", err),
		)

		os.Exit(1)
	}

	async.RegisterRunner(ctx, apiController)

	logger.InfoContext(ctx, "application start")
	defer logger.InfoContext(ctx, "application stop")

	err = async.Serve(ctx)
	if err != nil {
		logger.ErrorContext(
			ctx, "fail serve",
			slog.Any("error", err),
		)

		os.Exit(1)
	}
}
