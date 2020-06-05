package conf

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

///Users/nicker/go/etc/src/CenterSettlement-go/conf
var logconffilepath = "CenterSettlement-go/conf/app.conf"

//var logconffilepath = "CenterSettlement-go/conf/app.conf"

type LogConfig struct { //配置文件要通过tag来指定配置文件中的名称
	LogmaxAge       int    `ini:"log_maxAge"`
	LogrotationTime int    `ini:"log_rotationTime"`
	LogPath         string `ini:"log_Path"`
	LogFileName     string `ini:"log_FileName"`
}

//读取配置文件并转成结构体
func ReadlogConfig(path string) (LogConfig, error) {
	var logconfig LogConfig
	conf, err := ini.Load(path) //加载配置文件
	if err != nil {
		log.Println("log conf load config file fail!")
		return logconfig, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&logconfig) //解析成结构体
	if err != nil {
		log.Println("log conf mapto config file fail!")
		return logconfig, err
	}
	return logconfig, nil
}

func LogConfigInit() *LogConfig {
	//读配置文件
	config, err := ReadlogConfig(logconffilepath) //也可以通过os.arg或flag从命令行指定配置文件路径
	if err != nil {
		log.Fatal(err)
	}
	log.Println(config)
	return &config
}
