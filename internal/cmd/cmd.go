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

	// Use WaitGroup to wait for audio and timer (not mouse mover)
	var wg sync.WaitGroup
	wg.Add(2)

	// Start moving the mouse (don't wait for this one in WaitGroup)
	go util.StartMouseMover(stopMouseMover)

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

	// Wait for timer and audio to complete
	wg.Wait()

	// Then stop the mouse mover
	close(stopMouseMover)

	// Play alarm sound
	audio.PlayWav("alarm.wav")

	if taskExists {
		t.UpdateTask(taskName, minutes)
	} else {
		t.AddTask(taskName, minutes)
	}

	t.SaveTask()

	audio.SwitchToSpeakers()
}

