package main

import (
	"context"
	"fmt"
	"sync"

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
	// Run the application.
	initializer.Init(
		run,
		initializer.WithMsgf("awesome-service version:[%s] commit:[%s] date:[%s]", version, commit, date),
		initializer.WithOptionsLogz(
			logz.WithCaller(false),
			logz.WithLevel(zerolog.LevelDebugValue),
		))
}

func run(ctx context.Context, _ *sync.WaitGroup) error {
	// Do something here.
	<-ctx.Done()

	return fmt.Errorf("something went wrong")
}
