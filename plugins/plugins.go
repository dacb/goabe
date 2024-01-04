package plugins

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/dacb/goabe/logger"
	"github.com/spf13/viper"
	//"github.com/dacb/goabe/logger"
)

type Hook struct {
	SubStep     int
	Core        func(context.Context) error
	Thread      func(context.Context, int, string) error
	Description string
}

type PlugIn interface {
	Init(context.Context, string) error
	Name() string
	Version() (int, int, int)
	Description() string
	Filename() string
	GetHooks() []Hook
}

var LoadedPlugins []*PlugIn

func openPlugIn(_ context.Context, pluginFilename string) (*PlugIn, error) {
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

func LoadPlugins(ctx context.Context) error {
	log := ctx.Value("log").(*slog.Logger)

	// list of found plug in files, note these are just all files ending in .so
	var pluginFiles []string

	// find the plugin directory list from the config
	// these should be separated by ':'
	// iterate over the directories finding each .so file
	pluginDirsList := viper.GetString("plugin_dirs")
	pluginDirs := strings.Split(pluginDirsList, ":")
	for _, pluginDir := range pluginDirs {
		// if the dir is not real or is not a directory, this is not very graceful
		if stat, err := os.Stat(pluginDir); err != nil || !stat.IsDir() {
			log.Error(fmt.Sprintf("unable to open the plugin directory '%s'", pluginDir))
			return err
		}
		// go through the directory and find any files ending in .so
		filepath.WalkDir(pluginDir, func(str string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(d.Name()) == ".so" {
				pluginFiles = append(pluginFiles, str)
			}
			return nil
		})
		log.Info(fmt.Sprintf("plugin directory %s contained %d plugins", pluginDir, len(pluginFiles)))
	}

	// open each file and try to call the basic functions
	// incuding Init, Name, Version, Description
	for _, pluginFilename := range pluginFiles {
		log.Debug(fmt.Sprintf("loading plugin from file %s", pluginFilename))
		ctx := context.WithValue(ctx, "log", logger.Log.With("plugin", pluginFilename))
		plg, err := loadPlugIn(ctx, pluginFilename)
		if err != nil {
			log.Error(fmt.Sprintf("unable to load plugin in file %s", pluginFilename))
			return err
		}
		log.Debug(fmt.Sprintf("loaded plugin %s from file %s", (*plg).Name(), pluginFilename))
	}
	return nil
}

// given a context and a filename, load the plugin and initialize it
func loadPlugIn(ctx context.Context, pluginFilename string) (*PlugIn, error) {
	log := ctx.Value("log").(*slog.Logger)
	ctx = context.WithValue(ctx, "log", logger.Log.With("plugin", pluginFilename))

	plugin, err := openPlugIn(ctx, pluginFilename)
	if err != nil {
		return nil, err
	}
	LoadedPlugins = append(LoadedPlugins, plugin)

	// call the init function for the plug in
	err = (*plugin).Init(ctx, pluginFilename)
	if err != nil {
		return nil, err
	}

	// call the data functions to print some information and verify
	name := (*plugin).Name()
	description := (*plugin).Description()
	hooks := (*plugin).GetHooks()
	major, minor, patch := (*plugin).Version()
	filename := (*plugin).Filename()

	log.Info(fmt.Sprintf("plugin: %s v%d.%d.%d - %s - %s - %d hooks",
		name, major, minor, patch, description, filename, len(hooks),
	))

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
