package initializer

import (
	"context"
	"log/slog"
	"sync"

	"github.com/rs/zerolog/log"
)

type ShutdownHolder struct {
	ctxCancel context.CancelFunc
	funcs     []shutdownInfo

	mutex sync.Mutex
	isRun bool
}

type shutdownInfo struct {
	name string
	fn   func() error
}

var Shutdown ShutdownHolder

// Cancel is a function that cancels the root context.
//
// This helps to stop the application gracefully without any errors.
func (s *ShutdownHolder) CtxCancel() {
	if s.ctxCancel == nil {
		return
	}

	s.ctxCancel()
}

func (s *ShutdownHolder) Add(fn func() error, options ...OptionShutdownAdd) {
	option := optionShutdownAdd{}
	for _, opt := range options {
		opt(&option)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.funcs = append(s.funcs, shutdownInfo{
		name: option.name,
		fn:   fn,
	})
}

func (s *ShutdownHolder) Run(options ...OptionShutdownRun) {
	option := optionShutdownRun{}
	for _, opt := range options {
		opt(&option)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRun && option.once {
		return
	}

	// run opposite order
	for i := len(s.funcs) - 1; i >= 0; i-- {
		inf := s.funcs[i]

		if err := inf.fn(); err != nil {
			logFn(DefaultLogger, map[logger]func(){
				Zerolog: func() {
					log.Err(err).Str("name", inf.name).Msg("shutdown error")
				},
				Slog: func() {
					slog.Error("shutdown error", slog.String("name", inf.name), slog.String("error", err.Error()))
				},
			})
		}
	}

	s.isRun = true
}

type optionShutdownAdd struct {
	name string
}

type OptionShutdownAdd func(options *optionShutdownAdd)

func WithShutdownName(name string) OptionShutdownAdd {
	return func(options *optionShutdownAdd) {
		options.name = name
	}
}

type optionShutdownRun struct {
	once bool
}

type OptionShutdownRun func(options *optionShutdownRun)

func WithShutdownRunOnce() OptionShutdownRun {
	return func(options *optionShutdownRun) {
		options.once = true
	}
}
