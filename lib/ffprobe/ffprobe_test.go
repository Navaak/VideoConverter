package ffprobe

import "testing"

func TestGetDetail(t *testing.T) {
	_, err := GetDetail("/home/ehsan/Downloads/Alex.mp4")
	if err != nil {
		t.Error(err)
	}
}
