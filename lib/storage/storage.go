package storage

import (
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	deps "github.com/avito-tech/avito-pixel/lib"
	"github.com/avito-tech/avito-pixel/lib/config"
)

type Clickhouse struct {
	DB     clickhouse.Conn
	conf   config.Clickhouse
	logger deps.Logger
}

func NewClickhouse(
	conf config.Clickhouse,
	logger deps.Logger,
) *Clickhouse {
	db, err := newClient(conf)
	if err != nil {
		panic(err)
	}

	c := &Clickhouse{
		DB:     db,
		conf:   conf,
		logger: logger,
	}

	return c
}

func newClient(conf config.Clickhouse) (driver.Conn, error) {
	opts := clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", conf.Host, conf.Port)},
		Auth: clickhouse.Auth{
			Database: conf.Database,
			Username: conf.User,
			Password: conf.Password,
		},
		MaxOpenConns:    conf.MaxOpenConns,
		MaxIdleConns:    conf.MaxIdleConns,
		ConnMaxLifetime: conf.ConnLifetime,
		Settings: map[string]interface{}{
			"max_threads": conf.MaxThreads,
		},
	}
	db, err := clickhouse.Open(&opts)
	return db, err
}
