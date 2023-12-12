package settings

import "flag"

var (
	configPath string
)

func parseConfigFile() {
	flag.StringVar(&configPath, "conf", "./config.toml", "config filename")
}

func Parse() string {
	parseConfigFile()
	flag.Parse()
	return configPath
}
