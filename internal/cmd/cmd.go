package cmd

import (
	"sync"
	"vk-time/internal/audio"
	"vk-time/internal/storage"
	"vk-time/internal/util"
)

func StartTask(taskName string, minutes int) {
	audio.SwitchToHeadphones()

	t := storage.Tasks{}
	t.ReadFile("tasks.json")

	taskExists := t.PrintTask(taskName)

	// Channel to signal stopping the mouse mover
	stopMouseMover := make(chan struct{})

	// Use WaitGroup to wait for all goroutines
	var wg sync.WaitGroup
	wg.Add(3)

	// Start moving the mouse
	go func() {
		defer wg.Done()
		util.StartMouseMover(stopMouseMover)
	}()

	// Start playing MP3
	go func() {
		defer wg.Done()
		audio.PlayMP3("default_music.mp3", minutes)
	}()

	// Start the timer
	go func() {
		defer wg.Done()
		util.Timer(taskName, minutes)
	}()

	// Wait for the timer and audio to complete
	wg.Wait()

	// After task is complete, stop mouse mover
	close(stopMouseMover)

	// Play alarm sound after both processes complete
	audio.PlayWav("alarm.wav")

	if taskExists {
		t.UpdateTask(taskName, minutes)
	} else {
		t.AddTask(taskName, minutes)
	}

	t.SaveTask()

	audio.SwitchToSpeakers()
}
