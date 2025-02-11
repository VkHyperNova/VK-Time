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

func SaveTask(taskName string, minutes int) {

	localFilePath := "data/tasks.json"
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

	backupPath := "/media/veikko/VK DATA/DATABASES/TIME/tasks.json"
	err = os.WriteFile(backupPath, updatedData, 0644)
	if err != nil {
		fmt.Println("Failed to write BACKUPPATH tasks to file:", err)
	}
}
