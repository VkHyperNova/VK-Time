package audio

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"sync/atomic"
	"time"
	"vk-time/internal/storage"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

//go:embed music.mp3
var embeddedFiles embed.FS

func PlayMP3(duration time.Duration, paused *atomic.Bool, doneChan <-chan struct{}) {
	embeddedFile, err := embeddedFiles.Open("music.mp3")
	if err != nil {
		fmt.Println("❌ Failed to open embedded MP3:", err)
		return
	}
	defer embeddedFile.Close()

	streamer, format, err := mp3.Decode(embeddedFile)
	if err != nil {
		fmt.Println("❌ Failed to decode MP3:", err)
		return
	}
	defer streamer.Close()

	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(ctrl)

	end := time.After(duration)

	for {
		select {
		case <-doneChan:
			speaker.Clear()
			return
		case <-time.After(100 * time.Millisecond):
			speaker.Lock()
			ctrl.Paused = paused.Load()
			speaker.Unlock()
		case <-end:
			speaker.Clear()
			return
		}
	}
}

func CountdownTimer(task string, duration time.Duration, paused *atomic.Bool, doneChan <-chan struct{}) {
	
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var elapsed time.Duration
	lastTick := time.Now()

	for {
		select {
		case <-doneChan:
			return
		case <-ticker.C:
			if paused.Load() {
				lastTick = time.Now()
				continue
			}

			now := time.Now()
			elapsed += now.Sub(lastTick)
			lastTick = now

			remaining := duration - elapsed
			if remaining <= 0 {
				fmt.Printf("\r\033[K%s - Time passed: %s / %s", task, duration, duration)
				fmt.Println("\n⏰ Time’s up!")
				SwitchToSpeakers()
				time.Sleep(time.Second)
				tasks := storage.Tasks{}
				tasks.Save(task, duration)
				os.Exit(0)
			}

			fmt.Printf("\r\033[K%s - Time passed: %s / %s", task, elapsed.Truncate(time.Second), duration)
		}
	}
}

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
