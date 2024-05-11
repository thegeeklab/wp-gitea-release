package gitea

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"code.gitea.io/sdk/gitea"
	"github.com/rs/zerolog/log"
)

var (
	ErrReleaseNotFound = errors.New("release not found")
	ErrFileExists      = errors.New("asset file already exist")
)

const (
	FileExistsOverwrite FileExists = "overwrite"
	FileExistsFail      FileExists = "fail"
	FileExistsSkip      FileExists = "skip"
)

type Client struct {
	client  APIClient
	Release *Release
}

type Release struct {
	client APIClient
	Opt    ReleaseOptions
}

type ReleaseOptions struct {
	Owner      string
	Repo       string
	Tag        string
	Draft      bool
	Prerelease bool
	FileExists string
	Title      string
	Note       string
}

type FileExists string

// NewClient creates a new Client instance with the provided Gitea client.
func NewClient(url, key string, client *http.Client) (*Client, error) {
	c, err := gitea.NewClient(url, gitea.SetToken(key), gitea.SetHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: c,
		Release: &Release{
			client: c,
			Opt:    ReleaseOptions{},
		},
	}, nil
}

// Find retrieves the release with the specified tag name from the repository.
// If the release is not found, it returns an ErrReleaseNotFound error.
func (r *Release) Find() (*gitea.Release, error) {
	releases, _, err := r.client.ListReleases(r.Opt.Owner, r.Opt.Repo, gitea.ListReleasesOptions{})
	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		if release.TagName == r.Opt.Tag {
			log.Info().Msgf("found release: %s", r.Opt.Tag)

			return release, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrReleaseNotFound, r.Opt.Tag)
}

// Create creates a new release on the Gitea repository with the specified options.
// It returns the created release or an error if the creation failed.
func (r *Release) Create() (*gitea.Release, error) {
	opts := gitea.CreateReleaseOption{
		TagName:      r.Opt.Tag,
		IsDraft:      r.Opt.Draft,
		IsPrerelease: r.Opt.Prerelease,
		Title:        r.Opt.Title,
		Note:         r.Opt.Note,
	}

	release, _, err := r.client.CreateRelease(r.Opt.Owner, r.Opt.Repo, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create release: %w", err)
	}

	log.Info().Msgf("created release: %s", r.Opt.Tag)

	return release, nil
}

// AddAttachments uploads the specified files as attachments to the release with the given ID.
// It first checks for any existing attachments with the same names,
// and handles them according to the FileExists option:
//
// - "overwrite": overwrites the existing attachment
// - "fail": returns an error if the file already exists
// - "skip": skips uploading the file and logs a warning
//
// If there are no conflicts, it uploads the new files as attachments to the release.
func (r *Release) AddAttachments(releaseID int64, files []string) error {
	attachments, _, err := r.client.ListReleaseAttachments(
		r.Opt.Owner,
		r.Opt.Repo,
		releaseID,
		gitea.ListReleaseAttachmentsOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to fetch attachments: %w", err)
	}

	existingAttachments := make(map[string]bool)
	attachmentsMap := make(map[string]*gitea.Attachment)

	for _, attachment := range attachments {
		attachmentsMap[attachment.Name] = attachment
		existingAttachments[attachment.Name] = true
	}

	for _, file := range files {
		fileName := path.Base(file)
		if existingAttachments[fileName] {
			switch FileExists(r.Opt.FileExists) {
			case FileExistsOverwrite:
				_, err := r.client.DeleteReleaseAttachment(r.Opt.Owner, r.Opt.Repo, releaseID, attachmentsMap[fileName].ID)
				if err != nil {
					return fmt.Errorf("failed to delete artifact: %s: %w", fileName, err)
				}

				log.Info().Msgf("deleted artifact: %s", fileName)
			case FileExistsFail:
				return fmt.Errorf("%w: %s", ErrFileExists, fileName)
			case FileExistsSkip:
				log.Warn().Msgf("skip existing artifact: %s", fileName)

				continue
			}
		}

		if err := r.uploadFile(releaseID, file); err != nil {
			return err
		}
	}

	return nil
}

func (r *Release) uploadFile(releaseID int64, file string) error {
	handle, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to read artifact: %s: %w", file, err)
	}
	defer handle.Close()

	_, _, err = r.client.CreateReleaseAttachment(r.Opt.Owner, r.Opt.Repo, releaseID, handle, path.Base(file))
	if err != nil {
		return fmt.Errorf("failed to upload artifact: %s: %w", file, err)
	}

	log.Info().Msgf("uploaded artifact: %s", path.Base(file))

	return nil
}
