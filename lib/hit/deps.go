package hit

import (
	"context"

	"github.com/avito-tech/avito-pixel/lib/storage"
)

type Storage interface {
	SaveHits(ctx context.Context, hits []storage.HitTable) error
}
