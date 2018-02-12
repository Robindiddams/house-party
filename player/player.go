package player

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func control() {
	f, err := os.Open("dogsong.mp3")
	check(err)
	s, format, err := mp3.Decode(f)
	check(err)

	f, err = os.Open("ohmy.mp3")
	check(err)
	s2, fmt2, err := mp3.Decode(f)
	check(err)

	fmt.Println(s.Len(), s2.Len())
	fmt.Println(format.SampleRate, fmt2.SampleRate)
	fmt.Println(s.Len()/format.SampleRate.N(time.Second), s2.Len()/fmt2.SampleRate.N(time.Second))

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan struct{})
	ctrl := &beep.Ctrl{Streamer: beep.Seq(s, beep.Callback(func() {
		close(done)
	}))}
	speaker.Play(ctrl)

	time.AfterFunc(time.Second*3, func() {
		speaker.Lock()
		ctrl.Paused = true
		fmt.Println(s.Len())
		speaker.Unlock()
	})
	time.AfterFunc(time.Second*4, func() {
		speaker.Lock()
		fmt.Println(s.Len())
		ctrl.Paused = false
		speaker.Unlock()
	})

	time.AfterFunc(time.Second*10, func() {
		speaker.Lock()
		ctrl.Streamer = beep.Seq(s2, ctrl.Streamer)
		fmt.Println(s.Len())
		speaker.Unlock()
	})
	<-done
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
