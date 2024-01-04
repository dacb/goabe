package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dacb/goabe/plugins"
)

var log *slog.Logger

// the empty struct used for the plugin
type plugin struct {
}

// returns the version of plugin in major, minor, and patch
func (p *plugin) Init(ctx context.Context) error {
	mylog, ok := ctx.Value("log").(*slog.Logger)
	if !ok {
		panic(fmt.Errorf("unable to find the logger value in the current context"))
	}
	log = mylog
	log.Info("example plugin Init function was called")

	return nil
}

// levels as separate integers
func (p *plugin) Version() (int, int, int) {
	log.Info("example plugin Version function was called")
	return 0, 1, 0
}

// returns the short name of the module as a string
func (p *plugin) Name() string {
	log.Info("example plugin Name function was called")
	return "example"
}

// returns a short description of the module as a string
func (p *plugin) Description() string {
	log.Info("example plugin Description function was called")
	return "example plugin for code template"
}

func (p *plugin) GetHooks() []plugins.Hook {
	log.Info("example plugin GetHooks function was called")

	var hooks []plugins.Hook
	hooks = append(hooks, plugins.Hook{0, 0, CoreSubStep0, nil})
	hooks = append(hooks, plugins.Hook{0, 1, nil, ThreadSubStep1})

	return hooks
}

func CoreSubStep0() error {
	log.With("actor", "core").Info("core substep 0 hook called")
	return nil
}

func ThreadSubStep1(id int, name string) error {
	log.With("actor", name).Info("core substep 0 hook called")
	return nil
}

var PlugIn plugin
