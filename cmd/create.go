package cmd

import (
	"log/slog"

	"github.com/dacb/goabe/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a configuration file with default options.",
	Long:  `This saves a default configuration to the file for future modificaiton and usage`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log.With(
			slog.Group("cmd",
				slog.String("cmd", "config"),
				slog.String("sub", "create"),
				slog.String("config_file", viper.ConfigFileUsed()),
			),
		).Info("create called to save out the configuration; will not overwrite an existing file")
		err := viper.SafeWriteConfig()
		if err != nil {
			logger.Log.With(
				slog.Group("cmd",
					slog.String("cmd", "config"),
					slog.String("sub", "create"),
					slog.String("config_file", viper.ConfigFileUsed()),
				),
			).Error("unable to safe write configuration file, does it already exist?")
		}
	},
}

func init() {
	configCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
