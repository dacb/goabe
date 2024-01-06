package plugins

import (
	"context"
	"fmt"
	"log/slog"
)

type Hook struct {
	SubStep     int
	Core        func(context.Context) error
	Thread      func(context.Context, int, string) error
	Description string
}

type PluginInit func(context.Context) error
type PluginName func() string
type PluginVersion func() (int, int, int)
type PluginDescription func() string
type PluginGetHooks func() []Hook

type Plugin struct {
	Init        PluginInit
	Name        PluginName
	Version     PluginVersion
	Description PluginDescription
	GetHooks    PluginGetHooks
}

var LoadedPlugins []Plugin

func LoadPlugins(ctx context.Context) error {
	log := ctx.Value("log").(*slog.Logger)

	for _, plugin := range LoadedPlugins {
		plugin.Init(ctx)
		// call the data functions to print some information and verify
		name := plugin.Name()
		description := plugin.Description()
		hooks := plugin.GetHooks()
		major, minor, patch := plugin.Version()
		log.Info(fmt.Sprintf("plugin: %s v%d.%d.%d - %s - %d hooks",
			name, major, minor, patch, description, len(hooks),
		))
	}

	return nil
}
