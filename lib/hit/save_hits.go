package hit

import (
	"context"
	"fmt"
	"time"

	"github.com/avito-tech/avito-pixel/lib/storage"
	"github.com/mroth/jitter"
)

func (h *Handler) saveHits(ctx context.Context) error {
	h.mu.Lock()
	h.flushTicker.Stop()
	h.flushTicker = *jitter.NewTicker(h.config.FlushTimeout, 0.1)
	if len(h.HitBuffer) == 0 {
		h.mu.Unlock()
		return nil
	}
	sendTimer := h.metrics.Timer("save_hits.timer")
	h.logger.Debug(ctx, fmt.Sprintf("saving %d hits", len(h.HitBuffer)))
	hits := mapHits(h.HitBuffer)
	h.HitBuffer = []HitPayload{}
	err := h.storage.SaveHits(ctx, hits)
	if err != nil {
		h.metrics.Increment("save_hits.count.fail")
	} else {
		h.metrics.Increment("save_hits.count.success")
	}
	sendTimer()
	h.mu.Unlock()
	return err
}

func (h *Handler) saveHitsAsync(ctx context.Context) {
	err := h.saveHits(ctx)
	if err != nil {
		h.logger.Error(ctx, err)
	}
}

func mapHits(hits []HitPayload) []storage.HitTable {
	res := []storage.HitTable{}
	for _, h := range hits {
		meta := map[string]any{}
		if h.AppVersion != nil {
			meta["app_version"] = *h.AppVersion
		}
		res = append(res, storage.HitTable{
			SessionID: h.SessionID,
			EventTime: time.Now(),
			EventType: "page_view",
			Platform:  h.Platform,
			Meta:      meta,
		})
	}
	return res
}

func (h *Handler) runFlusher(ctx context.Context) {
	defer func() {
		h.flushTicker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			err := h.saveHits(ctx)
			if err != nil {
				h.logger.Error(ctx, err)
			}
			return
		case <-h.flushTicker.C:
			err := h.saveHits(context.Background())
			if err != nil {
				h.logger.Error(ctx, err)
			}
		}
	}
}
