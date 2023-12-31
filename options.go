package initializer

import (
	"context"
	"fmt"

	"github.com/worldline-go/logz"
)

type optionInit struct {
	msg string
	ctx context.Context //nolint:containedctx // temporary

	logzOptions []logz.Option
	// initLog is a flag that indicates if the init message should be logged.
	initLog bool
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

// WithContext is a function that sets the context to be used as parent context.
func WithContext(ctx context.Context) OptionInit {
	return func(options *optionInit) {
		options.ctx = ctx
	}
}

func WithOptionsLogz(logzOpts ...logz.Option) OptionInit {
	return func(options *optionInit) {
		options.logzOptions = logzOpts
	}
}

func WithInitLog(v bool) OptionInit {
	return func(options *optionInit) {
		options.initLog = v
	}
}

func optionInitRunner(options ...OptionInit) *optionInit {
	option := &optionInit{
		ctx:     context.Background(),
		initLog: true,
	}

	for _, opt := range options {
		opt(option)
	}

	return option
}
