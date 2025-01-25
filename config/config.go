package config

import "time"

type Config[T any] struct {
	API         API         `envconfig:"API" yaml:"api"`
	Application Application `envconfig:"APPLICATION" yaml:"application"`
	Parsers     *T          `envconfig:"PARSERS" yaml:"parsers"`
}

func DefaultConfig[T any](defaultParsers func() *T) Config[T] {
	return Config[T]{
		API:         DefaultAPI(),
		Parsers:     defaultParsers(),
		Application: DefaultApplication(),
	}
}

type Application struct {
	ExportPath    string        `envconfig:"EXPORT_PATH" yaml:"export_path"`
	Debug         bool          `envconfig:"DEBUG" yaml:"debug"`
	TraceEndpoint string        `envconfig:"TRACE_ENDPOINT" yaml:"trace_endpoint"`
	ClientTimeout time.Duration `envconfig:"CLIENT_TIMEOUT" yaml:"client_timeout"`
}

func DefaultApplication() Application {
	return Application{
		ExportPath:    "",
		Debug:         false,
		TraceEndpoint: "",
		ClientTimeout: time.Minute,
	}
}

type API struct {
	Addr  string `envconfig:"ADDR" yaml:"addr"`
	Token string `envconfig:"TOKEN" yaml:"token"`
}

func DefaultAPI() API {
	return API{
		Addr:  ":8080",
		Token: "",
	}
}

type Parsers struct {
	HG4Token string   `envconfig:"HG4_TOKEN" yaml:"hg4_token"`
	Enabled  []string `envconfig:"ENABLED" yaml:"enabled"`
}

func DefaultParsers() *Parsers {
	return &Parsers{
		HG4Token: "",
		Enabled: []string{
			"mock",
			"hgraber_local",
		},
	}
}
