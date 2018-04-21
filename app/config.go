package app

var DefaultConfig = Config{
	WatchPath: "./watch",
	DonePath:  "./done",
	MaxUseCPU: 3,
	LogPath:   "./logs",
}

type Config struct {
	WatchPath string `json:"watch_path"`
	DonePath  string `json:"done_path"`
	MaxUseCPU int    `json:"max_use_cpu"`
	LogPath   string `json:"log_path"`
}
