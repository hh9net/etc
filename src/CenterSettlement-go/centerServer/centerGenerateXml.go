package centerServer

import (
	"CenterSettlement-go/common"
	"encoding/xml"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"

	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

//生成记账包、争议包、清分包
func GenerateKeepAccountXml() {
	tiker := time.NewTicker(time.Second * 3)
	for {
		log.Println("执行   生成记账包", <-tiker.C)

		gerr := GenerateJZxml()
		if gerr != nil {
			log.Fatal(gerr)
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

func GenerateJZxml() error {
	db := NewDatabase()
	//查询数据

	//查询出需要打包的原始交易包序号
	msgerr, msgs := QueryKeepAccountMsgdata()
	if msgerr != nil {
		log.Println("查询可以记账的原始交易消息包出错", msgerr)
		return msgerr
	}
	for _, jsmsj := range *msgs {
		log.Println("查询可以记账的原始交易消息包 ", jsmsj.MessageId)
		msgjschulisj := make([]Jiessjchuli, 0)
		//同一个数据包 可以记账的数据
		qerr := db.orm.Where("f_nb_yuansjybxh=?", jsmsj.MessageId).And("f_nb_jizjg=?", 1).Find(&msgjschulisj)
		if qerr != nil {
			log.Fatalln("查询结算数据出错", qerr)
			return qerr
		}
		log.Println("此消息包总共查询出需要记账的数据数据", jsmsj.MessageId, len(msgjschulisj))

		if len(msgjschulisj) == jsmsj.Count {
			log.Println("此包没有争议数据", jsmsj.MessageId)
			//生成记账包xml
			fname := GenerateJZokXml(jsmsj.MessageId, &msgjschulisj, &jsmsj)
			if fname == "" {
				return errors.New("生成记账包xml 失败")
			}
			fnstr := strings.Split(fname, "_")
			fstr := strings.Split(fnstr[1], ".")
			f, _ := strconv.Atoi(fstr[0]) //记账包号
			//更新记账状态 清分状态、记账包号
			err1 := UpdatemsgJZ(jsmsj.MessageId)
			if err1 != nil {
				return err1
			}
			//更新记账处理
			err2 := UpdateJZclzt(jsmsj.MessageId, int64(f))
			if err2 != nil {
				return err1
			}
			break
		}
	}
	return nil
}

func GenerateJZokXml(msgid int64, jzsj *[]Jiessjchuli, msg *JieSuanMessage) string { //获取记账消息包的序号
	Messageid := GenerateMessageId()
	Filenameid := fmt.Sprintf("%020d", Messageid)
	log.Println("记账包消息的文件id", Filenameid)

	var t string
	t = common.DateFormat() //生成记账包时间
	//
	count := len(*jzsj)
	log.Println("全部都可以记账，没有争议数据")

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
			Amount:            msg.Amount,         //确认记帐总金额 交易总金额(元) 数据库为分【注意转换的小数问题】
			DisputedCount:     0,
		}}
	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(jzmsg, "  ", "  ")
	if err != nil {
		log.Printf("error: %v\n", err)
		return ""
	}

	fname := createxml("JZB-ok", outputxml, Filenameid)
	if fname == "" {
		log.Println("生成记账包失败")
		return ""
	}

	log.Println("生成记账包 成功：", fname)
	return fname
}

//创建xml文件
func createxml(lx string, outputxml []byte, Filenameid string) string {
	fw, f_werr := os.Create("../centerkeepaccount/" + lx + "_" + Filenameid + ".xml")
	if f_werr != nil {
		log.Fatal("创建xml文件 Read error :", f_werr)
		return ""
	}
	//加入XML头
	//headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	//xmlOutPutData := append(headerBytes, outputxml...)

	_, fwerr := fw.Write((outputxml))
	if fwerr != nil {
		log.Printf("Write xml file error: %v\n", fwerr)
		return ""
	}
	//更新消息包信息
	fw.Close()
	return lx + "_" + Filenameid + ".xml"
}
