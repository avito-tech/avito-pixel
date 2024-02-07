package report

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) JsonBuild() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// if r.Method != http.MethodPost {
		// 	err := ResponseFail(w, 404, "Route not found")
		// 	if err != nil {
		// 		h.logger.Error(ctx, err)
		// 	}
		// 	return
		// }
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

func JsonResponseOk(w http.ResponseWriter, metrics Metrics) error {
	resp := jsonResponsePayload{
		Data: metrics,
	}
	raw, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(raw)
	if err != nil {
		return err
	}
	return nil
}
