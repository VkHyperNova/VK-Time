package main

import (
	"flag"
	"fmt"
	"sync"
	"vk-time/internal/audio"
	"vk-time/internal/storage"
	"vk-time/internal/timer"
)

func main() {
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

	if *minutes < 0 {
		fmt.Println("Minutes must be non-negative.")
		return
	}

	// Use WaitGroup to wait for both goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Start playing MP3 in a goroutine
	go func() {
		defer wg.Done()
		audio.PlayMP3("default_song.mp3", *minutes)
	}()

	// Start the timer in a goroutine
	go func() {
		defer wg.Done()
		timer.Timer(*taskName, *minutes)
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	// Play alarm sound after both processes complete
	audio.PlaySound("alarm.wav")

	// Save task after completion
	t := storage.Tasks{}
	t.SaveTask(*taskName, *minutes)
}
