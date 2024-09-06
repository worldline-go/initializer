package initializer

import (
	"context"

	"github.com/rakunlabs/into"
	"github.com/worldline-go/logz"
)

// Init is a function that initializes the application.
//
// This function will initialize the logger and run the shutdown function on exit.
func Init(fn func(context.Context) error, options ...Option) {
	opt := optionInitRunner(options...)

	logzOpt := logz.ReadOptions(opt.logzOptions...)
	logz.InitializeLog(logz.WithOption(logzOpt))

	logzOpt.Caller = new(bool)
	logzLogger := logz.Logger(logz.WithOption(logzOpt))

	optionsInto := []into.Option{
		into.WithMsgf(opt.msg),
		into.WithLogger(logz.AdapterKV{
			Log: logzLogger,
		}),
	}
	optionsInto = append(optionsInto, opt.intoOptions...)

	into.Init(fn, optionsInto...)
}
