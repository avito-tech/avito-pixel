package report

import (
	"context"
)

type Storage interface {
	GetReport(ctx context.Context, payload ReportSettings) (Metrics, error)
}
