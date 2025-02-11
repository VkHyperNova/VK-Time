package timer

import (
	"fmt"
	"time"
	"vk-time/internal/util"
)

func Timer(taskName string, minutes int) {

	duration := time.Duration(minutes) * time.Minute

	t := time.NewTimer(duration)
	defer t.Stop()

	util.PrintCoundown(taskName, minutes)

	<-t.C
	fmt.Println("\nTimer expired!")
}