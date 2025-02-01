package config

import "time"

type Config[T any] struct {
	API         API         `envconfig:"API" yaml:"api"`
	Application Application `envconfig:"APPLICATION" yaml:"application"`
	Parsers     *T          `envconfig:"PARSERS" yaml:"parsers"`
	FSBase      FSBase      `envconfig:"FS_BASE" yaml:"fs_base"`
	Sqlite      Sqlite      `envconfig:"SQLITE" yaml:"sqlite"`
	ZipScanner  ZipScanner  `envconfig:"ZIP_SCANNER" yaml:"zip_scanner"`
	Highway     Highway     `envconfig:"HIGHWAY" yaml:"highway"`
}

func DefaultConfig[T any](defaultParsers func() *T) Config[T] {
	return Config[T]{
		API:         DefaultAPI(),
		Parsers:     defaultParsers(),
		Application: DefaultApplication(),
		FSBase:      DefaultFSBase(),
		Sqlite:      DefaultSqlite(),
		Highway:     DefaultHighway(),
	}
}

type Application struct {
	Debug           bool          `envconfig:"DEBUG" yaml:"debug"`
	TraceEndpoint   string        `envconfig:"TRACE_ENDPOINT" yaml:"trace_endpoint"`
	ClientTimeout   time.Duration `envconfig:"CLIENT_TIMEOUT" yaml:"client_timeout"`
	ServiceName     string        `envconfig:"SERVICE_NAME" yaml:"service_name"`
	UseUnsafeCloser bool          `envconfig:"USE_UNSAFE_CLOSER" yaml:"use_unsafe_closer"`
}

func DefaultApplication() Application {
	return Application{
		Debug:         false,
		TraceEndpoint: "",
		ClientTimeout: time.Minute,
		ServiceName:   "hgraber-next-agent",
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

type FSBase struct {
	ExportPath          string `envconfig:"EXPORT_PATH" yaml:"export_path"`
	FilePath            string `envconfig:"FILE_PATH" yaml:"file_path"`
	EnableDeduplication bool   `envconfig:"ENABLE_DEDUPLICATION" yaml:"enable_deduplication"`
	ExportLimitOnFolder int    `envconfig:"EXPORT_LIMIT_ON_FOLDER" yaml:"export_limit_on_folder"`
}

func DefaultFSBase() FSBase {
	return FSBase{}
}

type Sqlite struct {
	FilePath string `envconfig:"FILE_PATH" yaml:"file_path"`
}

func DefaultSqlite() Sqlite {
	return Sqlite{}
}

type ZipScanner struct {
	MasterAddr  string `envconfig:"MASTER_ADDR" yaml:"master_addr"`
	MasterToken string `envconfig:"MASTER_TOKEN" yaml:"master_token"`
}

func DefaultZipScanner() ZipScanner {
	return ZipScanner{}
}

type Highway struct {
	PrivateKey    string        `envconfig:"PRIVATE_KEY" yaml:"private_key"`
	TokenLifetime time.Duration `envconfig:"TOKEN_LIFETIME" yaml:"token_lifetime"`
}

func DefaultHighway() Highway {
	return Highway{}
}
