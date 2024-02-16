package storage

import (
	"context"
	"github.com/avito-tech/avito-pixel/lib/report"
)

func (c *Clickhouse) GetPlatforms(
	ctx context.Context,
) (report.PlatformList, error) {

	var data report.PlatformList
	var query string

	query = "SELECT if(platform = 'web', 'mav', platform) as platform FROM visitors_1_day_mv GROUP BY platform"

	if err := c.DB.Select(ctx, &data, query); err != nil {
		return data, err
	}

	return data, nil
}
