package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	fmt.Println("starting")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("page.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("url:", r.Form["url"])
	}
}

func play() {
	f, err := os.Open("dogsong.mp3")
	check(err)
	s, format, err := mp3.Decode(f)
	check(err)

	f, err = os.Open("ohmy.mp3")
	check(err)
	s2, _, err := mp3.Decode(f)
	check(err)

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan struct{})
	ctrl := &beep.Ctrl{Streamer: beep.Seq(s, beep.Callback(func() {
		close(done)
	}))}
	speaker.Play(ctrl)
	time.AfterFunc(time.Second*3, func() {
		speaker.Lock()
		ctrl.Paused = true
		speaker.Unlock()
	})
	time.AfterFunc(time.Second*4, func() {
		speaker.Lock()
		ctrl.Paused = false
		speaker.Unlock()
	})

	time.AfterFunc(time.Second*10, func() {
		speaker.Lock()
		ctrl.Streamer = beep.Seq(s2, ctrl.Streamer)
		speaker.Unlock()
	})
	<-done
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
