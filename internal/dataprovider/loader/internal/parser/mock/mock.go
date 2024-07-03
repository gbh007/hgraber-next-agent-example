package mock

import (
	"app/internal/dataprovider/loader/internal/parser/common"
	"app/internal/domain/hgraber"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Проверка соответствия базового типа
var (
	_ hgraber.BookParser = (*BookParser)(nil)
	_ hgraber.Parser     = (*Parser)(nil)

	ParserError = errors.New("parser mock")
)

type Parser struct {
	common.CoreParser
}

func New(r common.Requester) *Parser {
	return &Parser{
		CoreParser: common.NewCoreParser(r, []string{
			"http://localhost",
		}),
	}
}

func (p *Parser) Load(ctx context.Context, URL string) (hgraber.BookParser, error) {
	bookParser := BookParser{
		url: URL,
	}

	if len(bookParser.url) > 1 && bookParser.url[len(bookParser.url)-1] == '/' {
		bookParser.url = bookParser.url[:len(bookParser.url)-1]
	}

	body, err := p.Requester.RequestString(ctx, URL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ParserError, err)
	}

	bookParser.body = body

	return bookParser, nil
}

type BookParser struct {
	url, body string
}

func (p BookParser) Pages(ctx context.Context) ([]hgraber.Page, error) {
	result := make([]hgraber.Page, 0)

	rp := `(?sm)` + regexp.QuoteMeta(`<a href="`) + `(.+?)\.(.+?)` + regexp.QuoteMeta(`">`)
	for i, name := range regexp.MustCompile(rp).FindAllStringSubmatch(p.body, -1) {
		if len(name) > 1 {
			result = append(result, hgraber.Page{
				PageNumber: i + 1,
				URL:        p.url + "/" + strings.TrimSpace(name[1]) + "." + strings.TrimSpace(name[2]),
				Ext:        name[2],
			})
		}
	}

	return result, nil
}

func (p BookParser) Name(ctx context.Context) (string, error) {
	return "mock name", nil
}

func (p BookParser) Tags(ctx context.Context) ([]string, error) {
	return []string{"mock Tags"}, nil
}

func (p BookParser) Authors(ctx context.Context) ([]string, error) {
	return []string{"mock Authors"}, nil
}

func (p BookParser) Characters(ctx context.Context) ([]string, error) {
	return []string{"mock Characters"}, nil
}

func (p BookParser) Languages(ctx context.Context) ([]string, error) {
	return []string{"mock Languages"}, nil
}

func (p BookParser) Categories(ctx context.Context) ([]string, error) {
	return []string{"mock Categories"}, nil
}

func (p BookParser) Parodies(ctx context.Context) ([]string, error) {
	return []string{"mock Parodies"}, nil
}

func (p BookParser) Groups(ctx context.Context) ([]string, error) {
	return []string{"mock Groups"}, nil
}
