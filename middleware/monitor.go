package middleware

type MonitorConfig struct {
	Port int    `json:"port" yaml:"port"`
	Path string `json:"path" yaml:"path"`
}

func NewMonitor(cfg *MonitorConfig) {

}
