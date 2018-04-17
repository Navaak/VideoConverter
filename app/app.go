package app

import (
	"log"
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
		jobsCount := 0
		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Create {
					a.newFile(event.Name)
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

func (a *application) newFile(f string) chan bool {
	done := make(chan bool)

	return done
}
