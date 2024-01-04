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
		log := logger.Log.With("cmd", "plugin list")
		log.Info("plugin list called")
		ctx := context.WithValue(cmd.Context(), "log", log)

		err := plugins.LoadPlugins(ctx)
		if err != nil {
			log.Error("an error occurred loading the plugins")
			panic(err)
		}

		for _, plugin := range plugins.LoadedPlugins {
			name := (*plugin).Name()
			description := (*plugin).Description()
			hooks := (*plugin).GetHooks()
			major, minor, patch := (*plugin).Version()
			filename := (*plugin).Filename()
			log.Info(fmt.Sprintf("plugin: %s v%d.%d.%d - %s - %s - %d hooks",
				name, major, minor, patch, description, filename, len(hooks),
			))
			for _, hook := range hooks {
				log.Info(fmt.Sprintf("plugin: %s - hook '%s' at step %d substep %d (%0x, %0x)", name, hook.Description, hook.Step, hook.SubStep, hook.Core, hook.Thread))
			}
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
