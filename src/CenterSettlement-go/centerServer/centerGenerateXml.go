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

//生成争议处理包
func GenerateDisputeXml() {
	tiker := time.NewTicker(time.Second * 5)
	for {
		log.Println("执行   生成争议处理包", <-tiker.C)

		gerr := GenerateZYxml()
		if gerr != nil {
			log.Fatal(gerr)
		}
	}
}

//生成清分处理包
func GenerateClearlingtXml() {
	tiker := time.NewTicker(time.Second * 5)
	for {
		log.Println("执行   生成清分统计处理包", <-tiker.C)

		gerr := GenerateQFxml()
		if gerr != nil {
			log.Fatal(gerr)
		}
	}
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
			//更新记账状态：2 清分状态1
			err1 := UpdatemsgJZ(jsmsj.MessageId)
			if err1 != nil {
				return err1
			}
			//更新记账处理、记账状态：2 清分状态1 记账包号
			err2 := UpdateJZclzt(jsmsj.MessageId, int64(f))
			if err2 != nil {
				return err1
			}
			break
		}

		if len(msgjschulisj) != jsmsj.Count {
			log.Println("此包有争议数据", jsmsj.MessageId)

		}
	}
	return nil
}

//生成没有争议的记账xml包
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
	var pwd string
	switch lx {
	case "JZB-ok":
		pwd = "../centerkeepaccount/"
	case "JZB":
		pwd = "../centerkeepaccount/"

	case "QFB":
		pwd = "../centerClearing/"

	case "ZYB":
		pwd = "../centerdispute/"
	}

	fw, f_werr := os.Create(pwd + lx + "_" + Filenameid + ".xml")
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

//生成争议处理包
func GenerateZYxml() error {

	//获取数据

	//组织xml数据

	//生成xml
	//GenerateXml()
	return nil
}

func GenerateQFxml() error {
	db := NewDatabase()
	//获取需要的清分包数据
	msgerr, msgs := QueryCleariungMsgdata()
	if msgerr != nil {
		log.Println("查询可以记账的原始交易消息包出错", msgerr)
		return msgerr
	}
	if len(*msgs) == 0 {
		log.Println("没有需要清分的数据包")
		return nil
	}

	var msgids []int64
	Count := 0
	Amount := 0
	ClearTargetDate := ""
	msgcount := strconv.Itoa(len(*msgs))
	for _, jsmsj := range *msgs {
		log.Println("查询可以清分的原始交易消息包 ", jsmsj.MessageId)

		msgjschulisj := make([]Jiessjchuli, 0)
		//同一个数据包 可以记账的数据
		qerr := db.orm.Where("f_nb_yuansjybxh=?", jsmsj.MessageId).And("f_nb_qingfjg=?", 1).Find(&msgjschulisj)
		if qerr != nil {
			log.Fatalln("查询结算数据出错", qerr)
			return qerr
		}

		log.Println("此消息包总共查询出需要清分的数据", jsmsj.MessageId, len(msgjschulisj))

		if len(msgjschulisj) == jsmsj.Count {
			log.Println("此包没有争议数据，均可以清分", jsmsj.MessageId)
			msgids = append(msgids, jsmsj.MessageId)

			Count = Count + jsmsj.Count

			amountstr := strings.Split(jsmsj.Amount, ".")
			amount, _ := strconv.Atoi(amountstr[0] + amountstr[1])
			Amount = Amount + amount

			ClearTargetDate = jsmsj.ClearTargetDate
			continue
			//break
		}

		if len(msgjschulisj) != jsmsj.Count {
			log.Println("此包有争议数据", jsmsj.MessageId)
			//break
		}

	}

	log.Println("清分数据", Count, int64(Amount), msgcount, ClearTargetDate)

	//生成清分统计包xml
	fname := GenerateQFokXml(msgids, Count, int64(Amount), msgcount, ClearTargetDate)
	if fname == "" {
		return errors.New("生成记账包xml 失败")
	}
	fnstr := strings.Split(fname, "_")
	fstr := strings.Split(fnstr[1], ".")
	f, _ := strconv.Atoi(fstr[0]) //清分包号

	//更新清分状态 清分状态：已清分2、清分包号
	err1 := UpdatemsgQFzt(msgids)
	if err1 != nil {
		return err1
	}

	//更新清分处理记录
	err2 := UpdateQFclzt(msgids, int64(f))
	if err2 != nil {
		return err1
	}

	//GenerateXml()
	return nil
}

func GenerateQFokXml(msgids []int64, Count int, Amount int64, msgcount string, ClearTargetDate string) string { //获取记账消息包的序号
	fMessageid := GenerateMessageId()
	Filenameid := fmt.Sprintf("%020d", fMessageid)
	log.Println("清分包消息的文件id", Filenameid)
	amount := common.Fen2Yuan(Amount)
	var t string
	t = common.DateFormat() //生成清分包时间

	jzmsg := &ClearingOKMessage{
		Header: ClearingHeader{
			Version:      "00010000",
			MessageClass: 5,
			MessageType:  5,
			SenderId:     "0000000000000020",
			ReceiverId:   "00000000000000FD",
			MessageId:    fMessageid},
		Body: ClearingOKBody{
			ContentType:     2,
			ClearTargetDate: ClearTargetDate, //清分目标日
			Count:           Count,           //对应清分总金额的的交易记录数量，包含原始交易包中由发行方确认应付的交易数量和争议处理结果中确认应付的交易数量之和， 不包含争议处理结果中为坏帐的记录数量。
			Amount:          amount,          //清分总金额(确认付款金额)
			ProcessTime:     t,               //处理时间 清分统计处理时间
			//IssuerId        : "00000000000000FD",//发行服务机构Id，
			List: ListOK{
				MessageCount:      msgcount,
				FileCount:         "0",
				ServiceProviderId: "00000000000000FD",
				MessageId:         msgids,
			},
		}}

	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(jzmsg, "  ", "  ")
	if err != nil {
		log.Printf("error: %v\n", err)
		return ""
	}

	fname := createxml("QFB", outputxml, Filenameid)
	if fname == "" {
		log.Println("生成记账包失败")
		return ""
	}

	log.Println("生成记账包 成功：", fname)
	return fname
}
