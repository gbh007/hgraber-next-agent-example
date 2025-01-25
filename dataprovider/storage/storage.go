package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB

	logger *slog.Logger
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	path string,
) (*Storage, error) {
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	err = migrate(ctx, logger, db.DB)
	if err != nil {
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	s := Storage{
		db:     db,
		logger: logger,
	}

	return &s, nil
}

func URLToDB(u *url.URL) sql.NullString {
	if u == nil {
		return sql.NullString{}
	}

	return sql.NullString{
		String: u.String(),
		Valid:  true,
	}
}
