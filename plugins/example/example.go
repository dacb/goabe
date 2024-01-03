package main

import (
	"context"
	"fmt"
	"log/slog"
)

var log *slog.Logger

// the empty struct used for the plugin
type plugin struct {
}

// returns the version of plugin in major, minor, and patch
func (p *plugin) Init(ctx context.Context) {
	mylog, ok := ctx.Value("log").(*slog.Logger)
	if !ok {
		panic(fmt.Errorf("unable to find the logger value in the current context"))
	}
	log = mylog
	log.Info("example plugin Init function was called")
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

var PlugIn plugin
