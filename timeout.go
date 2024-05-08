package initializer

import (
	"log/slog"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type timeout struct {
	Duration time.Duration
}

func newTimeout(duration time.Duration) *timeout {
	return &timeout{Duration: duration}
}

func (t *timeout) wait(wg *sync.WaitGroup) {
	timerWait := time.NewTimer(t.Duration)

	wgTimeout := &sync.WaitGroup{}
	wgTimeout.Add(2)

	mutex := &sync.Mutex{}
	canceled := false

	go func() {
		defer func() {
			mutex.Lock()
			defer mutex.Unlock()

			if !canceled {
				wgTimeout.Add(-2)
			}

			canceled = true
		}()

		<-timerWait.C
		logFn(DefaultLogger, map[logger]func(){
			Zerolog: func() {
				log.Warn().Msg("timeout reached while waiting WaitGroup")
			},
			Slog: func() {
				slog.Warn("timeout reached while waiting WaitGroup")
			},
		})
	}()

	go func() {
		defer func() {
			mutex.Lock()
			defer mutex.Unlock()

			if !canceled {
				wgTimeout.Add(-2)
			}

			canceled = true
		}()

		wg.Wait()
	}()

	wgTimeout.Wait()

	if timerWait != nil {
		timerWait.Stop()
	}
}
