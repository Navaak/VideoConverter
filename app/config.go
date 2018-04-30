package app

var DefaultConfig = Config{
	WatchPath:     "./watch",
	ExportPath:    "./done",
	WorkPath:      "./work",
	MaxUseCPU:     3,
	LogPath:       "./logs",
	SnapshotsPath: "./snapshots",
}

type Config struct {
	WatchPath     string `json:"watch_path" valid:"required"`
	WorkPath      string `json:"work_path"`
	WebhookURL    string `json:"webhook_url" valid:"required"`
	WebhookToken  string `json:"webhook_token" valid:"required"`
	ExportPath    string `json:"export_path"`
	MaxUseCPU     int    `json:"max_use_cpu"`
	LogPath       string `json:"log_path"`
	SnapshotsPath string `json:"snapshots_path`
}
