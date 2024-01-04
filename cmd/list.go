package cmd

import (
	"context"
	"fmt"

	"github.com/dacb/goabe/logger"
	"github.com/dacb/goabe/plugins"
	"github.com/spf13/cobra"
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
		// find the plugin directory list from the config
		//pluginDirStr := viper.GetString("plugin_dir")

		// iterate over the directories finding each .so file
		// open each file and try to call the basic functions
		// incuding Init, Name, Version, Description
		fmt.Println("plugin called")
		pluginFilename := "plugins/example/example.so"
		plg, err := plugins.LoadPlugIn(pluginFilename)
		if err != nil {
			panic(err)
		}
		ctx := context.Background()
		ctx = context.WithValue(ctx, "log", logger.Log.With("plugin", pluginFilename))
		(*plg).Init(ctx)
		(*plg).Name()
		(*plg).Description()
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
