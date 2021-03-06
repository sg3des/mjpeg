package mjpeg

import (
	"fmt"
	"net/http"
	"sync"
)

const (
	boundary = "MJPEGBOUNDARY"
	headerf  = "\r\n" +
		"--" + boundary + "\r\n" +
		"Content-Type: image/jpeg\r\n" +
		"Content-Length: %d\r\n\r\n"
)

// MJPEG represents a single video feed.
type MJPEG struct {
	clients map[chan []byte]bool
	sync.RWMutex
}

// NewStream initializes and returns a new Stream.
func New() *MJPEG {
	return &MJPEG{
		clients: make(map[chan []byte]bool),
	}
}

// ServeHTTP responds to HTTP requests with the MJPEG stream, implementing the http.Handler interface.
func (s *MJPEG) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+boundary)

	c := make(chan []byte)

	s.Lock()
	s.clients[c] = true
	s.Unlock()

	for {
		img := <-c

		fmt.Fprintf(w, headerf, len(img))
		if _, err := w.Write(img); err != nil {
			break
		}
	}

	s.Lock()
	delete(s.clients, c)
	s.Unlock()
	close(c)
}

func (s *MJPEG) UpdateFrame(img []byte) {
	s.RLock()
	for c := range s.clients {
		select {
		case c <- img:
		default:
			// log.Warning("queue full")
		}
	}
	s.RUnlock()

	return
}
