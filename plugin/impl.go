package plugin

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/thegeeklab/wp-gitea-release/gitea"
	plugin_file "github.com/thegeeklab/wp-plugin-go/v4/file"
)

var (
	ErrPluginEventNotSupported = errors.New("event not supported")
	ErrFileExistInvalid        = errors.New("invalid file_exist value")
)

//nolint:revive
func (p *Plugin) run(ctx context.Context) error {
	if err := p.FlagsFromContext(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Execute(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	var err error

	fileExistsValues := map[string]bool{
		"overwrite": true,
		"fail":      true,
		"skip":      true,
	}

	if p.Settings.Event != "tag" {
		return fmt.Errorf("%w: %s", ErrPluginEventNotSupported, p.Metadata.Pipeline.Event)
	}

	if !fileExistsValues[p.Settings.FileExists] {
		return ErrFileExistInvalid
	}

	if p.Settings.Note != "" {
		if p.Settings.Note, _, err = plugin_file.ReadStringOrFile(p.Settings.Note); err != nil {
			return fmt.Errorf("error while reading %s: %w", p.Settings.Note, err)
		}
	}

	if p.Settings.Title != "" {
		if p.Settings.Title, _, err = plugin_file.ReadStringOrFile(p.Settings.Title); err != nil {
			return fmt.Errorf("error while reading %s: %w", p.Settings.Title, err)
		}
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	client, err := gitea.NewClient(p.Settings.baseURL.String(), p.Settings.APIKey, p.Network.Client)
	if err != nil {
		return fmt.Errorf("failed to create Gitea client: %w", err)
	}

	client.Release.Opt = gitea.ReleaseOptions{
		Owner:      p.Metadata.Repository.Owner,
		Repo:       p.Metadata.Repository.Name,
		Tag:        strings.TrimPrefix(p.Settings.CommitRef, "refs/tags/"),
		Draft:      p.Settings.Draft,
		Prerelease: p.Settings.PreRelease,
		FileExists: p.Settings.FileExists,
		Title:      p.Settings.Title,
		Note:       p.Settings.Note,
	}

	release, err := client.Release.Find()
	if err != nil && !errors.Is(err, gitea.ErrReleaseNotFound) {
		return fmt.Errorf("failed to retrieve release: %w", err)
	}

	// If no release was found by that tag, create a new one.
	if release == nil {
		release, err = client.Release.Create()
		if err != nil {
			return fmt.Errorf("failed to create release: %w", err)
		}
	}

	if err := client.Release.AddAttachments(release.ID, p.Settings.files); err != nil {
		return fmt.Errorf("failed to upload the files: %w", err)
	}

	return nil
}

func (p *Plugin) FlagsFromContext() error {
	var err error

	baseURL := p.Context.String("base-url")

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	p.Settings.baseURL, err = url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base url: %w", err)
	}

	var files []string

	rawFiles := p.Context.StringSlice("files")
	for _, glob := range rawFiles {
		globed, err := filepath.Glob(glob)
		if err != nil {
			return fmt.Errorf("failed to glob %s: %w", glob, err)
		}

		if globed != nil {
			files = append(files, globed...)
		}
	}

	if len(p.Settings.Checksum.Value()) > 0 {
		var err error

		files, err = WriteChecksums(files, p.Settings.Checksum.Value(), "")
		if err != nil {
			return fmt.Errorf("failed to write checksums: %w", err)
		}
	}

	p.Settings.files = files

	return nil
}
