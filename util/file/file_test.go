package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	baseFile = "/home/ehsan/Downloads/Alex.mp4"
)

func TestCopy(t *testing.T) {
	dest := filepath.Base(baseFile)
	if err := Copy(baseFile, dest); err != nil {
		t.Error(err)
	}
}

func TestMove(t *testing.T) {
	base := filepath.Base(baseFile)
	name := strings.Split(base, ".")[0]
	os.MkdirAll(name, 0777)
	destpath := filepath.Join(name, base)
	if err := Move(base, destpath); err != nil {
		t.Error(err)
	}
}
