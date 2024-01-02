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

//go:generate stringer -type=engineMsg
type engineMsg int

const (
	HALT engineMsg = iota
	CONTINUE
)

func runCore(threads int) {
	// this waitgroup is used to signal the close of the threads
	wgThreadsDone := new(sync.WaitGroup)
	wgThreadsDone.Add(threads)
	// channels
	syncChan := make(chan engineMsg)

	// spawn the threads
	for i := 0; i < threads; i++ {
		go runThread(wgThreadsDone, syncChan, fmt.Sprintf("thread_%d", i), i)
	}
	// iterate over steps
	for step := int64(0); step < runSteps; step++ {
		for threadI := 0; threadI < threads; threadI++ {
			// release the threads
			logger.Log.Debug(fmt.Sprintf("releasing thread %d on step %d", threadI, step))
			syncChan <- CONTINUE
		}
	}
	for threadI := 0; threadI < threads; threadI++ {
		syncChan <- HALT
	}
	// wait until the threads are done
	wgThreadsDone.Wait()
}

func runThread(wgDone *sync.WaitGroup, syncChan chan engineMsg, name string, id int) {
	defer wgDone.Done()
	logger.Log.Debug(fmt.Sprintf("thread %s started", name))
	state := <-syncChan
	for state != HALT {
		logger.Log.Debug(fmt.Sprintf("thread %s heartbeat %d", name, id))
		state = <-syncChan
	}
}
