package service

import (
	"CenterSettlement-go/common"
	"CenterSettlement-go/conf"
	storage "CenterSettlement-go/storage"
	"CenterSettlement-go/types"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

var (
	Filenameid string
	amountStr  string
)

//线程1
//处理原始数据的打包
func HandleGeneratexml() {
	log.Println("执行线程1，处理原始数据的打包")
	//tiker := time.NewTicker(time.Minute * 10)//每10分钟执行一下
	tiker := time.NewTicker(time.Second * 20) //每10秒执行一下

	for {
		log.Println(common.DateTimeFormat(<-tiker.C), "执行线程1，处理原始数据的打包")

		//其他省市地区    xml文件生成
		for _, Diqu := range types.Gl_network {

			//储值卡 cz xml文件生成
			czerr, czfn := Genaratexml(types.PRECARD, Diqu)
			if czerr != nil {
				log.Println("此地区储值卡原始交易数据打包 错误", Diqu, czerr)
				return
			}
			if czerr == nil && czfn == "no" {
				//没有储值卡原始数据
				log.Println("没有此地区储值卡原始数据", Diqu)
			}

			//记账卡 jz xml文件生成
			jzerr, jzfn := Genaratexml(types.CREDITCARD, Diqu)
			if jzerr != nil {
				log.Println("此地区记账卡原始交易数据打包 错误", Diqu, jzerr)
				return
			}
			if jzerr == nil && jzfn == "no" {
				log.Println("没有此地区记账卡原始数据", Diqu)
			}
		}

	}
}

//生成原始交易消息包
func Genaratexml(Kalx int, Diqu string) (error, string) {
	//从数据层获取准备的数据
	jiesuansj := *storage.QueryQitaJiessj(Kalx, Diqu)
	if Kalx == types.PRECARD && len(jiesuansj) == 0 {
		log.Println("数据库没有此地区要打包的储值卡的数据", Diqu)
		return nil, "no"
	}
	if Kalx == types.CREDITCARD && len(jiesuansj) == 0 {
		log.Println("数据库没有此地区要打包的记账卡的数据", Diqu)
		return nil, "no"
	}
	//生成原始交易消息包序号
	Messageid := conf.GenerateMessageId()
	Filenameid = fmt.Sprintf("%020d", Messageid)
	//数据赋值，获取xml交易数据
	jiaoyisj := TransAssignment(jiesuansj, Messageid)

	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(jiaoyisj, "  ", " ")
	if err != nil {
		log.Printf("打包原始记录消息包 xml.MarshalIndent error: %v\n", err)
		return err, ""
	}
	//更新结算数据为  打包中
	sjid := make([]string, 0)
	for _, v := range jiesuansj {
		sjid = append(sjid, v.FVcJiaoyjlid)
	}
	// 更新打包状态 为： 打包中
	uerr := storage.UpdatePackaging(sjid)
	if uerr != nil {
		log.Println("打包状态 为打包中 失败", uerr)
		return uerr, ""
	}

	var fname string
	//创建xml文件 cz 储值卡 、jz 记账卡
	fname = createxml(Diqu, Kalx, outputxml)
	if fname == "" {
		//把打包中的状态更新为未打包
		//更新结算数据为  未打包  jiesuansj
		sjid := make([]string, 0)
		for _, v := range jiesuansj {
			sjid = append(sjid, v.FVcJiaoyjlid)
		}
		//更新打包状态 为  初始
		uperr := storage.UpdatePackagingInit(sjid)
		if uperr != nil {
			log.Println("更新打包状态 为 初始、未打包状态 失败", uperr)
			return uperr, ""
		}

	}

	if fname != "" {
		//打包成功
		//原始记录消息包赋值
		yuansjyxx := YuanshiMsgAssignment(jiaoyisj, fname)
		//  新增原始交易数据打包记录【插入原始交易消息表】
		err1 := storage.PackagingRecordInsert(yuansjyxx)
		if err1 != nil {
			log.Println("新增消息包打包记录 error ", err1)
			return err1, ""
		}

		//   新增原始交易数据打包明细记录
		yuansjyxxmx := YuanshiMsgMXAssignment(jiaoyisj)
		err2 := storage.PackagingMXRecordInsert(yuansjyxxmx)
		if err2 != nil {
			log.Println("新增打包明细记录 error ", err2)
			return err2, ""
		}

		// 更新结算数据打包结果【打包状态：已打包、原始交易包号、包内序号、清分目标日】
		err3 := storage.UpdateDataPackagingResults(sjid, Messageid, jiaoyisj)
		if err3 != nil {
			log.Println("更新结算数据打包结果 error ", err3)
			return err3, ""
		}
		return nil, fname
	}
	return nil, fname
}

//原始交易xml赋值
func TransAssignment(jiesuansj []types.BJsJiessj, Messageid int64) *types.Message {

	Trans := make([]types.Transaction, 0)
	count := len(jiesuansj)
	var amount int64
	for i, v := range jiesuansj {
		//把数据库数据 准备为 xml需要的  赋值
		var Tran types.Transaction
		Tran.TransId = i + 1 //包内序号

		jiaoysj := common.DateTimeFormat(v.FDtJiaoysj) //交易时间 的格式转换

		Tran.Time = jiaoysj                //交易时间
		yuan := common.Fen2Yuan(v.FNbJine) //分转元
		Tran.Fee = yuan                    //交易金额 元
		amount = amount + v.FNbJine        //总金额 分   224行 转元

		Tran.Service.ServiceType = types.SERVICETYPE //交易服务类型 【写死2】
		//账单描述  南京南站南广场P3|11小时32分40秒
		//停车场名称
		name := storage.GetTingcc(v.FVcTingccbh)
		//停车时长
		tcsj := common.TimeDifference(v.FDtYonghrksj, v.FDtJiaoysj)
		//账单描述  南京南站南广场P3|11小时32分40秒
		Tran.Service.Description = name + "｜" + tcsj

		rksj := common.DateTimeFormat(v.FDtYonghrksj)

		//站编码转换
		if v.FVcTingccbh == "0025000011" || v.FVcTingccbh == "0025000010" {
			v.FVcTingccbh = "3201020001"
		}
		//if v.FVcTingccbh=="0025000010"{
		//	v.FVcTingccbh="3201020001"
		//}

		//cx:车型 ckz：出口站、入口站  ckcd：出口车道，入口车道cksj：出口时间 rksj：入口时间
		detail := common.Detail(v.FVcChex, v.FVcTingccbh, v.FVcChedid, jiaoysj, rksj)
		Tran.Service.Detail = detail     //交易详细信息 1|04|3201|3201000006|1105|20191204 211733|03|3201|320
		Tran.ICCard.CardType = v.FNbKalx //卡类型 22,23
		Tran.ICCard.NetNo = v.FVcKawlh   //卡网络号
		Tran.ICCard.CardId = v.FVcKah    //卡号
		ys := ChePYS(v.FVcObucpys)
		Tran.ICCard.License = v.FVcCheph + ys //卡内车牌号 皖EYG511蓝
		jiaoyqye := common.Fen2Yuan(v.FNbJiaoyqye)
		Tran.ICCard.PreBalance = jiaoyqye //交易前余额，以元为单位 Decimal
		jiaoyhye := common.Fen2Yuan(v.FNbJiaoyhye)
		Tran.ICCard.PostBalance = jiaoyhye //交易后余额，以元为单位 Decimal

		Tran.Validation.TAC = v.FVcTacm                 //交易TAG码
		Tran.Validation.TransType = types.TRANSTYPE     //交易标识，2位16进制数，PBOC定义，如06为传统交易，09为复合交【09写死】
		Tran.Validation.TerminalNo = v.FVcJiamkh        //终端机编号  加密卡号
		Tran.Validation.TerminalTransNo = v.FVcKajmjyxh //PSAM卡脱机交易序号，在MAC1计算过程中得到  加密加密序列号

		Tran.OBU.NetNo = v.FVcKawlh //OBU网络号 此处与卡网络号一致
		Tran.OBU.OBUId = v.FVcObuid
		Tran.OBU.OBEState = v.FVcObuzt //2字节的OBU状态
		obucpys := ChePYS(v.FVcObucpys)
		Tran.OBU.License = v.FVcObucp + obucpys   //OBU中记录的车牌号 皖EYG511蓝 【加颜色】
		jyje := common.ToHexFormat(v.FNbJine, 8)  //交易金额
		jyhje := common.ToHexFormat(v.FNbJine, 8) //交易前金额
		jyqje := common.ToHexFormat(v.FNbJine, 8) //交易后金额

		customizedData := common.CustomizedData(v.FVcTacm, jyje, types.TRANSTYPE, v.FVcJiamkh, v.FVcKajmjyxh, jiaoysj, jyhje, jyqje, v.FVcKajyxh)
		Tran.CustomizedData = customizedData
		//停车场消费交易编号(停车场编号+交易发生的时间+流水号【怎么取】)
		//流水号  FVcChedjyxh 车道交易序号  取后两位
		liush := common.GetLiush(v.FVcChedjyxh)
		ID := common.GetId(v.FVcTingccbh, jiaoysj, liush)
		Tran.Id = ID                   //停车场消费交易编号
		Tran.Name = name               //停车场名称
		Tran.ParkTime = v.FNbYonghtcsc //停车时常
		Tran.VehicleType = v.FVcChex   //车型

		Tran.AlgorithmIdentifier = types.ALGORITHMIDENTIFIER //算法标识 1-3DEX  2-SM4【取1】

		Trans = append(Trans, Tran)
	}
	amountStr = common.Fen2Yuan(amount) //总金额 元

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
			ClearTargetDate:   common.TodayFormat(),
			ServiceProviderId: "00000000000000FD",
			IssuerId:          "0000000000000020",
			MessageId:         Messageid,
			Count:             count,
			Amount:            amountStr,
		}}
	for _, T := range Trans {
		jiaoyisj.Body.Transaction = append(jiaoyisj.Body.Transaction, T)
	}
	return jiaoyisj
}

//原始记录消息包 赋值
func YuanshiMsgAssignment(jiaoyisj *types.Message, fname string) types.BJsYuansjyxx {
	var yuansjyxx types.BJsYuansjyxx
	yuansjyxx.FVcBanbh = "00010000"                       //版本号
	yuansjyxx.FNbXiaoxlb = 5                              //消息类别c
	yuansjyxx.FNbXiaoxlx = 7                              //消息类型Type
	yuansjyxx.FVcFaszid = "00000000000000FD"              //发送者ID
	yuansjyxx.FVcJieszid = "0000000000000020"             //接受者ID
	yuansjyxx.FNbXiaoxxh = jiaoyisj.Body.MessageId        //消息序号【消息包号】
	yuansjyxx.FDtDabsj = common.DateTimeNowFormat()       // 打包时间
	yuansjyxx.FVcQingfmbr = jiaoyisj.Body.ClearTargetDate //清分目标日
	yuansjyxx.FVcTingccqffid = "00000000000000FD"         //停车场清分方ID
	yuansjyxx.FVcFaxfwjgid = "0000000000000020"           //发行服务机构ID 0000000000000020
	yuansjyxx.FNbJilsl = jiaoyisj.Body.Count              //记录数量
	yuansjyxx.FNbZongje = jiaoyisj.Body.Amount            //总金额
	yuansjyxx.FVcXiaoxwjlj = "generatexml/" + fname       //消息文件路径
	return yuansjyxx
}

//原始交易消息明细 赋值
func YuanshiMsgMXAssignment(jiaoyisj *types.Message) []types.BJsYuansjymx {
	var mx []types.BJsYuansjymx
	var yuansjymx types.BJsYuansjymx
	yuansjymx.FVcXiaoxxh = jiaoyisj.Header.MessageId //消息序号
	for _, T := range jiaoyisj.Body.Transaction {
		yuansjymx.FNbBaonxh = T.TransId           //包内序号
		yuansjymx.FDtJiaoysj = T.Time             //交易时间
		yuansjymx.FNbJine = T.Fee                 //交易金额
		yuansjymx.FVcDingzjyxx = T.CustomizedData //定制交易信息
		yuansjymx.FVcJiaoybh = T.Id               //交易编号
		yuansjymx.FVcTingccmc = T.Name            //停车场名称
		yuansjymx.FNbTingfsc = T.ParkTime         //停放时长 分
		cx, _ := strconv.Atoi(T.VehicleType)      //车型
		yuansjymx.FNbShoufcx = cx
		yuansjymx.FNbSuanfbs = T.AlgorithmIdentifier  //算法标识
		yuansjymx.FNbFuwlx = T.Service.ServiceType    //服务类型
		yuansjymx.FVcZhangdsm = T.Service.Description //账单说明
		yuansjymx.FVcJiaoyxxxx = T.Service.Detail     //交易详细信息
		yuansjymx.FNbKalx = T.ICCard.CardType         //卡类型
		yuansjymx.FVcWanglbm = T.ICCard.NetNo         //网络编码 卡网络号
		yuansjymx.FVcKawlbh = T.ICCard.CardId         //卡物理编号
		yuansjymx.FVcKancph = T.ICCard.License        //卡内车牌号
		//yuansjymx.FVcKajyxh=T.ICCard.   //卡交易序号
		jyqye, _ := strconv.Atoi(T.ICCard.PreBalance)
		yuansjymx.FNbJiaoyqye = int64(jyqye) //交易前余额
		jyhye, _ := strconv.Atoi(T.ICCard.PostBalance)
		yuansjymx.FNbJiaoyhye = int64(jyhye)                  //交易后余额
		yuansjymx.FVcTacm = T.Validation.TAC                  //TAC码
		yuansjymx.FVcJiaoybs = T.Validation.TransType         //交易标识
		yuansjymx.FVcZongdjh = T.Validation.TerminalNo        //终端机号
		yuansjymx.FVcZongdjyxh = T.Validation.TerminalTransNo //终端交易序号
		yuansjymx.FVcObuwlbh = T.OBU.NetNo                    //obu物理编号
		yuansjymx.FVcObuzt = T.OBU.OBEState                   //obu状态
		yuansjymx.FVcObuncph = T.OBU.License                  //obu内车牌号
		mx = append(mx, yuansjymx)
	}

	return mx
}

//创建xml文件
func createxml(Kawlh string, Kalx int, outputxml []byte) string {
	var kalxstr string
	if Kalx == 22 {
		kalxstr = "CZ"
	}
	if Kalx == 23 {
		kalxstr = "JZ"
	}
	fw, f_werr := os.Create("./generatexml/" + kalxstr + "_" + Kawlh + "_" + Filenameid + ".xml") //go run main.go
	if f_werr != nil {
		log.Fatal("Read:", f_werr)
		return ""
	}
	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, outputxml...)

	_, ferr := fw.Write((xmlOutPutData))
	if ferr != nil {
		log.Printf("Write xml file error: %v\n", ferr)
		return ""
	}
	//更新消息包信息
	fw.Close()
	return kalxstr + "_" + Kawlh + "_" + Filenameid + ".xml"

}

//车牌号颜色转换
func ChePYS(in string) string {
	var ys string
	switch in {
	case "0":
		ys = "蓝"
	case "1":
		ys = "黄"
	case "2":
		ys = "黑"
	case "3":
		ys = "白"
	case "4":
		ys = "渐变绿"
	case "5":
		ys = "黄绿双拼"
	case "6":
		ys = "蓝白渐变"
	case "7":
		ys = "临时拍照"
	case "11":
		ys = "绿"
	case "12":
		ys = "红"
	case "10000":
		ys = "蓝"
	}
	return ys
}

//返回一个32位md5加密后的字符串
func GetMD5Encode(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
