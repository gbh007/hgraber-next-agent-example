package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
)

const defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36"

type Requester struct {
	client *http.Client

	logger *slog.Logger
}

func New(logger *slog.Logger, timeout time.Duration) *Requester {
	return &Requester{
		client: &http.Client{
			Timeout: timeout,
			Transport: otelhttp.NewTransport(
				http.DefaultTransport,
				otelhttp.WithPropagators(noopPropagator{}),
			),
		},
		logger: logger,
	}
}

type noopPropagator struct{}

func (noopPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {}
func (noopPropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return ctx
}
func (noopPropagator) Fields() []string { return nil }

// requestBuffer запрашивает данные по урле и возвращает их в виде буффера
func (r *Requester) requestBuffer(ctx context.Context, URL string, headers http.Header, body io.Reader) (*bytes.Buffer, string, error) {
	buff := new(bytes.Buffer)

	var (
		req *http.Request
		err error
	)

	if body != nil {
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, URL, body)
	} else {
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	}

	if err != nil {
		r.logger.ErrorContext(ctx, err.Error())
		return nil, "", err
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
		return nil, "", err
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

		return nil, "", err
	}

	_, err = buff.ReadFrom(response.Body)
	if err != nil {
		r.logger.ErrorContext(ctx, err.Error())

		return nil, "", err
	}

	if response.Request != nil &&
		response.Request.URL != nil &&
		response.Request.URL.String() != URL {
		URL = response.Request.URL.String()
	}

	return buff, URL, nil
}

// RequestString запрашивает данные по урле и возвращает их строкой
func (r *Requester) RequestString(ctx context.Context, URL string) (string, error) {
	buff, _, err := r.requestBuffer(ctx, URL, nil, nil)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

func (r *Requester) RequestStringWithRedirect(ctx context.Context, URL string) (string, string, error) {
	buff, resultURL, err := r.requestBuffer(ctx, URL, nil, nil)
	if err != nil {
		return "", "", err
	}

	return buff.String(), resultURL, nil
}

// RequestBytes запрашивает данные по урле и возвращает их массивом байт
func (r *Requester) RequestBytes(ctx context.Context, URL string) ([]byte, error) {
	buff, _, err := r.requestBuffer(ctx, URL, nil, nil)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (r *Requester) Request(ctx context.Context, URL string, headers http.Header) (io.ReadCloser, error) {
	// FIXME: работать с потоком напрямую
	buff, _, err := r.requestBuffer(ctx, URL, headers, nil)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(buff), nil
}

func (r *Requester) RequestPost(ctx context.Context, u string, headers http.Header, body io.Reader) ([]byte, error) {
	buff, _, err := r.requestBuffer(ctx, u, headers, body)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
