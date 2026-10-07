package config

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    Config
		wantErr []string // substrings that must all appear in the error
	}{
		{
			name: "defaults applied",
			env:  map[string]string{"DATABASE_URL": "postgres://x"},
			want: Config{DatabaseURL: "postgres://x", HTTPAddr: ":8080", LogFormat: "json"},
		},
		{
			name: "explicit values win",
			env: map[string]string{
				"DATABASE_URL": "postgres://x",
				"HTTP_ADDR":    ":9000",
				"LOG_FORMAT":   "text",
			},
			want: Config{DatabaseURL: "postgres://x", HTTPAddr: ":9000", LogFormat: "text"},
		},
		{
			name:    "missing database url",
			env:     map[string]string{},
			wantErr: []string{"DATABASE_URL is required"},
		},
		{
			name:    "all errors reported together",
			env:     map[string]string{"LOG_FORMAT": "xml"},
			wantErr: []string{"DATABASE_URL is required", `LOG_FORMAT must be json or text, got "xml"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(key string) string { return tt.env[key] }

			got, err := Load(getenv)

			if len(tt.wantErr) > 0 {
				if err == nil {
					t.Fatalf("expected error, got config %+v", got)
				}
				for _, want := range tt.wantErr {
					if !strings.Contains(err.Error(), want) {
						t.Errorf("error %q does not contain %q", err, want)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
