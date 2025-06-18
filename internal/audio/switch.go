package audio

import (
	"fmt"
	"os/exec"
)


func switchAudioSink(sinkName, label string) {

	cmd := exec.Command("pactl", "set-default-sink", sinkName)

	err := cmd.Run()

	if err != nil {

		fmt.Printf("Error switching to %s: %v\n", label, err)

		return
	}

	fmt.Printf("Audio output switched to %s\n", label)
}

func SwitchToHeadphones() {

	headphones := "alsa_output.usb-Corsair_CORSAIR_HS70_Pro_Wireless_Gaming_Headset-00.analog-stereo"

	switchAudioSink(headphones, "Headphones")
}

func SwitchToSpeakers() {

	speakers := "alsa_output.usb-Creative_Technology_Ltd_Sound_Blaster_Play__3_00301228-00.analog-stereo"
	
	switchAudioSink(speakers, "Speakers")
}