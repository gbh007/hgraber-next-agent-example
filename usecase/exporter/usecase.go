package exporter

import (
	"context"
	"io"
	"log/slog"

	"github.com/google/uuid"
)

type fileStorage interface {
	Create(ctx context.Context, bookID uuid.UUID, bookName string, body io.Reader) error
}

type UseCase struct {
	logger *slog.Logger

	fileStorage fileStorage
}

func New(
	logger *slog.Logger,
	fileStorage fileStorage,
) *UseCase {
	return &UseCase{
		logger:      logger,
		fileStorage: fileStorage,
	}
}
