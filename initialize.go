package initializer

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/logz"
)

var (
	DefaultExitCode = 0

	exitCode = exitCodeHolder{}
	mutext   = sync.Mutex{}
)

type exitCodeHolder struct {
	code int
	set  bool
}

// SetExitCode is a function that sets the exit code.
//
// If used this function, the application will exit with the given code and not the default one.
// This function is thread safe.
//
// If override is set to true, the exit code set again.
func SetExitCode(code int, override bool) {
	mutext.Lock()
	defer mutext.Unlock()

	if exitCode.set && !override {
		return
	}

	exitCode.code = code
	exitCode.set = true
}

// Init is a function that initializes the application.
//
// This function will initialize the logger and run the shutdown function on exit.
func Init(fn func(context.Context, *sync.WaitGroup) error, options ...OptionInit) {
	opts := optionInitRunner(options...)
	logz.InitializeLog(opts.logzOptions...)

	log.Log().Msgf("starting %s", opts.msg)

	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}

		exitC := DefaultExitCode
		if exitCode.set {
			exitC = exitCode.code
		}

		os.Exit(exitC)
	}()

	wg := sync.WaitGroup{}
	ctx, ctxCancel := context.WithCancel(context.Background())

	defer wg.Wait()
	defer ctxCancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
		case <-signalChan:
			SetExitCode(1, false)

			log.Warn().Msg("received shutdown signal")
		}

		Shutdown.Run(WithShutdownRunOnce())
	}()

	if err := fn(ctx, &wg); err != nil {
		SetExitCode(1, false)

		log.Error().Err(err).Msgf("failed to run service, closing: %s", opts.msg)
	}
}
