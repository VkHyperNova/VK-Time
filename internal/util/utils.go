package util

import (
	"fmt"
	"time"
)

func PrintCoundown(taskName string, minutes int) {

	totalSeconds := minutes * 60 // 20 minutes in seconds

	for remaining := totalSeconds; remaining >= 0; remaining-- {
		minutes := remaining / 60
		seconds := remaining % 60
		fmt.Printf("\r%s: %02d:%02d", taskName, minutes, seconds)
		time.Sleep(1 * time.Second)
	}
}
