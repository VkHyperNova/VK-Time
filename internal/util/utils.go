package util

import (
	"flag"
	"fmt"
	"time"
	"github.com/go-vgo/robotgo"
)

func PrintCoundown(taskName string, minutes int) {

	totalSeconds := minutes * 60 // 20 minutes in seconds

	for remaining := totalSeconds; remaining >= 0; remaining-- {
		minutes := remaining / 60
		seconds := remaining % 60
		fmt.Printf("\r%s: %02d:%02d", taskName, minutes, seconds)
		time.Sleep(1 * time.Second)
	}
}

func Timer(taskName string, minutes int) {

	duration := time.Duration(minutes) * time.Minute

	t := time.NewTimer(duration)
	defer t.Stop()

	PrintCoundown(taskName, minutes)

	<-t.C
	fmt.Println("\nTimer expired!")
}

func ParseFlags() (string, int) {
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

	return *taskName, *minutes
}

// StartMouseMover moves the mouse in a rectangle shape (up, right, down, left) every minute per side, until stopped.
func StartMouseMover(stop chan struct{}) {
	const movePixels = 38 // Approx. 1 cm
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Define movement directions: (dx, dy)
	moves := [][2]int{
		{0, -movePixels}, // Up
		{movePixels, 0},  // Right
		{0, movePixels},  // Down
		{-movePixels, 0}, // Left
	}

	moveIndex := 0

	for {
		select {
		case <-ticker.C:
			x, y := robotgo.Location()
			dx, dy := moves[moveIndex][0], moves[moveIndex][1]
			robotgo.Move(x+dx, y+dy)

			moveIndex = (moveIndex + 1) % len(moves) // cycle through moves
		case <-stop:
			return
		}
	}
}

