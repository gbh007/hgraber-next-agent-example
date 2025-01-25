package exportAPI

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gbh007/hgraber-next-agent-core/entities"
)

func (uc *UseCase) Create(ctx context.Context, data entities.ExportData) error {
	c, err := uc.storage.ExportedCountByID(ctx, data.BookID)
	if err != nil {
		return fmt.Errorf("check export count by id: %w", err)
	}

	if c > 0 {
		uc.logger.DebugContext(
			ctx, "export already exists",
			slog.String("book_id", data.BookID.String()),
		)

		return nil
	}

	if data.BookURL != nil {
		c, err := uc.storage.ExportedCountByURL(ctx, *data.BookURL)
		if err != nil {
			return fmt.Errorf("check export count by url: %w", err)
		}

		if c > 0 {
			uc.logger.DebugContext(
				ctx, "export already exists",
				slog.String("book_id", data.BookID.String()),
				slog.String("book_url", data.BookURL.String()),
			)

			return nil
		}
	}

	relativePath, err := uc.exportFS.CreateExport(ctx, data)
	if err != nil {
		return fmt.Errorf("create file in export fs: %w", err)
	}

	err = uc.storage.CreateExport(ctx, entities.ExportInfo{
		BookID:     data.BookID,
		BookURL:    data.BookURL,
		FSPath:     relativePath,
		ExportedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("create export info: %w", err)
	}

	return nil
}
