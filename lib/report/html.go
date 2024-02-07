package report

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"net/http"
)

//go:embed assets/report.html
var reportTemplate string

type htmlPayloadDay struct {
	Date  int    `json:"date"`
	Total uint64 `json:"total"`
}

type htmlPayloadMonth struct {
	Month string           `json:"month"`
	Year  int              `json:"year"`
	Total uint64           `json:"total"`
	Days  []htmlPayloadDay `json:"days"`
}

type htmlPayload struct {
	Total  uint64             `json:"total"`
	Months []htmlPayloadMonth `json:"months"`
}

func (h *Handler) HtmlBuild() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
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

		payload := htmlPayload{}
		var monthlyPayload *htmlPayloadMonth
		for _, dailyReport := range metrics {
			payload.Total += dailyReport.Value
			if monthlyPayload == nil {
				monthlyPayload = &htmlPayloadMonth{
					Month: dailyReport.EventTime.Month().String(),
					Year:  dailyReport.EventTime.Year(),
				}
			}

			if dailyReport.EventTime.Month().String() != monthlyPayload.Month ||
				dailyReport.EventTime.Year() != monthlyPayload.Year {
				payload.Months = append(payload.Months, *monthlyPayload)
				monthlyPayload = &htmlPayloadMonth{
					Month: dailyReport.EventTime.Month().String(),
					Year:  dailyReport.EventTime.Year(),
				}
			}
			monthlyPayload.Total += dailyReport.Value
			monthlyPayload.Days = append(monthlyPayload.Days, htmlPayloadDay{
				Date:  dailyReport.EventTime.Day(),
				Total: dailyReport.Value,
			})
		}

		if monthlyPayload != nil {
			payload.Months = append(payload.Months, *monthlyPayload)
		}

		var payloadBytes []byte
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			h.logger.Error(ctx, "Failed to marshal payload", err)
			err = ResponseFail(w, 500, "Internal error")
			if err != nil {
				h.logger.Error(ctx, err)
			}
			return
		}

		tmpl := template.New("report")
		if tmpl, err = tmpl.Parse(reportTemplate); err != nil {
			h.logger.Error(ctx, "Failed to parse template", err)
			err = ResponseFail(w, 500, "Internal error")
			if err != nil {
				h.logger.Error(ctx, err)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, map[string]string{
			"PayloadString": string(payloadBytes),
		})
		if err != nil {
			h.logger.Error(ctx, err)
			ResponseFail(w, 500, "Internal error")
		}
	})
}
