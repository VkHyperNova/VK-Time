package timer

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
	"vk-time/internal/audio"
)

func StartCountdownTimer(task string, duration time.Duration, paused *atomic.Bool, doneChan <-chan struct{}) {

	ticker := time.NewTicker(time.Second)

	defer ticker.Stop()

	var elapsed time.Duration

	lastTick := time.Now()

	for {

		select {

		case <-doneChan:

			return

		case <-ticker.C:

			if paused.Load() {

				lastTick = time.Now()

				continue

			}

			now := time.Now()

			elapsed += now.Sub(lastTick)

			lastTick = now

			remaining := duration - elapsed

			if remaining <= 0 {

				fmt.Printf("\r\033[K%s - Time passed: %s / %s", task, duration, duration)

				fmt.Println("\n⏰ Time’s up!")

				audio.SwitchToSpeakers()

				time.Sleep(time.Second)

				os.Exit(0)

			}

			fmt.Printf("\r\033[K%s - Time passed: %s / %s", task, elapsed.Truncate(time.Second), duration)
		}
	}
}

func StartCountdown(task string, paused *atomic.Bool, doneChan <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var elapsed time.Duration
	lastTick := time.Now()

	for {
		select {
		case <-doneChan:
			return
		case <-ticker.C:
			if paused.Load() {
				lastTick = time.Now()
				continue
			}

			now := time.Now()
			elapsed += now.Sub(lastTick)
			lastTick = now

			fmt.Printf("\r\033[K%s - Time passed: %s", task, elapsed.Truncate(time.Second))
		}
	}
}

