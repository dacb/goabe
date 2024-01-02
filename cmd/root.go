package cmd

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goabe.json)")
	rootCmd.PersistentFlags().IntVar(&Threads, "threads", 1, "concurrent threads (default is 1)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig sets default config values and reads in config from input, if possible.
// If the global cfgFile is set by cobra CLI handling (see init()), then this
// function will try to read the config from that file.  If not, it will
// try to find the config file in the home directory (i.e., .goabe.json). This
// may need to change in the future to use a different default or even a URL.
func initConfig() {
	if cfgFile != "" {
		// use config file specified by command line flag
		viper.SetConfigFile(cfgFile)
	} else {
		// identify home dir for user
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// search config in home directory with name ".goabe" w/ .json
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".goabe")
	}

	// setup a default environment that can be overridden
	log_level_text, err := slog.LevelInfo.MarshalText()
	if err != nil {
		panic(err)
	}
	viper.SetDefault("log_level", string(log_level_text))
	viper.SetDefault("log_file", "goabe.log.json")
	viper.SetDefault("substeps", 8)
	viper.SetDefault("random_seed", time.Now().UnixNano())

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
		logger.Log.Info("no configuration file found and/or specified; using defaults")
	}
	logger.Log.Info(fmt.Sprintf("using %d threads", Threads))

	// initialize the random seed
	random_seed := viper.GetInt64("random_seed")
	rand.Seed(random_seed)
	logger.Log.Info(fmt.Sprintf("using %d as the random seed", random_seed))
}
