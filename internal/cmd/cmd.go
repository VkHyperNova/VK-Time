package cmd

import (
	"sync"
	"vk-time/internal/audio"
	"vk-time/internal/storage"
	"vk-time/internal/util"
)

func StartTask(taskName string, minutes int) {

	t := storage.Tasks{}

	t.ReadFile("tasks.json")

	taskExists := t.PrintTask(taskName)

	// Use WaitGroup to wait for both goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Start playing MP3 in a goroutine
	go func() {
		defer wg.Done()
		audio.PlayMP3("default_music.mp3", minutes)
	}()

	// Start the timer in a goroutine
	go func() {
		defer wg.Done()
		util.Timer(taskName, minutes)
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	// Play alarm sound after both processes complete
	audio.PlayWav("alarm.wav")

	if taskExists {
		t.Update(taskName, minutes)
	} else {
		t.Add(taskName, minutes)
	}

	t.SaveTask()
}
