package app

var DefaultConfig = Config{
	WatchPath:  "./watch",
	ExportPath: "./done",
	WorkPath:   "./work",
	MaxUseCPU:  3,
	LogPath:    "./logs",
}

type Config struct {
	WatchPath  string `json:"watch_path"`
	WorkPath   string `json:"work_path"`
	ExportPath string `json:"export_path"`
	MaxUseCPU  int    `json:"max_use_cpu"`
	LogPath    string `json:"log_path"`
}
