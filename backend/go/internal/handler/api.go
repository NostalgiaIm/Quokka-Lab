package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/service"
	"github.com/NostalgiaIm/Quokka-Lab/backend/go/internal/ws"
)

// APIHandler owns the public REST and WebSocket endpoints of the Go gateway.
type APIHandler struct {
	service *service.QuokkaService
	hub     *ws.Hub
	logger  *slog.Logger
}

func (h APIHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]any{
		"service":   "quokka-go-gateway",
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h APIHandler) GenerateChords(w http.ResponseWriter, r *http.Request) {
	var request service.ChordRequest
	if err := decodeJSON(r, &request); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, h.service.GenerateChords(r.Context(), request))
}

func (h APIHandler) SmartMix(w http.ResponseWriter, r *http.Request) {
	var request service.MixRequest
	if err := decodeJSON(r, &request); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, h.service.SmartMix(r.Context(), request))
}

func (h APIHandler) StyleTransfer(w http.ResponseWriter, r *http.Request) {
	var request service.StyleTransferRequest
	if err := decodeJSON(r, &request); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusAccepted, h.service.StyleTransfer(r.Context(), request))
}

func (h APIHandler) ProcessAudio(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.ProcessAudio(r.Context(), r)
	if err != nil {
		h.logger.Error("audio processing failed", "error", err)
		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, result)
}

func (h APIHandler) LioraTrigger(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any
	if err := decodeJSON(r, &payload); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusAccepted, h.service.AcceptLioraTrigger(r.Context(), payload))
}

func (h APIHandler) CollaborationSocket(w http.ResponseWriter, r *http.Request) {
	roomID := strings.TrimPrefix(r.URL.Path, "/ws/collaboration/")
	if roomID == "" {
		h.writeError(w, http.StatusBadRequest, "room id is required")
		return
	}

	h.hub.ServeRoom(w, r, roomID)
}

func (h APIHandler) ProxyToRails(w http.ResponseWriter, r *http.Request) {
	h.service.ProxyToRails(w, r)
}

func (h APIHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("json response failed", "error", err)
	}
}

func (h APIHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
