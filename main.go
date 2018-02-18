package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Ch3ck/youtube-dl/api"
	"github.com/gorilla/websocket"
)

type state struct {
	Playing bool   `json:"playing"`
	Title   string `json:"title"`
}

type control struct {
	Action string `json:"action"`
	Meta   string `json:"meta"`
}

var generalState state

var controlChan chan (control)

func processControl() {
	for c := range controlChan {
		fmt.Println("data is good:", c.Action)
		switch c.Action {
		case "pause":
			generalState.Playing = false
		case "play":
			generalState.Playing = true
		case "queue":
			fmt.Println("got a url", c.Meta)
			err := downloadSong(c.Meta)
			if err != nil {
				fmt.Println(err)
			}
		default:
		}
	}
}

func downloadSong(url string) error {
	fmt.Println("there", url)
	ID, err := api.GetVideoId(url)
	if err != nil {
		return err
	}
	err = api.APIGetVideoStream("mp3", ID, "/tmp/", 192)
	if err != nil {
		return err
	}
	return nil
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
	controlChan = make(chan control)
	go processControl()
}

func main() {
	startChannels()

	generalState.Playing = false
	generalState.Title = "cannabis"
	fmt.Println("starting")
	http.HandleFunc("/", handler)
	// http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
