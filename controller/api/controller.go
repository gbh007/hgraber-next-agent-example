package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/url"
	"time"

	"github.com/gbh007/hgraber-next-agent-example/entities"
	"github.com/gbh007/hgraber-next-agent-example/open_api/agentAPI"
	"go.opentelemetry.io/otel/trace"

	"github.com/google/uuid"
)

type parsingUseCases interface {
	CheckBooks(ctx context.Context, urls []url.URL) ([]entities.AgentBookCheckResult, error)
	ParseBook(ctx context.Context, u url.URL) (entities.AgentBookDetails, error)
	DownloadPage(ctx context.Context, bookURL, imageURL url.URL) (io.Reader, error)
	CheckPages(ctx context.Context, pages []entities.AgentPageURL) ([]entities.AgentPageCheckResult, error)
	MultiHandle(ctx context.Context, multiUrl url.URL) ([]entities.AgentBookCheckResult, error)
}

type ExportUseCase interface {
	Create(ctx context.Context, data entities.ExportData) error
}

type FileUseCase interface {
	Create(ctx context.Context, fileID uuid.UUID, body io.Reader) error
	Delete(ctx context.Context, fileID uuid.UUID) error
	Get(ctx context.Context, fileID uuid.UUID) (io.Reader, error)
	IDs(ctx context.Context) ([]uuid.UUID, error)
}

type Controller struct {
	startAt time.Time
	logger  *slog.Logger
	tracer  trace.Tracer
	addr    string
	debug   bool

	ogenServer *agentAPI.Server

	exportUseCase   ExportUseCase
	fileUseCase     FileUseCase
	parsingUseCases parsingUseCases

	token       string
	parserCodes []string
}

func New(
	startAt time.Time,
	logger *slog.Logger,
	tracer trace.Tracer,
	parsingUseCases parsingUseCases,
	exportUseCase ExportUseCase,
	fileUseCase FileUseCase,
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

		parserCodes: parserCodes,

		parsingUseCases: parsingUseCases,
		exportUseCase:   exportUseCase,
		fileUseCase:     fileUseCase,
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
