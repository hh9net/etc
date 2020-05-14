package sysinit

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/gopkg.in/ini.v1"
	"log"
)

var filepath = "../conf/app.conf"

type Config struct { //配置文件要通过tag来指定配置文件中的名称
	MHostname     string `ini:"mysql_hostname"`
	MPort         string `ini:"mysql_port"`
	MUserName     string `ini:"mysql_user"`
	MPass         string `ini:"mysql_pass"`
	Mdatabasename string `ini:"mysql_databasename"`
	MKeepalive    int    `ini:"mysql_keepalive"`
	MTimeout      int    `ini:"mysql_timeout"`
}

//连接数据库
func NewEngine() (*xorm.Engine, error) {
	//读配置文件
	config, err := ReadConfig(filepath) //也可以通过os.arg或flag从命令行指定配置文件路径
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	params := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", config.MUserName, config.MPass, config.MHostname, config.Mdatabasename)
	x, err := xorm.NewEngine("mysql", params)
	if x == nil {
		log.Println("获取xorm为空")
		x = new(xorm.Engine)
	}
	if err != nil {
		log.Fatal("连接数据库error")
	}
	log.Println("连接数据库成功")
	return x, err
}

//读取配置文件并转成结构体
func ReadConfig(path string) (Config, error) {
	var config Config
	conf, err := ini.Load(path) //加载配置文件
	if err != nil {
		log.Println("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config) //解析成结构体
	if err != nil {
		log.Println("mapto config file fail!")
		return config, err
	}
	return config, nil
}
