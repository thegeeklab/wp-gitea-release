package plugin

import (
	"errors"
	"fmt"
	"os"
	"path"

	"code.gitea.io/sdk/gitea"
	"github.com/rs/zerolog"
)

var (
	ErrReleaseNotFound = errors.New("release not found")
	ErrFileExists      = errors.New("asset file already exist")
)

type GiteaClient struct {
	client  *gitea.Client
	log     zerolog.Logger
	Release *GiteaRelease
}

type GiteaRelease struct {
	client *gitea.Client
	log    zerolog.Logger
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
func NewGiteaClient(client *gitea.Client, logger zerolog.Logger) *GiteaClient {
	return &GiteaClient{
		client: client,
		Release: &GiteaRelease{
			client: client,
			log:    logger,
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
			r.log.Info().Msgf("successfully retrieved %s release", r.Opt.Tag)

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

	r.log.Info().Msgf("successfully created %s release", r.Opt.Tag)

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
		return fmt.Errorf("failed to fetch existing assets: %w", err)
	}

	var uploadFiles []string

files:
	for _, file := range files {
		for _, attachment := range attachments {
			r.log.Debug().Msgf("found existing %s artifact", attachment.Name)

			if attachment.Name == path.Base(file) {
				switch r.Opt.FileExists {
				case "overwrite":
					// do nothing
				case "fail":
					return fmt.Errorf("%w: %s", ErrFileExists, path.Base(file))
				case "skip":
					r.log.Warn().Msgf("skipping pre-existing %s artifact", attachment.Name)

					continue files
				}
			}
		}

		uploadFiles = append(uploadFiles, file)
	}

	for _, file := range uploadFiles {
		handle, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to read %s artifact: %w", file, err)
		}

		for _, attachment := range attachments {
			if attachment.Name == path.Base(file) {
				if _, err := r.client.DeleteReleaseAttachment(r.Opt.Owner, r.Opt.Repo, releaseID, attachment.ID); err != nil {
					return fmt.Errorf("failed to delete %s artifact: %w", file, err)
				}

				r.log.Info().Msgf("successfully deleted old %s artifact", attachment.Name)
			}
		}

		_, _, err = r.client.CreateReleaseAttachment(r.Opt.Owner, r.Opt.Repo, releaseID, handle, path.Base(file))
		if err != nil {
			return fmt.Errorf("failed to upload %s artifact: %w", file, err)
		}

		r.log.Info().Msgf("successfully uploaded %s artifact", file)
	}

	return nil
}
