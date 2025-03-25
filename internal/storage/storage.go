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

func (t *Tasks) AddTask(taskName string, taskTime int) bool {
	for i := range t.Tasks {
		if t.Tasks[i].Name == taskName {
			t.Tasks[i].TotalMinutes += taskTime
			fmt.Println(t.Tasks[i].Name+":", t.Tasks[i].TotalMinutes, "minutes")
			return true
		}
	}

	// If Task does not exist 
	t.Tasks = append(t.Tasks, Task{Name: taskName, TotalMinutes: taskTime})
	return false
}

func (t *Tasks) SaveTask(taskName string, minutes int) bool {

	localFilePath := "tasks.json"
	

	t.ReadFile(localFilePath)

	t.AddTask(taskName, minutes) 

	updatedData, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		fmt.Println("Failed to encode tasks to JSON:", err)
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
	return true
}
