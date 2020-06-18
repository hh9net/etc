package centerServer

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

//生成记账包、争议包、清分包
func GenerateKeepAccountXml() {

	tiker := time.NewTicker(time.Second * 3)
	for {
		log.Println("执行   生成记账包", <-tiker.C)
		//获取数据
		QueryKeepAccountdata()

		//组织xml数据

		//生成xml
		GenerateXml()

	}
}

func GenerateDisputeXml() {
	//获取数据

	//组织xml数据

	//生成xml
	GenerateXml()
}

func GenerateClearlingtXml() {
	//获取数据

	//组织xml数据

	//生成xml
	GenerateXml()
}

func GenerateXml() {
	//消息包序号
	Messageid := GenerateMessageId()
	Filenameid := fmt.Sprintf("%020d", Messageid)
	log.Println(Filenameid)
}
