package sse

import "strconv"

type Conn struct {
	Closed      chan error
	LastEventID int
	event       chan *Event
	isOpen      bool
}

func (c *Conn) Send(event string, data interface{}) error {
	c.LastEventID++
	return c.SendEvent(&Event{
		Event: event,
		ID:    strconv.Itoa(c.LastEventID),
		Data:  data,
	})
}

func (c *Conn) SendEvent(event *Event) error {
	if !c.isOpen {
		return ErrConnectionClosed
	}
	c.event <- event
	return nil
}

func (c *Conn) IsOpen() bool {
	return c.isOpen
}

func (c *Conn) Close() {
	c.Closed <- nil
}
