package cmd

import (
	"fmt",

	"github.com/dacb/goabe/cmd/plugins",

	"github.com/spf13/cobra",
)

// pluginCmd represents the plugin command
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "A set of tools for interacting with goabe plugins",
	Long: `Some tools, including 'list' to work with goabe plugins.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("plugin called")
		plg, err := plugins.LoadPlugIn("plugins/example.so")
		if err != nil {
			panic(err)
		}
		(*plg).Init()
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pluginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pluginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
