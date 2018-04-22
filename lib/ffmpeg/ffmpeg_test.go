package ffmpeg

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	v, err := NewVideo("/home/ehsan/Downloads/Alex.mp4", "", P480, P360, P240)
	if err != nil {
		t.Error(err)
	}
	v.SetWorkerCount(2)
	v.Run()
	v.ShowProgressBar()
	fmt.Println(v.Logger())
}

func TestSnapshots(t *testing.T) {
	v, err := NewVideo("/home/ehsan/Downloads/Alex.mp4", "", P480, P360, P240)
	if err != nil {
		t.Error(err)
	}
	v.Snapshots("/home/ehsan/convertor_tmp/snapshots")
}
