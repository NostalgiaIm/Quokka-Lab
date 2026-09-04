package config

import "os"

// Config contains runtime settings for the Go API gateway.
type Config struct {
	HTTPAddr        string
	RailsURL        string
	AudioSocketPath string
	AudioBinaryPath string
}

// Load reads environment variables and applies development-safe defaults.
func Load() Config {
	return Config{
		HTTPAddr:        getenv("QUOKKA_GO_ADDR", ":8080"),
		RailsURL:        getenv("QUOKKA_RAILS_URL", "http://rails:3000"),
		AudioSocketPath: getenv("QUOKKA_AUDIO_SOCKET", "/tmp/quokka-audio.sock"),
		AudioBinaryPath: getenv("QUOKKA_AUDIO_ENGINE", "/audio_engine/build/quokka_audio"),
	}
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
