package initializer

import (
	"context"
	"fmt"
	"time"

	"github.com/rakunlabs/logi"
	"github.com/worldline-go/logz"
)

type logger uint8

const (
	Zerolog logger = iota
	Slog
)

type optionInit struct {
	msg string
	ctx context.Context //nolint:containedctx // temporary

	logzOptions []logz.Option
	logiOptions []logi.Option
	// initLog is a flag that indicates if the init message should be logged.
	initLog bool

	logger logger

	wgWaitTimeout time.Duration
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

func WithOptionsLogi(logiOpts ...logi.Option) OptionInit {
	return func(options *optionInit) {
		options.logiOptions = logiOpts
	}
}

func WithLogger(l logger) OptionInit {
	return func(options *optionInit) {
		options.logger = l
	}
}

func WithInitLog(v bool) OptionInit {
	return func(options *optionInit) {
		options.initLog = v
	}
}

func WithWaitTimeout(duration time.Duration) OptionInit {
	return func(options *optionInit) {
		options.wgWaitTimeout = duration
	}
}

func optionInitRunner(options ...OptionInit) *optionInit {
	option := &optionInit{
		ctx:           context.Background(),
		logger:        DefaultLogger,
		initLog:       true,
		wgWaitTimeout: DefaultWgTimeout,
	}

	for _, opt := range options {
		opt(option)
	}

	return option
}
