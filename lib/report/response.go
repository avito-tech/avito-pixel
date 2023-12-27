package report

import (
	"encoding/json"
	"net/http"
)

type jsonResponsePayload struct {
	Data Metrics `json:"data"`
}

type errPayload struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
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

func ResponseFail(w http.ResponseWriter, code int, message string) error {
	resp := errPayload{
		Code:    code,
		Message: message,
	}
	raw, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	http.Error(w, string(raw), code)
	return nil
}
