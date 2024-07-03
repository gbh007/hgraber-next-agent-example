package common

import (
	"app/internal/domain/hgraber"
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// Проверка соответствия базового типа
var (
	_ hgraber.Parser = (*CoreParser)(nil)
)

type Requester interface {
	RequestString(ctx context.Context, URL string) (string, error)
	RequestPost(ctx context.Context, u string, headers http.Header, body io.Reader) ([]byte, error)
}

type CoreParser struct {
	Requester  Requester
	prefixes   []string
	collisions map[string][]string
}

func NewCoreParser(requester Requester, prefixes []string) CoreParser {
	collisions := make(map[string][]string, len(prefixes))

	if len(prefixes) > 1 { // Если префикс только 1, то коллизий быть не может.
		for i1, pref1 := range prefixes {
			values := make([]string, 0, len(prefixes)-1)

			for i2, pref2 := range prefixes {
				if i1 == i2 {
					continue
				}

				values = append(values, pref2)
			}

			collisions[pref1] = values
		}
	}

	return CoreParser{
		Requester:  requester,
		prefixes:   prefixes,
		collisions: collisions,
	}
}

func (cp CoreParser) Load(ctx context.Context, u string) (hgraber.BookParser, error) {
	return nil, fmt.Errorf("unimplemented in core parser")
}

func (cp CoreParser) Prefixes() []string {
	return cp.prefixes
}

func (cp CoreParser) Collisions() map[string][]string {
	return cp.collisions
}

func (cp CoreParser) Headers(u string) (http.Header, error) {
	return nil, nil
}

func (cp CoreParser) AllBooks(ctx context.Context, u string) ([]string, error) {
	return nil, nil
}

func TrimLastSlash(URL string, count int) string {
	c := 0

	ind := strings.LastIndexFunc(URL, func(r rune) bool {
		if r != rune('/') {
			return false
		}
		c++
		return c == count
	})

	return URL[:ind]
}

func OneMatch(rgx *regexp.Regexp, raw string) (string, bool) {
	res := rgx.FindAllStringSubmatch(raw, -1)
	if len(res) < 1 || len(res[0]) != 2 {
		return "", false
	}

	return strings.TrimSpace(html.UnescapeString(res[0][1])), true
}
