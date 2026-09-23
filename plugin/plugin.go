package plugin

import (
	"fmt"
	"net/url"
	"slices"

	plugin_base "github.com/thegeeklab/wp-plugin-go/v7/plugin"
	"github.com/urfave/cli/v3"
)

//go:generate go run ../hack/docs-gen/main.go -output=../docs/data/data.yaml

// Plugin implements provide the plugin.
type Plugin struct {
	*plugin_base.Plugin
	Settings *Settings
}

// Settings for the Plugin.
type Settings struct {
	APIKey     string
	FileExists string
	Checksum   []string
	Draft      bool
	PreRelease bool
	Title      string
	Note       string
	CommitRef  string
	Event      string

	baseURL *url.URL
	files   []string
}

func New(e plugin_base.ExecuteFunc, build ...string) *Plugin {
	p := &Plugin{
		Settings: &Settings{},
	}

	options := plugin_base.Options{
		Name:        "wp-gitea-release",
		Description: "Publish files and artifacts to Gitea releases",
		Flags: slices.Concat(
			plugin_base.LoggingFlags(plugin_base.FlagsPluginCategory),
			plugin_base.NetworkFlags(plugin_base.FlagsPluginCategory),
			Flags(p.Settings, plugin_base.FlagsPluginCategory),
		),
		Execute:             p.run,
		HideWoodpeckerFlags: true,
	}

	if len(build) > 0 {
		options.Version = build[0]
	}

	if len(build) > 1 {
		options.VersionMetadata = fmt.Sprintf("date=%s", build[1])
	}

	if e != nil {
		options.Execute = e
	}

	p.Plugin = plugin_base.New(options)

	return p
}

// Flags returns a slice of CLI flags for the plugin.
func Flags(settings *Settings, category string) []cli.Flag {
	return []cli.Flag{
		// Api key to access Gitea API.
		&cli.StringFlag{
			Name:        "api-key",
			Usage:       "api key to access Gitea API",
			Sources:     cli.EnvVars("PLUGIN_API_KEY", "GITEA_RELEASE_API_KEY", "GITEA_TOKEN"),
			Destination: &settings.APIKey,
			Category:    category,
			Required:    true,
		},
		// List of files to upload.
		&cli.StringSliceFlag{
			Name:     "files",
			Usage:    "list of files to upload",
			Sources:  cli.EnvVars("PLUGIN_FILES", "GITEA_RELEASE_FILES"),
			Category: category,
		},
		// What to do if file already exist.
		&cli.StringFlag{
			Name:        "file-exists",
			Value:       "overwrite",
			Usage:       "what to do if file already exist",
			Sources:     cli.EnvVars("PLUGIN_FILE_EXIST", "GITEA_RELEASE_FILE_EXIST"),
			Destination: &settings.FileExists,
			Category:    category,
		},
		// Generate specific checksums.
		&cli.StringSliceFlag{
			Name:        "checksum",
			Usage:       "generate specific checksums",
			Sources:     cli.EnvVars("PLUGIN_CHECKSUM", "GITEA_RELEASE_CHECKSUM"),
			Destination: &settings.Checksum,
			Category:    category,
		},
		// Create a draft release.
		&cli.BoolFlag{
			Name:        "draft",
			Usage:       "create a draft release",
			Sources:     cli.EnvVars("PLUGIN_DRAFT", "GITEA_RELEASE_DRAFT"),
			Destination: &settings.Draft,
			Category:    category,
		},
		// Set the release as prerelease.
		&cli.BoolFlag{
			Name:        "prerelease",
			Usage:       "set the release as prerelease",
			Sources:     cli.EnvVars("PLUGIN_PRERELEASE", "GITEA_RELEASE_PRERELEASE"),
			Destination: &settings.PreRelease,
			Category:    category,
		},
		// URL of the Gitea instance.
		&cli.StringFlag{
			Name:     "base-url",
			Usage:    "URL of the Gitea instance",
			Sources:  cli.EnvVars("PLUGIN_BASE_URL", "GITEA_RELEASE_BASE_URL"),
			Category: category,
			Required: true,
		},
		// File or string with notes for the release.
		&cli.StringFlag{
			Name:        "note",
			Usage:       "file or string with notes for the release",
			Sources:     cli.EnvVars("PLUGIN_NOTE", "GITEA_RELEASE_NOTE"),
			Destination: &settings.Note,
			Category:    category,
		},
		// File or string for the title shown in the Gitea release.
		&cli.StringFlag{
			Name:        "title",
			Usage:       "file or string for the title shown in the Gitea release",
			Sources:     cli.EnvVars("PLUGIN_TITLE", "GITEA_RELEASE_TITLE", "CI_COMMIT_TAG"),
			Destination: &settings.Title,
			DefaultText: "$CI_COMMIT_TAG",
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "event",
			Value:       "push",
			Usage:       "build event",
			Sources:     cli.EnvVars("CI_PIPELINE_EVENT"),
			Destination: &settings.Event,
			DefaultText: "$CI_PIPELINE_EVENT",
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "commit-ref",
			Value:       "refs/heads/main",
			Usage:       "git commit ref",
			Sources:     cli.EnvVars("CI_COMMIT_REF"),
			Destination: &settings.CommitRef,
			Category:    category,
		},
	}
}
