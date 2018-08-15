package ffmpeg

import "navaak/convertor/lib/ffprobe"

const (
	P240  = "240p"
	P360  = "360p"
	P480  = "480p"
	P720  = "720p"
	P1080 = "1080p"
	ext   = ".mp4"
)

var (
	scales = map[string]ffprobe.Resolution{
		P240:  {426, 240},
		P360:  {640, 360},
		P480:  {854, 480},
		P720:  {1280, 720},
		P1080: {1920, 1080},
	}
	scalesPreExt = map[string]string{
		P240:  "240",
		P360:  "360",
		P480:  "480",
		P720:  "720",
		P1080: "1080",
	}
	scalesByHeight = map[int]string{
		240:  P240,
		360:  P360,
		480:  P480,
		720:  P720,
		1080: P1080,
	}
	scalesProfiles = map[string]string{
		P240:  "baseline",
		P360:  "baseline",
		P480:  "main",
		P720:  "high",
		P1080: "high",
	}
	scalesBuffRates = map[string]string{
		P240:  "600k",
		P360:  "800k",
		P480:  "1000k",
		P720:  "3000k",
		P1080: "5000k",
	}
	scalesBV = map[string]string{
		P240:  "300k",
		P360:  "400k",
		P480:  "500k",
		P720:  "1500k",
		P1080: "3000k",
	}
	scalesBA = map[string]string{
		P240:  "96k",
		P360:  "96k",
		P480:  "128k",
		P720:  "196k",
		P1080: "196k",
	}
	scalesDescSort = []string{P1080, P720, P480, P360, P240}
)
