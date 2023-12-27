package storage

import (
	"context"
)

func (clickhouse *Clickhouse) SaveHits(
	ctx context.Context,
	hits []HitTable,
) error {
	batch, err := clickhouse.DB.PrepareBatch(ctx, "INSERT INTO hits")
	if err != nil {
		clickhouse.logger.Error(
			ctx,
			"failed to prepare clickhouse insert; ",
			err,
		)
		return err
	}

	for _, hit := range hits {
		err := batch.Append(
			hit.EventTime,
			hit.EventType,
			hit.SessionID,
			hit.Platform,
			hit.Meta,
		)
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
