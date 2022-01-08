package config

var Conf *Config

type Database struct {
	Type        string `toml:"type"`
	DSN         string `toml:"dsn"`
	TablePrefix string `toml:"table_prefix"`
	MaxIdleConn int    `toml:"max_idle_conn"`
	MaxOpenConn int    `toml:"max_open_conn"`
}
type Tracer struct {
	Type        string `toml:"type"`
	Host        string `toml:"host"`
	ServiceName string `toml:"service_name"`
}

type Application struct {
}
type Config struct {
	Database    Database    `toml:"database"`
	Tracer      Tracer      `toml:"tracer"`
	Application Application `toml:"application"`
}
