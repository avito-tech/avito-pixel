package storage

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"strings"

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

	var (
		tlsConf      *tls.Config
		clientCAs    *x509.CertPool
		certificates []tls.Certificate
	)

	if len(strings.TrimSpace(conf.CaCertFile)) > 0 {
		caCert, err := os.ReadFile(path.Join(conf.CaCertFile))
		if err != nil {
			panic(err)
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			panic("failed to append from pem")
		}
		clientCAs = caCertPool
	}

	if len(string(conf.Cert)) > 0 && len(string(conf.CertKey)) > 0 {
		cert, err := tls.LoadX509KeyPair(conf.Cert, conf.CertKey)
		if err != nil {
			panic(err)
		}
		certificates = []tls.Certificate{cert}
	}

	if conf.SkipVerify != nil || len(certificates) > 0 || clientCAs != nil {
		tlsConf = &tls.Config{
			Certificates: certificates,
			ClientCAs:    clientCAs,
		}
		if conf.SkipVerify != nil && *conf.SkipVerify {
			tlsConf.InsecureSkipVerify = true
		}
	}

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
		TLS:             tlsConf,
		Settings: map[string]interface{}{
			"max_threads": conf.MaxThreads,
		},
	}
	db, err := clickhouse.Open(&opts)
	return db, err
}
