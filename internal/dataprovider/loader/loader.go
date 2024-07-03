package loader

import (
	"app/internal/dataprovider/loader/internal/parser/hgraber_local"
	"app/internal/dataprovider/loader/internal/parser/mock"
	"app/internal/dataprovider/loader/internal/request"
	"app/internal/domain/hgraber"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type Loader struct {
	logger    *slog.Logger
	requester *request.Requester

	parsers []hgraber.Parser
}

func New(logger *slog.Logger, hgToken string) *Loader {
	requester := request.New(logger)

	return &Loader{
		logger:    logger,
		requester: requester,
		parsers: []hgraber.Parser{
			mock.New(requester),
			hgraber_local.New(requester, hgToken),
		},
	}
}

func (l *Loader) Prefixes() []string {
	prefixes := make([]string, 0, len(l.parsers))

	for _, p := range l.parsers {
		prefixes = append(prefixes, p.Prefixes()...)
	}

	return prefixes
}

func (l *Loader) getParser(u string) (hgraber.Parser, error) {
	for _, p := range l.parsers {
		for _, prefix := range p.Prefixes() {
			if strings.HasPrefix(u, prefix) {
				return p, nil
			}
		}
	}

	return nil, hgraber.InvalidLinkError
}

func (l *Loader) HasParser(ctx context.Context, u string) (bool, error) {
	_, err := l.getParser(u)
	if errors.Is(err, hgraber.InvalidLinkError) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("get parser: %w", err)
	}

	return true, nil
}

func (l *Loader) Load(ctx context.Context, u string) (hgraber.BookParser, error) {
	p, err := l.getParser(u)
	if err != nil {
		return nil, fmt.Errorf("get parser: %w", err)
	}

	bookParser, err := p.Load(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	return bookParser, nil
}

func (l *Loader) LoadImage(ctx context.Context, u string, bookUrl string) (io.ReadCloser, error) {
	var headers http.Header

	// FIXME: переписать получше
	p, _ := l.getParser(bookUrl)
	if p != nil {
		headers, _ = p.Headers(bookUrl)
	}

	data, err := l.requester.Request(ctx, u, headers)
	if err != nil {
		return nil, fmt.Errorf("load image: %w", err)
	}

	return data, nil
}

func (l *Loader) AllBooks(ctx context.Context, u string) ([]string, error) {
	p, err := l.getParser(u)
	if err != nil {
		return nil, err
	}

	data, err := p.AllBooks(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("load image: %w", err)
	}

	return data, nil
}

// FIXME: сделать методом парсера (включить в базовый парсер)
func (l *Loader) Collisions(ctx context.Context, u string) ([]string, error) {
	p, err := l.getParser(u)
	if err != nil {
		return nil, fmt.Errorf("get parser: %w", err)
	}

	for prefix, replacements := range p.Collisions() {
		if strings.HasPrefix(u, prefix) {
			res := make([]string, 0, len(replacements))

			for _, v := range replacements {
				res = append(res, strings.Replace(u, prefix, v, 1))
			}

			return res, nil
		}
	}

	return []string{}, nil // Коллизий нет
}
