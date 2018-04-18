package ffmpeg

import "navaak/convertor/lib/ffprobe"

type ExportLog struct {
	DestFile   string
	Resolution ffprobe.Resolution
	Success    bool
	ScaleTitle string
	Error      error
}

type Log struct {
	SourceFile       string
	SourceResolution ffprobe.Resolution
	Exports          []ExportLog
	Size             string
}
