package conf

import (
	"gopkg.in/ini.v1"
	"log"
	"strconv"
	"sync"
)

type MessageId struct { //配置文件要通过tag来指定配置文件中的名称
	Messageid int64 `ini:"messageid"`
}

var m *sync.RWMutex

//注意读取messageid时，要做加锁处理
func GenerateMessageId() int64 {
	m = new(sync.RWMutex)

	///Users/nicker/go/etc/src/CenterSettlement-go/conf
	//cfg, err := ini.Load("CenterSettlement-go/conf/id.conf") //读配置文件
	cfg, err := ini.Load("./conf/id.conf") //读配置文件  goland不能使用 ./ 方式  go run main.go 可以

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
	m.Lock()
	cfg.Section("").Key("messageid").SetValue(s) //  修改后值然后进行保存
	//Saveerr := cfg.SaveTo("CenterSettlement-go/conf/id.conf")
	Saveerr := cfg.SaveTo("./conf/id.conf")
	m.Unlock()
	if Saveerr != nil {
		log.Fatal("Fail to SaveTo file:", Saveerr)
	}

	log.Println(" file new messageid: ", s)
	log.Println(" xmlfile new messageid: ", id)
	return id
}
