package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36"

type Requester struct {
	client *http.Client

	logger *slog.Logger
}

func New(logger *slog.Logger) *Requester {
	return &Requester{
		client: &http.Client{
			Timeout: time.Minute,
		},
		logger: logger,
	}
}

// requestBuffer запрашивает данные по урле и возвращает их в виде буффера
func (r *Requester) requestBuffer(ctx context.Context, URL string, headers http.Header, body io.Reader) (*bytes.Buffer, error) {
	buff := new(bytes.Buffer)

	var (
		req *http.Request
		err error
	)

	if body != nil {
		req, err = http.NewRequest(http.MethodPost, URL, body)
	} else {
		req, err = http.NewRequest(http.MethodGet, URL, nil)
	}

	if err != nil {
		r.logger.ErrorContext(ctx, err.Error())
		return buff, err
	}

	if len(headers) > 0 {
		for key, values := range headers {
			for _, v := range values {
				req.Header.Add(key, v)
			}
		}
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", defaultUserAgent)
	}

	// выполняем запрос
	response, err := r.client.Do(req)

	if err != nil {
		r.logger.ErrorContext(ctx, err.Error())
		return buff, err
	}

	defer func() {
		closeErr := response.Body.Close()
		if closeErr != nil {
			r.logger.ErrorContext(ctx, closeErr.Error())
		}
	}()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		err = fmt.Errorf("%s ошибка %s", URL, response.Status)
		r.logger.ErrorContext(ctx, err.Error())

		return buff, err
	}

	_, err = buff.ReadFrom(response.Body)
	if err != nil {
		r.logger.ErrorContext(ctx, err.Error())

		return buff, err
	}

	return buff, nil
}

// RequestString запрашивает данные по урле и возвращает их строкой
func (r *Requester) RequestString(ctx context.Context, URL string) (string, error) {
	buff, err := r.requestBuffer(ctx, URL, nil, nil)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

// RequestBytes запрашивает данные по урле и возвращает их массивом байт
func (r *Requester) RequestBytes(ctx context.Context, URL string) ([]byte, error) {
	buff, err := r.requestBuffer(ctx, URL, nil, nil)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (r *Requester) Request(ctx context.Context, URL string, headers http.Header) (io.ReadCloser, error) {
	// FIXME: работать с потоком напрямую
	buff, err := r.requestBuffer(ctx, URL, headers, nil)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(buff), nil
}

func (r *Requester) RequestPost(ctx context.Context, u string, headers http.Header, body io.Reader) ([]byte, error) {
	buff, err := r.requestBuffer(ctx, u, headers, body)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
