package hit

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	lib "github.com/avito-tech/avito-pixel/lib"
	"github.com/avito-tech/avito-pixel/lib/config"
	"github.com/avito-tech/avito-pixel/lib/metrics"
	guuid "github.com/google/uuid"
	jitter "github.com/mroth/jitter"
)

type HitPayload struct {
	SessionID  string         `json:"sid"`
	Platform   string         `json:"string"`
	AppVersion *string        `json:"appVersion,omitempty"`
	Tags       map[string]any `json:"tags"`
}

type requestBody struct {
	Type      string          `json:"type"`
	SessionID string          `json:"sessionId"`
	Meta      requestBodyMeta `json:"meta"`
}

type requestBodyMeta struct {
	Platform string         `json:"platform"`
	Tags     map[string]any `json:"tags"`
}

type Handler struct {
	HitBuffer       []HitPayload
	batchSizeJitter int
	responses       responses
	mu              *sync.Mutex
	logger          lib.Logger
	metrics         metrics.Metrics
	storage         Storage
	flushTicker     jitter.Ticker
	config          config.Collector
}

func NewHandler(
	storage Storage,
	config config.Config,
	logger lib.Logger,
	baseMetrics lib.Metrics,
) *Handler {
	h := &Handler{
		storage: storage,
		batchSizeJitter: rand.Intn(
			int(
				math.Round(
					float64(config.Collector.BatchSize) / 10,
				),
			),
		),
		mu:          &sync.Mutex{},
		logger:      logger,
		config:      config.Collector,
		metrics:     metrics.NewMetrics(baseMetrics),
		flushTicker: *jitter.NewTicker(config.Collector.FlushTimeout, 1),
		responses:   buildResponses(),
	}
	go h.runFlusher(context.Background())
	return h
}

func (h *Handler) Build() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, string(h.responses.NotFound), http.StatusNotFound)
			return
		}
		ctx := r.Context()
		bodyRaw, err := io.ReadAll(r.Body)
		if err != nil {
			_, err := w.Write(h.responses.Failed)
			if err != nil {
				h.logger.Error(ctx, err)
			}
			return
		}

		body := requestBody{}
		err = json.Unmarshal(bodyRaw, &body)
		if err != nil {
			h.logger.Error(ctx, err)
			_, err = w.Write(h.responses.Failed)
			if err != nil {
				h.logger.Error(ctx, err)
			}
			return
		}

		sessionID := ""
		if len(body.SessionID) != 0 {
			sessionID = body.SessionID
		} else {
			sessionIDCookie, err := r.Cookie(h.config.SessionIDCookieName)
			if err != nil || sessionIDCookie.Value == "" {
				h.logger.Debug(ctx, h.config.SessionIDCookieName+" cookie is empty")
				sessionID = guuid.New().String()
				w.Header().Set(h.config.SessionIDCookieName, sessionID)
				expiration := time.Now().Add(365 * 24 * time.Hour)
				cookie := http.Cookie{
					Name:    h.config.SessionIDCookieName,
					Value:   sessionID,
					Expires: expiration,
				}
				http.SetCookie(w, &cookie)
			} else {
				sessionID = sessionIDCookie.Value
			}
		}
		hitPayload := HitPayload{
			SessionID: sessionID,
			Platform:  body.Meta.Platform,
		}
		h.HitBuffer = append(h.HitBuffer, hitPayload)
		if len(h.HitBuffer)-h.batchSizeJitter >= h.config.BatchSize {
			go h.saveHitsAsync(r.Context())
		}
		_, err = w.Write(h.responses.Success)
		if err != nil {
			h.logger.Error(ctx, err)
		}
	})
}
