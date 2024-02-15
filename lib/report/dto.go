package report

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"time"
)

type MetricReport struct {
	EventTime time.Time `json:"date" ch:"eventTime"`
	Platform  string    `json:"platform" ch:"platform"`
	Value     uint64    `json:"value" ch:"value"`
}

type Metrics = []MetricReport

func ToCsv(metrics Metrics) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"date", "value", "platform"}
	err := writer.Write(header)
	if err != nil {
		return nil, err
	}

	for _, v := range metrics {
		record := []string{v.EventTime.Format("2006-01-02 15:04:05"), strconv.Itoa(int(v.Value)), v.Platform}
		err = writer.Write(record)
		if err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}
