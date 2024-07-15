package mjpeg

import (
	"errors"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// ErrorEndOfStream signals the end of the Motion JPEG frames from a MJPEG
// stream
var ErrorEndOfStream = errors.New("end of Motion JPEG Stream")

// A Handler is an http.Handler that streams mjpeg using an image stream. Encoding
// quality can be controlled using the Options parameters.
type Handler struct {
	Next    func(string) (image.Image, error)
	Options *jpeg.Options
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("--------> MJPEG")
	w.Header().Add("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	streamName := mux.Vars(r)["streamName"]
	boundary := "\r\n--frame\r\nContent-Type: image/jpeg\r\n\r\n"
	for {
		img, err := h.Next(streamName)
		if err != nil {
			return
		}

		n, err := io.WriteString(w, boundary)
		if err != nil || n != len(boundary) {
			return
		}

		err = jpeg.Encode(w, img, h.Options)
		if err != nil {
			return
		}

		n, err = io.WriteString(w, "\r\n")
		if err != nil || n != 2 {
			return
		}
	}
}
