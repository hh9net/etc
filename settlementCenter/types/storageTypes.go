package types

import "time"

//  B_JS_JIESSJ【结算数据】`b_js_jiessj`
type BJsJiessj struct {
	FVcJiaoyjlid   string    `xorm:"pk"` //F_VC_JIAOYJLID	交易记录ID	VARCHAR(128)
	FVcTingccbh    string    //F_VC_TINGCCBH	停车场编号	VARCHAR(32)
	FVcChedid      string    //F_VC_CHEDID	车道ID	VARCHAR(32)
	FVcGongsjtid   string    //F_VC_GONGSJTID	公司/集团ID	VARCHAR(32)
	FNbTingcclx    int       //F_NB_TINGCCLX	停车场类型	INT
	FNbJiaoyzt     int       //F_NB_JIAOYZT	交易状态	INT
	FNbYuansjybxh  int64     //F_NB_YUANSJYBXH	原始交易包序号	BIGINT
	FNbJiaoybnxh   int       //F_NB_JIAOYBNXH	交易包内序号	INT
	FNbJizjg       int       //F_NB_JIZJG	记账结果	INT
	FNbZhengylx    int       //F_NB_ZHENGYLX	争议类型	INT
	FNbJizbxh      int       //F_NB_JIZBXH	记账包序号	INT
	FNbZhengyclbxh int64     //F_NB_ZHENGYCLBXH	争议处理包序号	BIGINT
	FNbQingfbxh    int64     //F_NB_QINGFBXH	清分包序号	BIGINT
	FVcXiaofjlbh   string    //F_VC_XIAOFJLBH	消费记录编号	VARCHAR(128)
	FVcJiamkh      string    //F_VC_JIAMKH	加密卡号	VARCHAR(32)
	FVcKajmjyxh    string    //F_VC_KAJMJYXH	加密卡交易序号	VARCHAR(32)
	FVcObuid       string    //F_VC_OBUID	Obuid	VARCHAR(32)
	FVcObufxf      string    //F_VC_OBUFXF	obu发行方	VARCHAR(32)
	FVcObucp       string    //F_VC_OBUCP	obu内车牌	VARCHAR(32)
	FVcObucpys     string    //F_VC_OBUCPYS	obu车牌颜色	VARCHAR(32)
	FVcKah         string    //F_VC_KAH	卡号	VARCHAR(32)
	FVcKawlh       string    //F_VC_KAWLH	卡网络号	VARCHAR(32)
	FVcKajyxh      string    //F_VC_KAJYXH	卡交易序号	VARCHAR(32)
	FVcKafxf       string    //F_VC_KAFXF	卡发行方	VARCHAR(32)
	FNbKalx        int       //F_NB_KALX	卡类型	INT
	FNbJiaoyqye    int       //F_NB_JIAOYQYE	交易前余额	INT
	FNbJiaoyhye    int       //F_NB_JIAOYHYE	交易后余额	INT
	FNbJine        int       //F_NB_JINE	金额	INT
	FVcTacm        string    //F_VC_TACM	TAC码	VARCHAR(32)
	FDtJiaoysj     time.Time //F_DT_JIAOYSJ	交易时间	DATETIME
	FDtJiaoylx     string    //F_DT_JIAOYLX	交易类型	VARCHAR(32)
	FVcChex        string    //F_VC_CHEX	车型	VARCHAR(32)
	FVcObuzt       string    //F_VC_OBUZT	OBu状态	VARCHAR(32)
	FVcSuanfbs     string    //F_VC_SUANFBS	算法标识	VARCHAR(32)
	FDtYonghrksj   time.Time //F_DT_YONGHRKSJ	用户入口时间	DATETIME
	FNbYonghtcsc   int       //F_NB_YONGHTCSC	用户停车时长(分)	INT
	FVcZhangdms    string    //F_VC_ZHANGDMS	账单描述（给用户通知的信息）	VARCHAR(32)
	FVcMiybbh      string    //F_VC_MIYBBH	密钥版本号	VARCHAR(32)
	FVcObuyyxlh    string    //F_VC_OBUYYXLH	obu应用序列号	VARCHAR(32)
	FVcChedjyxh    string    //F_VC_CHDJYXH	车道交易序号	VARCHAR(32)【数据库字段有误】

}

//   B_JS_YUANSJYXX【原始交易消息表】
type BJsYuansjyxx struct {
	FVcBanbh       string    //F_VC_BANBH	版本号	VARCHAR(32)
	FNbXiaoxlb     int       //F_NB_XIAOXLB	消息类别	INT
	FNbXiaoxlx     int       //F_NB_XIAOXLX	消息类型	INT
	FVcFaszid      string    //F_VC_FASZID	发送者ID	VARCHAR(32)
	FVcJieszid     string    //F_VC_JIESZID	接收者ID	VARCHAR(32)
	FNbXiaoxxh     int64     //F_NB_XIAOXXH	消息序号	BIGINT
	FDtDabsj       time.Time //F_DT_DABSJ	打包时间	DATETIME
	FNbFaszt       int       //F_NB_FASZT	发送状态	INT
	FDtFassj       time.Time //F_DT_FASSJ	发送时间	DATETIME
	FNbYingdzt     int       //F_NB_YINGDZT	应答状态	INT
	FVcQingfmbr    string    //F_VC_QINGFMBR	清分目标日	VARCHAR(32)
	FVcTingccqffid string    //F_VC_TINGCCQFFID	停车场清分方ID	VARCHAR(32)
	FVcFaxfwjgid   string    //F_VC_FAXFWJGID	发行服务机构ID	VARCHAR(32)
	FNbJilsl       int       //F_NB_JILSL	记录数量	INT
	FNbZongje      int       //F_NB_ZONGJE	总金额	INT
	FVcXiaoxwjlj   string    //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)
}

//   B_JS_YUANSJYMX【原始交易明细】
type BJsYuansjymx struct {
	FVcXiaoxxh   int64     //F_VC_XIAOXXH	消息序号	BIGINT
	FNbBaonxh    int       //F_NB_BAONXH	包内序号	INT
	FDtJiaoysj   time.Time //F_DT_JIAOYSJ	交易时间	DATETIME
	FNbJine      int       //F_NB_JINE	金额	INT
	FVcDingzjyxx string    //F_VC_DINGZJYXX	定制交易信息	VARCHAR(512)
	FVcJiaoybh   string    //F_VC_JIAOYBH	交易编号	VARCHAR(128)
	FVcTingccmc  string    //F_VC_TINGCCMC	停车场名称	VARCHAR(256)
	FNbTingfsc   int       //F_NB_TINGFSC	停放时长	INT
	FNbShoufcx   int       //F_NB_SHOUFCX	收费车型	INT
	FNbSuanfbs   int       //F_NB_SUANFBS	算法标识	INT
	FNbFuwlx     int       //F_NB_FUWLX	服务类型	INT
	FVcZhangdsm  string    //F_VC_ZHANGDSM	账单说明	VARCHAR(256)
	FVcJiaoyxxxx string    //F_VC_JIAOYXXXX	交易详细信息	VARCHAR(512)
	FNbKalx      int       //F_NB_KALX	卡类型	INT
	FVcWanglbm   string    //F_VC_WANGLBM	网络编码	VARCHAR(32)
	FVcKawlbh    string    //F_VC_KAWLBH	卡物理编号	VARCHAR(32)
	FVcKancph    string    //F_VC_KANCPH	卡内车牌号	VARCHAR(32)
	FVcKajtxh    string    //F_VC_KAJYXH	卡交易序号	VARCHAR(32)
	FNbJiaoyqye  int       //F_NB_JIAOYQYE	交易前余额	INT
	FNbJiaoyhye  int       //F_NB_JIAOYHYE	交易后余额	INT
	FVcTacm      string    //F_VC_TACM	TAC码	VARCHAR(32)
	FVcjiaoybs   string    //F_VC_JIAOYBS	交易标识	VARCHAR(32)
	FVcZongdjh   string    //F_VC_ZONGDJH	终端机号	VARCHAR(32)
	FVcZongdjyxh string    //F_VC_ZONGDJYXH	终端交易序号	VARCHAR(32)
	FVcObuwlbh   string    //F_VC_OBUWLBH	OBU物理编号	VARCHAR(32)
	FVcObuzt     string    //F_VC_OBUZT	OBU状态	VARCHAR(32)
	FVcObucph    string    //F_VC_OBUNCPH	OBU内车牌号	VARCHAR(32)
}

//   B_JS_YUANSJYYDXX【原始交易应答消息】
type BJsYuansjyydxx struct {
	FVcBanbh     string    //F_VC_BANBH	版本号	VARCHAR(32)
	FNbXiaoxlb   int       //F_NB_XIAOXLB	消息类别	INT
	FNbXiaoxlx   int       //F_NB_XIAOXLX	消息类型	INT
	FVcFaszid    string    //F_VC_FASZID	发送者ID	VARCHAR(32)
	FVcJieszid   string    //F_VC_JIESZID	接收者ID	VARCHAR(32)
	FNbXiaoxxh   int64     //F_NB_XIAOXXH	消息序号	BIGINT
	FNbQuerdxxxh int64     //F_NB_QUERDXXXH	确认的消息序号	BIGINT
	FDtChulsj    time.Time //F_DT_CHULSJ	处理时间	DATETIME
	FNbZhixjg    int       //F_NB_ZHIXJG	执行结果	INT
	FVcQingfmbr  string    //F_VC_QINGFMBR	清分目标日	VARCHAR(32)
	FVcXiaoxwjlj string    //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)
}

// B_JS_JIZCLXX【记账处理消息】
type BJsJizclxx struct {
	FVcBanbh       string    //F_VC_BANBH	版本号	VARCHAR(32)
	FNbXiaoxlb     int       //F_NB_XIAOXLB	消息类别	INT
	FNbXiaoxlx     int       //F_NB_XIAOXLX	消息类型	INT
	FVcFaszid      string    //F_VC_FASZID	发送者ID	VARCHAR(32)
	FVcJieszid     string    //F_VC_JIESZID	接收者ID	VARCHAR(32)
	FNbXiaoxxh     int64     //F_NB_XIAOXXH	消息序号	BIGINT
	FDtJiessj      time.Time //F_DT_JIESSJ	接收时间	DATETIME
	FNbYuansjyxxxx int64     //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
	FNbJilsl       int       //F_NB_JILSL	记录数量	INT
	FNbZongje      int       //F_NB_ZONGJE	总金额	INT
	FNbZhengysl    int       //F_NB_ZHENGYSL	争议数量	INT
	FNbZhixjg      int       //F_NB_ZHIXJG	执行结果	INT
	FDtChulsj      time.Time //F_DT_CHULSJ	处理时间	DATETIME
	FVcXiaoxwjlj   string    //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)
}

//   B_JS_JIZCLMX【记账处理明细】
type BJsJizclmx struct {
	FNbYuansjyxxxh int64 //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
	FNbbaonxh      int   //F_NB_BAONXH	包内序号	INT
	FNbChuljg      int   //F_NB_CHULJG	处理结果	INT
}

//  B_JS_JIZCLYDXX【记账处理应答消息】
type BJsJizclydxx struct {
	FVcbanbh     string    //F_VC_BANBH	版本号	VARCHAR(32)
	FNbXiaoxlb   int       //F_NB_XIAOXLB	消息类别	INT
	FNbXiaoxlx   int       //F_NB_XIAOXLX	消息类型	INT
	FVcFaszid    string    //F_VC_FASZID	发送者ID	VARCHAR(32)
	FVcJieszid   string    //F_VC_JIESZID	接收者ID	VARCHAR(32)
	FNbXiaoxxh   int64     //F_NB_XIAOXXH	消息序号	BIGINT
	FNbQuerdxxxh int64     //F_NB_QUERDXXXH	确认的消息序号	BIGINT
	FVcChulsj    time.Time //F_DT_CHULSJ	处理时间	DATETIME
	FNbZhixjg    int       //F_NB_ZHIXJG	执行结果	INT
	FVcXiaoxwjlj string    //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)
}

//   B_JS_ZHENGYCLXX【争议交易处理消息】
type BJsZhengyclxx struct {
	FVcBanbh        string    //F_VC_BANBH	版本号	VARCHAR(32)
	FNbXiaoxlb      int       //F_NB_XIAOXLB	消息类别	INT
	FNbXiaoxlx      int       //F_NB_XIAOXLX	消息类型	INT
	FVcFaszid       string    //F_VC_FASZID	发送者ID	VARCHAR(32)
	FVcJieszid      string    //F_VC_JIESZID	接收者ID	VARCHAR(32)
	FNbXiaoxxh      int64     //F_NB_XIAOXXH	消息序号	BIGINT
	FDtJiessj       time.Time //F_DT_JIESSJ	接收时间	DATETIME
	FVcAQingffid    string    //F_VC_QINGFFID	清分方ID	VARCHAR(32)
	FVclianwzxid    string    //F_VC_LIANWZXID	联网中心ID	VARCHAR(32)
	FVcFaxfid       string    //F_VC_FAXFID	发行方ID	VARCHAR(32)
	FVcZhengyjgwjid int       //F_VC_ZHENGYJGWJID	争议结果文件ID	INT
	FDtZhengyclsj   time.Time //F_DT_ZHENGYCLSJ	争议处理时间	DATETIME
	FNbZhengysl     int       //F_NB_ZHENGYSL	争议数量	INT
	FNbQuerxyjzdzje int       //F_NB_QUERXYJZDZJE	确认需要记账的总金额	INT
}

//     B_JS_ZHENGYJYCLMX【争议交易处理明细】
type BJsZhengyjyclmx struct {
	FNbZhengyjyclxxxh int64 //F_NB_ZHENGYJYCLXXXH	争议交易处理消息序号	BIGINT
	FNbFenzxh         int   //F_NB_FENZXH	分组序号	INT
	FNbYuansjyxxxh    int64 //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
	FNbZunjlsl        int   //F_NB_ZUNJLSL	组内记录数量	INT
	FNbZunjezh        int   //F_NB_ZUNJEZH	组内金额总和	INT
	FNbYuansbnxh      int   //F_NB_YUANSBNXH	原始包内序号	INT
	FNbChuljg         int   //F_NB_CHULJG	处理结果	INT
}

//    B_JS_ZHENGYJYCLYDXX【争议交易处理应答消息】
type BJsZhengyclygxx struct {
	FVcBanbh     string    //F_VC_BANBH	版本号	VARCHAR(32)
	FNbXiaoxlb   int       //F_NB_XIAOXLB	消息类别	INT
	FNbXiaoxlx   int       //F_NB_XIAOXLX	消息类型	INT
	FVcFaszid    string    //F_VC_FASZID	发送者ID	VARCHAR(32)
	FVcJieszid   string    //F_VC_JIESZID	接收者ID	VARCHAR(32)
	FNbXiaoxxh   int64     //F_NB_XIAOXXH	消息序号	BIGINT
	FNbQuerdxxxh int64     //F_NB_QUERDXXXH	确认的消息序号	BIGINT
	FVcChulsj    time.Time //F_DT_CHULSJ	处理时间	DATETIME
	FNbZhixjg    int       //F_NB_ZHIXJG	执行结果	INT

}

//    B_JS_QINGFTJXX【清分统计消息】
type BJsQingftjxx struct {
	FVcBanbh         string    //F_VC_BANBH	版本号	VARCHAR(32)
	FNbXiaoxlb       int       //F_NB_XIAOXLB	消息类别	INT
	FNbXiaoxlx       int       //F_NB_XIAOXLX	消息类型	INT
	FVcFaszid        string    //F_VC_FASZID	发送者ID	VARCHAR(32)
	FVcJieszid       string    //F_VC_JIESZID	接收者ID	VARCHAR(32)
	FNbXiaoxxh       int64     //F_NB_XIAOXXH	消息序号	BIGINT
	FDtJiessj        time.Time //F_DT_JIESSJ	接收时间	DATETIME
	FVcQingfmbr      time.Time //F_VC_QINGFMBR	清分目标日	DATE
	FNbQingfzje      int       //F_NB_QINGFZJE	清分总金额	INT
	FNbQingfsl       int       //F_NB_QINGFSL	清分数量	INT
	FDtQingftjclsj   time.Time //F_DT_QINGFTJCLSJ	清分统计处理时间	DATETIME
	FNbYuansjysl     int       //F_NB_YUANSJYSL	原始包交易数量	INT
	FNbZhengycljgbsl int       //F_NB_ZHENGYCLJGBSL	争议处理结果包数量	INT
	FVcXiaoxwjlj     string    //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)
}

//    B_JS_QINGFTONGJIMX【清分统计明细】
type BJsQingftongjimx struct {
	FNbQingftjxxxh    int64  //F_NB_QINGFTJXXXH	清分统计消息序号	BIGINT
	FNbFenzxh         int    //F_NB_FENZXH	分组序号	INT
	FVcTongxbzxxtid   string //F_VC_TONGXBZXXTID	通行宝中心系统ID	VARCHAR(32)
	FNbYuansjyxxxh    int64  //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
	FNbZhengycljgwjid int    //F_NB_ZHENGYCLJGWJID	争议处理结果文件ID	INT
}

//   B_JS_QINGFTJXXYD【清分统计消息应答】
type BJsQingftjxxyd struct {
	FVcBanbh     string    //F_VC_BANBH	版本号	VARCHAR(32)
	FNbXiaoxlb   int       //F_NB_XIAOXLB	消息类别	INT
	FNbXiaoxlx   int       //F_NB_XIAOXLX	消息类型	INT
	FVcFaszid    string    //F_VC_FASZID	发送者ID	VARCHAR(32)
	FVcJieszid   string    //F_VC_JIESZID	接收者ID	VARCHAR(32)
	FNbXiaoxxh   int64     //F_NB_XIAOXXH	消息序号	BIGINT
	FNbQuerdxxxh int64     //F_NB_QUERDXXXH	确认的消息序号	BIGINT
	FDtChulsj    time.Time //F_DT_CHULSJ	处理时间	DATETIME
	FNbZhixjg    int       //F_NB_ZHIXJG	执行结果	INT
	FVcXiaoxwjlj string    //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)

}
