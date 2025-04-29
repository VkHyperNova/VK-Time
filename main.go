package main

import (
	"sync"
	"vk-time/internal/audio"
	"vk-time/internal/storage"
	"vk-time/internal/util"
)

func main() {
	// Parse command-line flags for task name and duration in minutes
	taskName, minutes := util.ParseFlags()

	// Route audio output to headphones
	audio.SwitchToHeadphones()

	// Load existing tasks
	tasks := storage.Tasks{}
	tasks.ReadFile("tasks.json")

	// Check if the task already exists
	taskExists := tasks.PrintTask(taskName)

	// Channel to signal mouse mover goroutine to stop
	stopMouseMover := make(chan struct{})

	// Use WaitGroup to wait for audio and timer to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Start mouse movement (doesn't block main goroutine)
	go util.StartMouseMover(stopMouseMover)

	// Start background audio playback
	go func() {
		defer wg.Done()
		audio.PlayMP3("default_music.mp3", minutes)
	}()

	// Start countdown timer
	go func() {
		defer wg.Done()
		util.Timer(taskName, minutes)
	}()

	// Wait for both timer and audio to finish
	wg.Wait()

	// Stop the mouse mover
	close(stopMouseMover)

	// Play alarm sound
	audio.PlayWav("alarm.wav")

	// Update or add task based on existence
	if taskExists {
		tasks.UpdateTask(taskName, minutes)
	} else {
		tasks.AddTask(taskName, minutes)
	}

	// Save tasks to disk
	tasks.SaveTask()

	// Switch audio output back to speakers
	audio.SwitchToSpeakers()
}
