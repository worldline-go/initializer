package initializer

import (
	"context"

	"github.com/rakunlabs/into"
	"github.com/rs/zerolog/log"
	"github.com/worldline-go/logz"
)

// Init is a function that initializes the application.
//
// This function will initialize the logger and run the shutdown function on exit.
func Init(fn func(context.Context) error, options ...Option) {
	opt := optionInitRunner(options...)

	logz.InitializeLog(opt.logzOptions...)

	optionsInto := []into.Option{
		into.WithMsgf(opt.msg),
		into.WithLogger(logz.AdapterKV{Log: log.Logger}),
	}
	optionsInto = append(optionsInto, opt.intoOptions...)

	into.Init(fn, optionsInto...)
}
