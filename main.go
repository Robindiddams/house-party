package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/robindiddams/house-party/player"

	"github.com/gorilla/websocket"
	wrapper "github.com/robindiddams/house-party/youtube-dl-wrapper"
)

type state struct {
	Playing   bool            `json:"playing"`
	Title     string          `json:"title"`
	SongQueue []*player.Sound `json:"songQueue"`
}

type control struct {
	Action string `json:"action"`
	Meta   string `json:"meta"`
}

var generalState state

var controlChan = make(chan control, 100)
var fetchChan = make(chan string, 100)

// var soundChan = make(chan *player.Sound, 100)

func processControl() {
	var song *player.Sound

	for c := range controlChan {
		log.Println("\ndata is good:", c.Action)
		switch c.Action {
		case "pause":
			if song != nil {
				player.Pause()
				generalState.Playing = false
			}
		case "play":
			if song != nil {
				player.Resume()
				generalState.Playing = true
			} else {
				//if theres no song go to the next one
				controlChan <- control{Action: "next"}
			}
		case "queue":
			fmt.Println("got a url", c.Meta)
			fetchChan <- c.Meta
		case "next":
			fmt.Println("going to next song")
			//TODO: tell the frontend to spinner here

			if len(generalState.SongQueue) > 0 {
				// good theres a song in the tank!
				// if song != nil {
				// 	//song isn't nil, stop that sukka!
				// 	player.Pause()
				// }
				song, generalState.SongQueue = generalState.SongQueue[0], generalState.SongQueue[1:]
				generalState.Title = song.Name
				player.Play(song)
				generalState.Playing = true
			} else {
				fmt.Println("no next song in queue, doing nothing")
			}

		default:
			fmt.Println("huh:", c.Action)
		}
	}
}

func downloadSongs() {
	for url := range fetchChan {
		f, err := wrapper.DownloadMp3(url)
		if err != nil {
			fmt.Println("error:", err.Error())
			continue
		}

		//debug
		fmt.Println("mp3 got")
		info, _ := f.Stat()
		fmt.Println(info.Name(), info.Size())

		song, err := player.NewSound(f)
		if err != nil {
			fmt.Println("error:", err.Error())
			continue
		}

		fmt.Println("song sound created")

		generalState.SongQueue = append(generalState.SongQueue, song)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("incomming connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("normal err", err)
		return
	}
	fmt.Println("connection established", conn.RemoteAddr().String())
	// connDex[conn.RemoteAddr().String()] = conn
	fmt.Println("writing state")

	err = conn.WriteJSON(generalState)
	if err != nil {
		fmt.Println("error writing to connection", err)
	}

	// detect changes to state and update
	go func() {
		oldState := generalState
		for {
			if oldState.Playing != generalState.Playing {
				fmt.Println("updating state...", oldState.Playing, generalState.Playing)
				err := conn.WriteJSON(generalState)
				if err != nil {
					if strings.Contains(err.Error(), "use of closed network connection") || strings.Contains(err.Error(), "close sent") {
						return
					}
					fmt.Println("error writing to connection", err)
				}
				oldState = generalState
			}
		}
	}()
	for {
		var c control
		err = conn.ReadJSON(&c)
		fmt.Println("incomming data")
		if err != nil {
			if err.Error() == "websocket: close 1006 (abnormal closure): unexpected EOF" {
				// fmt.Println("removing conn", len(connDex))
				conn.Close()
				return
			}
			fmt.Println("read err", err)
			return
		}
		controlChan <- c
	}
}

func checkorigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: checkorigin,
}

func startChannels() {
	// controlChan = make(chan control, 100)
	// fetchChan = make(chan string, 100)
	// soundChan = make(chan *player.Sound, 100)

	go processControl()
	go downloadSongs()
}

func main() {
	startChannels()

	player.RegisterStopHandler(func() {
		// when the song completes tell the player to move on
		fmt.Println("song end")
		controlChan <- control{Action: "next"}
	})
	generalState.Playing = false
	generalState.Title = "nothing"
	fmt.Println("starting")
	http.HandleFunc("/", handler)
	// http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
