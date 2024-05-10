package plugin

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"code.gitea.io/sdk/gitea"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thegeeklab/wp-gitea-release/plugin/mocks"
)

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
		mockClient := mocks.NewMockIGiteaClient(t)
		r := &GiteaRelease{
			Opt:    tt.opt,
			client: mockClient,
		}

		mockClient.
			On("ListReleases", mock.Anything, mock.Anything, mock.Anything).
			Return([]*gitea.Release{
				{
					ID:           1,
					TagName:      "v1.0.0",
					Title:        "Release v1.0.0",
					Note:         "This is the release notes for v1.0.0",
					IsDraft:      false,
					IsPrerelease: false,
				},
			}, nil, nil)

		t.Run(tt.name, func(t *testing.T) {
			release, err := r.Find()

			if tt.wantErr != nil {
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
		mockClient := mocks.NewMockIGiteaClient(t)
		r := &GiteaRelease{
			Opt:    tt.opt,
			client: mockClient,
		}

		mockClient.
			On("CreateRelease", mock.Anything, mock.Anything, mock.Anything).
			Return(&gitea.Release{
				ID:           1,
				TagName:      tt.opt.Tag,
				Title:        tt.opt.Title,
				Note:         tt.opt.Note,
				IsDraft:      tt.opt.Draft,
				IsPrerelease: tt.opt.Prerelease,
			}, nil, nil)

		t.Run(tt.name, func(t *testing.T) {
			release, err := r.Create()

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
		logBuffer.Reset()

		mockClient := mocks.NewMockIGiteaClient(t)
		r := &GiteaRelease{
			Opt:    tt.opt,
			client: mockClient,
		}

		mockClient.
			On("ListReleaseAttachments", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]*gitea.Attachment{
				{
					Name: "file1.txt",
				},
			}, nil, nil)

		if FileExists(tt.opt.FileExists) == FileExistsOverwrite {
			mockClient.
				On("DeleteReleaseAttachment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(nil, nil)
		}

		if tt.wantErr == nil {
			mockClient.
				On("CreateReleaseAttachment", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(nil, nil, nil)
		}

		t.Run(tt.name, func(t *testing.T) {
			err := r.AddAttachments(1, tt.files)

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
