package gitea

import (
	"io"

	"code.gitea.io/sdk/gitea"
)

//nolint:lll
type APIClient interface {
	ListReleases(owner, repo string, opt gitea.ListReleasesOptions) ([]*gitea.Release, *gitea.Response, error)
	CreateRelease(owner, repo string, opt gitea.CreateReleaseOption) (*gitea.Release, *gitea.Response, error)
	ListReleaseAttachments(user, repo string, release int64, opt gitea.ListReleaseAttachmentsOptions) ([]*gitea.Attachment, *gitea.Response, error)
	CreateReleaseAttachment(user, repo string, release int64, file io.Reader, filename string) (*gitea.Attachment, *gitea.Response, error)
	DeleteReleaseAttachment(user, repo string, release, id int64) (*gitea.Response, error)
}
