package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// Task represents a task with a name and total duration in seconds
type Task struct {
	Name         string `json:"name"`
	TotalSeconds int    `json:"total_seconds"`
}

func Timer(taskName string, minutes int, seconds int) {
	totalSeconds := minutes*60 + seconds
	if totalSeconds <= 0 {
		fmt.Println("Invalid duration. Please enter a positive time.")
		return
	}

	duration := time.Duration(totalSeconds) * time.Second
	t := time.NewTimer(duration)
	defer t.Stop()

	fmt.Printf("Task '%s': Waiting for %v...\n", taskName, duration)
	<-t.C
	fmt.Println("Timer expired!")
	playSound()

	// Save or update the task
	saveTask(taskName, totalSeconds)
}

func playSound() {
	filePath := "alarm.wav"
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open sound file (%s): %v\n", filePath, err)
		return
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		fmt.Println("Failed to decode sound file:", err)
		return
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		fmt.Println("Failed to initialize speaker:", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		fmt.Println("Sound finished playing")
		wg.Done()
	})))

	wg.Wait() // Wait for the sound to finish playing
}

func saveTask(taskName string, additionalSeconds int) {

	filePath := "tasks.json"
	var tasks []Task

	// Read existing tasks from the file
	data, err := os.ReadFile(filePath)
	if err == nil {
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			fmt.Println("Failed to parse JSON file:", err)
			return
		}
	}

	// Update or add the task
	found := false
	for i := range tasks {
		if tasks[i].Name == taskName {
			tasks[i].TotalSeconds += additionalSeconds
			found = true
			break
		}
	}
	if !found {
		tasks = append(tasks, Task{Name: taskName, TotalSeconds: additionalSeconds})
	}

	// Write updated tasks back to the file
	updatedData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("Failed to encode tasks to JSON:", err)
		return
	}

	err = os.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		fmt.Println("Failed to write tasks to file:", err)
	}
}



func main() {
	var taskName string
	var minutes, seconds int

	fmt.Print("Enter task name and time (format: name minutes seconds): ")
	_, err := fmt.Scanf("%s %d %d", &taskName, &minutes, &seconds)
	if err != nil {
		fmt.Println("Invalid input. Please use the format: name minutes seconds")
		return
	}

	if minutes < 0 || seconds < 0 {
		fmt.Println("Minutes and seconds must be non-negative numbers.")
		return
	}

	Timer(taskName, minutes, seconds)
}

