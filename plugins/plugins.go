package plugins

import (
	"context"
	"errors"
	"fmt"
	"plugin"
)

type PlugIn interface {
	Init(ctx context.Context)
	Name() string
	Version() (int, int, int)
	Description() string
}

func LoadPlugIn(filename string) (*PlugIn, error) {
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
