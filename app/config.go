package app

var DefaultConfig = Config{
	WatchPath: "./watch",
	WorkPath:  "./tmp",
	DonePath:  "./done",
	MaxUseCPU: 20,
	MaxUseRam: 20,
}

type Config struct {
	WatchPath string `json:"watch_path`
	WorkPath  string `json:"work_path"`
	DonePath  string `json:"done_path"`
	MaxUseCPU int    `json:"max_use_cpu`
	MaxUseRam int    `json:"max_use_ram`
}
