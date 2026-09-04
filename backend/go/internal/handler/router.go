package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/service"
	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/ws"
)

// RouterConfig wires gateway handlers to their service dependencies.
type RouterConfig struct {
	Service *service.QuokkaService
	Hub     *ws.Hub
	Logger  *slog.Logger
}

// NewRouter builds the public HTTP surface exposed to the Vue application.
func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()
	api := APIHandler{service: cfg.Service, hub: cfg.Hub, logger: cfg.Logger}

	mux.HandleFunc("GET /api/v1/health", api.Health)
	mux.HandleFunc("POST /api/v1/ai/chords", api.GenerateChords)
	mux.HandleFunc("POST /api/v1/ai/mix", api.SmartMix)
	mux.HandleFunc("POST /api/v1/ai/style_transfer", api.StyleTransfer)
	mux.HandleFunc("POST /api/v1/audio/process", api.ProcessAudio)
	mux.HandleFunc("POST /internal/liora_trigger", api.LioraTrigger)
	mux.HandleFunc("/ws/collaboration/", api.CollaborationSocket)

	// Rails remains the authority for persistence-heavy resources during the MVP.
	mux.HandleFunc("/api/", api.ProxyToRails)
	mux.HandleFunc("/rails/active_storage/", api.ProxyToRails)

	return recoverMiddleware(requestLogMiddleware(cfg.Logger, corsMiddleware(mux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func requestLogMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/v1/health") {
			logger.Info("request", "method", r.Method, "path", r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
