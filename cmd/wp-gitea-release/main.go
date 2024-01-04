package main

import (
	"fmt"

	"github.com/thegeeklab/wp-gitea-release/plugin"

	wp "github.com/thegeeklab/wp-plugin-go/plugin"
)

//nolint:gochecknoglobals
var (
	BuildVersion = "devel"
	BuildDate    = "00000000"
)

func main() {
	settings := &plugin.Settings{}
	options := wp.Options{
		Name:            "wp-gitea-release",
		Description:     "Add comments to GitHub Issues and Pull Requests",
		Version:         BuildVersion,
		VersionMetadata: fmt.Sprintf("date=%s", BuildDate),
		Flags:           settingsFlags(settings, wp.FlagsPluginCategory),
	}

	plugin.New(options, settings).Run()
}
