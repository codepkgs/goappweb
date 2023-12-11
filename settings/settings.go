package settings

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	*AppConfig   `mapstructure:"app"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database int    `mapstructure:"database"`
	Password string `mapstructure:"password"`
}

var Conf Config

var ErrInvalidConfigFilename = fmt.Errorf("invalid config filename, config filename must be in format of <name>.<type>")
var ErrUnsupportedConfigType = fmt.Errorf("unsupported config type, supported config types are: toml, json, yaml, hcl, ini")

func InitConfig(filename string) (err error) {
	parts := strings.Split(filename, ".")
	if len(parts) != 2 {
		return ErrInvalidConfigFilename
	}

	configName := parts[0]
	configType := parts[1]

	switch configType {
	case "toml", "json", "yaml", "yml", "hcl", "ini":
	default:
		return ErrUnsupportedConfigType
	}

	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigName(configName)
	v.SetConfigType(configType)
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&Conf); err != nil {
		return err
	}
	return nil
}
