package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"time"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management tools",
	Long: `These tools include creating a default configuration file, validating
a configuration file, and other aspects of working with configuraiton data.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("config called")
	//},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")	// setup a default environment that can be overridden

	log_level_text, err := slog.LevelInfo.MarshalText()
	if err != nil {
		panic(err)
	}
	viper.SetDefault("log_level", string(log_level_text))
	viper.SetDefault("log_file", "goabe.log.json")
	viper.SetDefault("substeps", 10)
	viper.SetDefault("random_seed", time.Now().UnixNano())
}
