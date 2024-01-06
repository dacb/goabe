package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	"github.com/dacb/goabe/example"
	"github.com/dacb/goabe/life"
	"github.com/dacb/goabe/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// the name of the config file on the file system from the user (from cobra)
var cfgFile string

// the number of concurrent processes (threads) to try to use (from cobra)
var Threads int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goabe",
	Short: "Go Agent Based Engine (goabe)",
	Long: `A scalable, parallel engine for simulation
and analysis of agent based computations.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.SetContext(context.WithValue(cmd.Context(), "threads", Threads))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./goabe.json)")
	rootCmd.PersistentFlags().IntVar(&Threads, "threads", 1, "concurrent threads (default is 1)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// register the plugins
	// this should be done with a go generate
	example.Register()
	life.Register()
}

// initConfig sets default config values and reads in config from input, if possible.
// If the global cfgFile is set by cobra CLI handling (see init()), then this
// function will try to read the config from that file.  If not, it will
// try to find the config file in the current directory (i.e., goabe.json). This
// may need to change in the future to use a different default or even a URL.
func initConfig() {
	if cfgFile != "" {
		// use config file specified by command line flag
		viper.SetConfigFile(cfgFile)
	} else {
		// search for config file in current directory w/ name goabe.json
		viper.AddConfigPath(".")
		viper.SetConfigType("json")
		viper.SetConfigName("goabe")
	}

	// read in environment variables
	viper.AutomaticEnv()

	// if a config file is found, read it in
	configFromFile := false
	if err := viper.ReadInConfig(); err == nil {
		configFromFile = true
	}

	// set up the logger, this has to be done after the config is read in because
	// it contains the name of the log output
	logger.InitLogger()
	if configFromFile {
		logger.Log.With("config_file", viper.ConfigFileUsed()).Info("loaded config from file")
	} else {
		if cfgFile != "" {
			logger.Log.With("config_file", viper.ConfigFileUsed()).Error("unable to read configuration file")
			panic(fmt.Errorf("unable to read configuration from %s", viper.ConfigFileUsed()))
		}
		logger.Log.Info("no configuration file found and/or specified; using defaults")
	}
	logger.Log.Info(fmt.Sprintf("using %d threads", Threads))

	// initialize the random seed
	random_seed := viper.GetInt64("random_seed")
	rand.Seed(random_seed)
	logger.Log.Info(fmt.Sprintf("using %d as the random seed", random_seed))
}
