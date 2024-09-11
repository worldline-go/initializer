package initializer

import "github.com/rakunlabs/into"

var (
	// WaitGroup returns global wait group.
	WaitGroup = into.WaitGroup
	// ShutdownAdd add function will be called when the context is done.
	ShutdownAdd = into.ShutdownAdd
	// CtxCancel is a function that cancels the root context.
	CtxCancel = into.CtxCancel
)
