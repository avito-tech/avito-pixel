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
