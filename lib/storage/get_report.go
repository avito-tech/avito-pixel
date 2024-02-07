package storage

import (
	"context"

	"github.com/avito-tech/avito-pixel/lib/report"
)

func (c *Clickhouse) GetReport(
	ctx context.Context,
	payload report.ReportSettings,
) (report.Metrics, error) {
	query := "SELECT eventTime, sum(totalHits) as value FROM visitors_1_day_mv GROUP BY eventTime HAVING eventTime >= ? AND eventTime <= ? ORDER BY eventTime ASC"
	var data report.Metrics
	if err := c.DB.Select(ctx, &data, query, payload.From, payload.To); err != nil {
		return data, err
	}
	return data, nil
}
