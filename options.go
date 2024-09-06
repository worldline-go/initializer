package initializer

import (
	"fmt"

	"github.com/rakunlabs/into"
	"github.com/worldline-go/logz"
)

type option struct {
	msg string

	logzOptions []logz.Option
	intoOptions []into.Option
}

type Option func(options *option)

// WithMsg is a function that sets the message to be logged when the application starts.
//
// This will override the default message.
func WithMsgf(format string, a ...any) Option {
	return func(options *option) {
		options.msg = fmt.Sprintf(format, a...)
	}
}

func WithOptionsLogz(logzOpts ...logz.Option) Option {
	return func(options *option) {
		options.logzOptions = logzOpts
	}
}

func WithOptionsInto(intoOpts ...into.Option) Option {
	return func(options *option) {
		options.intoOptions = intoOpts
	}
}

func optionInitRunner(options ...Option) *option {
	option := &option{}

	for _, opt := range options {
		opt(option)
	}

	return option
}
