package logging

import (
	"log/slog"
	"testing"
)

func TestLevelFromEnv(t *testing.T) {
	tests := []struct {
		name string
		want slog.Level
	}{
		{
			name: "DEBUG",
			want: slog.LevelDebug,
		},
		{
			name: "WARNING",
			want: slog.LevelWarn,
		},
		{
			name: "ERROR",
			want: slog.LevelError,
		},
		{
			name: "INFO",
			want: slog.LevelInfo,
		},
		{
			name: "",
			want: slog.LevelInfo,
		},
		{
			name: "other",
			want: slog.LevelInfo,
		},
	}
	for _, tt := range tests {
		t.Setenv("LOG_LEVEL", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			if got := LevelFromEnv(); got != tt.want {
				t.Errorf("LevelFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
