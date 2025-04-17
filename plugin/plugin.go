package plugin

import (
	"fmt"
	"net/url"

	plugin_base "github.com/thegeeklab/wp-plugin-go/v6/plugin"
	"github.com/urfave/cli/v2"
)

//go:generate go run ../internal/doc/main.go -output=../docs/data/data-raw.yaml

// Plugin implements provide the plugin.
type Plugin struct {
	*plugin_base.Plugin
	Settings *Settings
}

// Settings for the Plugin.
type Settings struct {
	APIKey     string
	FileExists string
	Checksum   cli.StringSlice
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
		Name:                "wp-gitea-release",
		Description:         "Publish files and artifacts to Gitea releases",
		Flags:               Flags(p.Settings, plugin_base.FlagsPluginCategory),
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
		&cli.StringFlag{
			Name:        "api-key",
			Usage:       "api key to access Gitea API",
			EnvVars:     []string{"PLUGIN_API_KEY", "GITEA_RELEASE_API_KEY", "GITEA_TOKEN"},
			Destination: &settings.APIKey,
			Category:    category,
			Required:    true,
		},
		&cli.StringSliceFlag{
			Name:     "files",
			Usage:    "list of files to upload",
			EnvVars:  []string{"PLUGIN_FILES", "GITEA_RELEASE_FILES"},
			Category: category,
		},
		&cli.StringFlag{
			Name:        "file-exists",
			Value:       "overwrite",
			Usage:       "what to do if file already exist",
			EnvVars:     []string{"PLUGIN_FILE_EXIST", "GITEA_RELEASE_FILE_EXIST"},
			Destination: &settings.FileExists,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "checksum",
			Usage:       "generate specific checksums",
			EnvVars:     []string{"PLUGIN_CHECKSUM", "GITEA_RELEASE_CHECKSUM"},
			Destination: &settings.Checksum,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "draft",
			Usage:       "create a draft release",
			EnvVars:     []string{"PLUGIN_DRAFT", "GITEA_RELEASE_DRAFT"},
			Destination: &settings.Draft,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "prerelease",
			Usage:       "set the release as prerelease",
			EnvVars:     []string{"PLUGIN_PRERELEASE", "GITEA_RELEASE_PRERELEASE"},
			Destination: &settings.PreRelease,
			Category:    category,
		},
		&cli.StringFlag{
			Name:     "base-url",
			Usage:    "URL of the Gitea instance",
			EnvVars:  []string{"PLUGIN_BASE_URL", "GITEA_RELEASE_BASE_URL"},
			Category: category,
			Required: true,
		},
		&cli.StringFlag{
			Name:        "note",
			Usage:       "file or string with notes for the release",
			EnvVars:     []string{"PLUGIN_NOTE", "GITEA_RELEASE_NOTE"},
			Destination: &settings.Note,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "title",
			Usage:       "file or string for the title shown in the Gitea release",
			EnvVars:     []string{"PLUGIN_TITLE", "GITEA_RELEASE_TITLE", "CI_COMMIT_TAG"},
			Destination: &settings.Title,
			DefaultText: "$CI_COMMIT_TAG",
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "event",
			Value:       "push",
			Usage:       "build event",
			EnvVars:     []string{"CI_PIPELINE_EVENT"},
			Destination: &settings.Event,
			DefaultText: "$CI_PIPELINE_EVENT",
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "commit-ref",
			Value:       "refs/heads/main",
			Usage:       "git commit ref",
			EnvVars:     []string{"CI_COMMIT_REF"},
			Destination: &settings.CommitRef,
			Category:    category,
		},
	}
}
