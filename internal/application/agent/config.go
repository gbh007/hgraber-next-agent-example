package agent

import (
	"flag"
	"log/slog"
	"os"
)

type configRaw struct {
	Addr  string
	Token string

	HG4Token string

	ExportPath string

	Debug bool
}

func parseFlag() configRaw {
	hgAddr := flag.String("addr", ":8080", "Адрес сервера API")
	hgToken := flag.String("token", "", "Токен для доступа к API")
	debug := flag.Bool("debug", false, "Режим отладки")
	hg4Token := flag.String("hg", "", "hgraber v4 token")
	exportPath := flag.String("export-path", "", "Путь для экспорта")

	flag.Parse()

	cfg := configRaw{
		Addr:       *hgAddr,
		Token:      *hgToken,
		Debug:      *debug,
		HG4Token:   *hg4Token,
		ExportPath: *exportPath,
	}

	return cfg
}

func initLogger(cfg configRaw) *slog.Logger {
	slogOpt := &slog.HandlerOptions{
		AddSource: cfg.Debug,
		Level:     slog.LevelInfo,
	}

	if cfg.Debug {
		slogOpt.Level = slog.LevelDebug
	}

	return slog.New(slog.NewJSONHandler(
		os.Stderr,
		slogOpt,
	))
}
