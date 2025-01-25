package exporter

import (
	"context"
	"io"

	"github.com/google/uuid"
)

func (uc *UseCase) ExportBook(ctx context.Context, bookID uuid.UUID, bookName string, body io.Reader) error {
	return uc.fileStorage.Create(ctx, bookID, bookName, body)
}
