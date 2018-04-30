package main

import (
	"encoding/json"
	"fmt"
	"log"
	"navaak/convertor/app"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
)

var env = map[string]string{
	"watch_path":     "VC_WATCH",
	"export_path":    "VC_EXPORTS",
	"work_path":      "VC_TMP",
	"max_use_cpu":    "VC_WORKER",
	"log_path":       "VC_LOGS",
	"snapshots_path": "VC_SNAPSHOTS_BASE_PATH",
	"webhook_url":    "VC_WEBHOOK_URL",
	"webhook_token":  "VC_WEBHOOK_TOKEN",
}

func main() {
	config := loadConfig()
	cdata, _ := json.Marshal(&config)
	println("service stating with :")
	fmt.Println(string(cdata))
	println()
	time.Sleep(time.Second)
	a, _ := app.New(config)
	a.Run()
}

func loadConfig() app.Config {
	var config = app.Config{
		WatchPath:     os.Getenv(env["watch_path"]),
		WorkPath:      os.Getenv(env["work_path"]),
		ExportPath:    os.Getenv(env["export_path"]),
		LogPath:       os.Getenv(env["log_path"]),
		SnapshotsPath: os.Getenv(env["snapshots_path"]),
		WebhookURL:    os.Getenv(env["webhook_url"]),
		WebhookToken:  os.Getenv(env["webhook_token"]),
	}
	cpu, _ := strconv.Atoi(os.Getenv(env["max_use_cpu"]))
	if cpu < 1 {
		cpu = 1
	}
	if _, err := govalidator.ValidateStruct(&config); err != nil {
		log.Fatal(err)
	}
	config.MaxUseCPU = cpu
	return config
}
