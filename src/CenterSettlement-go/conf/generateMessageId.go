package conf

import (
	"github.com/gopkg.in/ini.v1"
	"log"
	"strconv"
)

type MessageId struct { //配置文件要通过tag来指定配置文件中的名称
	Messageid int64 `ini:"messageid"`
}

//注意读取messageid时，要做加锁处理
func GenerateMessageId() int64 {

	cfg, err := ini.Load("../conf/app.conf") //读配置文件
	if err != nil {
		log.Fatal("Fail to read file:", err)
	}

	id, rerr := cfg.Section("").Key("messageid").Int64() //取值messageid
	if rerr != nil {
		log.Fatal("Fail to read messageid:", rerr)
	}
	log.Println("read conffile messageid: ", id)

	newid := id + 1
	s := strconv.Itoa(int(newid))

	cfg.Section("").Key("messageid").SetValue(s) //  修改后值然后进行保存
	Saveerr := cfg.SaveTo("../conf/app.conf")
	if Saveerr != nil {
		log.Fatal("Fail to SaveTo file:", Saveerr)
	}

	log.Println(" file new messageid: ", s)
	log.Println(" xmlfile new messageid: ", id)
	return id
}