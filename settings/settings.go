package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		// 读取配置文件失败
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了。。。")
	})
	return

}
