package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/AlexEidt/aio"
	"github.com/benitogf/gondi"
)

var (
	videoFrames     int64   = 0
	videoFPS        float64 = 0
	videoFPSFixed   float64 = 0
	videoFrameRateN int32   = 30000
	videoFrameRateD int32   = 1001

	audioFrames     int64  = 0
	audioFormat     string = ""
	audioSampleRate int32  = 0
	audioChannels   int32  = 0
	audioNumSamples int32  = 0
)

func streamNDI(fileName string, sender *gondi.SendInstance) {
	go audioToNDI(fileName, sender)
	go videoToNDI(fileName, sender)
}

func audioToNDI(fileName string, sender *gondi.SendInstance) {
	for {
		audio, err := aio.NewAudio(fileName, &aio.Options{
			Format: "f32",
		})
		if err != nil {
			log.Panic("failed to get audio", err)
		}

		for audio.Read() {
			audioInput := gondi.NewAudioFrameV3()
			audioInput.SampleRate = int32(audio.SampleRate())
			audioInput.NumChannels = int32(audio.Channels())
			audioInput.NumSamples = int32(len(audio.Samples().([]float32)) / 2)
			audioInput.Data = &audio.Buffer()[0]
			// audioInput.ChannelStride = 0

			audioFormat = audio.Format()
			audioSampleRate = audioInput.SampleRate
			audioChannels = audioInput.NumChannels
			audioNumSamples = audioInput.NumSamples

			sender.SendAudioFrame32f(audioInput)

			audioFrames++
		}
	}
}

func videoToNDI(fileName string, sender *gondi.SendInstance) {
	for {
		video, err := vidio.NewVideo(fileName)
		if err != nil {
			log.Panic("videoToNDI: failed to read video", err)
		}

		videoFrames = 0

		for video.Read() {
			videoFrame := gondi.NewVideoFrameV2()
			videoFrame.FourCC = gondi.FourCCTypeRGBA
			videoFrame.FrameFormatType = gondi.FrameFormatProgressive
			videoFrame.Xres = int32(video.Width())
			videoFrame.Yres = int32(video.Height())
			// videoFrame.LineStride = 0 // 2 bytes per pixel

			// The most common fractional fps values are 23.976 fps (24000 / 1001), 29.97 fps (30000 / 1001) or 59.94 fps (60000 / 1001)
			videoFPS = video.FPS()
			videoFPSFixed = float64(int(videoFPS*100)) / 100
			switch videoFPSFixed {
			case 23.97:
				videoFrame.FrameRateN = 24000
				videoFrame.FrameRateD = 1001
			case 29.97:
				videoFrame.FrameRateN = 30000
				videoFrame.FrameRateD = 1001
			case 30.00:
				videoFrame.FrameRateN = 30000
				videoFrame.FrameRateD = 1000
			case 59.94:
				videoFrame.FrameRateN = 60000
				videoFrame.FrameRateD = 1001
			case 60.00:
				videoFrame.FrameRateN = 60000
				videoFrame.FrameRateD = 3000
			default:
				videoFrame.FrameRateN = 30000
				videoFrame.FrameRateD = 1001
			}
			videoFrameRateN = videoFrame.FrameRateN
			videoFrameRateD = videoFrame.FrameRateD
			videoFrame.Data = &video.FrameBuffer()[0]

			sender.SendVideoFrame(videoFrame)

			videoFrames++
		}
	}
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	gondi.InitLibrary("")

	version := gondi.GetVersion()
	log.Printf("NDI version: %s\n", version)

	// Set up sender, block on both audio and video as we are using separate threads for audio and video
	sender, err := gondi.NewSendInstance("mock", "", true, true)
	if err != nil {
		log.Panic("failed to create ndi sender", err)
	}

	defer sender.Destroy()
	streamNDI("mock.mp4", sender)

	// Show info
	for {
		clear()
		log.Println("-- Video")
		log.Println("Frames: ", videoFrames)
		log.Println("FPS: ", videoFPS)
		log.Println("FPS[fixed]: ", videoFPSFixed)
		log.Println("FrameRateN: ", videoFrameRateN)
		log.Println("FrameRateD: ", videoFrameRateD)
		log.Println("-- Audio:")
		log.Println("Frames: ", audioFrames)
		log.Println("Format: ", audioFormat)
		log.Println("Rate: ", audioSampleRate)
		log.Println("Channels: ", audioChannels)
		log.Println("Samples: ", audioNumSamples)
		time.Sleep(1 * time.Second)
	}
}
