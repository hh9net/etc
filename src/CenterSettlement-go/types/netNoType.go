package types

//卡网络号
const (
	JS_NETWORK  string = "3201" // 江苏
	SH_NETWORK  string = "3101" // 上海
	ZJ_NETWORK  string = "3301" // 浙江
	AH_NETWORK  string = "3401" // 安徽
	FJ_NETWORK  string = "3501" // 福建
	JX_NETWORK  string = "3601" // 江西
	SD_NETWORK  string = "3701" // 山东
	SD_NETWORK2 string = "3702" // 山东

	/* 华北区路网代码定义*/
	BJ_NETWORK  string = "1101" // 北京
	TJ_NETWORK  string = "1201" // 天津
	HEB_NETWORK string = "1301" // 河北
	SX_NETWORK  string = "1401" // 山西
	NM_NETWORK  string = "1501" // 内蒙古

	/* 东北区路网代码定义*/
	LN_NETWORK  string = "2101" // 辽宁
	JL_NETWORK  string = "2201" // 吉林
	HLJ_NETWORK string = "2301" // 黑龙江

	/*华中、华南区路网代码定义*/
	HEN_NETWORK  string = "4101" // 河南
	HUB_NETWORK  string = "4201" // 湖北
	HUB_NETWORK2 string = "4202" // 湖北

	HUN_NETWORK  string = "4301" // 湖南
	GD_NETWORK   string = "4401" // 广东
	GX_NETWORK   string = "4501" // 广西
	HAIN_NETWORK string = "4601" // 海南

	/*西南区路网代码定义*/
	CQ_NETWORK  string = "5001" // 重庆
	SC_NETWORK  string = "5101" // 四川
	SC_NETWORK2 string = "5102" // 四川
	SC_NETWORK3 string = "5103" // 四川
	SC_NETWORK4 string = "5104" // 四川
	SC_NETWORK5 string = "5105" // 四川
	GZ_NETWORK  string = "5201" // 贵州
	YN_NETWORK  string = "5301" // 云南
	XZ_NETWORK  string = "5401" // 西藏

	/*西北区路网代码定义*/
	SHANXI_NETWORK  string = "6101" // 陕西
	SHANXI_NETWORK2 string = "6102" // 陕西
	SHANXI_NETWORK3 string = "6103" // 陕西
	SHANXI_NETWORK4 string = "6104" // 陕西
	SHANXI_NETWORK5 string = "6105" // 陕西
	SHANXI_NETWORK6 string = "6106" // 陕西
	SHANXI_NETWORK7 string = "6107" // 陕西

	GS_NETWORK string = "6201" // 甘肃
	QH_NETWORK string = "6301" // 青海
	NX_NETWORK string = "6401" // 宁夏
	XJ_NETWORK string = "6501" // 新疆

	ARMY_CARDNETWORK string = "501" // 军车卡的网络编号

	NETWORK_CODE_CNT = 40 // 联网省份数量（发行商）

	PRECARD                    = 22   //储值卡
	CREDITCARD                 = 23   //记账卡
	TRANSTYPE           string = "09" //交易标识
	SERVICETYPE         int    = 2    //原始交易服务类型
	ALGORITHMIDENTIFIER int    = 1    //algorithmIdentifier"` //算法标识 1-3DEX  2-SM4

)

//各省市发行方(除江苏)
var Gl_network = []string{
	SH_NETWORK, ZJ_NETWORK, AH_NETWORK, FJ_NETWORK, JX_NETWORK,
	BJ_NETWORK, TJ_NETWORK, HEB_NETWORK, SHANXI_NETWORK, SHANXI_NETWORK2, SHANXI_NETWORK3, SHANXI_NETWORK4, SHANXI_NETWORK5, SHANXI_NETWORK6, SHANXI_NETWORK7, LN_NETWORK, SD_NETWORK, SD_NETWORK2, HUN_NETWORK, SX_NETWORK,
	HEN_NETWORK, HUB_NETWORK2, GD_NETWORK, GZ_NETWORK, ARMY_CARDNETWORK,
	JL_NETWORK, CQ_NETWORK, SC_NETWORK, SC_NETWORK2, SC_NETWORK3, SC_NETWORK4, SC_NETWORK5, YN_NETWORK,
	GS_NETWORK, QH_NETWORK, NX_NETWORK,
	NM_NETWORK, HLJ_NETWORK, GX_NETWORK, XJ_NETWORK}
