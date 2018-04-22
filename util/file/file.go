package file

import (
	"io"
	"os"
)

func Copy(src, dest string) error {
	srcIO, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcIO.Close()
	destIO, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destIO.Close()
	if _, err := io.Copy(destIO, srcIO); err != nil {
		return err
	}
	destIO.Sync()
	return nil
}

func Move(src, dest string) error {
	if err := Copy(src, dest); err != nil {
		return err
	}
	return os.Remove(src)
}
