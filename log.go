package initializer

func logFn(l logger, mapFn map[logger]func()) {
	if fn, ok := mapFn[l]; ok {
		if fn != nil {
			fn()
		}
	}
}
