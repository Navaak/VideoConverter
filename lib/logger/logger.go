package logger

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const fileExt = ".log.json"

type Logger struct {
	sync.Mutex
	supplementaryConst map[string]interface{}
	datetimeFormat     string
	loggerPath         string
}

func New(path string) *Logger {
	os.MkdirAll(path, 0777)
	return &Logger{
		loggerPath:     path,
		datetimeFormat: time.RFC3339Nano,
	}
}

func (l *Logger) SetDateTimeFormat(f string) {
	l.Lock()
	defer l.Unlock()
	l.datetimeFormat = f
}

func (l *Logger) SetSupplementry(s map[string]interface{}) {
	l.Lock()
	defer l.Unlock()
	l.supplementaryConst = s
}

func (l *Logger) Log(title string, s interface{}) {
	res := map[string]interface{}{}
	res["data"] = s
	for key, val := range l.supplementaryConst {
		res[key] = val
	}
	res["datetime"] = time.Now().Format(l.datetimeFormat)
	path := filepath.Join(l.loggerPath, title+fileExt)
	data, err := json.Marshal(&res)
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(path, data, 0777); err != nil {
		log.Fatal(err)
	}
}
