package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// Embed the alarm.wav file
//go:embed alarm.wav
var alarmFile embed.FS

type Task struct {
	Name         string `json:"name"`
	TotalMinutes int    `json:"total_minutes"`
}

func Timer(taskName string, minutes int) {

	duration := time.Duration(minutes) * time.Minute

	t := time.NewTimer(duration)
	defer t.Stop()

	fmt.Printf("Task '%s': Waiting for %v...\n", taskName, duration)

	<-t.C
	fmt.Println("Timer expired!")

	playSound()

	saveTask(taskName, minutes)

	printTask(taskName)
}

func playSound() {
	// Open the embedded alarm.wav file
	alarmData, err := alarmFile.Open("alarm.wav")
	if err != nil {
		fmt.Println("Failed to open embedded sound file:", err)
		return
	}
	defer alarmData.Close()

	// Decode the WAV file
	streamer, format, err := wav.Decode(alarmData)
	if err != nil {
		fmt.Println("Failed to decode sound file:", err)
		return
	}
	defer streamer.Close()

	// Initialize the speaker
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		fmt.Println("Failed to initialize speaker:", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Play the sound
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		fmt.Println("Sound finished playing")
		wg.Done()
	})))

	wg.Wait()
}

func saveTask(taskName string, minutes int) {

	localFilePath := "tasks.json"
	backupPath := "/media/veikko/VK DATA/DATABASES/TIME/tasks.json"
	var tasks []Task

	data, err := os.ReadFile(localFilePath)
	if err == nil {
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			fmt.Println("Failed to parse JSON file:", err)
			return
		}
	}

	found := false
	for i := range tasks {
		if tasks[i].Name == taskName {
			tasks[i].TotalMinutes += minutes
			found = true
			break
		}
	}
	if !found {
		tasks = append(tasks, Task{Name: taskName, TotalMinutes: minutes})
	}

	updatedData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("Failed to encode tasks to JSON:", err)
		return
	}

	err = os.WriteFile(localFilePath, updatedData, 0644)
	if err != nil {
		fmt.Println("Failed to write LOCALPATH tasks to file:", err)
	}

	err = os.WriteFile(backupPath, updatedData, 0644)
	if err != nil {
		fmt.Println("Failed to write BACKUPPATH tasks to file:", err)
	}
}

func printTask(taskName string) {
	filePath := "tasks.json"
	var tasks []Task

	data, err := os.ReadFile(filePath)
	if err == nil {
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			fmt.Println("Failed to parse JSON file:", err)
			return
		}
	}

	for i := range tasks {
		if tasks[i].Name == taskName {
			fmt.Printf("[%s %d]\n", tasks[i].Name, tasks[i].TotalMinutes)
		}
	}
}

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

	printTask(*taskName)

	Timer(*taskName, *minutes)
}
