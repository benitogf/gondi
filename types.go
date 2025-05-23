package gondi

import "math"

type VideoFrameV2 struct {
	// The resolution of this frame.
	Xres, Yres int32

	// What FourCC describing the type of data for this frame.
	FourCC FourCCType

	// What is the frame rate of this frame.
	// For instance NTSC is 30000,1001 = 30000/1001 = 29.97 fps.
	FrameRateN, FrameRateD int32

	// What is the picture aspect ratio of this frame.
	// For instance 16.0/9.0 = 1.778 is 16:9 video
	// 0 means square pixels.
	PictureAspectRatio float32

	// Is this a fielded frame, or is it progressive.
	FrameFormatType FrameFormat

	// The timecode of this frame in 100ns intervals.
	Timecode int64

	// The video data itself.
	Data *byte

	// The inter line stride of the video data, in bytes.
	LineStride int32

	//Per frame metadata for this frame. This is a NULL terminated UTF8 string that should be
	//in XML format. If you do not want any metadata then you may specify NULL here.
	Metadata *byte

	//This is only valid when receiving a frame and is specified as a 100ns time that was the exact
	//moment that the frame was submitted by the sending side and is generated by the SDK. If this
	//value is NDIlib_recv_timestamp_undefined then this value is not available and is NDIlib_recv_timestamp_undefined.
	Timestamp int64
}

// Video frame format
type FrameFormat int32

const (
	FrameFormatInterleaved FrameFormat = iota //A fielded frame with the field 0 being on the even lines and field 1 being on the odd lines.
	FrameFormatProgressive                    //A progressive frame.

	//Individual fields.
	FrameFormatField0
	FrameFormatField1
)

// An enumeration to specify the type of a packet returned by the capture functions
type FrameType int32

const (
	FrameTypeNone FrameType = iota
	FrameTypeVideo
	FrameTypeAudio
	FrameTypeMetadata
	FrameTypeError

	//This indicates that the settings on this input have changed.
	//For instamce, this value will be returned from NDIlib_recv_capture_v2 and NDIlib_recv_capture
	//when the device is known to have new settings, for instance the web-url has changed ot the device
	//is now known to be a PTZ camera.
	FrameTypeStatusChange FrameType = 100
)

type FourCCType [4]byte

var (
	FourCCTypeRGBA = [4]byte{'R', 'G', 'B', 'A'}
	FourCCTypeRGBX = [4]byte{'R', 'G', 'B', 'X'}
	// YCbCr color space using 4:2:2.
	FourCCTypeUYVY FourCCType = [4]byte{'U', 'Y', 'V', 'Y'}
	// Planar 8bit, 4:4:4:4 video format.
	// Color ordering in memory is blue, green, red, alpha
	FourCCTypeBGRA FourCCType = [4]byte{'B', 'G', 'R', 'A'}
	// Planar 8bit, 4:4:4 video format, packed into 32bit pixels.
	// Color ordering in memory is blue, green, red, 255
	FourCCTypeBGRX FourCCType = [4]byte{'B', 'G', 'R', 'X'}
	// YCbCr + Alpha color space, using 4:2:2:4.
	// In memory there are two separate planes. The first is a regular
	// UYVY 4:2:2 buffer. Immediately following this in memory is a
	// alpha channel buffer.
	FourCCTypeUYVA FourCCType = [4]byte{'U', 'Y', 'V', 'A'}
)

const (
	SendTimecodeSynthesize int64 = math.MaxInt64
	SendTimecodeEmpty      int64 = 0
)

type AudioFrameV2 struct {
	//The sample-rate of this buffer.
	SampleRate int32

	//The number of audio channels.
	NumChannels int32

	//The number of audio samples per channel.
	NumSamples int32

	// The timecode of this frame in 100-nanosecond intervals.
	Timecode int64

	// The audio data as float32 samples
	Data *float32

	// The inter channel stride of the audio channels, in bytes.
	ChannelStride int32

	// Per frame metadata for this frame. This is a NULL terminated UTF8 string that should be in XML format.
	// If you do not want any metadata then you may specify NULL here.
	Metadata *byte

	// This is only valid when receiving a frame and is specified as a 100-nanosecond time that was the exact
	// moment that the frame was submitted by the sending side and is generated by the SDK. If this value is
	// NDIlib_recv_timestamp_undefined then this value is not available and is NDIlib_recv_timestamp_undefined.
	Timestamp int64
}

type AudioFrameV3 struct {
	//The sample-rate of this buffer.
	SampleRate int32

	//The number of audio channels.
	NumChannels int32

	//The number of audio samples per channel.
	NumSamples int32

	// The timecode of this frame in 100-nanosecond intervals.
	Timecode int64

	// The audio data as int16 samples
	Data *byte

	// The inter channel stride of the audio channels, in bytes.
	ChannelStride int32

	// Per frame metadata for this frame. This is a NULL terminated UTF8 string that should be in XML format.
	// If you do not want any metadata then you may specify NULL here.
	Metadata *byte

	// This is only valid when receiving a frame and is specified as a 100-nanosecond time that was the exact
	// moment that the frame was submitted by the sending side and is generated by the SDK. If this value is
	// NDIlib_recv_timestamp_undefined then this value is not available and is NDIlib_recv_timestamp_undefined.
	Timestamp int64
}

/* Borrowed from ndi-go/ndi.go */
type RecvColorFormat int32

const (
	RecvColorFormatBGRXBGRA RecvColorFormat = 0 //No alpha channel: BGRX, Alpha channel: BGRA
	RecvColorFormatUYVYBGRA RecvColorFormat = 1 //No alpha channel: UYVY, Alpha channel: BGRA
	RecvColorFormatRGBXRGBA RecvColorFormat = 2 //No alpha channel: RGBX, Alpha channel: RGBA
	RecvColorFormatUYVYRGBA RecvColorFormat = 3 //No alpha channel: UYVY, Alpha channel: RGBA

	//Read the SDK documentation to understand the pros and cons of this format.
	RecvColorFormatFastest RecvColorFormat = 100
)

/* Borrowed from ndi-go/ndi.go */
type RecvBandwidth int32

const (
	RecvBandwidthMetadataOnly RecvBandwidth = -10 //Receive metadata.
	RecvBandwidthAudioOnly    RecvBandwidth = 10  //Receive metadata, audio.
	RecvBandwidthLowest       RecvBandwidth = 0   //Receive metadata, audio, video at a lower bandwidth and resolution.
	RecvBandwidthHighest      RecvBandwidth = 100 //Receive metadata, audio, video at full resolution.
)

type RecvPerformance struct {
	//The current number of video frames
	VideoFrames int64

	//The current number of audio frames
	AudioFrames int64

	//The current number of metadata frames
	MetadataFrames int64
}

type MetadataFrame struct {
	// The length of the string in UTF8 characters. This includes the NULL terminating character.
	// If this is 0, then the length is assume to be the length of a null terminated string.
	Length int32

	// The timecode of this frame in 100ns intervals.
	Timecode int64

	// The metadata as a UTF8 XML string. This is a NULL terminated string.
	Data *byte
}

type Tally struct {
	// Program tally, usually indicated by red
	Program bool
	// Preview tally, usually indicated by green
	Preview bool
}

type Source struct {
	name    *byte
	address *byte
}

type routingCreateSettings struct {
	// The name of the ndi source to create.
	name *byte

	// The groups that you would like the source to be in
	groups *byte
}

type NewRecvInstanceSettings struct {
	// Source to connect to
	SourceToConnectTo *Source

	// Your preferred colorspace, default is gondi.RecvColorFormatUYVYBGRA
	ColorFormat RecvColorFormat

	// The bandwidth setting that you wish to use for this video source. Bandwidth
	// controlled by changing both the compression level and the resolution of the source.
	// Default value, and for full quality and all frame types: gondi.RecvBandwidthHighest
	Bandwidth RecvBandwidth

	// When this flag is FALSE, all video that you receive will be progressive. For sources
	// that provide fields, this is de-interlaced on the receiving side (because we cannot change
	// what the up-stream source was actually rendering. This is provided as a convenience to
	// down-stream sources that do not wish to understand fielded video. There is almost no
	// performance impact of using this function.
	// Default is true
	AllowVideoFields bool

	// The name of the ndi receiver to create. This should be named the same way that you
	// would like the source on the network to be named. If this is empty then it will use the filename of your application
	// indexed with the number of the instance number of this receiver.
	Name string
}

type recvCreateSettings struct {
	sourceToConnectTo Source
	colorFormat       RecvColorFormat
	bandwidth         RecvBandwidth
	allowVideoFields  bool
	name              *byte
}

type findCreateSettings struct {
	// Do we want to incluide the list of NDI sources that are running
	// on the local machine ?
	// If TRUE then local sources will be visible, if FALSE then they
	// will not.
	showLocalSources bool

	// Which groups do you want to search in for sources
	groups *byte

	// The list of additional IP addresses that exist that we should query for
	// sources on. For instance, if you want to find the sources on a remote machine
	// that is not on your local sub-net then you can put a comma seperated list of
	// those IP addresses here and those sources will be available locally even though
	// they are not mDNS discoverable. An example might be "12.0.0.8,13.0.12.8".
	// When none is specified the registry is used.
	// Default = nil;
	extraIPs *byte
}

// Sender creation settings
type sendCreateSettings struct {
	// The name of the NDI source to create.
	name *byte

	// What groups should this source be part of. nil or empty string means default.
	groups *byte

	// Do you want audio and video to "clock" themselves. When they are clocked then by adding video frames,
	// they will be rate limited to match the current frame rate that you are submitting at. The same is true
	// for audio. In general if you are submitting video and audio off a single thread then you should only
	// clock one of them (video is probably the better of the two to clock off). If you are submitting audio
	// and video of separate threads then having both clocked can be useful.
	clockVideo, clockAudio bool
}

// Sender instance struct
type SendInstance struct {
	ndiInstance    uintptr
	createSettings *sendCreateSettings
}

// Finder instance struct
type FindInstance struct {
	ndiInstance    uintptr
	createSettings *findCreateSettings
}

// Receiver instance struct
type RecvInstance struct {
	ndiInstance    uintptr
	createSettings *recvCreateSettings
}

// ROuting instance struct
type RoutingInstance struct {
	ndiInstance    uintptr
	createSettings *routingCreateSettings
	name           string
	groups         string
}
