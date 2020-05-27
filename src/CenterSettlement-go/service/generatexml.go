package service

import (
	"CenterSettlement-go/common"
	"CenterSettlement-go/conf"
	storage "CenterSettlement-go/storages"
	"CenterSettlement-go/types"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	Messageid  int64
	Filenameid string
	count      int
	amount     int64
	amountStr  string
)

//线程1
//处理原始数据的打包
func HandleGeneratexml() {

	tiker := time.NewTicker(time.Second * 5)
	for {
		log.Println("执行线程1")
		<-tiker.C
		//储值卡 cz xml文件生成
		czfname := Generatexml(types.PRECARD)
		if czfname != "" {
			//没有本省的储值卡原始数据
		}

		//记账卡 jz xml文件生成
		jzfname := Generatexml(types.CREDITCARD)
		if jzfname != "" {
			//没有本省的记账卡原始数据
		}
	}
}

//
func Generatexml(Kalx int) string {
	//从数据层获取准备的数据
	Trans := make([]types.Transaction, 0)
	//获取本省数据
	jiesuansj := *storage.QueryJiessj(Kalx)
	if Kalx == types.PRECARD && len(jiesuansj) == 0 {
		log.Println("数据库没有要打包的本省的储值卡的数据")
		return ""
	}
	if Kalx == types.CREDITCARD && len(jiesuansj) == 0 {
		log.Println("数据库没有要打包的本省的记账卡的数据")
		return ""
	}
	count = len(jiesuansj)
	//消息包序号
	Messageid = conf.GenerateMessageId()
	Filenameid = fmt.Sprintf("%020d", Messageid)
	for i, v := range jiesuansj {
		//把数据库数据 准备为 xml需要的  赋值
		var Tran types.Transaction
		Tran.TransId = i + 1                                           //包内序号
		Tran.Time = v.FDtJiaoysj                                       //交易时间
		yuan := common.Fen2Yuan(v.FNbJine)                             //分转元
		Tran.Fee = yuan                                                //交易金额 yuan
		amount = amount + v.FNbJine                                    //总金额
		Tran.CustomizedData = time.Now().Format("2006-01-02 15:04:05") //【清分目标日】 当前日期
		Tran.Service.ServiceType = 2                                   //交易服务类型 【写死2】
		//获取停车场名称
		tName := storage.GetTingcc(v.FVcTingccbh)
		//获取停车时长

		Tran.Service.Description = tName + "|" + "" //账单描述  南京南站南广场P3|11小时32分40秒
		Tran.Service.Description = v.FVcZhangdms
		//Tran.Service.Detail=//交易详细信息 1|04|3201|3201000006|1105|20191204 211733|03|3201|320
		Tran.ICCard.CardType = v.FNbKalx //卡类型
		Tran.ICCard.NetNo = v.FVcKawlh   //卡网络号
		Tran.ICCard.CardId = v.FVcKah    //卡号
		//Tran.ICCard.License=v.//卡内车牌号
		Tran.ICCard.PostBalance = v.FNbJiaoyhye / 100 //交易后余额，以元为单位 Decimal
		Tran.ICCard.PreBalance = v.FNbJiaoyqye        //交易前余额，以元为单位 Decimal
		Tran.Validation.TAC = v.FVcTacm               //交易TAG码
		//Tran.Validation.TransType=//交易标识，2位16进制数，PBOC定义，如06为传统交易，09为复合交
		//Tran.Validation.TerminalNo//终端机编号
		//Tran.Validation.TerminalTransNo//PSAM卡脱机交易序号，在MAC1计算过程中得到
		Tran.OBU.License = v.FVcObucp //OBU中记录的车牌号
		//Tran.OBU.NetNo=//OBU网络号
		Tran.OBU.OBEState = v.FVcObuzt //2字节的OBU状态
		Tran.OBU.OBUId = v.FVcObuid
		//Trans[Ti].Service.Description=	v.FNbYonghtcsc//  	停车场名｜停车时常

		Trans = append(Trans, Tran)
		//log.Println(Trans)
	}
	amountStr = common.Fen2Yuan(amount)

	//赋值
	//把原始交易数据转化成xml文件
	jiaoyisj := &types.Message{
		Header: types.Header{
			Version:      "00010000",
			MessageClass: 5,
			MessageType:  7,
			SenderId:     "00000000000000FD",
			ReceiverId:   "0000000000000020",
			MessageId:    Messageid},
		Body: types.Body{
			ContentType: 1,
			//清分目标日 当前日期
			ClearTargetDate:   time.Now().Format("2006-01-02"),
			ServiceProviderId: "00000000000000FD",
			IssuerId:          "0000000000000020",
			MessageId:         Messageid,
			Count:             count,
			Amount:            amountStr,
		}}
	for _, T := range Trans {
		jiaoyisj.Body.Transaction = append(jiaoyisj.Body.Transaction, T)
	}
	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(jiaoyisj, " ", " ")
	//使用Marshal函数，生成的XML格式无缩进
	//outputxml,err:=xml.Marshal(v)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	//log.Println(outputxml)
	var fname string
	//创建xml文件 cz 储值卡
	if Kalx == types.PRECARD {
		fname = createxml(Kalx, outputxml)
	}
	//创建xml文件 jz 记账卡
	if Kalx == types.CREDITCARD {
		fname = createxml(Kalx, outputxml)
	}
	//打包成功
	//
	return fname
}

func description() {

}

func createxml(Kalx int, outputxml []byte) string {
	var kalxstr string
	if Kalx == 22 {
		kalxstr = "CZ"
	}
	if Kalx == 23 {
		kalxstr = "JZ"
	}
	//
	fw, f_werr := os.Create("../generatexml/" + kalxstr + "_3201_" + Filenameid + ".xml")
	if f_werr != nil {
		log.Fatal("Read:", f_werr)
	}
	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, outputxml...)
	//这里可以不写，直接使用channel发送给线程2
	//写入文件
	//ioutil.WriteFile("../generatexml/"+kalxstr+"_3201_"+Filenameid+".xml", xmlOutPutData, os.ModeAppend)

	_, ferr := fw.Write((xmlOutPutData))
	if ferr != nil {
		log.Printf("Write xml file error: %v\n", ferr)
	}
	//更新消息包信息
	fw.Close()
	//return "../generatexml/"+kalxstr+"_3201_"+Filenameid+".xml"
	return kalxstr + "_3201_" + Filenameid + ".xml"

}

//返回一个32位md5加密后的字符串
func GetMD5Encode(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

//封装一个函数，处理xml数据的准备
func xmldata() {
	//原始交易数据

}
