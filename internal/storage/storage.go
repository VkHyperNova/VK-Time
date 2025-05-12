package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	Name         string `json:"name"`
	TotalMinutes int    `json:"total_minutes"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

func (t *Tasks) ReadFile(name string) {
	data, err := os.ReadFile(name)
	if err == nil {
		err = json.Unmarshal(data, t)
		if err != nil {
			fmt.Println("Failed to parse JSON file:", err)
			return
		}
	}

}

func (t *Tasks) AddTask(taskName string, taskTime int) {
	t.Tasks = append(t.Tasks, Task{Name: taskName, TotalMinutes: taskTime})
}

func (t *Tasks) UpdateTask(taskName string, taskTime int) {
	for i := range t.Tasks {
		if t.Tasks[i].Name == taskName {
			t.Tasks[i].TotalMinutes += taskTime
		}
	}

}

func (t *Tasks) SaveTask() bool {

	updatedData, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		fmt.Println("Failed to encode tasks to JSON:", err)
	}

	err = os.WriteFile("tasks.json", updatedData, 0644)
	if err != nil {
		fmt.Println("Failed to write LOCALPATH tasks to file:", err)
	}

	backupPath := "/media/veikko/VK DATA/DATABASES/TIME/tasks.json"
	err = os.WriteFile(backupPath, updatedData, 0644)
	if err != nil {
		fmt.Println("Failed to write BACKUPPATH tasks to file:", err)
	}
	return true
}

func (t *Tasks) PrintTask(taskName string) bool {
	for i := range t.Tasks {
		if t.Tasks[i].Name == taskName {
			hours, minutesLeft := CalculateMinutesAndHours(t.Tasks[i].TotalMinutes)
			fmt.Println(t.Tasks[i].Name+": ", hours, "h ", minutesLeft, " minutes")
			return true
		}
	}
	return false
}

func CalculateMinutesAndHours(minutes int) (int, int) {
	hours := minutes/60
	minutesLeft := minutes - (hours * 60)

	return hours, minutesLeft
}
