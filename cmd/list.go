package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dacb/goabe/logger"
	"github.com/dacb/goabe/plugins"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log.With("cmd", "plugin").Info("plugin list called")
		var pluginFiles []string
		// find the plugin directory list from the config
		// these should be separated by ':'
		// iterate over the directories finding each .so file
		pluginDirsList := viper.GetString("plugin_dirs")
		pluginDirs := strings.Split(pluginDirsList, ":")
		for _, pluginDir := range pluginDirs {
			if stat, err := os.Stat(pluginDir); err != nil || !stat.IsDir() {
				logger.Log.Error(fmt.Sprintf("unable to open the plugin directory '%s'", pluginDir))
				panic(err)
			}
			filepath.WalkDir(pluginDir, func(str string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if filepath.Ext(d.Name()) == ".so" {
					pluginFiles = append(pluginFiles, str)
				}
				return nil
			})
			logger.Log.With("cmd", "plugin").Info(fmt.Sprintf("plugin directory %s contained %d plugins", pluginDir, len(pluginFiles)))
		}

		// open each file and try to call the basic functions
		// incuding Init, Name, Version, Description
		for _, pluginFilename := range pluginFiles {
			logger.Log.With("cmd", "plugin").Info(fmt.Sprintf("loading plugin from file %s", pluginFilename))
			ctx := context.WithValue(cmd.Context(), "log", logger.Log.With("plugin", pluginFilename))
			plg, err := plugins.LoadPlugIn(ctx, pluginFilename)
			if err != nil {
				panic(err)
			}
			logger.Log.With("cmd", "plugin").Debug(fmt.Sprintf("loaded plugin %s from file %s", (*plg).Name(), pluginFilename))
		}
	},
}

func init() {
	pluginCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
