package example

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dacb/goabe/plugins"
)

var log *slog.Logger

var threads int

func Register() {
	plugins.LoadedPlugins = append(plugins.LoadedPlugins, plugins.Plugin{Init, Name, Version, Description, GetHooks})
}

// main initiailization function for the plugin
func Init(ctx context.Context) error {
	mylog, ok := ctx.Value("log").(*slog.Logger)
	if !ok {
		return errors.New("no logger found on the current context")
	}
	log = mylog.With("plugin", Name())
	log.Info("example plugin Init function was called")

	threadCount, ok := ctx.Value("threads").(int)
	if !ok {
		return errors.New("missing number of threads in current context")
	}
	threads = threadCount

	return nil
}

// major, minor, patch
func Version() (int, int, int) {
	return 1, 0, 0
}

// returns the short name of the module as a string
func Name() string {
	return "example"
}

// returns a short description of the module as a string
func Description() string {
	return "example plugin for code template"
}

func GetHooks() []plugins.Hook {
	log.Debug("example plugin GetHooks function was called")

	var hooks []plugins.Hook
	hooks = append(hooks, plugins.Hook{0, nil, ThreadSubStep0, "thread sum"})
	hooks = append(hooks, plugins.Hook{1, CoreSubStep1, nil, "core sum"})

	return hooks
}

// note this logs through the context
func CoreSubStep1(ctx context.Context) error {
	log := ctx.Value("log").(*slog.Logger).With("plugin", Name())
	log.Debug("core substep 1 hook called")
	return nil
}

// note this logs through the context
func ThreadSubStep0(ctx context.Context, id int, name string) error {
	log := ctx.Value("log").(*slog.Logger).With("actor", name).With("plugin", Name())
	log.Debug("thread substep 0 hook called")
	return nil
}
