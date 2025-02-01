package config

import (
	"log"
	"path/filepath"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	ServerConfig  ServerConfig  `json:"server" yaml:"server"`
	DBConfig      DBConfig      `json:"db"`
	LoggingConfig LoggingConfig `json:"logging" yaml:"logging"`
	TraceConfig   TraceConfig   `json:"trace" yaml:"trace"`
}

type LoggingConfig struct {
	Level       int                `json:"level"`
	Encoding    string             `json:"encoding"`
	Development bool               `json:"development"`
	Info        InfoLoggingConfig  `json:"infoLevel" yaml:"infoLevel"`
	Error       ErrorLoggingConfig `json:"errorLevel" yaml:"errorLevel"`
}

type InfoLoggingConfig struct {
	Filename   string `json:"filename"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
	Compress   bool   `json:"compress"`
}

type ErrorLoggingConfig struct {
	Filename   string `json:"filename"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
	Compress   bool   `json:"compress"`
}

type ServerConfig struct {
	Port             int           `json:"port"`
	ReadTimeout      time.Duration `json:"readTimeout"`
	WriteTimeout     time.Duration `json:"writeTimeout"`
	GracefulShutdown time.Duration `json:"gracefulShutdown"`
	GoroutineTimeout time.Duration `json:"goroutineTimeout"`
	Host             string        `json:"host"`
}

type DBConfig struct {
	DataSourceName string `json:"dataSourceName"`
	LogLevel       int    `json:"logLevel"`
	Pool           struct {
		MaxOpen     int           `json:"maxOpen"`
		MaxIdle     int           `json:"maxIdle"`
		MaxLifetime time.Duration `json:"maxLifetime"`
	} `json:"pool"`
}

type TraceConfig struct {
	IsTraceEnabled    bool   `json:"is_enabled"`
	Url               string `json:"trace_url"`
	ServiceName       string `json:"trace_service_name"`
	IsHttpBodyEnabled bool   `json:"trace_is_http_body_enabled"`
}

func Load() (*Config, error) {
	k := koanf.New(".")

	path, err := filepath.Abs("./config/config.local.yaml")
	if err != nil {
		log.Printf("failed to get absoulute config path. configPath:%s, err: %v", "./config/config.local.yaml", err)
		return nil, err
	}

	log.Printf("load config file from %s", path)
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		log.Printf("failed to load config from file. err: %v", err)
		return nil, err
	}

	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "json", FlatPaths: false}); err != nil {
		log.Printf("failed to unmarshal with conf. err: %v", err)
		return nil, err
	}
	return &cfg, err
}
