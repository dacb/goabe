package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dacb/goabe/plugins"
)

var pluginFilename string // when this is populated, the plugin has alread been initialized
var log *slog.Logger

// the empty struct used for the plugin
type plugin struct {
}

var PlugIn plugin

// main initiailization function for the plugin
func (p *plugin) Init(ctx context.Context, pluginFname string) error {
	mylog, ok := ctx.Value("log").(*slog.Logger)
	if !ok {
		return errors.New("no logger found on the current context")
	}
	log = mylog
	log.Info("example plugin Init function was called")

	if pluginFilename != "" {
		log.Error("plugin has already been initialized? refusing to load the plugin twice")
		return errors.New("this plug in has already been loaded!")
	}
	pluginFilename = pluginFname

	return nil
}

// major, minor, patch
func (p *plugin) Version() (int, int, int) {
	log.Debug("example plugin Version function was called")
	return 1, 0, 0
}

// returns the short name of the module as a string
func (p *plugin) Name() string {
	log.Debug("example plugin Name function was called")
	return "example"
}

// returns a short description of the module as a string
func (p *plugin) Description() string {
	log.Debug("example plugin Description function was called")
	return "example plugin for code template"
}

func (p *plugin) GetHooks() []plugins.Hook {
	log.Debug("example plugin GetHooks function was called")

	var hooks []plugins.Hook
	hooks = append(hooks, plugins.Hook{0, nil, ThreadSubStep0, "thread sum"})
	hooks = append(hooks, plugins.Hook{1, CoreSubStep1, nil, "core sum"})

	return hooks
}

func (p *plugin) Filename() string {
	log.Debug("example plugin Filename function was called")
	return pluginFilename
}

// note this logs through the context
func CoreSubStep1(ctx context.Context) error {
	log := ctx.Value("log").(*slog.Logger).With("plugin", pluginFilename)
	log.Debug("core substep 1 hook called")
	return nil
}

// note this logs through the context
func ThreadSubStep0(ctx context.Context, id int, name string) error {
	log := ctx.Value("log").(*slog.Logger).With("actor", name).With("plugin", pluginFilename)
	log.Debug("thread substep 0 hook called")
	return nil
}
