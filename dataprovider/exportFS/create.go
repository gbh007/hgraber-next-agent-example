package exportFS

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/gbh007/hgraber-next-agent-core/entities"
)

func (s *Storage) Create(ctx context.Context, data entities.ExportData) error {
	_, err := s.create(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateExport(ctx context.Context, data entities.ExportData) (string, error) {
	return s.create(ctx, data)
}

func (s *Storage) create(ctx context.Context, data entities.ExportData) (string, error) {
	var (
		relativePath, absolutePath string
		err                        error
	)

	if s.limitOnFolder > 0 {
		s.fsLimitMutex.Lock()
		// TODO: подумать над многопоточным способом (включая необходимость)
		defer s.fsLimitMutex.Unlock()

		relativePath, absolutePath, err = s.filepathWithLimits(data.BookID, data.BookName)
		if err != nil {
			return "", fmt.Errorf("get filepath with limits: %w", err)
		}
	} else {
		relativePath, absolutePath = s.filepath(data.BookID, data.BookName)
	}

	f, err := os.Create(absolutePath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}

	_, err = io.Copy(f, data.Body)
	if err != nil {
		fileCloseErr := f.Close()
		if fileCloseErr != nil {
			s.logger.ErrorContext(ctx, "close on write error", slog.Any("error", fileCloseErr))
		}

		return "", fmt.Errorf("write file: %w", err)
	}

	err = f.Close()
	if err != nil {
		return "", fmt.Errorf("close file: %w", err)
	}

	return relativePath, nil
}
