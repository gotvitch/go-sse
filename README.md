# go-sse
go-sse is a [Go](http://golang.org/) library that provides a simple way to handle SSE connections in a http handler.

## Features
- Connection timeout
- Detection and automatic incrementation of LastEventID

## Examples

```go
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

	http.ListenAndServe(":8080", nil)
}

```


## The MIT License (MIT)

MIT License

Copyright (c) 2016 Seb Gotvitch

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
