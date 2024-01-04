package plugin

import (
	"net/url"

	wp "github.com/thegeeklab/wp-plugin-go/plugin"
	"github.com/urfave/cli/v2"
)

// Plugin implements provide the plugin.
type Plugin struct {
	*wp.Plugin
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

func New(options wp.Options, settings *Settings) *Plugin {
	p := &Plugin{}

	options.Execute = p.run

	p.Plugin = wp.New(options)
	p.Settings = settings

	return p
}
