package main

import (
	"app/internal/dataprovider/loader"
	"app/internal/domain/hgraber"
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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
	flag.Parse()

	l := slog.New(slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
	))
	loader := loader.New(l, *hgToken)

	if *url != "" {
		err := handleBook(ctx, l, loader, *url, *withPages)
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
