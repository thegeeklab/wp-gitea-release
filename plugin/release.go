package plugin

import (
	"errors"
	"fmt"
	"os"
	"path"

	"code.gitea.io/sdk/gitea"
	"github.com/rs/zerolog/log"
)

var (
	ErrReleaseNotFound = errors.New("release not found")
	ErrFileExists      = errors.New("asset file already exist")
)

type GiteaClient struct {
	client  *gitea.Client
	Release *GiteaRelease
}

type GiteaRelease struct {
	client *gitea.Client
	Opt    GiteaReleaseOpt
}

type GiteaReleaseOpt struct {
	Owner      string
	Repo       string
	Tag        string
	Draft      bool
	Prerelease bool
	FileExists string
	Title      string
	Note       string
}

// NewGiteaClient creates a new GiteaClient instance with the provided Gitea client.
func NewGiteaClient(client *gitea.Client) *GiteaClient {
	return &GiteaClient{
		client: client,
		Release: &GiteaRelease{
			client: client,
			Opt:    GiteaReleaseOpt{},
		},
	}
}

// Find retrieves the release with the specified tag name from the repository.
// If the release is not found, it returns an ErrReleaseNotFound error.
func (r *GiteaRelease) Find() (*gitea.Release, error) {
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
func (r *GiteaRelease) Create() (*gitea.Release, error) {
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
func (r *GiteaRelease) AddAttachments(releaseID int64, files []string) error {
	attachments, _, err := r.client.ListReleaseAttachments(
		r.Opt.Owner,
		r.Opt.Repo,
		releaseID,
		gitea.ListReleaseAttachmentsOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to fetch attachments: %w", err)
	}

	var uploadFiles []string

files:
	for _, file := range files {
		for _, attachment := range attachments {
			if attachment.Name == path.Base(file) {
				switch r.Opt.FileExists {
				case "overwrite":
					// do nothing
				case "fail":
					return fmt.Errorf("%w: %s", ErrFileExists, path.Base(file))
				case "skip":
					log.Warn().Msgf("skip existing artifact: %s", path.Base(file))

					continue files
				}
			}
		}

		uploadFiles = append(uploadFiles, file)
	}

	for _, file := range uploadFiles {
		handle, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to read artifact: %s: %w", file, err)
		}

		for _, attachment := range attachments {
			if attachment.Name == path.Base(file) {
				if _, err := r.client.DeleteReleaseAttachment(r.Opt.Owner, r.Opt.Repo, releaseID, attachment.ID); err != nil {
					return fmt.Errorf("failed to delete artifact: %s: %w", file, err)
				}

				log.Info().Msgf("deleted artifact: %s", attachment.Name)
			}
		}

		_, _, err = r.client.CreateReleaseAttachment(r.Opt.Owner, r.Opt.Repo, releaseID, handle, path.Base(file))
		if err != nil {
			return fmt.Errorf("failed to upload artifact: %s: %w", file, err)
		}

		log.Info().Msgf("uploaded artifact: %s", path.Base(file))
	}

	return nil
}
