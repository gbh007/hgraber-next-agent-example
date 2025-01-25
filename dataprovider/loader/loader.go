package loader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gbh007/hgraber-next-agent-example/domain/hgraber"
	"github.com/gbh007/hgraber-next-agent-example/parser/hgraber_local"
	"github.com/gbh007/hgraber-next-agent-example/parser/mock"
	"github.com/gbh007/hgraber-next-agent-example/request"
)

type Loader struct {
	logger    *slog.Logger
	requester *request.Requester

	parsers []hgraber.Parser
}

func NewDefaultParsers(
	logger *slog.Logger,
	hgToken string,
	timeout time.Duration,
	enabledParsers []string,
) []hgraber.Parser {
	requester := request.New(logger, timeout)

	parsers := make([]hgraber.Parser, 0, len(enabledParsers))

	for _, code := range enabledParsers {
		var p hgraber.Parser

		switch code {
		case "mock":
			p = mock.New(requester)

		case "hgraber_local":
			p = hgraber_local.New(requester, hgToken)

		default:
			logger.Warn(
				"unknown parser code",
				slog.String("code", code),
			)

			continue
		}

		parsers = append(parsers, p)
	}

	return parsers
}

func New(
	logger *slog.Logger,
	timeout time.Duration,
	parsers []hgraber.Parser,
) *Loader {
	requester := request.New(logger, timeout)

	return &Loader{
		logger:    logger,
		requester: requester,
		parsers:   parsers,
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
		return nil, fmt.Errorf("load books: %w", err)
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
