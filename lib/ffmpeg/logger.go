package ffmpeg

import (
	"encoding/json"
	"navaak/convertor/lib/ffprobe"
)

type ExportLog struct {
	DestFile   string             `json:"dest_file"`
	Resolution ffprobe.Resolution `json:"resolution"`
	Success    bool               `json:"success"`
	ScaleTitle string             `json:"scale_title"`
	Error      error              `json:"error"`
}

type Log struct {
	SourceFile       string             `json:"source_file"`
	SourceResolution ffprobe.Resolution `json:"source_resolution"`
	Exports          []ExportLog        `json:"exports"`
	Size             int                `json:"size"`
	Duration         int                `json:"duration"`
}

func (l *Log) JSON() []byte {
	d, _ := json.Marshal(l)
	return d
}
