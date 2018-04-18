package app

import (
	"log"
	"navaak/convertor/lib/ffmpeg"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
)

type application struct {
	config Config
}

func New(config Config) (*application, error) {
	a := new(application)
	a.config = config
	runtime.GOMAXPROCS(a.config.MaxUseCPU)
	return a, nil
}

func (a *application) Run() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Create {
					log.Println("new file detected -- >",
						event.Name)
					a.newVid(event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	watcher.Add(a.config.WatchPath)
	if err != nil {
		return err
	}
	<-done
	return nil
}

func (a *application) newVid(f string) {
	if filepath.Ext(f) != ".mp4" {
		return
	}
	v, err := ffmpeg.NewVideo(f, a.config.DonePath,
		ffmpeg.P1080,
		ffmpeg.P720,
		ffmpeg.P360,
		ffmpeg.P480,
		ffmpeg.P240)
	if err != nil {
		log.Fatal(err)
	}
	v.SetWorkerCount(a.config.MaxUseCPU)
	v.Run()
	v.Wait()
}
