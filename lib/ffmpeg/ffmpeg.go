package ffmpeg

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"navaak/convertor/lib/ffprobe"
)

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
		p720:  {1280, 720},
		p1080: {1920, 1080},
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
	scalesDescSort = []string{P1080, P720, P480, P360, P240}
)

type Video struct {
	src            string
	destDir        string
	scales         []string
	details        FileDetail
	worker         int
	exports        []Export
	sourceDuration time.Duration
}

type Export struct {
	dest           string
	resolution     ffprobe.Resolution
	err            error
	progress       float32
	sourceDuration time.Duration
}

func NewVideo(src, destDir string, scales ...string) (Video, error) {
	details, err := ffprobe.GetDetail(src)
	if err != nil {
		return nil, err
	}
	v := new(Video)
	v.src = src
	v.destDir = destDir
	v.sourceDuration = strDurationToTime(v.details.Format.Duration)
	for _, scale := range scales {
		if _, err := v.newExp(scale); err != nil {
			return err
		}
	}
	return v, nil
}

func (v *Video) Progress() chan float32 {
	p := make(chan float32)
	go v.calculateProgress(p)
	return p
}

func (v *Video) calculateProgress(p chan float32) {
	for {
		time.Sleep(time.Second)
	}
}

func (v *Video) JobsCount() int {
	return len(v.exports)
}

func (v *Video) SetWorkerCount(n int) {
	v.worker = n
}

func (v *Video) Run() {
	if v.worker < 1 {
		v.SetWorkerCount(2)
	}
	var (
		job       sync.WaitGroup
		jobsCount int
	)
	for _, export := range v.exports {
		job.Add(1)
		jobsCount++
		go v.exec(export, job)
		if jobsCount >= v.worker {
			job.Wait()
		}
	}
}

func (v *Video) newExp(scale string) (*Export, error) {
	resolution, ok := scales[scale]
	if !ok {
		return nil, errors.New("ffmpeg: " + scale + " is undefined")
	}
	dest, err := v.makeFilepath(scale)
	if err != nil {
		return nil, err
	}
	e := new(Export)
	e.dest = dest
	e.resolution = resolution
	v.exports = append(v.exports, *e)
	e.sourceDuration = v.sourceDuration
	return e, nil
}

func (v *Video) exec(e Export, job sync.WaitGroup) error {
	defer job.Done()
	scale := fmt.Sprintf("scale=%d:%d",
		e.resolution.Width, e.resolution.Height)
	cmd := exec.Command("ffmpeg", "-y", "-i", v.src, "-vf", scale,
		"-codec:v", "libx264", e.dest)
	stdout, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	cmd.Start()
	e.progress = make(chan float)
	go func() {
		e.readout(stdout)
	}()

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (v *Video) makeFilepath(scale string) (string, error) {
	base := filepath.Base(v.src)
	ex := filepath.Ext(v.src)
	splits := strings.Split(base, ex, -1)
	if len(splits) < 2 {
		return "", errors.New("error source file path")
	}
	name := splits[0]
	filename = name + scalesPreExt[scale] + ext
	return filepath.Join(v.destDir, filepath), nil
}

func (e *Export) readout(r io.Reader) {
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
	return 20
}

func parseDurationFromReader(s string) time.Duration {
	re := regexp.MustCompile("time=([0-9]+):([0-9]+):([0-9]+)")
	submatches := re.FindAllStringSubmatch(s, -1)
	if len(submatches) < 4 {
		return time.Minute * 15
	}
	hour, _ := strconv.Atoi(submatches[1])
	min, _ := strconv.Atoi(submatches[2])
	sec, _ := strconv.Atoi(submatches[3])
	return time.Duration(int(time.Hour)*hour) +
		time.Duration(int(time.Minute)*hour) +
		time.Duration(int(time.Second)*hour)

}

func strDurationToTime(s string) time.Duration {
	n, _ := strconv.Atoi(s)
	return time.Duration(int(time.Second) * n)
}
