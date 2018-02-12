package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

// var connDex = make(map[string]*websocket.Conn)
var generalState state

// var conMux sync.Mutex

func parseControl(c control) {
	fmt.Println("data is good:", c.Action)
	switch c.Action {
	case "pause":
		generalState.Playing = false
	case "play":
		generalState.Playing = true
	default:
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
		parseControl(c)
	}
}

func checkorigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: checkorigin,
}

func main() {
	// player.control()

	generalState.Playing = false
	generalState.Title = "cannabis"
	fmt.Println("starting")
	http.HandleFunc("/", handler)
	// http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
