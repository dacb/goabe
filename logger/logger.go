package logger

import (
	"fmt"
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/viper"
)

var Log *slog.Logger

func InitLogger() {
	fmt.Println(viper.GetString("log_file"), "log_file")
	// initialize the system using the config data from viper
	logfile, err := os.OpenFile(viper.GetString("log_file"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
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
