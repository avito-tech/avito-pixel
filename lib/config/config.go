package config

import (
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Clickhouse Clickhouse
	Collector  Collector
}

type Clickhouse struct {
	User         string        `env:"CLICKHOUSE_USERNAME,required"`
	Password     string        `env:"CLICKHOUSE_PASSWORD,required"`
	Host         string        `env:"PROXY_CLICKHOUSE_HOST,required"`
	Port         int           `env:"PROXY_CLICKHOUSE_PORT,required"`
	Database     string        `env:"CLICKHOUSE_DATABASE,required"`
	MaxOpenConns int           `env:"CLICKHOUSE_MAX_OPEN_CONNS,required"`
	MaxIdleConns int           `env:"CLICKHOUSE_MAX_IDLE_CONNS,required"`
	ConnLifetime time.Duration `env:"CLICKHOUSE_CONN_LIFETIME,required"`
	MaxThreads   int           `env:"CLICKHOUSE_MAX_THREADS,required"`
	SkipVerify   bool          `env:"CLICKHOUSE_SKIP_VERIFY"`
	CaCertFile   string        `env:"CA_CERT_FILE"`
	Cert         string        `env:"CERT_FILE"`
	CertKey      string        `env:"KEY_FILE"`
}

type Collector struct {
	BatchSize           int           `env:"COLLECTOR_BATCH_SIZE,required"`
	FlushTimeout        time.Duration `env:"COLLECTOR_FLUSH_TIMEOUT,required"`
	SessionIDCookieName string        `env:"SESSION_ID_COOKIE_NAME,required"`
}

func Init() (Config, error) {
	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
