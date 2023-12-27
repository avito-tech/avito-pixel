package report

import (
	"errors"
	"net/http"
	"strconv"

	lib "github.com/avito-tech/avito-pixel/lib"
	"github.com/avito-tech/avito-pixel/lib/config"
	"github.com/avito-tech/avito-pixel/lib/metrics"
)

type ReportSettings struct {
	Metric   string
	Interval int64
	From     string
	To       string
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

func (h *Handler) JsonBuild() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method != http.MethodPost {
			err := ResponseFail(w, 404, "Route not found")
			if err != nil {
				h.logger.Error(ctx, err)
			}
			return
		}
		reportSettings, err := parseReportSettingsFromQueryParams(r)
		if err != nil {
			h.logger.Error(ctx, err)
			err = ResponseFail(w, 400, "Bad request: could not parse request body")
			if err != nil {
				h.logger.Error(ctx, err)
			}
			return
		}

		metrics, err := h.storage.GetReport(ctx, reportSettings)
		if err != nil {
			h.logger.Error(ctx, err)
			err = ResponseFail(w, 500, "Internal server error: could not retrieve metrics")
			if err != nil {
				h.logger.Error(ctx, err)
			}
			return
		}
		err = JsonResponseOk(w, metrics)
		if err != nil {
			h.logger.Error(ctx, err)
		}
	})
}

// TODO: fix duplicating request parsing and validation
func (h *Handler) CsvBuild() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method != http.MethodGet {
			err := ResponseFail(w, 404, "Route not found")
			if err != nil {
				h.logger.Error(ctx, err)
			}
			return
		}
		reportSettings, err := parseReportSettingsFromQueryParams(r)
		if err != nil {
			h.logger.Error(ctx, err)
			err = ResponseFail(w, 400, "Bad request: could not parse request body")
			if err != nil {
				h.logger.Error(ctx, err)
			}
			return
		}

		metrics, err := h.storage.GetReport(ctx, reportSettings)
		if err != nil {
			h.logger.Error(ctx, err)
			err = ResponseFail(w, 500, "Internal server error: could not retrieve metrics")
			if err != nil {
				h.logger.Error(ctx, err)
			}
			return
		}
		err = CsvResponseOk(w, metrics)
		if err != nil {
			h.logger.Error(ctx, err)
		}
	})
}

func parseReportSettingsFromQueryParams(r *http.Request) (ReportSettings, error) {
	var settings ReportSettings
	queryParams := r.URL.Query()

	settings.Metric = queryParams.Get("metric")
	if settings.Metric == "" {
		return settings, errors.New("metric is required")
	}
	settings.From = queryParams.Get("from")
	if settings.From == "" {
		return settings, errors.New("from is required")
	}
	settings.To = queryParams.Get("to")
	if settings.To == "" {
		return settings, errors.New("to is required")
	}

	intervalParam := queryParams.Get("interval")
	if intervalParam == "" {
		return settings, errors.New("interval is required")
	}

	interval, err := strconv.ParseInt(intervalParam, 10, 64)
	if err != nil {
		return settings, err
	}
	settings.Interval = interval
	return settings, nil
}
