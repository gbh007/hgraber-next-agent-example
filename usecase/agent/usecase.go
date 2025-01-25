package agent

import (
	"context"
	"io"
	"log/slog"

	"github.com/gbh007/hgraber-next-agent-example/domain/hgraber"
)

type loader interface {
	HasParser(ctx context.Context, u string) (bool, error)
	Load(ctx context.Context, URL string) (hgraber.BookParser, error)
	LoadImage(ctx context.Context, u string, bookUrl string) (io.ReadCloser, error)
	Collisions(ctx context.Context, u string) ([]string, error)
	AllBooks(ctx context.Context, u string) ([]string, error)
}

type UseCase struct {
	logger *slog.Logger

	loader loader
}

func New(
	logger *slog.Logger,
	loader loader,
) *UseCase {
	return &UseCase{
		logger: logger,
		loader: loader,
	}
}
