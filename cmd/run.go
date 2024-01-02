package cmd

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/dacb/goabe/logger"

	"github.com/spf13/cobra"
)

var runSteps int64

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Initialize the engine and run it",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

This is the core of the Go Agent Based Engine toolkit.
This application runs agent based models and analyzes
them`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log.With(
			slog.Group("cmd",
				slog.String("cmd", "run"),
				slog.Int64("runSteps", runSteps),
			),
		).Info("run was called")

		runCore(Threads)

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().Int64VarP(&runSteps, "steps", "", 0, "Number of steps to run the engine")
}

func runCore(threads int) {
	wgThreadsDone := new(sync.WaitGroup)
	wgThreadsDone.Add(threads)
	for i := 0; i < threads; i++ {
		go runThread(wgThreadsDone, fmt.Sprintf("thread_%d", i), i+1)
	}
	for step := int64(0); step < runSteps; step++ {

	}
	wgThreadsDone.Wait()
}

func runThread(wgDone *sync.WaitGroup, name string, actions int) {
	defer wgDone.Done()
	logger.Log.Info(fmt.Sprintf("thread %s started", name))
	for i := 0; i < actions; i++ {
		logger.Log.Debug(fmt.Sprintf("thread %s heartbeat %d", name, i+1))
		//time.Sleep(time.Millisecond * time.Duration(500))
	}
}
