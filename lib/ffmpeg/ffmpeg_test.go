package ffmpeg

import "testing"

func TestRun(t *testing.T) {
	if err := Run("/home/ehsan/Downloads/Alex.mp4", ""); err != nil {
		t.Error(err)
	}
}
