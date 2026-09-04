package service

// ChordRequest describes the melody context used for chord suggestions.
type ChordRequest struct {
	KeySignature string         `json:"key_signature"`
	Style        string         `json:"style"`
	Notes        []MelodyNote   `json:"notes"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

type MelodyNote struct {
	Name      string  `json:"name"`
	Frequency float64 `json:"frequency"`
	StartBeat float64 `json:"start_beat"`
	Duration  float64 `json:"duration"`
}

type ChordResponse struct {
	KeySignature string   `json:"key_signature"`
	Style        string   `json:"style"`
	Progression  []string `json:"progression"`
	Confidence   float64  `json:"confidence"`
	Source       string   `json:"source"`
}

type MixRequest struct {
	Tracks []TrackPeak `json:"tracks"`
}

type TrackPeak struct {
	TrackID string  `json:"track_id"`
	Peak    float64 `json:"peak"`
}

type TrackGain struct {
	TrackID string  `json:"track_id"`
	Gain    float64 `json:"gain"`
}

type MixResponse struct {
	HeadroomDB  float64     `json:"headroom_db"`
	Adjustments []TrackGain `json:"adjustments"`
}

type StyleTransferRequest struct {
	CompositionID int64  `json:"composition_id"`
	Style         string `json:"style"`
}

type StyleTransferResponse struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	Style     string `json:"style"`
	Transport string `json:"transport"`
}
