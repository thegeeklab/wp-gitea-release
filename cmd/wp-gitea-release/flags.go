package main

import (
	"github.com/thegeeklab/wp-gitea-release/plugin"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
//
//go:generate go run docs.go flags.go
func settingsFlags(settings *plugin.Settings, category string) []cli.Flag {
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
