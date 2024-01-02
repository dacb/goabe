package logger

import (
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/viper"
)

// This is the system wide log interface.
var Log *slog.Logger

// InitLogger reads configuration to initialize logging for goabe.
// The function needs the viper configuration system to be available
// which also needs the cobra command line parser to be availabel.
// These two dependencies means this function exists instead of doing
// this with init().
// When this function is run, the configuration will already have been
// read from whatever file specified by the command line arguments.
// This tooling runs dual logging streams: one on the stdout for text
// based logging and the other on the json file.
// This function will clobber (i.e., overwrite and truncate) any
// existing contents of the log file.
func InitLogger() {
	// initialize the system using the config data from viper
	logfile, err := os.OpenFile(viper.GetString("log_file"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	// fanout over the stdout in text and output log in json
	log_level_text := []byte(viper.GetString("log_level"))
	var log_level slog.Level
	log_level.UnmarshalText(log_level_text)
	opts := &slog.HandlerOptions{
		Level: log_level,
	}
	Log = slog.New(
		slogmulti.Fanout(
			slog.NewJSONHandler(logfile, opts),
			slog.NewTextHandler(os.Stdout, opts),
		),
	)
	Log.Info("logging initialized")
}
