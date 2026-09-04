package service

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/audio"
	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/client/rails"
)

// QuokkaService coordinates public requests across Rails, Redis, and audio workers.
type QuokkaService struct {
	rails  *rails.Client
	audio  *audio.Client
	logger *slog.Logger
}

// NewQuokkaService creates the application service used by HTTP handlers.
func NewQuokkaService(railsClient *rails.Client, audioClient *audio.Client, logger *slog.Logger) *QuokkaService {
	return &QuokkaService{rails: railsClient, audio: audioClient, logger: logger}
}

// ProxyToRails forwards persistence-heavy API calls to the internal Rails service.
func (s *QuokkaService) ProxyToRails(w http.ResponseWriter, r *http.Request) {
	s.rails.Proxy(w, r)
}

// GenerateChords returns a deterministic scaffold until an external AI model is connected.
func (s *QuokkaService) GenerateChords(_ context.Context, request ChordRequest) ChordResponse {
	key := defaultString(request.KeySignature, "C")
	style := defaultString(request.Style, "pop")

	return ChordResponse{
		KeySignature: key,
		Style:        style,
		Progression:  []string{key, "Am", "F", "G"},
		Confidence:   0.74,
		Source:       "go-gateway-scaffold",
	}
}

// SmartMix proposes track gain changes that leave headroom for browser playback.
func (s *QuokkaService) SmartMix(_ context.Context, request MixRequest) MixResponse {
	adjustments := make([]TrackGain, 0, len(request.Tracks))
	for _, track := range request.Tracks {
		gain := 1.0
		if track.Peak > 0.85 {
			gain = 0.85 / track.Peak
		}

		adjustments = append(adjustments, TrackGain{TrackID: track.TrackID, Gain: gain})
	}

	return MixResponse{HeadroomDB: -3.0, Adjustments: adjustments}
}

// StyleTransfer reserves the async task shape for future model-backed rendering.
func (s *QuokkaService) StyleTransfer(_ context.Context, request StyleTransferRequest) StyleTransferResponse {
	return StyleTransferResponse{
		TaskID:    "style-" + time.Now().UTC().Format("20060102150405"),
		Status:    "queued",
		Style:     defaultString(request.Style, "electronic"),
		Transport: "websocket",
	}
}

// ProcessAudio delegates upload handling to Rails until the UDS worker protocol is completed.
func (s *QuokkaService) ProcessAudio(_ context.Context, r *http.Request) (map[string]any, error) {
	s.logger.Info("audio process request routed through go gateway", "socket", s.audio.SocketPath())
	return s.rails.ForwardMultipart(r, "/api/v1/audio/process")
}

// AcceptLioraTrigger records the integration boundary for Liora sound-module events.
func (s *QuokkaService) AcceptLioraTrigger(_ context.Context, payload map[string]any) map[string]any {
	s.logger.Info("liora trigger accepted", "payload", payload)
	return map[string]any{
		"status":      "accepted",
		"integration": "liora",
		"handled_by":  "go-gateway",
		"received_at": time.Now().UTC().Format(time.RFC3339),
	}
}

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
