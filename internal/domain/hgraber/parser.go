package hgraber

import (
	"context"
	"errors"
	"net/http"
)

var (
	InvalidLinkError      = errors.New("invalid link")
	UnknownAttributeError = errors.New("unknown attribute")
)

// Parser интерфейс для реализации парсеров для различных сайтов
type Parser interface {
	Load(ctx context.Context, u string) (BookParser, error)
	Prefixes() []string
	Collisions() map[string][]string
	Headers(u string) (http.Header, error)
	AllBooks(ctx context.Context, u string) ([]string, error)
}

type BookParser interface {
	Name(ctx context.Context) (string, error)
	Pages(ctx context.Context) ([]Page, error)
	Tags(ctx context.Context) ([]string, error)
	Authors(ctx context.Context) ([]string, error)
	Characters(ctx context.Context) ([]string, error)
	Languages(ctx context.Context) ([]string, error)
	Categories(ctx context.Context) ([]string, error)
	Parodies(ctx context.Context) ([]string, error)
	Groups(ctx context.Context) ([]string, error)
}

func ParseBookAttr(ctx context.Context, p BookParser, attr Attribute) ([]string, error) {
	switch attr {
	case AttrAuthor:
		return p.Authors(ctx)

	case AttrCategory:
		return p.Categories(ctx)

	case AttrCharacter:
		return p.Characters(ctx)

	case AttrGroup:
		return p.Groups(ctx)

	case AttrLanguage:
		return p.Languages(ctx)

	case AttrParody:
		return p.Parodies(ctx)

	case AttrTag:
		return p.Tags(ctx)

	default:
		return []string{}, UnknownAttributeError
	}
}
