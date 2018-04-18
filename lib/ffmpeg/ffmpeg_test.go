package ffmpeg

import (
	"testing"

	"github.com/gosuri/uiprogress"
)

func TestRun(t *testing.T) {
	v, err := NewVideo("/home/ehsan/Downloads/Alex.mp4", "", P480, P360, P240)
	if err != nil {
		t.Error(err)
	}
	v.SetWorkerCount(2)
	v.Run()
	p := v.Progress()
	uiprogress.Start()
	bar := uiprogress.AddBar(100)
	bar.AppendCompleted()
	bar.PrependElapsed()
	last := 0
	for ch := range p {
		diff := int(ch) - last
		for i := 0; i < diff; i++ {
			bar.Incr()
		}
		last = int(ch)
	}
}
