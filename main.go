package main

import (
	"bufio"
	"fmt"
	"os"
	"sync/atomic"
	"vk-time/internal/audio"
	"vk-time/internal/mouse"
	"vk-time/internal/timer"
	"vk-time/internal/util"
)

func main() {

	taskName, duration := util.ParseFlags()

	var paused atomic.Bool
	done := make(chan struct{})

	if duration == 0 {
		go timer.StartCountdown(taskName, &paused, done)
	} else {
		go timer.StartCountdownTimer(taskName, duration, &paused, done)
	}
	go audio.PlayMP3(&paused, done)
	go mouse.StartMouseMover(&paused, done)

	audio.SwitchToHeadphones()
	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("\nType 'p' to pause, 'r' to resume, 'q' to quit: ")
		scanner.Scan()
		input := scanner.Text()

		switch input {

		case "p":
			audio.SwitchToSpeakers()
			paused.Store(true)
			fmt.Println("⏸️  Paused.")

		case "r":
			audio.SwitchToHeadphones()
			paused.Store(false)
			fmt.Println("▶️  Resumed.")

		case "q":
			close(done)
			fmt.Println("⏹️  Stopped.")
			audio.SwitchToSpeakers()
			os.Exit(0)
		}
	}
}
