# Initializer

Initializer is a simple library to initialize main go service.

```sh
go get github.com/worldline-go/initializer
```

> This library customized to use https://github.com/rakunlabs/into

## Usage

Default logger is __zerolog__

```go
var (
	version = "v0.0.0"
	commit  = "-"
	date    = "-"
)

func main() {
	// run the application
	initializer.Init(
		run,
		initializer.WithMsgf("awesome-service version:[%s] commit:[%s] date:[%s]", version, commit, date),
		// if you want to close the application after a certain time without waiting the waitgroup
		// initializer.WithWaitTimeout(0),  // 0 means no timeout as default (time.Duration)
		initializer.WithOptionsLogz(
			logz.WithLevel(zerolog.LevelDebugValue),
		))
}

func run(_ context.Context) error {
	// Do something here.
	log.Info().Msg("Hello World!")

	return fmt.Errorf("something went wrong")
}
```

Add shutdown function, it will be called when the context is done.

```go
initializer.ShutdownAdd(server.Close, "server")
```
