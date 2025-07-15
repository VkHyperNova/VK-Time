package main

import (
	"bufio"
	"fmt"
	"os"
	"sync/atomic"
	"time"
	"vk-time/internal/audio"
	"vk-time/internal/mouse"
	"vk-time/internal/storage"
	"vk-time/internal/timer"
	"vk-time/internal/util"
)

func main() {

	taskName, duration := util.ParseFlags()

	var paused atomic.Bool

	doneChan := make(chan struct{})

	go timer.StartCountdownTimer(taskName, duration, &paused, doneChan)

	go audio.PlayMP3(duration, &paused, doneChan)

	go mouse.StartMouseMover(duration, &paused, doneChan)

	start := time.Now()

	audio.SwitchToHeadphones()

	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("\nType 'p' to pause, 'r' to resume, 'q' to quit: ")

		scanner.Scan()

		input := scanner.Text()

		switch input {

		case "p":

			paused.Store(true)

			fmt.Println("⏸️  Paused.")

		case "r":

			paused.Store(false)

			fmt.Println("▶️  Resumed.")

		case "q":

			close(doneChan)

			fmt.Println("⏹️  Stopped.")

			audio.SwitchToSpeakers()

			elapsed := time.Since(start)

			tasks := storage.Tasks{}

			tasks.Save(taskName, elapsed)

			os.Exit(0)
		}
	}
}
