package audio

import (
	"embed"
	"fmt"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)


//go:embed alarm.wav
var alarmFile embed.FS


func PlaySound() {
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
