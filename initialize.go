package initializer

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rakunlabs/logi"
	"github.com/rs/zerolog/log"
	"github.com/worldline-go/logz"
)

var (
	DefaultExitCode = 0
	DefaultLogger   = Zerolog

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
	opt := optionInitRunner(options...)
	DefaultLogger = opt.logger

	logz.InitializeLog(opt.logzOptions...)
	logi.InitializeLog(opt.logiOptions...)

	if opt.initLog {
		logFn(opt.logger, map[logger]func(){
			Zerolog: func() {
				log.Log().Msg("starting " + opt.msg)
			},
			Slog: func() {
				// without level check
				_ = slog.Default().Handler().Handle(context.Background(), slog.Record{
					Time:    time.Now(),
					Message: "starting " + opt.msg,
				})
			},
		})
	}

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
	ctx, ctxCancel := context.WithCancel(opt.ctx)

	Shutdown = ShutdownHolder{
		ctxCancel: ctxCancel,
	}

	defer func() {
		if opt.wgWaitTimeout > 0 {
			newTimeout(opt.wgWaitTimeout).wait(&wg)
		} else {
			wg.Wait()
		}
	}()
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

			logFn(opt.logger, map[logger]func(){
				Zerolog: func() {
					log.Warn().Msg("received shutdown signal")
				},
				Slog: func() {
					slog.Warn("received shutdown signal")
				},
			})

			ctxCancel()
		}

		Shutdown.Run(WithShutdownRunOnce())
	}()

	if err := fn(ctx, &wg); err != nil {
		SetExitCode(1, false)

		logFn(opt.logger, map[logger]func(){
			Zerolog: func() {
				log.Error().Err(err).Msgf("service closing: %s", opt.msg)
			},
			Slog: func() {
				slog.Error("service closing: "+opt.msg, slog.String("error", err.Error()))
			},
		})
	}
}
