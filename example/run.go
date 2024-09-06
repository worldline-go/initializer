package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rakunlabs/into"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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
		initializer.WithOptionsLogz(
			// logz.WithCaller(false),
			logz.WithLevel(zerolog.LevelDebugValue),
		),
		initializer.WithOptionsInto(
			into.WithWaitTimeout(5*time.Second),
		),
	)
}

func run(ctx context.Context) error {
	log.Warn().Msg("this is a warning message")
	<-ctx.Done()

	wg := initializer.WaitGroup(ctx)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// above than WithWaitTimeout
		time.Sleep(3 * time.Second)
	}()

	return fmt.Errorf("something went wrong")
}
