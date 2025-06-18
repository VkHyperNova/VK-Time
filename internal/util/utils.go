package util

import (
	"flag"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-vgo/robotgo"
)

func ParseFlags() (string, time.Duration) {
	// Define the flags
	taskName := flag.String("name", "", "a string flag (optional)")
	minutes := flag.Int("age", 0, "an integer flag (optional)")

	// Parse the flags
	flag.Parse()

	// Check for positional arguments (remaining arguments)
	args := flag.Args()
	if len(args) > 0 {
		*taskName = args[0] // Use the first positional argument as the string flag
	}

	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", minutes) // Use the second positional argument as the integer flag
	}

	// Convert to time.Duration
	duration := time.Duration(*minutes) * time.Minute

	return *taskName, duration
}

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


