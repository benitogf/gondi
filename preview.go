package gondi

import (
	"errors"
	"image"
	"image/color"
)

const EMPTY_X = 1920
const EMPTY_Y = 1080

type Preview struct {
	StreamName string `json:"streamName"`
	IMG        *image.RGBA
	Width      int
	Height     int
}

var Previews []Preview

func GetPreview(streamName string) (*image.RGBA, error) {
	for _, preview := range Previews {
		if preview.StreamName == streamName {
			return preview.IMG, nil
		}
	}

	return GenerateAlpha(), errors.New("preview not found")
}

func GetPreviewIndex(streamName string) (int, error) {
	for index, preview := range Previews {
		if preview.StreamName == streamName {
			return index, nil
		}
	}

	return -1, errors.New("preview not found")
}

func ClearPreview(streamName string) {
	index, err := GetPreviewIndex(streamName)
	if err != nil {
		Previews = append(Previews, Preview{
			StreamName: streamName,
			IMG:        GenerateAlpha(),
			Width:      EMPTY_X,
			Height:     EMPTY_Y,
		})
		return
	}

	Previews[index].IMG = GenerateAlpha()
}

func SetPreviewFrame(streamName string, frame []byte, width, height int) {
	index, err := GetPreviewIndex(streamName)
	if err != nil {
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		img.Pix = frame
		Previews = append(Previews, Preview{
			StreamName: streamName,
			IMG:        img,
			Width:      width,
			Height:     height,
		})
		return
	}

	if Previews[index].Width != width || Previews[index].Height != height {
		Previews[index].IMG = image.NewRGBA(image.Rect(0, 0, width, height))
		Previews[index].Width = width
		Previews[index].Height = height
	}

	Previews[index].IMG.Pix = frame
}

func generateStatic(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			gray := color.RGBA{0, 0, 0, 0}
			img.SetRGBA(i, j, gray)
		}
	}

	return img.Pix
}

func GenerateAlpha() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, EMPTY_X, EMPTY_Y))
	img.Pix = generateStatic(EMPTY_X, EMPTY_Y)
	return img
}
