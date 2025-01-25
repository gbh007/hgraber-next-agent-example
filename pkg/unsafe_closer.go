package pkg

import (
	"errors"
	"io"
)

// UnsafeCloser оболочка для закрытия потока, после полного чтения.
// Является костылем, отказаться с реализацией https://github.com/ogen-go/ogen/issues/1023.
type UnsafeCloser struct {
	isClosed bool

	Body io.ReadCloser
}

func (c *UnsafeCloser) Read(p []byte) (n int, err error) {
	if c.isClosed || c.Body == nil {
		return 0, io.EOF
	}

	n, err = c.Body.Read(p)

	if errors.Is(err, io.EOF) {
		_ = c.Body.Close()
		c.isClosed = true
	}

	return
}
