package player

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	mplayer "github.com/robindiddams/go-mplayer"
)

// Sound can be played, Paused and data pulled out
type Sound struct {
	Name string
	Path string
}

func init() {
	mplayer.StartSlave()
}

// Play plays a song from the beginning
func Play(song *Sound) {
	mplayer.SendCommand(fmt.Sprintf("loadfile %s", song.Path))
}

// Resume sound
func Resume() {
	fmt.Println("resuming")
	mplayer.Paused = false
	mplayer.SendCommand("pause")
}

// Pause sound
func Pause() {
	fmt.Println("pausing")
	mplayer.Paused = true
	mplayer.SendCommand("pause")
}

func RegisterStopHandler(callback func()) {
	mplayer.RegisterStopHandler(callback)
}

// NewSound will mint a sound struct from a filename
// you must give it a callback to call when it finishes playing
func NewSound(file *os.File) (*Sound, error) {
	path, err := catalogueSong(file)
	if err != nil {
		return nil, fmt.Errorf("error cataloguing file %s", err.Error())
	}

	// fmt.Println(path)
	song := Sound{
		Path: path,
	}
	arr := strings.Split(path, ".mp3")
	if len(arr) == 2 {
		song.Name = strings.Replace(arr[1], "_", " ", -1)
	}

	return &song, nil
}

func catalogueSong(f *os.File) (string, error) {
	//put song in temp dir and return where it is
	name := f.Name()
	name = strings.Replace(name, " ", "_", -1)
	newhome, err := ioutil.TempFile("", name)
	if err != nil {
		return "", err
	}
	defer newhome.Close()
	_, err = io.Copy(newhome, f)
	if err != nil {
		return "", err
	}
	// if we ever have a db register the path here
	err = f.Close()
	if err != nil {
		return "", err
	}
	err = os.Remove(f.Name())
	if err != nil {
		return "", err
	}
	name = newhome.Name()
	return name, nil
}
