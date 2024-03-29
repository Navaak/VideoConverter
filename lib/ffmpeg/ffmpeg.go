package ffmpeg

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"

	"navaak/convertor/lib/ffprobe"
)

// Video caching reference for convertor jobs
type Video struct {
	src            string
	destDir        string
	scales         []string
	details        *ffprobe.FileDetail
	worker         int
	exports        []*export
	done           chan bool
	sourceDuration time.Duration
	errs           []error
}

type export struct {
	dest           string
	resolution     ffprobe.Resolution
	err            error
	progress       float32
	sourceDuration time.Duration
	scale          string
	done           bool
}

// New Video instant
func NewVideo(src, destDir string, scales ...string) (*Video, error) {
	details, err := ffprobe.GetDetail(src)
	if err != nil {
		return nil, err
	}
	v := new(Video)
	v.src = src
	v.destDir = destDir
	v.details = details
	v.sourceDuration = strDurationToTime(v.details.Format.Duration)
	for _, scale := range scales {
		if err := v.newExp(scale); err != nil {
			return nil, err
		}
	}
	return v, nil
}

// Run convertor scaling video
func (v *Video) Run() {
	for _, export := range v.exports {
		v.exec(export)
	}
}

func (v *Video) JobsCount() int {
	return len(v.exports)
}

func (v *Video) SetWorkerCount(n int) {
	v.worker = n
}
func (v *Video) Logger() Log {
	size, _ := strconv.Atoi(v.details.Format.Size)
	log := Log{
		SourceFile:       v.src,
		SourceResolution: v.details.Resolution,
		Size:             size,
		Duration:         int(v.sourceDuration.Seconds()),
		Errors:           v.errs,
	}
	for _, e := range v.exports {
		log.Exports = append(log.Exports, ExportLog{
			DestFile:   e.dest,
			Resolution: e.resolution,
			ScaleTitle: e.scale,
			Success:    e.done,
			Error:      e.err,
		})
	}
	return log
}

func (v *Video) Snapshots(path string) {
	println("snapshoting on ", path)
	cmd := exec.Command("ffmpeg", "-i",
		v.src, "-f", "image2", "-bt", "20M",
		"-vf", "fps=1/20", filepath.Join(path,
			"shot%02d.jpg"))
	if err := cmd.Run(); err != nil {
		errMsg := err.Error() + " on running command : " + cmd.Args[0]
		v.errs = append(v.errs, errors.New("snapshot: "+errMsg))
	}
}

func (v *Video) newExp(scale string) error {
	resolution, ok := scales[scale]
	if !ok {
		return errors.New("ffmpeg: " + scale + " is undefined")
	}
	dest, err := v.makeFilepath(scale)
	if err != nil {
		return err
	}
	if (resolution.Height > v.details.Resolution.Height || resolution.Width > v.details.Resolution.Width) &&
		scale != P240 {
		return nil
	}
	e := new(export)
	e.dest = dest
	e.resolution = resolution
	e.scale = scale
	v.exports = append(v.exports, e)
	e.sourceDuration = v.sourceDuration
	return nil
}

func (v *Video) exec(e *export) {
	scale := fmt.Sprintf("scale=%d:%d",
		e.resolution.Width, e.resolution.Height)
	cmd := exec.Command("ffmpeg", "-y",
		"-i", v.src,
		"-vf", scale,
		"-threads", "7",
		"-codec:v", "libx264",
		"-preset", "slow",
		"-b:v", scalesBV[e.scale],
		"-b:a", scalesBA[e.scale],
		"-maxrate", scalesBuffRates[e.scale],
		"-bufsize", scalesBuffRates[e.scale],
		"-profile:v", scalesProfiles[e.scale],
		e.dest)
	command := strings.Join(cmd.Args, " ")
	println("command :  ", command, "   has executed successfully!")
	if err := cmd.Run(); err != nil {
		e.err = errors.New(err.Error() + " on running command : " + command)
		logrus.Error(e.err.Error())
	}
}

func (v *Video) makeFilepath(scale string) (string, error) {
	base := filepath.Base(v.src)
	splits := strings.Split(base, ".")
	if len(splits) < 2 {
		return "", errors.New("error source file path")
	}
	name := splits[0]
	filename := name + "." + scalesPreExt[scale] + ext
	path := filepath.Join(v.destDir, filename)
	return path, nil
}

func (e *export) readout(r io.Reader) {
	buf := make([]byte, 1024, 1024)
	counter := 0
	for {
		n, err := r.Read(buf[:])
		counter++
		if counter < 50 {
			continue
		}
		if n > 0 {
			d := buf[:n]
			current := parseDurationFromReader(string(d))
			e.progress = getProgress(current, e.sourceDuration)
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
	}
}

func getProgress(current, total time.Duration) float32 {
	p := float32(current.Seconds()/total.Seconds()) * 100
	if p >= 100 {
		return 100
	}
	return p
}

func parseDurationFromReader(s string) time.Duration {
	re := regexp.MustCompile("time=([0-9]+):([0-9]+):([0-9]+)")
	submatches := re.FindAllStringSubmatch(s, -1)
	if len(submatches) < 1 {
		return time.Minute * 15
	}
	if len(submatches[0]) < 4 {
		return time.Minute * 15
	}
	hour, _ := strconv.Atoi(submatches[0][1])
	min, _ := strconv.Atoi(submatches[0][2])
	sec, _ := strconv.Atoi(submatches[0][3])
	return time.Duration(int(time.Hour)*hour) +
		time.Duration(int(time.Minute)*min) +
		time.Duration(int(time.Second)*sec)
}

func strDurationToTime(s string) time.Duration {
	n, _ := strconv.ParseFloat(s, 32)
	return time.Duration(int(time.Second) * int(n))
}

func testRun() error {
	cmd := exec.Command("ffmpeg", "-version")
	return cmd.Run()
}

func init() {
	if err := testRun(); err != nil {
		log.Fatal("ffmpeg is missing could not execute")
	}
}
