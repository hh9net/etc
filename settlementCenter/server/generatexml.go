package server

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	storage "settlementCenter/storages"
	"settlementCenter/types"
	"time"
)

var (
	Messageid      int64
	Filename       string
	count          int
	amount         int
	customizedData string
	serviceType    string
)

//线程1
func Generatexml() string {

	Trans := make([]types.Transaction, 0)
	//获取本省数据
	jiesuansj := *storage.QueryJiessj()
	//消息包起始序号
	Messageid = 999999
	Filename = fmt.Sprintf("%020d", Messageid)
	count = len(jiesuansj)
	log.Println(count)

	for i, v := range jiesuansj {
		//把数据库数据 准备xml需要的  赋值
		var Tran types.Transaction
		Tran.TransId = i + 1                         //包内序号
		Tran.Time = v.FDtJiaoysj                     //交易时间
		Tran.Fee = v.FNbJine                         //交易金额
		Tran.CustomizedData = "===customizedData===" //【清分目标日？？？】
		//Tran.Service.ServiceType = v.FDtJiaoylx                    //交易类型
		Tran.Service.Description = v.FVcZhangdms //账单描述  南京南站南广场P3|11小时32分40秒
		//Tran.Service.Detail=//交易详细信息 1|04|3201|3201000006|1105|20191204 211733|03|3201|320
		Tran.ICCard.CardType = v.FNbKalx //卡类型
		Tran.ICCard.NetNo = v.FVcKawlh   //卡网络号
		Tran.ICCard.CardId = v.FVcKah    //卡号
		//Tran.ICCard.License=v.//卡内车牌号
		//Tran.ICCard.PostBalance=v.FNbJiaoyhye//交易后余额，以元为单位 Decimal
		//Tran.ICCard.PreBalance=v.FNbJiaoyqye  //交易前余额，以元为单位 Decimal
		Tran.Validation.TAC = v.FVcTacm //交易TAG码
		//Tran.Validation.TransType=//交易标识，2位16进制数，PBOC定义，如06为传统交易，09为复合交
		//Tran.Validation.TerminalNo//终端机编号
		//Tran.Validation.TerminalTransNo//PSAM卡脱机交易序号，在MAC1计算过程中得到
		Tran.OBU.License = v.FVcObucp //OBU中记录的车牌号
		//Tran.OBU.NetNo=//OBU网络号
		Tran.OBU.OBEState = v.FVcObuzt //2字节的OBU状态
		Tran.OBU.OBUId = v.FVcObuid
		//Trans[Ti].Service.Description=	v.FNbYonghtcsc//  	停车场名｜停车时常

		Trans = append(Trans, Tran)
		log.Println(Trans)
	}

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
			ContentType:       1,
			ClearTargetDate:   time.Now(),
			ServiceProviderId: "00000000000000FD",
			IssuerId:          "0000000000000020",
			MessageId:         Messageid,
			Count:             count,
			Amount:            amount,
		}}
	for _, T := range Trans {
		jiaoyisj.Body.Transaction = append(jiaoyisj.Body.Transaction, T)
	}
	log.Println(jiaoyisj)
	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(jiaoyisj, "  ", "  ")
	//使用Marshal函数，生成的XML格式无缩进
	//outputxml,err:=xml.Marshal(v)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	log.Println(outputxml)

	//创建文件
	fw, f_werr := os.Create("../generatexml/" + "CZ_3201_" + Filename + ".xml")
	if f_werr != nil {
		log.Fatal("Read:", f_werr)
	}
	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, outputxml...)
	//写入文件
	ioutil.WriteFile("CZ_3201_"+Filename+".xml", xmlOutPutData, os.ModeAppend)
	//
	defer fw.Close()
	_, ferr := fw.Write((xmlOutPutData))
	if ferr != nil {
		fmt.Printf("Write xml file error: %v\n", err)
	}
	//更新消息包信息

	return "CZ_3201_" + Filename + ".xml"
}

//返回一个32位md5加密后的字符串
func GetMD5Encode(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
