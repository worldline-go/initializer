package initializer

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type ShutdownHolder struct {
	funcs []shutdownInfo
	mutex sync.Mutex

	isRun bool
}

type shutdownInfo struct {
	name string
	fn   func() error
}

var Shutdown = ShutdownHolder{}

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

	for _, inf := range s.funcs {
		if err := inf.fn(); err != nil {
			log.Err(err).Str("name", inf.name).Msg("shutdown error")
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
