package entities

import (
	"io"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type ExportInfo struct {
	BookID     uuid.UUID
	BookURL    *url.URL
	FSPath     string
	ExportedAt time.Time
}

type ExportData struct {
	BookID   uuid.UUID
	BookName string
	BookURL  *url.URL
	Body     io.Reader
}
