package storage

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/gbh007/hgraber-next-agent-example/entities"
	"github.com/google/uuid"
)

func (s *Storage) CreateExport(ctx context.Context, info entities.ExportInfo) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO export_infos (book_id, book_url, relative_path, exported_at) VALUES (?,?,?,?) ON CONFLICT DO NOTHING;`,
		info.BookID,
		URLToDB(info.BookURL),
		info.FSPath,
		info.ExportedAt.Unix(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ExportedCountByID(ctx context.Context, bookID uuid.UUID) (int, error) {
	var c sql.NullInt64

	err := s.db.GetContext(ctx, &c, `SELECT COUNT(*) FROM export_infos WHERE book_id = ?;`, bookID)
	if err != nil {
		return 0, err
	}

	return int(c.Int64), nil
}

func (s *Storage) ExportedCountByURL(ctx context.Context, u url.URL) (int, error) {
	var c sql.NullInt64

	err := s.db.GetContext(ctx, &c, `SELECT COUNT(*) FROM export_infos WHERE book_url = ?;`, u.String())
	if err != nil {
		return 0, err
	}

	return int(c.Int64), nil
}

func (s *Storage) ExportedCountByRelativePath(ctx context.Context, path string) (int, error) {
	var c sql.NullInt64

	err := s.db.GetContext(ctx, &c, `SELECT COUNT(*) FROM export_infos WHERE relative_path = ?;`, path)
	if err != nil {
		return 0, err
	}

	return int(c.Int64), nil
}
