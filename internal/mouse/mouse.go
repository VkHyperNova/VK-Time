package mouse

import (
	"sync/atomic"
	"time"

	"github.com/go-vgo/robotgo"
)

func StartMouseMover(duration time.Duration, paused *atomic.Bool, doneChan <-chan struct{}) {

	const movePixels = 38

	ticker := time.NewTicker(60 * time.Second)

	defer ticker.Stop()

	timeout := time.After(duration)

	moves := [][2]int{
		
		{0, -movePixels}, // Up

		{movePixels, 0},  // Right

		{0, movePixels},  // Down

		{-movePixels, 0}, // Left
	}

	moveIndex := 0

	for {

		select {

		case <-doneChan:

			return

		case <-timeout:

			return

		case <-ticker.C:

			if paused.Load() {
				continue
			}

			x, y := robotgo.Location()

			dx, dy := moves[moveIndex][0], moves[moveIndex][1]

			robotgo.Move(x+dx, y+dy)

			moveIndex = (moveIndex + 1) % len(moves)
		}
	}
}


