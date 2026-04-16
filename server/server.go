package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	m "github.com/jsutton9/iron-orbit/match"
)

func handleTracking(match *m.Match, w http.ResponseWriter, r *http.Request) {
	tracking := make(chan m.TrackingUpdate)
	match.TrackingBroker.Sub <- tracking
	for u := range tracking {
		_, err := fmt.Fprintf(w, `{"id":%d,"p":{"x":%f,"y":%f},"v":{"x":%f,"y":%f}}
`, u.Id, u.P.X, u.P.Y, u.V.X, u.V.Y)
		if err != nil {
			panic(err)
		}
	}
	match.TrackingBroker.Unsub <- tracking
	for range tracking {}
}

func handleQuit(match *m.Match, w http.ResponseWriter, r *http.Request) {
	match.QuitChannel <- struct{}{}
}

func handlePause(match *m.Match, w http.ResponseWriter, r *http.Request) {
	match.TimeChannel <- m.Pause
}

func handlePlay(match *m.Match, w http.ResponseWriter, r *http.Request) {
	match.TimeChannel <- m.RealTime
}

func handle(match *m.Match, path string, handler func(match *m.Match, w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		handler(match, w, r)
	})
}

func handleDebugPage(match *m.Match, w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("./debug.html")
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprint(w, string(body[:]))
	if err != nil {
		panic(err)
	}
}

func Serve(match *m.Match) {
	handle(match, "/tracking", handleTracking)
	handle(match, "/pause", handlePause)
	handle(match, "/play", handlePlay)
	handle(match, "/quit", handleQuit)
	handle(match, "/debug.html", handleDebugPage)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
