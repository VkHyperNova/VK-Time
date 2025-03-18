package audio

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

//go:embed alarm.wav meditation1.mp3
var embeddedFiles embed.FS

// wrapper for *bytes.Reader to implement io.ReadCloser
type readCloser struct {
	*bytes.Reader
}

func (rc *readCloser) Close() error {
	return nil // No actual closing needed
}

func PlaySound(name string) {
	// Open the embedded alarm.wav file
	alarmData, err := embeddedFiles.Open(name)
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

func PlayMP3(name string, minutes int) {
	// Convert minutes to duration
	duration := time.Duration(minutes) * time.Minute

	var fileData []byte
	var err error

	// Try to load the file from embedded storage
	fileData, err = embeddedFiles.ReadFile(name)
	if err != nil {
		// If not found in embed, try from disk
		fileData, err = os.ReadFile(name)
		if err != nil {
			fmt.Println("Error opening MP3 file:", err)
			return
		}
	}

	// Create a reader that implements io.ReadCloser
	mp3Reader := &readCloser{bytes.NewReader(fileData)}

	// Decode MP3
	streamer, format, err := mp3.Decode(mp3Reader)
	if err != nil {
		fmt.Println("Error decoding MP3:", err)
		return
	}
	defer streamer.Close()

	// Initialize speaker
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Create a done channel
	done := make(chan bool)

	// Play the MP3 file in a loop until duration is reached
	startTime := time.Now()
	go func() {
		for time.Since(startTime) < duration {
			streamer.Seek(0) // Restart the stream
			speaker.Play(beep.Seq(streamer, beep.Callback(func() {
				done <- true
			})))
			<-done
		}
		speaker.Clear()
		fmt.Println("Stopped playing after", minutes, "minutes")
	}()

	// Block until playback is complete
	time.Sleep(duration)
}
