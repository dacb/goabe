package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/dacb/goabe/logger"
	"github.com/dacb/goabe/plugins"

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
		log := logger.Log.With(
			slog.Group("cmd",
				slog.String("cmd", "run"),
				slog.Int64("runSteps", runSteps),
			),
		)
		log.Info("run command called")
		ctx := context.WithValue(cmd.Context(), "log", log)

		err := plugins.LoadPlugins(ctx)
		if err != nil {
			log.Error("an error occurred loading the plugins")
			panic(err)
		}
		setupPluginHooks(ctx)

		runCore(ctx, Threads)
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

var pluginHooks map[int][]plugins.Hook

func setupPluginHooks(ctx context.Context) {
	pluginHooks = make(map[int][]plugins.Hook)
	for _, plugin := range plugins.LoadedPlugins {
		hooks := (*plugin).GetHooks()
		for _, hook := range hooks {
			subStep := hook.SubStep
			pluginHooks[subStep] = append(pluginHooks[subStep], hook)
		}
	}
}

func runCore(ctx context.Context, threads int) {
	// this is the logger from the command context
	log := ctx.Value("log").(*slog.Logger)

	subSteps := viper.GetInt("substeps")
	// this waitgroup is used to signal the close of the threads
	wgThreadsDone := new(sync.WaitGroup)
	wgThreadsDone.Add(threads)
	// channels
	syncChan := make([]chan engineMsg, threads)

	// spawn the threads
	for threadI := 0; threadI < threads; threadI++ {
		syncChan[threadI] = make(chan engineMsg)
		threadName := fmt.Sprintf("thread_%d", threadI)
		tctx := context.WithValue(ctx, "log", log.With("actor", threadName))
		go runThread(tctx, wgThreadsDone, syncChan[threadI], threadName, threadI)
	}
	// release the threads
	stepStartTime := time.Now()
	for threadI := 0; threadI < threads; threadI++ {
		syncChan[threadI] <- CONTINUE
	}

	// setup the log for the rest of the core's activities
	log = log.With("actor", "core")
	ctx = context.WithValue(ctx, "log", log)

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
					log.Info("received HALT message from thread, shutting down core")
					panic("unimplemented graceful termination")
				}
			}

			// do atomic stuff at end of substep
			// for each plugin that is registered for the core at this step and substep
			{
				hooks := pluginHooks[subStep]
				for _, hook := range hooks {
					if hook.Core != nil {
						err := hook.Core(ctx)
						if err != nil {
							// can this be made to report the plug in as well?
							log.Error(fmt.Sprintf("error occurred calling plugin hook %s", hook.Description))
							panic(err)
						}
					}
				}
			}
			if subStep == subSteps-1 {
				// do atomic stuff at end of step
				runTime := time.Now().Sub(stepStartTime)
				log.With("step", step).With("run_time", runTime).Info("finished")
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

func runThread(ctx context.Context, wgDone *sync.WaitGroup, syncChan chan engineMsg, name string, id int) {
	defer wgDone.Done()
	log := ctx.Value("log").(*slog.Logger)
	log.Debug("started")

	// configure the thread
	subSteps := viper.GetInt("substeps")

	// wait until released
	state := <-syncChan
	for step := int64(0); step < runSteps && state != HALT; step++ {
		//logger.Log.With("cmd", "run").With("actor", name).
		//	With("step", step).Debug("starting")
		for subStep := 0; subStep < subSteps && state != HALT; subStep++ {
			//logger.Log.With("cmd", "run").With("actor", name).
			//	With("step", step).With("substep", subStep).
			//	With("workTimeMS", workTimeMS).Debug("working")
			{
				hooks := pluginHooks[subStep]
				for _, hook := range hooks {
					if hook.Thread != nil {
						// make the thread call for this substep
						err := hook.Thread(ctx, id, name)
						if err != nil {
							// can this be made to report the plug in as well?
							log.Error(fmt.Sprintf("error occurred calling plugin hook %s", hook.Description))
							panic(err)
						}
					}
				}
			}
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
