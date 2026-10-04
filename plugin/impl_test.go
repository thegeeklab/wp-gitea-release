package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name       string
		event      string
		fileExists string
		wantErr    error
	}{
		{
			name:       "valid tag event",
			event:      "tag",
			fileExists: "overwrite",
		},
		{
			name:       "valid tag event with skip",
			event:      "tag",
			fileExists: "skip",
		},
		{
			name:       "non tag event is rejected",
			event:      "push",
			fileExists: "overwrite",
			wantErr:    ErrPluginEventNotSupported,
		},
		{
			name:       "invalid file exists value",
			event:      "tag",
			fileExists: "bogus",
			wantErr:    ErrFileExistInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PLUGIN_API_KEY", "test-token")
			t.Setenv("PLUGIN_BASE_URL", "https://gitea.example.com")

			p := New(func(_ context.Context) error { return nil })
			_ = p.App.Run(t.Context(), []string{"wp-gitea-release"})

			p.Settings = &Settings{
				Event:      tt.event,
				FileExists: tt.fileExists,
			}

			err := p.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)

				return
			}

			assert.NoError(t, err)
		})
	}
}
