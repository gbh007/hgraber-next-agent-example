package exportAPI

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/gbh007/hgraber-next-agent-example/entities"
	"github.com/google/uuid"
)

type storage interface {
	CreateExport(ctx context.Context, info entities.ExportInfo) error
	ExportedCountByID(ctx context.Context, bookID uuid.UUID) (int, error)
	ExportedCountByURL(ctx context.Context, u url.URL) (int, error)
}

type exportFS interface {
	CreateExport(ctx context.Context, data entities.ExportData) (string, error)
}

type UseCase struct {
	logger *slog.Logger

	storage  storage
	exportFS exportFS
}

func New(
	logger *slog.Logger,
	storage storage,
	exportFS exportFS,
) *UseCase {
	return &UseCase{
		logger:   logger,
		storage:  storage,
		exportFS: exportFS,
	}
}
