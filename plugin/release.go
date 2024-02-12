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

// Release holds ties the Woodpecker env data and gitea client together.
type releaseClient struct {
	*gitea.Client
	Owner      string
	Repo       string
	Tag        string
	Draft      bool
	Prerelease bool
	FileExists string
	Title      string
	Note       string
}

func (rc *releaseClient) buildRelease() (*gitea.Release, error) {
	// first attempt to get a release by that tag
	release, err := rc.getRelease()

	if err != nil && release == nil {
		fmt.Println(err)
	} else if release != nil {
		return release, nil
	}

	// if no release was found by that tag, create a new one
	release, err = rc.newRelease()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve or create a release: %w", err)
	}

	return release, nil
}

func (rc *releaseClient) getRelease() (*gitea.Release, error) {
	releases, _, err := rc.Client.ListReleases(rc.Owner, rc.Repo, gitea.ListReleasesOptions{})
	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		if release.TagName == rc.Tag {
			log.Info().Msgf("successfully retrieved %s release", rc.Tag)

			return release, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrReleaseNotFound, rc.Tag)
}

func (rc *releaseClient) newRelease() (*gitea.Release, error) {
	r := gitea.CreateReleaseOption{
		TagName:      rc.Tag,
		IsDraft:      rc.Draft,
		IsPrerelease: rc.Prerelease,
		Title:        rc.Title,
		Note:         rc.Note,
	}

	release, _, err := rc.Client.CreateRelease(rc.Owner, rc.Repo, r)
	if err != nil {
		return nil, fmt.Errorf("failed to create release: %w", err)
	}

	log.Info().Msgf("successfully created %s release", rc.Tag)

	return release, nil
}

func (rc *releaseClient) uploadFiles(releaseID int64, files []string) error {
	attachments, _, err := rc.Client.ListReleaseAttachments(
		rc.Owner,
		rc.Repo,
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
			if attachment.Name == path.Base(file) {
				switch rc.FileExists {
				case "overwrite":
					// do nothing
				case "fail":
					return fmt.Errorf("%w: %s", ErrFileExists, path.Base(file))
				case "skip":
					log.Warn().Msgf("skipping pre-existing %s artifact", attachment.Name)

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
				if _, err := rc.Client.DeleteReleaseAttachment(rc.Owner, rc.Repo, releaseID, attachment.ID); err != nil {
					return fmt.Errorf("failed to delete %s artifact: %w", file, err)
				}

				log.Info().Msgf("successfully deleted old %s artifact", attachment.Name)
			}
		}

		if _, _, err = rc.Client.CreateReleaseAttachment(rc.Owner, rc.Repo, releaseID, handle, path.Base(file)); err != nil {
			return fmt.Errorf("failed to upload %s artifact: %w", file, err)
		}

		log.Info().Msgf("successfully uploaded %s artifact", file)
	}

	return nil
}
