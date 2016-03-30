package sse

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/manucorporat/sse"
)

var (
	ErrStreamingNotSupported = errors.New("Streaming not supported")
	ErrConnectionClosed      = errors.New("Connection already closed")
	ErrConnectionTimeout     = errors.New("Connection timeout")
)

type Event struct {
	ID    string
	Event string
	Data  interface{}
}

type Options struct {
	Timeout   int
	RetryTime int
}

var DefaultOptions = Options{
	Timeout:   30,
	RetryTime: 0,
}

func Upgrade(w http.ResponseWriter, r *http.Request, options Options) (*Conn, error) {

	f, ok := w.(http.Flusher)
	if !ok {
		return nil, ErrStreamingNotSupported
	}

	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")

	w.WriteHeader(http.StatusOK)
	// Write a empty string to avoid the warning "multiple response.WriteHeader calls" cause by Flush()
	w.Write([]byte(""))

	conn := &Conn{
		Closed: make(chan error),
		event:  make(chan *Event),
		isOpen: true,
	}

	if lastEventID, err := strconv.Atoi(r.Header.Get("Last-Event-ID")); err != nil {
		conn.LastEventID = lastEventID
	} else {
		conn.LastEventID = 0
	}

	if options.RetryTime > 0 {
		fmt.Fprintf(w, "retry: %d\n", options.RetryTime)
	}

	timeoutChannel := make(chan bool)
	closeNotifyChannel := w.(http.CloseNotifier).CloseNotify()

	if options.Timeout > 0 {
		go func() {
			time.Sleep(time.Duration(options.Timeout) * time.Second)
			timeoutChannel <- true
		}()
	}

	f.Flush()

	go func() {
		for {
			select {
			case <-timeoutChannel:
				conn.isOpen = false
				conn.Closed <- ErrConnectionTimeout
				return

			case <-closeNotifyChannel:
				conn.isOpen = false
				conn.Closed <- nil
				return

			case <-conn.Closed:
				conn.isOpen = false
				conn.Closed <- nil
				return

			case msg := <-conn.event:
				event := sse.Event{
					Id:    msg.ID,
					Event: msg.Event,
					Data:  msg.Data,
				}
				sse.Encode(w, event)
				f.Flush()
			}
		}
	}()

	return conn, nil
}
