package cmd

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/dacb/goabe/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	subSteps := viper.GetInt("substeps")
	// this waitgroup is used to signal the close of the threads
	wgThreadsDone := new(sync.WaitGroup)
	wgThreadsDone.Add(threads)
	// channels
	syncChan := make([]chan engineMsg, threads)

	// spawn the threads
	for threadI := 0; threadI < threads; threadI++ {
		syncChan[threadI] = make(chan engineMsg)
		go runThread(wgThreadsDone, syncChan[threadI], fmt.Sprintf("thread_%d", threadI), threadI)
	}
	// release the threads
	stepStartTime := time.Now()
	for threadI := 0; threadI < threads; threadI++ {
		syncChan[threadI] <- CONTINUE
	}
	// iterate over steps
	for step := int64(0); step < runSteps; step++ {
		//logger.Log.With("cmd", "run").With("actor", "core").
		//	With("step", step).Debug("starting")
		for subStep := 0; subStep < subSteps; subStep++ {
			//logger.Log.With("cmd", "run").With("actor", "core").
			//	With("step", step).With("substep", subStep).
			//	Debug("waiting")
			for threadI := 0; threadI < threads; threadI++ {
				cont := <-syncChan[threadI]
				if cont == HALT {
					logger.Log.With("cmd", "run").Info("received HALT message from thread, shutting down core")
					panic("unimplemented graceful termination")
				}
			}
			// do atomic stuff at end of substep
			{
			}
			if subStep == subSteps-1 {
				// do atomic stuff at end of step
				runTime := time.Now().Sub(stepStartTime)
				logger.Log.With("cmd", "run").With("actor", "core").
					With("step", step).With("run_time", runTime).Info("finished")
				stepStartTime = time.Now()
			}

			// release the threads
			for threadI := 0; threadI < threads; threadI++ {
				syncChan[threadI] <- CONTINUE
			}
		}
	}
	// wait until the threads are done
	logger.Log.With("cmd", "run").With("actor", "core").Debug("waiting for threads")
	wgThreadsDone.Wait()
	logger.Log.With("cmd", "run").With("actor", "core").Debug("done")
}

func runThread(wgDone *sync.WaitGroup, syncChan chan engineMsg, name string, id int) {
	defer wgDone.Done()
	logger.Log.With("cmd", "run").With("actor", name).Debug("started")

	// configure the thread
	subSteps := viper.GetInt("substeps")

	// wait until released
	state := <-syncChan
	for step := int64(0); step < runSteps && state != HALT; step++ {
		//logger.Log.With("cmd", "run").With("actor", name).
		//	With("step", step).Debug("starting")
		for subStep := 0; subStep < subSteps && state != HALT; subStep++ {
			workTimeMS := 10 // + rand.Intn(10)
			//logger.Log.With("cmd", "run").With("actor", name).
			//	With("step", step).With("substep", subStep).
			//	With("workTimeMS", workTimeMS).Debug("working")
			time.Sleep((time.Duration)(workTimeMS) * time.Millisecond)
			//logger.Log.With("cmd", "run").With("actor", name).
			//	With("step", step).With("substep", subStep).
			//	Debug("letting core know we are done")
			// send back our message that we are ready to continue
			syncChan <- CONTINUE
			//logger.Log.With("cmd", "run").With("actor", name).
			//	With("step", step).With("substep", subStep).
			//	Debug("waiting for signal to continue")
			// wait for go ahead to continue
			state = <-syncChan
		}
	}
}
