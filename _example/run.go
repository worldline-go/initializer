package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rakunlabs/logi"
	"github.com/rs/zerolog"

	"github.com/worldline-go/initializer"
	"github.com/worldline-go/logz"
)

var (
	version = "v0.0.0"
	commit  = "-"
	date    = "-"
)

func main() {
	// run the application
	initializer.Init(
		run,
		initializer.WithMsgf("awesome-service version:[%s] commit:[%s] date:[%s]", version, commit, date),
		initializer.WithWaitTimeout(5*time.Second),
		// initializer.WithInitLog(false),
		initializer.WithOptionsLogz(
			logz.WithCaller(false),
			logz.WithLevel(zerolog.LevelDebugValue),
		),
		initializer.WithOptionsLogi(
			logi.WithCaller(true),
			logi.WithLevel(zerolog.LevelDebugValue),
		),
		initializer.WithLogger(initializer.Slog),
	)
}

func run(ctx context.Context, wg *sync.WaitGroup) error {
	<-ctx.Done()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// above than WithWaitTimeout
		time.Sleep(10 * time.Second)
	}()

	return fmt.Errorf("something went wrong")
}
