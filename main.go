package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	data, err := ioutil.ReadFile(filePath)
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

	err = ioutil.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		fmt.Println("Failed to write tasks to file:", err)
	}
}



func main() {
	var taskName string
	var minutes, seconds int

	fmt.Print("Enter task name: ")
	fmt.Scanln(&taskName)

	fmt.Print("Enter minutes: ")
	_, err := fmt.Scan(&minutes)
	if err != nil || minutes < 0 {
		fmt.Println("Invalid input for minutes. Please enter a non-negative number.")
		return
	}

	fmt.Print("Enter seconds: ")
	_, err = fmt.Scan(&seconds)
	if err != nil || seconds < 0 {
		fmt.Println("Invalid input for seconds. Please enter a non-negative number.")
		return
	}

	Timer(taskName, minutes, seconds)
}
