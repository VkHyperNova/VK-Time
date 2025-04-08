package main

import (
	"vk-time/internal/cmd"
	"vk-time/internal/util"
)

func main() {
	taskName, minutes := util.ParseFlags()
	cmd.StartTask(taskName, minutes)
}


