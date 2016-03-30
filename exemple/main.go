package main

import (
	"net/http"
	"time"

	"github.com/gotvitch/go-sse"
)

func main() {

	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		sseConnection, err := sse.Upgrade(w, r, sse.DefaultOptions)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for {
			select {
			case <-time.After(time.Second):
				sseConnection.Send("time", time.Now())
			case <-sseConnection.Closed:
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe(":8080", nil)
}
