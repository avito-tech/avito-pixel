package report

import (
	"errors"
	"net/http"
	"strconv"

	lib "github.com/avito-tech/avito-pixel/lib"
	"github.com/avito-tech/avito-pixel/lib/config"
	"github.com/avito-tech/avito-pixel/lib/metrics"
)

type QueryParams struct {
	Metric   string
	Interval string
	From     string
	To       string
	Platform string
}

type ReportSettings struct {
	Metric   string
	Interval int64
	From     string
	To       string
	Platform string
}

type Handler struct {
	logger  lib.Logger
	metrics metrics.Metrics
	storage Storage
	config  config.Collector
}

func NewHandler(
	storage Storage,
	config config.Config,
	logger lib.Logger,
	baseMetrics lib.Metrics,
) *Handler {
	h := &Handler{
		storage: storage,
		logger:  logger,
		config:  config.Collector,
		metrics: metrics.NewMetrics(baseMetrics),
	}
	return h
}

func getReportSettingsFromQueryParams(r *http.Request) QueryParams {
	var settings QueryParams
	queryParams := r.URL.Query()

	settings.Metric = queryParams.Get("metric")
	settings.From = queryParams.Get("from")
	settings.To = queryParams.Get("to")
	settings.Platform = queryParams.Get("platform")
	settings.Interval = queryParams.Get("interval")

	return settings
}

func validateQueryParams(q QueryParams) (ReportSettings, error) {
	var settings ReportSettings

	settings.Metric = q.Metric

	if q.From == "" {
		return settings, errors.New("from is required")
	}
	settings.From = q.From

	if q.To == "" {
		return settings, errors.New("to is required")
	}
	settings.To = q.To

	if q.Interval == "" {
		return settings, errors.New("interval is required")
	}
	var interval, err = strconv.ParseInt(q.Interval, 10, 64)
	if err != nil {
		return settings, err
	}
	settings.Interval = interval
	settings.Platform = q.Platform

	return settings, nil
}

func parseReportSettingsFromQueryParams(r *http.Request) (ReportSettings, error) {
	query := getReportSettingsFromQueryParams(r)
	settings, err := validateQueryParams(query)
	if err != nil {
		return settings, err
	}

	return settings, nil
}
