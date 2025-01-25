package main

import (
	"archive/zip"
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gbh007/hgraber-next-agent-example/config"
	"github.com/gbh007/hgraber-next-agent-example/dataprovider/loader"
	"github.com/gbh007/hgraber-next-agent-example/domain/hgraber"
	"github.com/gbh007/hgraber-next-agent-example/external"
	"github.com/gbh007/hgraber-next-agent-example/pkg"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	url := flag.String("u", "", "debug url")
	listBooks := flag.String("l", "", "books handle list")
	hgToken := flag.String("hg", "", "hgraber v4 token")
	withPages := flag.Bool("pages", false, "download pages")
	asZip := flag.Bool("zip", false, "download book as zip")
	printCfg := flag.String("print-config", "", "generate example config")
	flag.Parse()

	l := slog.New(slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
	))

	if *printCfg != "" {
		err := config.ExportToFile(config.DefaultConfig[config.Parsers](config.DefaultParsers), *printCfg)
		if err != nil {
			l.ErrorContext(ctx, err.Error())

			return
		}

		return
	}

	loader := loader.New(
		l,
		time.Minute,
		loader.NewDefaultParsers(
			l,
			*hgToken,
			time.Minute,
			[]string{
				"mock",
				"hgraber_local",
			},
		),
	)

	var err error

	if *url != "" {
		if *asZip {
			err = handleBookToZip(ctx, loader, *url)
		} else {
			err = handleBook(ctx, l, loader, *url, *withPages)
		}

		if err != nil {
			l.ErrorContext(ctx, err.Error())

			return
		}
	}

	if *listBooks != "" {
		urls, err := loader.AllBooks(ctx, *listBooks)
		if err != nil {
			l.ErrorContext(ctx, err.Error())

			return
		}

		l.InfoContext(ctx, "list", slog.Int("count", len(urls)), slog.Any("urls", urls))

		for _, u := range urls {
			fmt.Println(u)
		}
	}
}

func handleBook(ctx context.Context, l *slog.Logger, loader *loader.Loader, bookUrl string, withPages bool) error {
	parser, err := loader.Load(ctx, bookUrl)
	if err != nil {
		return err
	}

	pages, err := parser.Pages(ctx)
	l.InfoContext(ctx, "Pages", slog.Any("data", pages), slog.Any("error", err))

	name, err := parser.Name(ctx)
	l.InfoContext(ctx, "Name", slog.Any("data", name), slog.Any("error", err))

	tags, err := parser.Tags(ctx)
	l.InfoContext(ctx, "Tags", slog.Any("data", tags), slog.Any("error", err))

	authors, err := parser.Authors(ctx)
	l.InfoContext(ctx, "Authors", slog.Any("data", authors), slog.Any("error", err))

	characters, err := parser.Characters(ctx)
	l.InfoContext(ctx, "Characters", slog.Any("data", characters), slog.Any("error", err))

	languages, err := parser.Languages(ctx)
	l.InfoContext(ctx, "Languages", slog.Any("data", languages), slog.Any("error", err))

	categories, err := parser.Categories(ctx)
	l.InfoContext(ctx, "Categories", slog.Any("data", categories), slog.Any("error", err))

	parodies, err := parser.Parodies(ctx)
	l.InfoContext(ctx, "Parodies", slog.Any("data", parodies), slog.Any("error", err))

	groups, err := parser.Groups(ctx)
	l.InfoContext(ctx, "Groups", slog.Any("data", groups), slog.Any("error", err))

	if !withPages {
		return nil
	}

	for _, p := range pages {
		err = loadPage(ctx, loader, p, bookUrl)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadPage(ctx context.Context, loader *loader.Loader, p hgraber.Page, bookUrl string) error {
	r, err := loader.LoadImage(ctx, p.URL, bookUrl)
	if err != nil {
		return fmt.Errorf("page load %d %s:%w", p.PageNumber, p.URL, err)
	}

	defer r.Close()

	f, err := os.Create(fmt.Sprintf("page_test_%d.%s", p.PageNumber, p.Ext))
	if err != nil {
		return fmt.Errorf("file create %d %s:%w", p.PageNumber, p.URL, err)
	}

	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return fmt.Errorf("page write %d %s:%w", p.PageNumber, p.URL, err)
	}

	return nil
}

func handleBookToZip(ctx context.Context, loader *loader.Loader, bookUrl string) error {
	parser, err := loader.Load(ctx, bookUrl)
	if err != nil {
		return err
	}

	pages, err := parser.Pages(ctx)
	if err != nil {
		return err
	}

	name, err := parser.Name(ctx)
	if err != nil {
		return err
	}

	tags, err := parser.Tags(ctx)
	if err != nil {
		return err
	}

	authors, err := parser.Authors(ctx)
	if err != nil {
		return err
	}

	characters, err := parser.Characters(ctx)
	if err != nil {
		return err
	}

	languages, err := parser.Languages(ctx)
	if err != nil {
		return err
	}

	categories, err := parser.Categories(ctx)
	if err != nil {
		return err
	}

	parodies, err := parser.Parodies(ctx)
	if err != nil {
		return err
	}

	groups, err := parser.Groups(ctx)
	if err != nil {
		return err
	}

	f, err := os.Create("dump.zip")
	if err != nil {
		return err
	}

	defer f.Close()

	pageUrls := pkg.SliceToMap(pages, func(p hgraber.Page) (int, string) {
		return p.PageNumber, p.URL
	})

	zipWriter := zip.NewWriter(f)

	err = external.WriteArchive(
		ctx, zipWriter,
		func(ctx context.Context, pageNumber int) (io.Reader, error) {
			u, ok := pageUrls[pageNumber]
			if !ok {
				return nil, fmt.Errorf("missing page %d", pageNumber)
			}

			return loader.LoadImage(ctx, u, bookUrl)
		},
		external.Info{
			Version: "1.0.0",
			Meta: external.Meta{
				Exported:    time.Now().UTC(),
				ServiceName: "hgraber next agent",
			},
			Data: external.Book{
				Name:             name,
				OriginURL:        bookUrl,
				PageCount:        len(pages),
				CreateAt:         time.Now(),
				AttributesParsed: true,
				Attributes: []external.Attribute{
					{
						Code:   external.AttributeCodeTag,
						Values: tags,
					},
					{
						Code:   external.AttributeCodeAuthor,
						Values: authors,
					},
					{
						Code:   external.AttributeCodeCharacter,
						Values: characters,
					},
					{
						Code:   external.AttributeCodeLanguage,
						Values: languages,
					},
					{
						Code:   external.AttributeCodeCategory,
						Values: categories,
					},
					{
						Code:   external.AttributeCodeParody,
						Values: parodies,
					},
					{
						Code:   external.AttributeCodeGroup,
						Values: groups,
					},
				},
				Pages: pkg.Map(pages, func(p hgraber.Page) external.Page {
					ext := p.Ext
					if !strings.HasPrefix(ext, ".") {
						ext = "." + ext
					}

					return external.Page{
						PageNumber: p.PageNumber,
						Ext:        ext,
						OriginURL:  p.URL,
						CreateAt:   time.Now(),
						Downloaded: true,
						LoadAt:     time.Now(),
					}
				}),
			},
		},
	)
	if err != nil {
		return err
	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	return nil
}
