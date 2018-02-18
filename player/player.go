package player

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// Sound can be played, Paused and data pulled out
type Sound struct {
	Duration int
	Name     string
	Filepath string
	Burn     int
	ctrl     *beep.Ctrl
}

// Play sound
func (s *Sound) Play() {
	speaker.Lock()
	fmt.Println("playing", s.ctrl.Paused)
	if s.ctrl.Paused {
		s.ctrl.Paused = false
	}
	speaker.Unlock()
}

// Pause sound
func (s *Sound) Pause() {
	speaker.Lock()
	fmt.Println("pausing")
	if !s.ctrl.Paused {
		s.ctrl.Paused = true
	}
	speaker.Unlock()
}

// OpenSound will mint a sound struct from a filename
// you must give it a callback to call when it finishes playing
func OpenSound(filename string, callback func()) (*Sound, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return OpenSoundFile(f, callback)
}

// OpenSoundFile will mint a sound struct from an os.file
// you must give it a callback to call when it finishes playing
func OpenSoundFile(f *os.File, callback func()) (*Sound, error) {
	s, format, err := mp3.Decode(f)
	if err != nil {
		return nil, err
	}
	d := s.Len() / format.SampleRate.N(time.Second)
	ctrl := &beep.Ctrl{Streamer: beep.Seq(s, beep.Callback(callback))}
	song := Sound{
		Duration: d,
		Name:     f.Name(),
		Burn:     0,
		ctrl:     ctrl,
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	song.ctrl.Paused = true
	speaker.Play(song.ctrl)
	return &song, nil
}

// Control will test
func Control() {
	s, err := OpenSound("dogsong.mp3", func() {
		fmt.Println("done")
	})
	if err != nil {
		panic(err)
	}
	s.Play()

	time.AfterFunc(time.Second*1, func() {
		s.Pause()
	})

	time.AfterFunc(time.Second*2, func() {
		s.Play()
	})
}
