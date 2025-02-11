package main

import (
	"flag"
	"fmt"
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
		fmt.Println("Minutes and seconds must be non-negative numbers.")
		return
	}

	timer.Timer(*taskName, *minutes)

	audio.PlaySound()

	storage.SaveTask(*taskName, *minutes)
}
