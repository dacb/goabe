package plugins

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"plugin"

	"github.com/dacb/goabe/logger"
	//"github.com/dacb/goabe/logger"
)

type Hook struct {
	Step        int
	SubStep     int
	Core        func() error
	Thread      func(id int, name string) error
	Description string
}

type PlugIn interface {
	Init(context.Context) error
	Name() string
	Version() (int, int, int)
	Description() string
	GetHooks() []Hook
}

func OpenPlugIn(_ context.Context, pluginFilename string) (*PlugIn, error) {
	plg, err := plugin.Open(pluginFilename)
	if err != nil {
		return nil, err
	}

	pluginStruct, err := lookUpSymbol[PlugIn](plg, "PlugIn")
	if err != nil {
		return nil, err
	}

	return pluginStruct, nil
}

func LoadPlugIn(ctx context.Context, pluginFilename string) (*PlugIn, error) {
	log := ctx.Value("log").(*slog.Logger)
	ctx = context.WithValue(ctx, "log", logger.Log.With("plugin", pluginFilename))

	plugin, err := OpenPlugIn(ctx, pluginFilename)
	if err != nil {
		return nil, err
	}

	// call the init function for the plug in
	err = (*plugin).Init(ctx)
	if err != nil {
		return nil, err
	}

	// call the data functions to print some information and verify
	name := (*plugin).Name()
	description := (*plugin).Description()
	callbacks := (*plugin).GetHooks()

	log.With("cmd", "plugin").With("description", description).
		Info(fmt.Sprintf("plugin %s had %d call backs", name, len(callbacks)))
	for _, callback := range callbacks {
		log.With("cmd", "plugin").Info(fmt.Sprintf("callback '%s' at step %d substep %d (%0x, %0x)", callback.Description, callback.Step, callback.SubStep, callback.Core, callback.Thread))
	}

	return plugin, nil
}

// this is based on the example in the go docs
func lookUpSymbol[M any](plugin *plugin.Plugin, symbolName string) (*M, error) {
	symbol, err := plugin.Lookup(symbolName)
	if err != nil {
		return nil, err
	}
	switch symbol.(type) {
	case *M:
		return symbol.(*M), nil
	case M:
		result := symbol.(M)
		return &result, nil
	default:
		return nil, errors.New(fmt.Sprintf("unhandled type from module symbol: %T", symbol))
	}
}
