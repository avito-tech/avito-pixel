package storage

import (
	"context"
	"fmt"

	"github.com/avito-tech/avito-pixel/lib/report"
)

func (c *Clickhouse) GetReport(
	ctx context.Context,
	payload report.ReportSettings,
) (report.Metrics, error) {

	var data report.Metrics
	var arguments []any
	var query string

	template := "SELECT eventTime, count(*) as value FROM visitors_1_day_mv WHERE eventTime >= ? AND eventTime <= ? %s GROUP BY eventTime ORDER BY eventTime ASC"

	if payload.Platform != "" {
		query = fmt.Sprintf(template, "AND platform = ?")
		arguments = append(arguments, payload.From, payload.To, payload.Platform)
	} else {
		query = fmt.Sprintf(template, "")
		arguments = append(arguments, payload.From, payload.To)
	}

	if err := c.DB.Select(ctx, &data, query, arguments...); err != nil {
		return data, err
	}

	return data, nil
}
