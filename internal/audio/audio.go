package audio

import (
	"embed"
	"fmt"
	"sync/atomic"
	"time"

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

	end := time.Now().Add(duration)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-doneChan:
			speaker.Clear()
			return
		case <-ticker.C:
			speaker.Lock()
			ctrl.Paused = paused.Load()
			speaker.Unlock()

			if time.Now().After(end) {
				speaker.Clear()
				return
			}
		}
	}
}



