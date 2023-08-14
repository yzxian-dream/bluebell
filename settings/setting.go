package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf *AppConfig

type AppConfig struct {
	Name       string `mapstructure:"name"`
	Mode       string `mapstructure:"mode"`
	Port       int    `mapstructure:"port"`
	Version    string `mapstructure:"version"`
	*LogConf   `mapstructure:"log"`
	*MysqlConf `mapstructure:"mysql"`
	*RedisConf `mapstructure:"redis"`
}

type MysqlConf struct {
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	DBname         string `mapstructure:"dbname"`
	Max_Conns      int    `mapstructure:"max_conns"`
	Max_Idle_Conns int    `mapstructure:"max_idle_conns"`
}

type RedisConf struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Pool_Size int    `mapstructure:"pool_size"`
	DB        int    `mapstructure:"db"`
	Password  string `mapstructure:"password"`
}

type LogConf struct {
	Level      string `mapstructure:",level"`
	FileName   string `mapstructure:",filename"`
	MaxSize    int    `mapstructure:",max_size"`
	MaxAge     int    `mapstructure:",max_age"`
	MaxBackups int    `mapstructure:",max_backups"`
}

func Init() error {
	viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	//viper.SetConfigName("config") // 配置文件名称(无扩展名)
	//viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".") // 查找配置文件所在的路径
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("fatal error config file: %s \n", err)
		return err
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		fmt.Printf("Unmarshal error config file: %s \n", err)
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:#{err}\n")
		}
	})
	return nil
}
