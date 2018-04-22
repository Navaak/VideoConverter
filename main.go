package main

import (
	"navaak/convertor/app"
	"os"
	"strconv"
)

var env = map[string]string{
	"watch_path":  "VC_WATCH",
	"export_path": "VC_EXPORTS",
	"work_path":   "VC_TMP",
	"max_use_cpu": "VC_WORKER",
	"log_path":    "VC_LOGS",
}

func main() {
	config := loadConfig()
	a, _ := app.New(config)
	a.Run()
}

func loadConfig() app.Config {
	var config = app.Config{
		WatchPath:  os.Getenv(env["watch_path"]),
		WorkPath:   os.Getenv(env["work_path"]),
		ExportPath: os.Getenv(env["export_path"]),
		LogPath:    os.Getenv(env["log_path"]),
	}
	cpu, _ := strconv.Atoi(os.Getenv(env["max_use_cpu"]))
	if cpu < 1 {
		cpu = 1
	}
	config.MaxUseCPU = cpu
	return config
}
