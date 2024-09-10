package initializer

import "github.com/rakunlabs/into"

var WaitGroup = into.WaitGroup

func ShutdownAdd(fn func() error, name string) {
	into.ShutdownAdd(fn, into.WithShutdownName(name))
}
