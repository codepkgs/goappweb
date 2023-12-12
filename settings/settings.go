package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

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
	Mode string `mapstructure:"mode"`
	Port int    `mapstructure:"port"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Output     string `mapstructure:"output"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxOpenConns string `mapstructure:"max_open_conns"`
	MaxIdleConns string `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database int    `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
}

var Conf Config

func Init(filename string) (err error) {
	v := viper.New()
	v.SetConfigFile(filename)

	if err = v.ReadInConfig(); err != nil {
		return err
	}

	// 将配置信息反序列化到Conf中
	if err = v.Unmarshal(&Conf); err != nil {
		return
	}

	// 监控配置文件的变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&Conf); err != nil {
			fmt.Printf("config file changed, but  failed to unmarshal: %v\n", err)
			return
		}
	})

	return
}
