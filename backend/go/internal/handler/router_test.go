package handler

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/audio"
	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/client/rails"
	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/service"
	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/ws"
)

func TestHealthEndpoint(t *testing.T) {
	logger := slog.Default()
	railsClient := rails.NewClient("http://rails.internal", logger)
	audioClient := audio.NewClient("/tmp/quokka-audio.sock", "/audio_engine/build/quokka_audio", logger)
	router := NewRouter(RouterConfig{
		Service: service.NewQuokkaService(railsClient, audioClient, logger),
		Hub:     ws.NewHub(logger),
		Logger:  logger,
	})

	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, response.Code)
	}
}
