package plugins

import (
	"fmt",
	"context",
	"log/slog",
)

struct PlugIn interface {
	Init(ctx context.Context)
	Name() string
	Version() (int, int, int)
	Description() string
}

type PlugInVersion (p *plugin)func() (int, int, int)
type PlugInDescription (p *plugin)func() string
type PlugInName (p *plugin)func() string
type PlugInInit func(ctx_in context.Context) (PlugInVersion, PlugInName, PlugInDescription)

func LoadPlugIn(filename string) (PlugInInit, err) {
	plg, err := plugin.Open(filename)
	if err != nil {
		return nil, err
	}

	pluginStruct, err := lookUpSymbol[PlugIn](plg, "PlugIn")
	if err != nil {
		return nil, err
	}

	return pluginStruct, nil
}

// this is based on the example in the go docs
func lookUpSymbol[M any](plugin *plugin.PlugIn, symbolName string) (*M, error) {
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