package exportFS

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/gbh007/hgraber-next-agent-core/pkg"
)

func (s *Storage) Get(ctx context.Context, relativePath string) (io.Reader, error) {
	f, err := os.Open(path.Join(s.fsPath, relativePath))
	if err != nil {
		return nil, fmt.Errorf("export fs: open: %w", err)
	}

	if s.useUnsafeCloser {
		return &pkg.UnsafeCloser{
			Body: f,
		}, nil
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("export fs: read all: %w", err)
	}

	return bytes.NewReader(data), nil
}
