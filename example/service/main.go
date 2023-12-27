package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/avito-tech/avito-pixel/example/service/extras"
	cfg "github.com/avito-tech/avito-pixel/lib/config"
	"github.com/avito-tech/avito-pixel/lib/hit"
	"github.com/avito-tech/avito-pixel/lib/report"
	"github.com/avito-tech/avito-pixel/lib/storage"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("1111")
	logger := extras.NewLogger()
	err := godotenv.Load()
	if err != nil {
		logger.Error(context.Background(), "Error loading .env file")
	}

	mux := http.NewServeMux()
	appConfig, err := cfg.Init()
	if err != nil {
		panic(err)
	}

	metrics := extras.NewMetrics()

	storage := storage.NewClickhouse(appConfig.Clickhouse, logger)
	hitHandler := hit.NewHandler(storage, appConfig, logger, metrics)
	reportHandler := report.NewHandler(storage, appConfig, logger, metrics)

	mux.Handle("/hit/", hitHandler.Build())
	mux.Handle("/report/csv/", reportHandler.CsvBuild())
	mux.Handle("/report/json/", reportHandler.JsonBuild())
	http.ListenAndServe(":3000", mux)

	fmt.Println("azazax")
}
