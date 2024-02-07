package report

import "net/http"

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

func CsvResponseOk(w http.ResponseWriter, metrics Metrics) error {
	raw, err := ToCsv(metrics)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/csv")
	_, err = w.Write(raw)
	if err != nil {
		return err
	}
	return nil
}
