package util

import (
	"flag"
	"fmt"
	"time"

)

func ParseFlags() (string, time.Duration) {

	taskName := flag.String("name", "", "a string flag (optional)")

	minutes := flag.Int("age", 0, "an integer flag (optional)")

	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {

		*taskName = args[0] 
	}

	if len(args) > 1 {

		fmt.Sscanf(args[1], "%d", minutes) 
	}

	duration := time.Duration(*minutes) * time.Minute

	return *taskName, duration
}

