package hit

import (
	"encoding/json"
)

type responsePayload struct {
	Ok bool `json:"ok"`
}

type notFoundPayload struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type responses struct {
	Success  []byte
	Failed   []byte
	NotFound string
}

func buildResponses() responses {
	res := responses{}
	successResponseRaw, err := json.Marshal(
		responsePayload{Ok: true},
	)
	if err != nil {
		panic("Error to marshal success response")
	}
	res.Success = []byte(successResponseRaw)

	failedResponseRaw, err := json.Marshal(
		responsePayload{Ok: false},
	)
	if err != nil {
		panic("Error to marshal failed response")
	}
	res.Failed = []byte(failedResponseRaw)

	notFoundResponseRaw, err := json.Marshal(
		notFoundPayload{
			Code:    404,
			Message: "route not found",
		},
	)
	if err != nil {
		panic("Error to marshal notFound response")
	}
	res.NotFound = string(notFoundResponseRaw)
	return res
}
