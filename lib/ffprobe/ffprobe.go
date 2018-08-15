package ffprobe

import (
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"strings"
	"time"
)

type Resolution struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Format struct {
	Filename string `json:"filename"`
	Size     string `json:"size"`
	Duration string `json:"duration"`
	BitRate  string `json:"bit_rate"`
}

type FileDetail struct {
	Streams    []Resolution `json:"streams"`
	Format     Format       `json:"format"`
	Resolution Resolution   `json:"-"`
}

func GetDetail(path string) (*FileDetail, error) {
	return getDetail(path, 1)
}

func getDetail(path string, try int) (*FileDetail, error) {
	cmd := exec.Command("ffprobe", "-show_format", "-print_format", "json",
		"-show_entries", "stream=width,height", path)
	out, err := cmd.Output()
	if err != nil && try < 3 {
		time.Sleep(time.Second)
		try++
		return getDetail(path, try)
	}
	command := strings.Join(cmd.Args, " ")
	if err != nil && try > 3 {
		err = errors.New(err.Error() + " on execute : " + command)
		return nil, err
	}
	f := new(FileDetail)
	if err := json.Unmarshal(out, f); err != nil {
		err = errors.New(err.Error() + "on get marshal out put : " + command)
		return nil, err
	}
	if err := f.parse(); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *FileDetail) parse() error {
	if len(f.Streams) < 1 {
		return errors.New("ffprobe: parsing streams error")
	}
	f.Resolution = f.Streams[0]
	return nil
}

func testRun() error {
	cmd := exec.Command("ffprobe", "-version")
	return cmd.Run()
}

func init() {
	if err := testRun(); err != nil {
		log.Fatal("ffprobe is missing could not execute")
	}
}
