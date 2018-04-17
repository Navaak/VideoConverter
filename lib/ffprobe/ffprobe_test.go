package ffprobe

import (
	"fmt"
	"testing"
)

func TestGetDetail(t *testing.T) {
	detail, err := GetDetail("/home/ehsan/Downloads/Alex.mp4")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(detail)
}
