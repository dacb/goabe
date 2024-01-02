package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	slogmulti "github.com/samber/slog-multi"
)

func main() {
	threads := 2

	// initialize logging subsystem
	logfile, err := os.OpenFile("goabe.log.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	// fanout over the stdout in text and goabe.log.json as json
	defer logfile.Close()
	log := slog.New(
		slogmulti.Fanout(
			slog.NewJSONHandler(logfile, &slog.HandlerOptions{}),
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}),
		),
	)

	// welcome
	log.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now()),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("A message")

	// wait group to manage main engine threads
	wgThreads := new(sync.WaitGroup)
	wgThreads.Add(threads)
	for i := 0; i < threads; i++ {
		go goabeThread(wgThreads, fmt.Sprintf("thread_%03d", i), (i+1)*2, (i+1)*100)
	}
	wgThreads.Wait()
}

func goabeThread(wg *sync.WaitGroup, name string, actions int, waitTimeMs int) {
	defer wg.Done()
	fmt.Printf("thread %s started\n", name)
	for i := 0; i < actions; i++ {
		fmt.Println(name, i)
		time.Sleep(time.Millisecond * time.Duration(waitTimeMs))
	}
}
