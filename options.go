package initializer

import (
	"fmt"

	"github.com/worldline-go/logz"
)

type optionInit struct {
	msg string

	logzOptions []logz.Option
}

type OptionInit func(options *optionInit)

// WithMsg is a function that sets the message to be logged when the application starts.
//
// This will override the default message.
func WithMsgf(format string, a ...any) OptionInit {
	return func(options *optionInit) {
		options.msg = fmt.Sprintf(format, a...)
	}
}

func WithOptionsLogz(logzOpts ...logz.Option) OptionInit {
	return func(options *optionInit) {
		options.logzOptions = logzOpts
	}
}

func optionInitRunner(options ...OptionInit) *optionInit {
	option := &optionInit{}
	for _, opt := range options {
		opt(option)
	}

	return option
}
