package centerServer

import (
	"CenterSettlement-go/common"
	"encoding/xml"
	"os"

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
		qerr, jzshuju, msg := QueryKeepAccountdata()
		if qerr != nil {
			log.Fatal(qerr)
		}
		//遍历map

		for msgid, jzsj := range *jzshuju {
			//for _,mid:=range *msg{
			//
			//
			//}

			//组织xml数据  生成xml
			GenerateJZokXml(msgid, jzsj, msg)
		}

	}
}

func GenerateDisputeXml() {
	//获取数据

	//组织xml数据

	//生成xml
	//GenerateXml()
}

func GenerateClearlingtXml() {
	//获取数据

	//组织xml数据

	//生成xml
	//GenerateXml()
}

func GenerateJZokXml(msgid int64, jzsj *[]Jiessjchuli, msg *[]JieSuanMessage) {
	//获取记账消息包的序号
	Messageid := GenerateMessageId()
	Filenameid := fmt.Sprintf("%020d", Messageid)
	log.Println(Filenameid)
	//new(KeepAccountokMessage)
	var t string
	t = common.DateFormat()
	//
	count := len(*jzsj)
	var amount string
	for _, v := range *msg {
		if v.MessageId == msgid {
			//全部记账
			if v.Count == count {
				log.Println("全部都可以记账，没有争议数据")
				amount = v.Amount
				jzmsg := &KeepAccountokMessage{
					Header: KeepAccountHeader{
						Version:      "00020000",
						MessageClass: 5,
						MessageType:  5,
						SenderId:     "0000000000000020",
						ReceiverId:   "00000000000000FD",
						MessageId:    Filenameid},
					Body: KeepAccountOkBody{
						ContentType:       1,
						ServiceProviderId: "00000000000000FD", //通行宝中心系统Id，
						IssuerId:          "0000000000000020", //发行服务机构Id， 记账消息哪一个发行方
						MessageId:         msgid,              //交易消息包Id。原始交易包消息中的messageid
						ProcessTime:       t,                  //处理时间
						Count:             count,              //本消息对应的原始交易包中交易记录的数量
						Amount:            amount,             //确认记帐总金额 交易总金额(元) 数据库为分【注意转换的小数问题】
						DisputedCount:     0,
					}}
				//使用MarshalIndent函数，生成的XML格式有缩进
				outputxml, err := xml.MarshalIndent(jzmsg, "  ", "  ")
				if err != nil {
					log.Printf("error: %v\n", err)
					//return err, ""
				}

				fname := createxml("JZB-ok", outputxml, Filenameid)
				if fname == "" {
					log.Println("生成记账包失败")
				}
			}

			//有争议

		}

	}

}

//创建xml文件
func createxml(lx string, outputxml []byte, Filenameid string) string {
	fw, f_werr := os.Create("../centerkeepaccount/" + lx + "_" + Filenameid + ".xml")
	if f_werr != nil {
		log.Fatal("Read:", f_werr)
		return ""
	}
	//加入XML头
	//headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	//xmlOutPutData := append(headerBytes, outputxml...)

	_, ferr := fw.Write((outputxml))
	if ferr != nil {
		log.Printf("Write xml file error: %v\n", ferr)
		return ""
	}
	//更新消息包信息
	fw.Close()
	return lx + "_" + Filenameid + ".xml"

}
