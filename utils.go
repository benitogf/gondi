package gondi

import "regexp"

var rgx = regexp.MustCompile(`\((.*?)\)`)

func ExtractSourceName(fullName string) string {
	rs := rgx.FindStringSubmatch(fullName)
	if len(rs) < 2 {
		return ""
	}
	return rs[1]
}

// send an alpha frame
func SendAlphaFrame(sender *SendInstance) {
	// send alpha on stop
	img := GenerateAlpha()
	videoFrame := NewVideoFrameV2()
	videoFrame.FourCC = FourCCTypeRGBA
	videoFrame.FrameFormatType = FrameFormatProgressive
	videoFrame.Xres = int32(EMPTY_X)
	videoFrame.Yres = int32(EMPTY_Y)
	videoFrame.LineStride = 0 // 2 bytes per pixel
	videoFrame.FrameRateN = 30000
	videoFrame.FrameRateD = 1001
	videoFrame.Data = &img.Pix[0]
	sender.SendVideoFrameAsync(videoFrame)
}
