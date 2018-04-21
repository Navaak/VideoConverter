package app

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"

	"navaak/convertor/lib/ffmpeg"
	"navaak/convertor/lib/logger"
)

type application struct {
	config Config
	logger *logger.Logger
}

func New(config Config) (*application, error) {
	a := new(application)
	a.config = config
	runtime.GOMAXPROCS(a.config.MaxUseCPU)
	a.logger = logger.New(a.config.LogPath)
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
	watchpath, err := filepath.Abs(a.config.WatchPath)
	if err != nil {
		log.Fatal(err)
	}
	println("watching path selected to ---> ", watchpath)
	watcher.Add(watchpath)
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
		return
	}
	v.SetWorkerCount(a.config.MaxUseCPU)
	v.Run()
	loggs := v.Logger()
	a.logger.Log(filepath.Base(f), loggs)
}
