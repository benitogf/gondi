package main

import (
	"errors"
	"flag"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
	"unsafe"

	"github.com/benitogf/gondi"
	"github.com/benitogf/gondi/mjpeg"
	"github.com/gorilla/mux"
)

var inputFlag = flag.String("input", "", "input to copy")
var outputFlag = flag.String("output", "copy", "output stream name")

var (
	NDIversion        string
	InputStreamName   string
	InputStreamAdress string
	FramesSentCount   int64
	FrameRateN        int32 = 30000
	FrameRateD        int32 = 1001
	CaptureCount      int64
	CaptureErrorCount int64
	CaptureNoneCount  int64
	NDISources        []*gondi.Source
)

func ndiToNDI(receiver *gondi.RecvInstance, sender *gondi.SendInstance) {
	for {
		CaptureCount++
		videoInput := gondi.NewVideoFrameV2()
		frametype := receiver.CaptureV3(videoInput, nil, nil, 1000)
		if frametype == gondi.FrameTypeNone {
			CaptureNoneCount++
		}
		if frametype == gondi.FrameTypeError {
			CaptureErrorCount++
		}
		if frametype == gondi.FrameTypeVideo {
			size := videoInput.LineStride * videoInput.Yres
			videoInputSlice := unsafe.Slice(videoInput.Data, size)
			frame := make([]byte, len(videoInputSlice))
			copy(frame, videoInputSlice)

			// if videoInput.FourCC == gondi.FourCCTypeUYVY {
			// 	log.Println("preview UYVY", len(frame))
			// }
			// if videoInput.FourCC == gondi.FourCCTypeBGRA {
			// 	log.Println("preview BGRA", len(frame))
			// }
			// if videoInput.FourCC == gondi.FourCCTypeUYVA {
			// 	log.Println("preview UYVA", len(frame))
			// }
			// if videoInput.FourCC == gondi.FourCCTypeBGRX {
			// 	log.Println("preview BGRX", len(frame))
			// }
			if videoInput.FourCC == gondi.FourCCTypeRGBA {
				gondi.SetPreviewFrame(*outputFlag, frame, int(videoInput.Xres), int(videoInput.Yres))
			}
			if videoInput.FourCC == gondi.FourCCTypeRGBX {
				gondi.SetPreviewFrame(*outputFlag, frame, int(videoInput.Xres), int(videoInput.Yres))
			}

			videoOutput := gondi.NewVideoFrameV2()
			videoOutput.FourCC = videoInput.FourCC
			videoOutput.FrameFormatType = videoInput.FrameFormatType
			videoOutput.Xres = videoInput.Xres
			videoOutput.Yres = videoInput.Yres
			videoOutput.FrameRateN = videoInput.FrameRateN
			videoOutput.FrameRateD = videoInput.FrameRateD
			FrameRateN = videoOutput.FrameRateN
			FrameRateD = videoOutput.FrameRateD
			videoOutput.Data = videoInput.Data

			sender.SendVideoFrame(videoOutput)
			receiver.FreeVideoV2(videoInput)
			FramesSentCount++
		}
	}
}

func getNDISources(finder *gondi.FindInstance) ([]*gondi.Source, error) {
	// Wait for sources to appear
	for {
		more := finder.WaitForSources(1000)
		if !more {
			break
		}
	}

	// Fetch the sources
	result := finder.GetCurrentSources()
	if len(result) == 0 {
		return nil, errors.New("no sources found")
	}

	return result, nil
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	flag.Parse()
	log.Println("Initializing NDI", *inputFlag)
	gondi.InitLibrary("")

	NDIversion = gondi.GetVersion()

	finder, err := gondi.NewFindInstance(true, "", "")
	if err != nil {
		log.Println("failed to create ndi finder", err)
		panic(err)
	}
	defer finder.Destroy()

	NDISources, err := getNDISources(finder)
	if err != nil {
		log.Println("failed to get ndi sources", err)
		panic(err)
	}

	input := NDISources[0]
	if *inputFlag != "" {
		found := false
		for _, source := range NDISources {
			// log.Println("===============", gondi.ExtractSourceName(source.Name()))
			if gondi.ExtractSourceName(source.Name()) == *inputFlag {
				found = true
				input = source
			}
		}

		if !found {
			log.Println("failed to find ndi feed: ", *inputFlag)
			log.Println("NDI sources found: ", len(NDISources))
			for _, source := range NDISources {
				log.Println("-- ", source.Name())
			}
			panic(errors.New("failed to find ndi feed"))
		}
	}

	InputStreamName = input.Name()
	InputStreamAdress = input.Address()

	// Set up receiver
	receiver, err := gondi.NewRecvInstance(&gondi.NewRecvInstanceSettings{
		SourceToConnectTo: input,
		ColorFormat:       gondi.RecvColorFormatRGBXRGBA,
		Bandwidth:         gondi.RecvBandwidthHighest,
		AllowVideoFields:  true,
	})
	receiver.Connect(input)
	if err != nil {
		log.Println("failed to receive ndi", err)
		panic(err)
	}
	defer receiver.Destroy()

	// Set up sender, block on both audio and video as we are using separate threads for audio and video
	sender, err := gondi.NewSendInstance(*outputFlag, "", true, true)
	if err != nil {
		log.Println("failed to send ndi", err)
		panic(err)
	}
	gondi.ClearPreview(*outputFlag)
	defer sender.Destroy()

	// Set up threads
	go ndiToNDI(receiver, sender)

	// Show info
	go func() {
		for {
			clear()
			totals, dropped := receiver.GetPerformance()
			log.Printf("version: %s\n", NDIversion)
			log.Println("input name: ", InputStreamName)
			log.Println("output name: ", *outputFlag)
			log.Println("output connections: ", sender.GetNumberOfConnections(10))
			log.Println("sources", len(NDISources))
			for _, source := range NDISources {
				log.Println("-- ", source.Name())
			}
			log.Println("input address: ", InputStreamAdress)
			log.Println("frames sent: ", FramesSentCount)
			log.Println("frame rate N: ", FrameRateN)
			log.Println("frame rate D: ", FrameRateD)
			log.Println("captures: ", CaptureCount)
			log.Println("capture error: ", CaptureErrorCount)
			log.Println("capture none: ", CaptureNoneCount)
			log.Println("received video frames: ", totals.VideoFrames)
			log.Println("dropped frames: ", dropped.VideoFrames)
			time.Sleep(1 * time.Second)
		}
	}()

	previewStream := mjpeg.Handler{
		Next: func(streamName string) (image.Image, error) {
			return gondi.GetPreview(streamName)
		},
		Options: &jpeg.Options{Quality: 20},
	}

	router := mux.NewRouter()
	router.Handle("/preview/{streamName}", previewStream)
	log.Fatal(http.ListenAndServe("0.0.0.0:8086", router))
}
