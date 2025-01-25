package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/url"
	"time"

	"github.com/gbh007/hgraber-next-agent-example/controller/api/internal/server"
	"github.com/gbh007/hgraber-next-agent-example/entities"

	"github.com/google/uuid"
)

type parsingUseCases interface {
	CheckBooks(ctx context.Context, urls []url.URL) ([]entities.AgentBookCheckResult, error)
	ParseBook(ctx context.Context, u url.URL) (entities.AgentBookDetails, error)
	DownloadPage(ctx context.Context, bookURL, imageURL url.URL) (io.Reader, error)
	CheckPages(ctx context.Context, pages []entities.AgentPageURL) ([]entities.AgentPageCheckResult, error)
	MultiHandle(ctx context.Context, multiUrl url.URL) ([]entities.AgentBookCheckResult, error)
}

type ExportUseCases interface {
	ExportBook(ctx context.Context, bookID uuid.UUID, bookName string, body io.Reader) error
}

type Controller struct {
	startAt time.Time
	logger  *slog.Logger
	addr    string
	debug   bool

	ogenServer *server.Server

	parsingUseCases parsingUseCases
	exportUseCases  ExportUseCases

	token       string
	parserCodes []string
}

func New(
	startAt time.Time,
	logger *slog.Logger,
	parsingUseCases parsingUseCases,
	exportUseCases ExportUseCases,
	addr string,
	debug bool,
	token string,
	parserCodes []string,
) (*Controller, error) {
	c := &Controller{
		startAt: startAt,
		logger:  logger,
		addr:    addr,
		debug:   debug,
		token:   token,

		parserCodes: parserCodes,

		parsingUseCases: parsingUseCases,
		exportUseCases:  exportUseCases,
	}

	ogenServer, err := server.NewServer(c, c)
	if err != nil {
		return nil, err
	}

	c.ogenServer = ogenServer

	return c, nil
}

var errorAccessForbidden = errors.New("access forbidden")

func (c *Controller) HandleHeaderAuth(ctx context.Context, operationName string, t server.HeaderAuth) (context.Context, error) {
	if c.token == "" {
		return ctx, nil
	}

	if c.token != t.APIKey {
		return ctx, errorAccessForbidden
	}

	return ctx, nil
}
