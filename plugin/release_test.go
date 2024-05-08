package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog/log"

	"code.gitea.io/sdk/gitea"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func giteaMockHandler(t *testing.T, opt GiteaReleaseOpt) func(http.ResponseWriter, *http.Request) {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Println(r.RequestURI)

		switch r.RequestURI {
		case "/api/v1/version":
			_, err := io.WriteString(w, `{"version":"1.21.0"}`)
			if err != nil {
				t.Fail()
			}
		case "/api/v1/repos/test-owner/test-repo/releases?limit=0&page=1":
			_, err := io.WriteString(w, `[{
				"id": 1,
				"tag_name": "v1.0.0",
				"name": "Release v1.0.0",
				"body": "This is the release notes for v1.0.0",
				"draft": false,
				"prerelease": false,
				"created_at": "2023-05-01T12:00:00Z",
				"published_at": "2023-05-01T12:30:00Z"
			}]`)
			if err != nil {
				t.Fail()
			}
		case "/api/v1/repos/test-owner/test-repo/releases":
			_, err := io.WriteString(w, fmt.Sprintf(`{
				"id": 1,
				"tag_name": "%s",
				"name": "Release %s",
				"body": "This is the release notes for %s",
				"draft": %t,
				"prerelease": %t,
				"created_at": "2023-05-01T12:00:00Z",
				"published_at": "2023-05-01T12:30:00Z"
			}`, opt.Tag, opt.Tag, opt.Tag, opt.Draft, opt.Prerelease))
			if err != nil {
				t.Fail()
			}
		case "/api/v1/repos/test-owner/test-repo/releases/1/assets?limit=0&page=1":
			_, err := io.WriteString(w, `[{
				"id": 1,
				"name": "file1.txt",
				"size": 1024,
				"created_at": "2023-05-01T12:30:00Z"
			}]`)
			if err != nil {
				t.Fail()
			}
		case "/api/v1/repos/test-owner/test-repo/releases/1/assets":
			_, err := io.WriteString(w, `{
				"id": 1,
				"name": "file1.txt",
				"size": 1024,
				"created_at": "2023-05-01T12:30:00Z"
			}`)
			if err != nil {
				t.Fail()
			}
		}
	}
}

func TestGiteaReleaseFind(t *testing.T) {
	tests := []struct {
		name    string
		opt     GiteaReleaseOpt
		want    *gitea.Release
		wantErr error
	}{
		{
			name: "find release by tag",
			opt: GiteaReleaseOpt{
				Owner: "test-owner",
				Repo:  "test-repo",
				Tag:   "v1.0.0",
			},
			want: &gitea.Release{
				TagName: "v1.0.0",
			},
		},
		{
			name: "release not found",
			opt: GiteaReleaseOpt{
				Owner: "test-owner",
				Repo:  "test-repo",
				Tag:   "v1.1.0",
			},
			want:    nil,
			wantErr: ErrReleaseNotFound,
		},
	}

	for _, tt := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			giteaMockHandler(t, tt.opt)(w, r)
		}))
		defer ts.Close()

		g, _ := gitea.NewClient(ts.URL)
		client := NewGiteaClient(g)

		t.Run(tt.name, func(t *testing.T) {
			client.Release.Opt = tt.opt
			release, err := client.Release.Find()

			if tt.want == nil {
				assert.Error(t, err)
				assert.Nil(t, release)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.TagName, release.TagName)
		})
	}
}

func TestGiteaReleaseCreate(t *testing.T) {
	tests := []struct {
		name    string
		opt     GiteaReleaseOpt
		want    *gitea.Release
		wantErr error
	}{
		{
			name: "create release",
			opt: GiteaReleaseOpt{
				Owner:      "test-owner",
				Repo:       "test-repo",
				Tag:        "v1.1.0",
				Title:      "Release v1.1.0",
				Note:       "This is the release notes for v1.1.0",
				Draft:      false,
				Prerelease: false,
			},
			want: &gitea.Release{
				TagName:      "v1.1.0",
				Title:        "Release v1.1.0",
				Note:         "This is the release notes for v1.1.0",
				IsDraft:      false,
				IsPrerelease: false,
			},
		},
		{
			name: "create draft release",
			opt: GiteaReleaseOpt{
				Owner:      "test-owner",
				Repo:       "test-repo",
				Tag:        "v1.2.0",
				Title:      "Release v1.2.0",
				Note:       "This is the release notes for v1.2.0",
				Draft:      true,
				Prerelease: false,
			},
			want: &gitea.Release{
				TagName:      "v1.2.0",
				Title:        "Release v1.2.0",
				Note:         "This is the release notes for v1.2.0",
				IsDraft:      true,
				IsPrerelease: false,
			},
		},
		{
			name: "create prerelease",
			opt: GiteaReleaseOpt{
				Owner:      "test-owner",
				Repo:       "test-repo",
				Tag:        "v1.3.0-rc1",
				Title:      "Release v1.3.0-rc1",
				Note:       "This is the release notes for v1.3.0-rc1",
				Draft:      false,
				Prerelease: true,
			},
			want: &gitea.Release{
				TagName:      "v1.3.0-rc1",
				Title:        "Release v1.3.0-rc1",
				Note:         "This is the release notes for v1.3.0-rc1",
				IsDraft:      false,
				IsPrerelease: true,
			},
		},
	}

	for _, tt := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			giteaMockHandler(t, tt.opt)(w, r)
		}))
		defer ts.Close()

		g, _ := gitea.NewClient(ts.URL)
		client := NewGiteaClient(g)

		t.Run(tt.name, func(t *testing.T) {
			client.Release.Opt = tt.opt
			release, err := client.Release.Create()

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, release)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.TagName, release.TagName)
			assert.Equal(t, tt.want.Title, release.Title)
			assert.Equal(t, tt.want.Note, release.Note)
			assert.Equal(t, tt.want.IsDraft, release.IsDraft)
			assert.Equal(t, tt.want.IsPrerelease, release.IsPrerelease)
		})
	}
}

func TestGiteaReleaseAddAttachments(t *testing.T) {
	logBuffer := &bytes.Buffer{}
	logger := zerolog.New(logBuffer)
	log.Logger = logger

	tests := []struct {
		name       string
		opt        GiteaReleaseOpt
		files      []string
		fileExists string
		wantErr    error
		wantLogs   []string
	}{
		{
			name: "add new attachments",
			opt: GiteaReleaseOpt{
				Owner:      "test-owner",
				Repo:       "test-repo",
				Tag:        "v2.0.0",
				Title:      "Release v2.0.0",
				FileExists: "overwrite",
			},
			files:    []string{createTempFile(t, "file1.txt"), createTempFile(t, "file2.txt")},
			wantLogs: []string{"uploaded artifact: file1.txt", "uploaded artifact: file2.txt"},
		},
		{
			name: "fail on existing attachments",
			opt: GiteaReleaseOpt{
				Owner:      "test-owner",
				Repo:       "test-repo",
				Tag:        "v2.0.0",
				Title:      "Release v2.0.0",
				FileExists: "fail",
			},
			files:   []string{createTempFile(t, "file1.txt"), createTempFile(t, "file2.txt")},
			wantErr: ErrFileExists,
		},
		{
			name: "overwrite on existing attachments",
			opt: GiteaReleaseOpt{
				Owner:      "test-owner",
				Repo:       "test-repo",
				Tag:        "v2.0.0",
				Title:      "Release v2.0.0",
				FileExists: "overwrite",
			},
			files:    []string{createTempFile(t, "file1.txt"), createTempFile(t, "file2.txt")},
			wantErr:  nil,
			wantLogs: []string{"deleted artifact: file1.txt", "uploaded artifact: file1.txt"},
		},
		{
			name: "skip on existing attachments",
			opt: GiteaReleaseOpt{
				Owner:      "test-owner",
				Repo:       "test-repo",
				Tag:        "v2.0.0",
				Title:      "Release v2.0.0",
				FileExists: "skip",
			},
			files:    []string{createTempFile(t, "file1.txt"), createTempFile(t, "file2.txt")},
			wantErr:  nil,
			wantLogs: []string{"skip existing artifact: file1"},
		},
		{
			name: "fail on invalid file",
			opt: GiteaReleaseOpt{
				Owner:      "test-owner",
				Repo:       "test-repo",
				Tag:        "v2.0.0",
				Title:      "Release v2.0.0",
				FileExists: "overwrite",
			},
			files:   []string{"testdata/file1.txt", "testdata/invalid.txt"},
			wantErr: errors.New("no such file or directory"),
		},
	}

	for _, tt := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			giteaMockHandler(t, tt.opt)(w, r)
		}))
		defer ts.Close()

		logBuffer.Reset()

		g, _ := gitea.NewClient(ts.URL)
		client := NewGiteaClient(g)

		t.Run(tt.name, func(t *testing.T) {
			client.Release.Opt = tt.opt
			release, _ := client.Release.Create()

			err := client.Release.AddAttachments(release.ID, tt.files)

			// Assert log output.
			for _, l := range tt.wantLogs {
				assert.Contains(t, logBuffer.String(), l)
			}

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr.Error())

				return
			}

			assert.NoError(t, err)
		})
	}
}

func createTempFile(t *testing.T, name string) string {
	t.Helper()

	name = filepath.Join(t.TempDir(), name)
	_ = os.WriteFile(name, []byte("hello"), 0o600)

	return name
}
