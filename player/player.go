package player

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
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

// RegisterStopHandler register a func for when
// theres nothing playing, ie when the player is stopped
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
	song := Sound{
		Path: path,
		Name: parseSongName(path),
	}
	return &song, nil
}

func parseSongName(path string) string {
	var name string
	dir := os.TempDir()
	fmt.Println(dir)
	re := regexp.MustCompile(dir + `(.*)\.mp3`)
	match := re.FindStringSubmatch(path)
	if len(match) == 2 {
		name = match[1]
		name = name[:strings.LastIndex(name, "-")]
		name = strings.Replace(name, "_", " ", -1)
	}
	return name
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
