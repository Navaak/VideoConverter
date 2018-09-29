package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"navaak/convertor/lib/ffmpeg"
	"navaak/convertor/lib/logger"
	"navaak/convertor/util/file"
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
					syncFile(event.Name)
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
	base := filepath.Base(f)
	if err := os.MkdirAll(a.config.WorkPath, 0777); err != nil {
		println(a.config.WorkPath, err.Error(), "-----------------")
	}
	v, err := ffmpeg.NewVideo(f, a.config.WorkPath,
		ffmpeg.P1080,
		ffmpeg.P720,
		ffmpeg.P480,
		ffmpeg.P360,
		ffmpeg.P240)
	if err != nil {
		go a.logger.Log(base, map[string]string{
			"error": err.Error(),
		})
		return
	}
	v.SetWorkerCount(a.config.MaxUseCPU)
	name := strings.Split(base, ".")[0]
	snapshotsPath := filepath.Join(a.config.SnapshotsPath, name,
		"snapshots")
	os.MkdirAll(snapshotsPath, 0777)
	v.Snapshots(snapshotsPath)
	hookerr := a.hookDone(name + "/snapshots")
	v.Run()
	loggs := v.Logger()
	if hookerr != nil {
		loggs.Errors = append(loggs.Errors, hookerr)
	}
	exportpath := filepath.Join(a.config.ExportPath, name)
	os.MkdirAll(exportpath, 0777)
	orgfile := filepath.Join(exportpath, base)
	if err := file.Move(loggs.SourceFile, orgfile); err != nil {
		loggs.Errors = append(loggs.Errors, err)
	}
	syncFile(exportpath)
	for i, export := range loggs.Exports {
		base := filepath.Base(export.DestFile)
		dest := filepath.Join(exportpath, base)
		loggs.Exports[i].DestFile = dest
		if err := file.Move(export.DestFile, dest); err != nil {
			loggs.Exports[i].Error = err
		}
	}
	destfilename := filepath.Join(exportpath, name)
	a.smil(destfilename, loggs)
	a.json(destfilename, orgfile, loggs)
	hookerr = a.hookDone(name)
	if hookerr != nil {
		loggs.Errors = append(loggs.Errors, hookerr)
	}
	logdest := destfilename + ".log.json"
	go a.logger.LogTo(logdest, loggs)
}

func (a *application) smil(dest string, logg ffmpeg.Log) {
	dest += ".smil"
	res := smilHead
	for _, ex := range logg.Exports {
		base := filepath.Base(ex.DestFile)
		vid := fmt.Sprintf(smilQualities[ex.Resolution.Height], base)
		res += vid
	}
	res += smilFooter
	if err := ioutil.WriteFile(dest, []byte(res), 0777); err != nil {
		logg.Errors = append(logg.Errors, err)
	}
}

func (a *application) json(dest, org string, logg ffmpeg.Log) {
	dest = dest + ".json"
	id := strings.Split(filepath.Base(org), ".")[0]
	qualities := []map[string]interface{}{}
	for _, ex := range logg.Exports {
		data := map[string]interface{}{
			"quality": ex.Resolution.Height,
			"size":    getFileSize(ex.DestFile),
			"bitRate": bitRate[ex.Resolution.Height],
		}
		qualities = append(qualities, data)
	}
	res := map[string]interface{}{
		"videoId":   id,
		"fullpath":  org,
		"duration":  logg.Duration,
		"size":      logg.Size,
		"qualities": qualities,
	}
	data, _ := json.Marshal(&res)
	if err := ioutil.WriteFile(dest, data, 0777); err != nil {
		logg.Errors = append(logg.Errors, err)
		// log.Fatal(err)
	}
}

func getFileSize(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	details, err := f.Stat()
	if err != nil {
		return 0
	}
	return int(details.Size())
}

func syncFile(path string) {
	var cmd *exec.Cmd
	if path != "" {
		cmd = exec.Command("sync", "-d", path)
	} else {
		cmd = exec.Command("sync")
	}
	if err := cmd.Run(); err != nil {
		log.Println("syncing error")
	}
	time.Sleep(time.Second)
}

func (a application) hookDone(path string) error {
	url := a.config.WebhookURL + path
	c := new(http.Client)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(""))
	if err != nil {
		println("request to ", url, "error!", err.Error())
		return err
	}
	req.Header.Set("token", a.config.WebhookToken)
	res, err := c.Do(req)
	if err != nil {
		println("request to ", url, "error! ", err.Error())
		return err
	}
	if res.StatusCode == 200 {
		println("request to ", url, "successfully")
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		println("request to ", url, "error! ", res.StatusCode, err.Error())
		return err
	}
	err = errors.New(string(body))
	println("request to ", url, "error! ", res.StatusCode, err.Error())
	return err
}
