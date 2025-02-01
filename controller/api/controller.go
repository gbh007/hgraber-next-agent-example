package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/url"
	"time"

	"github.com/gbh007/hgraber-next-agent-core/entities"
	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
	"go.opentelemetry.io/otel/trace"

	"github.com/google/uuid"
)

type ParsingUseCases interface {
	CheckBooks(ctx context.Context, urls []url.URL) ([]entities.AgentBookCheckResult, error)
	ParseBook(ctx context.Context, u url.URL) (entities.AgentBookDetails, error)
	DownloadPage(ctx context.Context, bookURL, imageURL url.URL) (io.Reader, error)
	CheckPages(ctx context.Context, pages []entities.AgentPageURL) ([]entities.AgentPageCheckResult, error)
	MultiHandle(ctx context.Context, multiUrl url.URL) ([]entities.AgentBookCheckResult, error)
}

type ExportUseCases interface {
	Create(ctx context.Context, data entities.ExportData) error
}

type FileUseCases interface {
	Create(ctx context.Context, fileID uuid.UUID, body io.Reader) error
	Delete(ctx context.Context, fileID uuid.UUID) error
	Get(ctx context.Context, fileID uuid.UUID) (io.Reader, error)
	State(ctx context.Context, includeFileIDs, includeFileSizes bool) (entities.FSState, error)
}

type HighwayUseCases interface {
	NewToken(ctx context.Context) (string, time.Time, error)
	ValidateToken(ctx context.Context, token string) error
	Get(ctx context.Context, fileID uuid.UUID) (io.Reader, error)
}

type Controller struct {
	startAt time.Time
	logger  *slog.Logger
	tracer  trace.Tracer
	addr    string
	debug   bool

	ogenServer *agentAPI.Server

	exportUseCase   ExportUseCases
	fileUseCase     FileUseCases
	parsingUseCases ParsingUseCases
	highwayUseCase  HighwayUseCases

	token          string
	parserCodes    []string
	enabledModules []string
}

func New(
	startAt time.Time,
	logger *slog.Logger,
	tracer trace.Tracer,
	parsingUseCases ParsingUseCases,
	exportUseCase ExportUseCases,
	fileUseCase FileUseCases,
	highwayUseCase HighwayUseCases,
	addr string,
	debug bool,
	token string,
	parserCodes []string,
) (*Controller, error) {
	c := &Controller{
		startAt: startAt,
		logger:  logger,
		tracer:  tracer,
		addr:    addr,
		debug:   debug,
		token:   token,

		parserCodes:    parserCodes,
		enabledModules: make([]string, 0, 3),

		parsingUseCases: parsingUseCases,
		exportUseCase:   exportUseCase,
		fileUseCase:     fileUseCase,
		highwayUseCase:  highwayUseCase,
	}

	if c.parsingUseCases != nil {
		c.enabledModules = append(c.enabledModules, "parsing")
	}

	if c.exportUseCase != nil {
		c.enabledModules = append(c.enabledModules, "export")
	}

	if c.fileUseCase != nil {
		c.enabledModules = append(c.enabledModules, "file system")
	}

	if c.highwayUseCase != nil {
		c.enabledModules = append(c.enabledModules, "highway")
	}

	ogenServer, err := agentAPI.NewServer(c, c)
	if err != nil {
		return nil, err
	}

	c.ogenServer = ogenServer

	return c, nil
}

var errorAccessForbidden = errors.New("access forbidden")

func (c *Controller) HandleHeaderAuth(ctx context.Context, operationName string, t agentAPI.HeaderAuth) (context.Context, error) {
	if c.token == "" {
		return ctx, nil
	}

	if c.token != t.APIKey {
		return ctx, errorAccessForbidden
	}

	return ctx, nil
}
