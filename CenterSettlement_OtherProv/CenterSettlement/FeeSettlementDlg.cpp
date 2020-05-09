
// FeeSettlementDlg.cpp : 实现文件
//

#include "stdafx.h"
#include "FeeSettlement.h"
#include "FeeSettlementDlg.h"
#include "afxdialogex.h"
#include "time.h"


#include <stdio.h>
#include <WinSock2.h>
#include <WS2tcpip.h>
#include <windows.h>
#include <iostream> 
#include <fstream>
#include <string>
#include <io.h>
#include <direct.h> 
#include <vector>
#include <stdarg.h>
#include <winbase.h>
#include <tchar.h>
#include <initializer_list>


#include "ocilib.h"
#include "md5.h"
#include "tinyxml2.h"
#include "Lz77.h"
#include "iconv.h"

//自定义时间
#define ONE_SEC		1000
#define ONE_MIN		(ONE_SEC * 60)
#define THREE_MIN	(ONE_MIN * 3)
#define FIVE_MIN	(ONE_MIN * 5)
#define HALF_HOUR   (ONE_MIN*30)

//定义ini配置文件
#define  CONFIG_NUM	          "config"
#define  KEYRECV              "num_of_recv" 
#define  KEYSEND              "num_of_send"
#define  KEYDEAL              "num_of_deal"
#define  DB_INI_PATH          ".\\db.ini" 
#define  LOG_DIR              "LogDir"

//oracle数据库连接
#define DB_USER		"admin"
#define DB_PWD		"123"
//#define DB_SERVICE	"221.226.132.158:1521/TXB"
#define DB_SERVICE	"192.168.4.144:1521/TXB"

//socket通讯
#define CLIENT_CONNECT_PORT	 6750
#define SERVER_LISTEN_PORT   6751

//最大交易数目，打包文件内最大记录条数
#define TRANS_MAX_COUNT 100

//记账结算xml文件大小，200K
#define JZJS_XML_FILE_LEN			(1024 * 200)

//记账结算压缩文件大小，120K
#define JZJS_COMPRESS_FILE_LEN		(1024 * 120)

#define ACCESS(fileName,accessMode) _access(fileName,accessMode)
#define MKDIR(path) _mkdir(path)
#define GETCWD(path) _getcwd(path, sizeof(path))

//消息句柄
#define MYMSG1 WM_USER+200
#define MYMSG2 WM_USER+201
#define WM_SHOWTASK WM_USER+202

#ifdef _DEBUG
#define new DEBUG_NEW
#endif

//加载lib库
#pragma comment(lib, "ociliba.lib")
#pragma comment(lib, "ws2_32.lib")
#pragma comment(lib, "libiconv.lib")

//卡网络号
#define JS_NETWORK				   3201		    // 江苏
#define SH_NETWORK				   3101		    // 上海
#define ZJ_NETWORK				   3301		    // 浙江
#define AH_NETWORK				   3401		    // 安徽
#define FJ_NETWORK				   3501		    // 福建
#define JX_NETWORK				   3601		    // 江西
#define SD_NETWORK                 3701         // 山东
#define SD_NETWORK2                3702         // 山东

/* 华北区路网代码定义*/
#define BJ_NETWORK                 1101         // 北京
#define TJ_NETWORK                 1201         // 天津
#define HEB_NETWORK                1301         // 河北
#define SX_NETWORK                 1401         // 山西
#define NM_NETWORK                 1501         // 内蒙古

/* 东北区路网代码定义*/
#define LN_NETWORK                 2101         // 辽宁
#define JL_NETWORK                 2201         // 吉林
#define HLJ_NETWORK                2301         // 黑龙江

/*华中、华南区路网代码定义*/
#define HEN_NETWORK                4101         // 河南
#define HUB_NETWORK                4201         // 湖北
#define HUB_NETWORK2               4202         // 湖北

#define HUN_NETWORK                4301         // 湖南
#define GD_NETWORK                 4401         // 广东
#define GX_NETWORK                 4501         // 广西
#define HAIN_NETWORK               4601         // 海南

/*西南区路网代码定义*/
#define CQ_NETWORK                 5001         // 重庆
#define SC_NETWORK                 5101         // 四川
#define SC_NETWORK2                5102         // 四川
#define SC_NETWORK3                5103         // 四川
#define SC_NETWORK4                5104         // 四川
#define SC_NETWORK5                5105         // 四川
#define GZ_NETWORK                 5201         // 贵州
#define YN_NETWORK                 5301         // 云南
#define XZ_NETWORK                 5401         // 西藏

/*西北区路网代码定义*/
#define SHANXI_NETWORK             6101          // 陕西
#define SHANXI_NETWORK2            6102          // 陕西
#define SHANXI_NETWORK3            6103          // 陕西
#define SHANXI_NETWORK4            6104          // 陕西
#define SHANXI_NETWORK5            6105          // 陕西
#define SHANXI_NETWORK6            6106          // 陕西
#define SHANXI_NETWORK7            6107          // 陕西

#define GS_NETWORK                 6201          // 甘肃
#define QH_NETWORK                 6301          // 青海
#define NX_NETWORK                 6401          // 宁夏
#define XJ_NETWORK                 6501          // 新疆

#define ARMY_CARDNETWORK	      	501			 // 军车卡的网络编号

#define NETWORK_CODE_CNT	        40			 // 联网省份数量（发行商）

#define PRECARD                     22           //储值卡
#define CREDITCARD                  23           //记账卡

//各省市发行方(除江苏)
static WORD gl_network[NETWORK_CODE_CNT] = {
	SH_NETWORK,ZJ_NETWORK,AH_NETWORK,FJ_NETWORK,JX_NETWORK,\
	BJ_NETWORK,TJ_NETWORK,HEB_NETWORK,SHANXI_NETWORK,SHANXI_NETWORK2,SHANXI_NETWORK3,SHANXI_NETWORK4,SHANXI_NETWORK5,SHANXI_NETWORK6,SHANXI_NETWORK7,LN_NETWORK,SD_NETWORK,SD_NETWORK2,HUN_NETWORK,SX_NETWORK,\
	HEN_NETWORK,HUB_NETWORK2,GD_NETWORK,GZ_NETWORK,ARMY_CARDNETWORK,\
	JL_NETWORK,CQ_NETWORK,SC_NETWORK,SC_NETWORK2,SC_NETWORK3,SC_NETWORK4,SC_NETWORK5,YN_NETWORK,\
	GS_NETWORK,QH_NETWORK,NX_NETWORK,\
	NM_NETWORK,HLJ_NETWORK,GX_NETWORK,XJ_NETWORK };

//////////////////////////////////////////////////////////////////////////
// 程序运行参数配置 
//////////////////////////////////////////////////////////////////////////
struct RUN_CONFIGE
{
	int settlerun;				//结算记账线程运行标记
	int sendrun;				//发送文件线程运行标记
	int recvrun;                //接收文件线程运行标记
	int analyzerun;             //解析xml文件线程运行标记
	char serverip[50];          //服务器ip地址

	char dbservice[50];			//通行宝数据库服务
	char dbuser[15];			//通行宝数据库用户
	char dbpwd[15];				//通行宝数据库密码

	char dirpath[150];			//程序当前运行路径
};

//////////////////////////////////////////////////////////////////////////
// 与西公所通讯消息-socket消息 
//////////////////////////////////////////////////////////////////////////
struct SEND_STRU 
{
	char messageid[20 + 1];							//消息报文顺序号 （SSSSSSS） 20字节
	char xmllen[6 + 1];								//压缩后的XML消息长度(CCCCCC) 6字节
	char md5str[32 + 1];							//32字节MD5 (16进制显示字符串)
	char xml_msg[JZJS_COMPRESS_FILE_LEN + 1];		//二进制压缩后的XML消息包( AAAAAAAAAAAAAAA…)
};

struct REPLY_STRU
{
	char messageid[20 + 1];				//消息报文顺序号 （SSSSSSS） 20字节
	char result;
};

//////////////////////////////////////////////////////////////////////////
// 与西公所通讯消息-xml文件标准消息头 
//////////////////////////////////////////////////////////////////////////
struct ETC_XML_HEADER
{
	int version;				//本次版本号统一使用0x00010000，表示版本1.0
	int messageclass;			//使用所接收消息的MessageClass
	int messagetype;			//使用与其对应的Response值
	char senderid[16 + 1];		//当前参与方Id
	char receiverid[16 + 1];	//准备接收确认消息的参与方Id
	DWORD messageid;			//消息序号，从1开始，逐1递增 ，8字节
};

//////////////////////////////////////////////////////////////////////////
// 发送的结算记账消息体 
//////////////////////////////////////////////////////////////////////////
struct ETC_XML_SEND_BODY_DESC
{
	int contenttype;				//1为记账包；9为调账包
	char cleartargetdate[10 + 1];	//清分目标日
	char serviceproviderid[16 + 1];	//通行宝中心Id，表示消息包中的交易是由哪个公路收费方产生的。
	char issuerid[16 + 1];			//发行方Id，表示产生交易记录的电子介质所属的发行方。			
	DWORD messageid;				//消息序号，从1开始，逐1递增 ，8字节
	int	count;						//记录条数
	double amount;					//交易总金额
};

struct ETC_XML_SEND_TRANSACTION_SERVICE
{
	int servicetype;				//交易服务类型
	char description[100];			//交易描述
	char detail[500];				//交易详细信息
};

struct ETC_XML_SEND_TRANSACTION_ICCARD
{
	int cardtype;					//卡类型
	char netno[4 + 1];				//网络编码
	char cardid[16 + 1];			//卡号
	char license[20];				//车牌 
	double prebalance;				//交易前余额
	double postbalance;				//交易后余额
};

struct ETC_XML_SEND_TRANSACTION_VALIDATION
{
	char tac[8 + 1];					//tac码
	char transtype[2 + 1];				//交易标识
	char terminalno[12 + 1];			//psam卡序号
	char terminaltransno[8 + 1];		//终端交易序号
};

struct ETC_XML_SEND_TRANSACTION_OBU
{
	char netno[4 + 1];					//网络编码
	char obuid[20];					    //OBU物理编号
	char obestate[8];				    //OBU状态
	char license[20];					//OBU车牌
};

struct ETC_XML_SEND_TRANSACTION
{
	int transid;									//是由通行宝中心产生的该包内顺序Id，从1开始递增。在通行宝中心、联网中心、发行方三方的交易通讯过程中均采用此Id表示包内唯一的交易记录。
	char time[20];									//交易时间
	double fee;										//交易金额
	//20171117 begin 
	long long int stationId;                        //收费站编号
	int lane;                                       //车道号
	//20171117 end
	ETC_XML_SEND_TRANSACTION_SERVICE service;		//服务信息
	ETC_XML_SEND_TRANSACTION_ICCARD iccard;			//IC卡相关信息
	ETC_XML_SEND_TRANSACTION_VALIDATION validation;	//校验相关信息
	ETC_XML_SEND_TRANSACTION_OBU obu;				//OBU信息
	char customizeddata[200];						//包含所有校验信息在内的原始信息
};

//////////////////////////////////////////////////////////////////////////
// 应答的结算记账消息体 
//////////////////////////////////////////////////////////////////////////
struct ETC_XML_RECV_BODY_DESC
{
	int contenttype;				//1为记账包；9为调账包
	DWORD messsageid;				//消息序号
	char processtime[20];			//处理时间
	int result;						//执行结果
	char description[100];			//执行结果描述
};

using namespace std;
using namespace tinyxml2;


UINT SettleWorkerThread(LPVOID lparam);
UINT SendFileWorkerThread(LPVOID lparam);

UINT RecvCompressFileThread(LPVOID lparam);
UINT AnalyzeXmlFileThread(LPVOID lparam);

// CFeeSettlementDlg 对话框

//////////////////////////////////////////////////////////////////////////
// 全局参数变量，存储程序运行各类参数
//////////////////////////////////////////////////////////////////////////
RUN_CONFIGE g_config;
CString g_msg1;
CString g_msg2;

CFeeSettlementDlg::CFeeSettlementDlg(CWnd* pParent /*=NULL*/)
	: CDialogEx(IDD_FEESETTLEMENT_DIALOG, pParent)
{
	m_hIcon = AfxGetApp()->LoadIcon(IDR_MAINFRAME);
}

void CFeeSettlementDlg::DoDataExchange(CDataExchange* pDX)
{
	CDialogEx::DoDataExchange(pDX);
	DDX_Control(pDX, IDC_LIST_MSG1, m_lstMsg1);
	DDX_Control(pDX, IDC_LIST_MSG2, m_lstMsg2);
}

void CFeeSettlementDlg::ToTray()
{
	NOTIFYICONDATA nid;
	nid.cbSize = (DWORD)sizeof(NOTIFYICONDATA);
	nid.hWnd = this->m_hWnd;
	nid.uID = IDR_MAINFRAME;
	nid.uFlags = NIF_ICON | NIF_MESSAGE | NIF_TIP;
	nid.uCallbackMessage = WM_SHOWTASK;//自定义的消息名称 
	nid.hIcon = LoadIcon(AfxGetInstanceHandle(), MAKEINTRESOURCE(IDR_MAINFRAME));
	lstrcpy(nid.szTip, _T("通行宝"));//信息提示条 
	Shell_NotifyIcon(NIM_ADD, &nid);//在托盘区添加图标
	ShowWindow(SW_HIDE);//隐藏主窗口
}

BEGIN_MESSAGE_MAP(CFeeSettlementDlg, CDialogEx)
	ON_WM_PAINT()
	ON_WM_QUERYDRAGICON()
	ON_MESSAGE(MYMSG1, OnMsg1Back)
	ON_MESSAGE(MYMSG2, OnMsg2Back)
	ON_MESSAGE(WM_SHOWTASK, OnShowTask)
	ON_BN_CLICKED(IDC_BTN_MIN, &CFeeSettlementDlg::OnBnClickedBtnMin)
	ON_BN_CLICKED(IDCANCEL, &CFeeSettlementDlg::OnBnClickedCancel)
END_MESSAGE_MAP()


// CFeeSettlementDlg 消息处理程序

BOOL CFeeSettlementDlg::OnInitDialog()
{
	CDialogEx::OnInitDialog();

	// 设置此对话框的图标。  当应用程序主窗口不是对话框时，框架将自动
	//  执行此操作
	SetIcon(m_hIcon, TRUE);			// 设置大图标
	SetIcon(m_hIcon, FALSE);		// 设置小图标

	// TODO: 在此添加额外的初始化代码
	//////////////////////////////////////////////////////////////////////////
	// 控制程序只有一个实例
	//////////////////////////////////////////////////////////////////////////
//	int result = OCI_Initialize(NULL, "D:\\Oracle\\Instant Client\\bin", OCI_ENV_THREADED);
	int result = OCI_Initialize(NULL, NULL, OCI_ENV_THREADED);
	if (!result)
	{
		return FALSE;
	}

	HANDLE hMutex = ::CreateMutex(NULL, TRUE, _T("FeeSettlement_OtherProv"));  
	if (hMutex != NULL)
	{
		if (GetLastError() == ERROR_ALREADY_EXISTS)
		{
			if (1 == MessageBox(_T("已经有一个程序运行"), _T("通行宝-数据结算组件"), MB_ICONEXCLAMATION | MB_OK))
			{
				OCI_Cleanup();
				exit(0);
			}
		}
	}

  //////////////////////////////////////////////////////////////////////////
	// 通过配置文件来控制线程是否运行
	//////////////////////////////////////////////////////////////////////////
	//读取配置文件
	if (!this->GetConfig())
	{
		if (1 == MessageBox(_T("读取程序运行配置文件出错"), _T("通行宝-数据结算组件"), MB_ICONEXCLAMATION | MB_OK))
		{
			exit(0);
		}
	}

	// 判断配置项
	
	if (g_config.settlerun)
	{
		//数据打包线程
		::AfxBeginThread(SettleWorkerThread, this);
	}

	//发送文件线程
	if (g_config.sendrun)
	{
		//数据包发送线程
		::AfxBeginThread(SendFileWorkerThread, this);
	}

	if (g_config.recvrun)
	{
	    //接收数据线程
		::AfxBeginThread(RecvCompressFileThread, this);
	}

	if (g_config.analyzerun)
	{
		//解析xml文件线程
		::AfxBeginThread(AnalyzeXmlFileThread, this);
	}

	return TRUE;  // 除非将焦点设置到控件，否则返回 TRUE
}

// 如果向对话框添加最小化按钮，则需要下面的代码
//  来绘制该图标。  对于使用文档/视图模型的 MFC 应用程序，
//  这将由框架自动完成。

//从托盘区恢复界面的函数
LRESULT CFeeSettlementDlg::OnShowTask(WPARAM wParam, LPARAM lParam)
{
	if (wParam != IDR_MAINFRAME)
		return 1;
	switch (lParam)
	{
	case WM_RBUTTONUP://右键起来时弹出快捷菜单，这里只有一个“关闭” 
	{
		LPPOINT lpoint = new tagPOINT;
		::GetCursorPos(lpoint);//得到鼠标位置
		CMenu menu;
		menu.CreatePopupMenu();//声明一个弹出式菜单
		menu.AppendMenu(MF_STRING, WM_DESTROY, _T("关闭"));//增加菜单项“关闭”，点击则发送消息WM_DESTROY给主窗口（已隐藏），将程序结束。               
		menu.TrackPopupMenu(TPM_LEFTALIGN, lpoint->x, lpoint->y, this);//确定弹出式菜单的位置                 
		HMENU hmenu = menu.Detach();
		menu.DestroyMenu();//资源回收 
		DeleteTray();
		delete lpoint;
	}  break;
	case WM_LBUTTONDBLCLK://双击左键的处理 
	{
		this->ShowWindow(SW_SHOW);//简单的显示主窗口
		DeleteTray();
	}  break;
	default:   break;
	}
	return 0;
}

//删除托盘图标的函数
void CFeeSettlementDlg::DeleteTray()
{
	NOTIFYICONDATA nid;
	nid.cbSize = (DWORD)sizeof(NOTIFYICONDATA);
	nid.hWnd = this->m_hWnd;
	nid.uID = IDR_MAINFRAME;
	nid.uFlags = NIF_ICON | NIF_MESSAGE | NIF_TIP;
	nid.uCallbackMessage = WM_SHOWTASK;//自定义的消息名称 
	nid.hIcon = LoadIcon(AfxGetInstanceHandle(), MAKEINTRESOURCE(IDR_MAINFRAME));
	lstrcpy(nid.szTip, _T("通行宝"));//信息提示条  
	Shell_NotifyIcon(NIM_DELETE, &nid);//在托盘区删除图标 
}

void CFeeSettlementDlg::OnPaint()
{
	if (IsIconic())
	{
		CPaintDC dc(this); // 用于绘制的设备上下文

		SendMessage(WM_ICONERASEBKGND, reinterpret_cast<WPARAM>(dc.GetSafeHdc()), 0);

		// 使图标在工作区矩形中居中
		int cxIcon = GetSystemMetrics(SM_CXICON);
		int cyIcon = GetSystemMetrics(SM_CYICON);
		CRect rect;
		GetClientRect(&rect);
		int x = (rect.Width() - cxIcon + 1) / 2;
		int y = (rect.Height() - cyIcon + 1) / 2;

		// 绘制图标
		dc.DrawIcon(x, y, m_hIcon);
	}
	else
	{
		CDialogEx::OnPaint();
	}
}

//当用户拖动最小化窗口时系统调用此函数取得光标
//显示。
HCURSOR CFeeSettlementDlg::OnQueryDragIcon()
{
	return static_cast<HCURSOR>(m_hIcon);
}

/*************************************************
函数名称:
函数描述:
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
void CFeeSettlementDlg::OnBnClickedBtnMin()
{
	// TODO: 在此添加控件通知处理程序代码
	//CWnd::ShowWindow(SW_SHOWMINIMIZED);
	//this->ShowWindow(HIDE_WINDOW);
	ToTray();//最小化到系统托盘
}

/*************************************************
函数名称:
函数描述:
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
void CFeeSettlementDlg::OnBnClickedCancel()
{
	// TODO: 在此添加控件通知处理程序代码
	if (1 == MessageBox(_T("确认要退出本程序吗"), _T("通行宝-数据结算组件"), MB_ICONEXCLAMATION | MB_OKCANCEL))
	{
		exit(0);
	}
}

/*************************************************
函数名称:
函数描述:取当前日期时间
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
CString CFeeSettlementDlg::GetNow()
{
	CTime time = CTime::GetCurrentTime();
	CString localtm = time.Format(_T("%Y年%m月%d日 %H:%M:%S "));
	return localtm;
}

/*************************************************
函数名称:
函数描述:生成日志文件
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
BOOL CreateLogFile(CString logstr, int logtype)
{
	char LogFileName[20] = { 0 };
	char Logdir[200] = { 0 };
	char LogFileDir[300] = { 0 };
	char LogFile[300] = { 0 };

	SYSTEMTIME st;
	GetLocalTime(&st);
	sprintf_s(LogFileName, sizeof(LogFileName),"%04d%02d%02d", st.wYear, st.wMonth, st.wDay);
	sprintf_s(LogFileDir, sizeof(LogFileDir),"%s\\Logfile", g_config.dirpath);
	if (ACCESS(LogFileDir, 0) != 0)
	{
		//路径不存在，就创建
		MKDIR(LogFileDir);
	}
	if (logtype == 1)
	{
		sprintf_s(LogFile, sizeof(LogFile),"%s\\%s_settle.log", LogFileDir, LogFileName);
	}
	else if (logtype == 2)
	{
		sprintf_s(LogFile, sizeof(LogFile), "%s\\%s_send.log", LogFileDir, LogFileName);
	}
	else if (logtype == 3)
	{
		sprintf_s(LogFile, sizeof(LogFile), "%s\\%s_recv.log", LogFileDir, LogFileName);
	}
	else if (logtype == 4)
	{
		sprintf_s(LogFile, sizeof(LogFile), "%s\\%s_analyze.log", LogFileDir, LogFileName);
	}
	else if (logtype == 5)
	{
		sprintf_s(LogFile, sizeof(LogFile), "%s\\%s_err.log", LogFileDir, LogFileName);
	}
	else if (logtype == 6)
	{
		sprintf_s(LogFile, sizeof(LogFile), "%s\\%s_ThisProvDebit.log", LogFileDir, LogFileName);
	}

	CTime time = CTime::GetCurrentTime();
	CString localtm = time.Format(_T("%Y-%m-%d %H:%M:%S "));

	//CString logstr_1 = localtm + logstr + "\n";
	//CFile file(LogFile, CFile::modeCreate | CFile::modeWrite);
	//file.Write(logstr_1, logstr_1.GetLength());
	//file.Flush();
	//file.Close();

	FILE *fp;
	if ((fopen_s(&fp, LogFile, "a+")) != 0)
	{
		return FALSE;
	}
	CString logstr_1 = localtm + logstr + "\n";
	//显示中文必须加下面的设置函数接口，否则中文显示乱码
	_wsetlocale(0, L"chs");
	fwprintf_s(fp,logstr_1);
	fclose(fp);
	return TRUE;
}


/*************************************************
函数名称:
函数描述:生成日志文件
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
void WriteLog(const char *fm, ...)
{

	char LogFileName[20] = { 0 };
	char Logdir[200] = { 0 };
	char LogFileDir[300] = { 0 };
	char LogFile[300] = { 0 };
	char LocalTime[50] = { 0 };

	int buflen = 5120;
    char buf[5120];
	int i = 0;
	memset(buf, 0, buflen);
	va_list args;
	va_start(args, fm);
	_vsnprintf_s(buf, buflen, fm, args);
	va_end(args);

	printf("%s\n", buf);

	SYSTEMTIME st;
	GetLocalTime(&st);
	sprintf_s(LogFileName, sizeof(LogFileName), "%04d%02d%02d", st.wYear, st.wMonth, st.wDay);
	sprintf_s(LogFileDir, sizeof(LogFileDir), "%s\\Logfile", g_config.dirpath);
	if (ACCESS(LogFileDir, 0) != 0)
	{
		//路径不存在，就创建
		MKDIR(LogFileDir);
	}

	sprintf_s(LogFile, sizeof(LogFile), "%s\\%s_settle.log", LogFileDir, LogFileName);

	sprintf_s(LocalTime, sizeof(LocalTime), "%04d-%02d-%02d %02d:%02d:%02d", st.wYear, st.wMonth, st.wDay ,st.wHour, st.wMinute, st.wSecond);

	FILE* logfile = NULL;
	if ((fopen_s(&logfile, LogFile, "a+")) != 0)
	{
		return;
	}

	const char* pTemp = LocalTime;
	fwrite(pTemp, 1, strlen(pTemp), logfile);
	fwrite(" ", 1, 1, logfile);
	//内容
	fwrite(buf, 1, strlen(buf), logfile);
	fwrite(" \r\n", 1, 3, logfile);
	fclose(logfile);

}

/*************************************************
函数名称:
函数描述:读取配置文件
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
int u2g(char *inbuf, int inlen, char *outbuf, int outlen);
BOOL CFeeSettlementDlg::GetConfig()
{
	memset((void*)&g_config, 0, sizeof(RUN_CONFIGE));

	//获取当前路径
	GETCWD(g_config.dirpath);

	//查找配置文件是否存在
	char configfile[120] = { 0 };
	sprintf_s(configfile, sizeof(configfile),"%s\\config.xml", g_config.dirpath);
	if (ACCESS(configfile, 0) != 0)
	{
		return FALSE;
	}

	tinyxml2::XMLDocument doc;
	doc.LoadFile(configfile);
	XMLElement *root = doc.RootElement();
	//配置文件检查
	if (strlen(root->FirstChildElement("WorkerThread")->FirstChildElement("Settle")->GetText()) == 0)
	{
		CreateLogFile(_T("config.xml 错误! 错误节点 ->Settle"), 5);
		return FALSE;
	}
	if (strlen(root->FirstChildElement("WorkerThread")->FirstChildElement("SendFile")->GetText()) == 0)
	{
		CreateLogFile(_T("config.xml 出错! 错误节点 ->SendFile"), 5);
		return FALSE;
	}
	if (strlen(root->FirstChildElement("WorkerThread")->FirstChildElement("RecvFile")->GetText()) == 0)
	{
		CreateLogFile(_T("config.xml 出错! 错误节点 ->RecvFile"), 5);
		return FALSE;
	}
	if (strlen(root->FirstChildElement("WorkerThread")->FirstChildElement("AnalyzeFile")->GetText()) == 0)
	{
		CreateLogFile(_T("config.xml 出错! 错误节点 ->AnalyzeFile"), 5);
		return FALSE;
	}
	if (strlen(root->FirstChildElement("WorkerThread")->FirstChildElement("ServerIp")->GetText()) == 0)
	{
		CreateLogFile(_T("config.xml 出错! 错误节点 ->ServerIp"), 5);
		return FALSE;
	}
	if (strlen(root->FirstChildElement("TxbDb")->FirstChildElement("DbService")->GetText()) == 0)
	{
		CreateLogFile(_T("config.xml 出错! 错误节点 ->DbService"), 5);
		return FALSE;
	}
	if (strlen(root->FirstChildElement("TxbDb")->FirstChildElement("DbUser")->GetText()) == 0)
	{
		CreateLogFile(_T("config.xml 出错! 错误节点 ->DbUser"), 5);
		return FALSE;
	}
	if (strlen(root->FirstChildElement("TxbDb")->FirstChildElement("DbPwd")->GetText()) == 0)
	{
		CreateLogFile(_T("config.xml 出错! 错误节点 ->DbPwd"), 5);
		return FALSE;
	}

	//配置文件读取
	g_config.settlerun = atoi(root->FirstChildElement("WorkerThread")->FirstChildElement("Settle")->GetText());
	g_config.sendrun = atoi(root->FirstChildElement("WorkerThread")->FirstChildElement("SendFile")->GetText());
	g_config.recvrun = atoi(root->FirstChildElement("WorkerThread")->FirstChildElement("RecvFile")->GetText());
	g_config.analyzerun = atoi(root->FirstChildElement("WorkerThread")->FirstChildElement("AnalyzeFile")->GetText());
	sprintf_s(g_config.serverip, sizeof(g_config.serverip),"%s", root->FirstChildElement("WorkerThread")->FirstChildElement("ServerIp")->GetText());
	//获取通行宝数据库连接配置
	sprintf_s(g_config.dbservice, sizeof(g_config.dbservice),"%s", root->FirstChildElement("TxbDb")->FirstChildElement("DbService")->GetText());
	sprintf_s(g_config.dbuser, sizeof(g_config.dbuser),"%s", root->FirstChildElement("TxbDb")->FirstChildElement("DbUser")->GetText());
	sprintf_s(g_config.dbpwd, sizeof(g_config.dbpwd),"%s", root->FirstChildElement("TxbDb")->FirstChildElement("DbPwd")->GetText());

	doc.Clear();

	return TRUE;
}



/*************************************************
函数名称:
函数描述:自定义消息
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
LRESULT CFeeSettlementDlg::OnMsg1Back(WPARAM wParam, LPARAM lParam)
{
	CString* str = (CString*)wParam;
	CString newstr = GetNow() + *str;
	m_lstMsg1.AddString(newstr);

	//function(CString) >> 
	//CreateLogFile(newstr);

	//定位最后一行
	int counts = m_lstMsg1.GetCount();
	if (counts > 0)
		m_lstMsg1.SetCurSel(counts - 1);

	//清空
	if (counts > 1000)
		m_lstMsg1.ResetContent();

	return 0;
}

LRESULT CFeeSettlementDlg::OnMsg2Back(WPARAM wParam, LPARAM lParam)
{
	CString* str = (CString*)wParam;
	CString newstr = GetNow() + *str;
	m_lstMsg2.AddString(newstr);

	//定位最后一行
	int counts = m_lstMsg2.GetCount();
	if (counts > 0)
		m_lstMsg2.SetCurSel(counts - 1);

	//清空
	if (counts > 1000)
		m_lstMsg2.ResetContent();

	return 0;
}

/*************************************************
函数名称:
函数描述:取当前日期
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
void GetDate(char* dt)
{
	CTime time = CTime::GetCurrentTime();
	sprintf_s(dt, 11, "%d-%02d-%02d", time.GetYear(), time.GetMonth(), time.GetDay());
	return;
}

/*************************************************
函数名称:
函数描述:取当前时间（精确到秒）
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
void GetDateTime(char* datetime)
{
	CTime time = CTime::GetCurrentTime();
	sprintf_s(datetime, 20, "%04d%02d%02d%02d%02d%02d", time.GetYear(), time.GetMonth(), time.GetDay(), time.GetHour(), time.GetMinute(), time.GetSecond());
	return;
}


/*************************************************
函数名称:
函数描述:读取目录下所有文件
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
void GetFiles(string path, vector<string>& files, vector<string>& filenames)
{
	//文件句柄
	long hFile = 0;
	//文件信息
	struct _finddata_t fileinfo;
	string p;

	if ((hFile = _findfirst(p.assign(path).append("\\*").c_str(), &fileinfo)) != -1)
	{
		do
		{
			//如果是目录,迭代之
			//如果不是,加入列表
			if ((fileinfo.attrib & _A_SUBDIR))
			{
				if (strcmp(fileinfo.name, ".") != 0 && strcmp(fileinfo.name, "..") != 0)
					GetFiles(p.assign(path).append("\\").append(fileinfo.name), files, filenames);
			}
			else 
			{
				filenames.push_back(fileinfo.name);
				files.push_back(p.assign(path).append("\\").append(fileinfo.name));
			}
		} while (_findnext(hFile, &fileinfo) == 0);
		_findclose(hFile);
	}
}

/*************************************************
函数名称:
函数描述:字符转换
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
int code_convert(const char *from_charset, const char *to_charset, char *inbuf, int inlen, char *outbuf, int outlen)
{
	iconv_t cd;
	//int rc;
	char **pin = &inbuf;
	char **pout = &outbuf;
	cd = iconv_open(to_charset, from_charset);
	if (cd == 0) return -1;
	memset(outbuf, 0, outlen);
	if ((int)iconv(cd, (const char**)pin, (size_t *)&inlen, pout, (size_t *)&outlen) == -1) return -1;
	iconv_close(cd);
	return 0;
}

int g2u(char *inbuf, size_t inlen, char *outbuf, size_t outlen)
{
	return code_convert("gb2312", "utf-8", inbuf, inlen, outbuf, outlen);
}

int u2g(char *inbuf, int inlen, char *outbuf, int outlen)
{
	return code_convert("utf-8", "gb2312", inbuf, inlen, outbuf, outlen);
}

/*************************************************
函数名称:
函数描述:显示窗口日志
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
void ShowLstMsg1(LPVOID lparam, CString msg)
{
	Sleep(5);
	g_msg1 = msg;
	CFeeSettlementDlg* dlg = (CFeeSettlementDlg*)lparam;
	PostMessage(dlg->m_hWnd, MYMSG1, WPARAM(&g_msg1), NULL);
	
}

void ShowLstMsg1(LPVOID lparam, CString msg, DWORD val)
{
	Sleep(5);
	CFeeSettlementDlg* dlg = (CFeeSettlementDlg*)lparam;
	g_msg1 = msg;
	g_msg1 += to_string(val).c_str();
	PostMessage(dlg->m_hWnd, MYMSG1, WPARAM(&g_msg1), NULL);
	
}

void ShowLstMsg1(LPVOID lparam, CString msg, char* ch)
{
	Sleep(5);
	CFeeSettlementDlg* dlg = (CFeeSettlementDlg*)lparam;
	g_msg1 = msg;
	g_msg1 += ch;
	PostMessage(dlg->m_hWnd, MYMSG1, WPARAM(&g_msg1), NULL);
	
}

void ShowLstMsg1(LPVOID lparam, CString msg, CString str)
{
	Sleep(5);
	CFeeSettlementDlg* dlg = (CFeeSettlementDlg*)lparam;
	g_msg1 = msg;
	g_msg1 += str;
	PostMessage(dlg->m_hWnd, MYMSG1, WPARAM(&g_msg1), NULL);
	
}

void ShowLstMsg2(LPVOID lparam, CString msg)
{
	Sleep(5);
	g_msg2 = msg;
	CFeeSettlementDlg* dlg = (CFeeSettlementDlg*)lparam;
	PostMessage(dlg->m_hWnd, MYMSG2, WPARAM(&g_msg2), NULL);
	
}

void ShowLstMsg2(LPVOID lparam, CString msg, DWORD val)
{
	Sleep(5);
	CFeeSettlementDlg* dlg = (CFeeSettlementDlg*)lparam;
	g_msg2 = msg;
	g_msg2 += to_string(val).c_str();
	PostMessage(dlg->m_hWnd, MYMSG2, WPARAM(&g_msg2), NULL);
	
}

void ShowLstMsg2(LPVOID lparam, CString msg, char* ch)
{
	Sleep(5);
	CFeeSettlementDlg* dlg = (CFeeSettlementDlg*)lparam;
	g_msg2 = msg;
	g_msg2 += ch;
	PostMessage(dlg->m_hWnd, MYMSG2, WPARAM(&g_msg2), NULL);
	
}

void ShowLstMsg2(LPVOID lparam, CString msg, CString str)
{
	Sleep(5);
	CFeeSettlementDlg* dlg = (CFeeSettlementDlg*)lparam;
	g_msg2 = msg;
	g_msg2 += str;
	PostMessage(dlg->m_hWnd, MYMSG2, WPARAM(&g_msg2), NULL);
	
}

/*************************************************
函数名称:
函数描述:创建xml文件
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
BOOL CreateXmlFile(LPVOID lparam, ETC_XML_HEADER* header, ETC_XML_SEND_BODY_DESC* body, ETC_XML_SEND_TRANSACTION* trans, char* filepath)
{
	if (header == NULL || body == NULL || trans == NULL || filepath == NULL)
	{
		ShowLstMsg1(lparam, _T("CreateXmlFile: 参数传入不正确..."));
		return FALSE;
	}

	//用来转换金额，保留2位小数位
	char toll[20] = { 0 };

	char utf8str[100] = { 0 };

	tinyxml2::XMLDocument doc;

	const char* declaration = "<?xml version=\"1.0\" encoding=\"utf-8\"?>";
	doc.Parse(declaration);		//会覆盖xml所有内容

	//根节点
	XMLElement* root = doc.NewElement("Message");
	doc.InsertEndChild(root);

	//消息头
	XMLElement* heardernode = doc.NewElement("Header");
	root->InsertEndChild(heardernode);
	XMLElement* hearderchildren = doc.NewElement("Version");
	hearderchildren->InsertEndChild(doc.NewText("00010000"));
	heardernode->InsertEndChild(hearderchildren);

	hearderchildren = doc.NewElement("MessageClass");
	hearderchildren->InsertEndChild(doc.NewText(to_string(header->messageclass).c_str()));
	heardernode->InsertEndChild(hearderchildren);

	hearderchildren = doc.NewElement("MessageType");
	hearderchildren->InsertEndChild(doc.NewText(to_string(header->messagetype).c_str()));
	heardernode->InsertEndChild(hearderchildren);
	
	hearderchildren = doc.NewElement("SenderId");
	hearderchildren->InsertEndChild(doc.NewText(header->senderid));
	heardernode->InsertEndChild(hearderchildren);

	hearderchildren = doc.NewElement("ReceiverId");
	hearderchildren->InsertEndChild(doc.NewText(header->receiverid));
	heardernode->InsertEndChild(hearderchildren);

	hearderchildren = doc.NewElement("MessageId");
	hearderchildren->InsertEndChild(doc.NewText(to_string(header->messageid).c_str()));
	heardernode->InsertEndChild(hearderchildren);

	//消息体
	XMLElement* bodynode = doc.NewElement("Body");
	bodynode->SetAttribute("ContentType", body->contenttype);
	root->InsertEndChild(bodynode);

	XMLElement* bodychildren = doc.NewElement("ClearTargetDate");
	bodychildren->InsertEndChild(doc.NewText(body->cleartargetdate));
	bodynode->InsertEndChild(bodychildren);

	bodychildren = doc.NewElement("ServiceProviderId");
	bodychildren->InsertEndChild(doc.NewText(body->serviceproviderid));
	bodynode->InsertEndChild(bodychildren);

	bodychildren = doc.NewElement("IssuerId");
	bodychildren->InsertEndChild(doc.NewText(body->issuerid));
	bodynode->InsertEndChild(bodychildren);

	bodychildren = doc.NewElement("MessageId");
	bodychildren->InsertEndChild(doc.NewText(to_string(body->messageid).c_str()));
	bodynode->InsertEndChild(bodychildren);

	bodychildren = doc.NewElement("Count");
	bodychildren->InsertEndChild(doc.NewText(to_string(body->count).c_str()));
	bodynode->InsertEndChild(bodychildren);

	bodychildren = doc.NewElement("Amount");
	sprintf_s(toll,sizeof(toll), "%.2f", body->amount);
	bodychildren->InsertEndChild(doc.NewText(toll));
	bodynode->InsertEndChild(bodychildren);

	//消息体-交易明细
	for (int i = 0; i < body->count; i++)
	{
		XMLElement* transnode = doc.NewElement("Transaction");
		bodynode->InsertEndChild(transnode);

		XMLElement* transchildren = doc.NewElement("TransId");
		transchildren->InsertEndChild(doc.NewText(to_string(trans[i].transid).c_str()));
		transnode->InsertEndChild(transchildren);

		transchildren = doc.NewElement("Time");
		transchildren->InsertEndChild(doc.NewText(trans[i].time));
		transnode->InsertEndChild(transchildren);

		transchildren = doc.NewElement("Fee");
		sprintf_s(toll, sizeof(toll), "%.2f", trans[i].fee);
		transchildren->InsertEndChild(doc.NewText(toll));
		transnode->InsertEndChild(transchildren);

		//服务信息子节点
		XMLElement* servicenode = doc.NewElement("Service");
		transnode->InsertEndChild(servicenode);
		
		XMLElement* servicechildren = doc.NewElement("ServiceType");
		servicechildren->InsertEndChild(doc.NewText(to_string(trans[i].service.servicetype).c_str()));
		servicenode->InsertEndChild(servicechildren);

		servicechildren = doc.NewElement("Description");
		memset(utf8str, 0, sizeof(utf8str));
		g2u(trans[i].service.description, strlen(trans[i].service.description) + 1, utf8str, 100);
		servicechildren->InsertEndChild(doc.NewText(utf8str));
		//servicechildren->InsertEndChild(doc.NewText(trans[i].service.description));
		servicenode->InsertEndChild(servicechildren);

		servicechildren = doc.NewElement("Detail");
		memset(utf8str, 0, sizeof(utf8str));
		g2u(trans[i].service.detail, strlen(trans[i].service.detail) + 1, utf8str, 100);
		servicechildren->InsertEndChild(doc.NewText(utf8str));
		servicenode->InsertEndChild(servicechildren);

		//IC卡节点
		XMLElement* iccardnode = doc.NewElement("ICCard");
		transnode->InsertEndChild(iccardnode);

		XMLElement* iccardchildren = doc.NewElement("CardType");
		iccardchildren->InsertEndChild(doc.NewText(to_string(trans[i].iccard.cardtype).c_str()));
		iccardnode->InsertEndChild(iccardchildren);

		iccardchildren = doc.NewElement("NetNo");
		iccardchildren->InsertEndChild(doc.NewText(trans[i].iccard.netno));
		iccardnode->InsertEndChild(iccardchildren);

		iccardchildren = doc.NewElement("CardId");
		iccardchildren->InsertEndChild(doc.NewText(trans[i].iccard.cardid));
		iccardnode->InsertEndChild(iccardchildren);

		iccardchildren = doc.NewElement("License");
		memset(utf8str, 0, sizeof(utf8str));
	
		g2u(trans[i].iccard.license, strlen(trans[i].iccard.license) + 1, utf8str, 100);
		iccardchildren->InsertEndChild(doc.NewText(utf8str));
		//iccardchildren->InsertEndChild(doc.NewText(trans[i].iccard.license));
		iccardnode->InsertEndChild(iccardchildren);

		iccardchildren = doc.NewElement("PreBalance");
		sprintf_s(toll, sizeof(toll), "%.2f", trans[i].iccard.prebalance*0.01);
		iccardchildren->InsertEndChild(doc.NewText(toll));
		iccardnode->InsertEndChild(iccardchildren);

		iccardchildren = doc.NewElement("PostBalance");
		sprintf_s(toll, sizeof(toll), "%.2f", trans[i].iccard.postbalance*0.01);
		iccardchildren->InsertEndChild(doc.NewText(toll));
		iccardnode->InsertEndChild(iccardchildren);
		
		//清分信息节点
		XMLElement* validationnode = doc.NewElement("Validation");
		transnode->InsertEndChild(validationnode);

		XMLElement* validationnodechildren = doc.NewElement("TAC");
		validationnodechildren->InsertEndChild(doc.NewText(trans[i].validation.tac));
		validationnode->InsertEndChild(validationnodechildren);

		validationnodechildren = doc.NewElement("TransType");
		validationnodechildren->InsertEndChild(doc.NewText(trans[i].validation.transtype));
		validationnode->InsertEndChild(validationnodechildren);

		validationnodechildren = doc.NewElement("TerminalNo");
		validationnodechildren->InsertEndChild(doc.NewText(trans[i].validation.terminalno));
		validationnode->InsertEndChild(validationnodechildren);

		validationnodechildren = doc.NewElement("TerminalTransNo");
		validationnodechildren->InsertEndChild(doc.NewText(trans[i].validation.terminaltransno));
		validationnode->InsertEndChild(validationnodechildren);

		//OBU子节点信息
		XMLElement* obunode = doc.NewElement("OBU");
		transnode->InsertEndChild(obunode);

		XMLElement* obuchildren = doc.NewElement("NetNo");
		obuchildren->InsertEndChild(doc.NewText(trans[i].obu.netno));
		obunode->InsertEndChild(obuchildren);

		obuchildren = doc.NewElement("OBUId");
		obuchildren->InsertEndChild(doc.NewText(trans[i].obu.obuid));
		obunode->InsertEndChild(obuchildren);

		obuchildren = doc.NewElement("OBEState");
		obuchildren->InsertEndChild(doc.NewText(trans[i].obu.obestate));
		obunode->InsertEndChild(obuchildren);

		obuchildren = doc.NewElement("License");
		memset(utf8str, 0, sizeof(utf8str));
		g2u(trans[i].obu.license, strlen(trans[i].obu.license) + 1, utf8str, 100);
		obuchildren->InsertEndChild(doc.NewText(utf8str));
		obunode->InsertEndChild(obuchildren);

		//描述
		transchildren = doc.NewElement("CustomizedData");
		transchildren->InsertEndChild(doc.NewText(trans[i].customizeddata));
		transnode->InsertEndChild(transchildren);
	}

	doc.SaveFile(filepath);
	doc.Clear();

	return TRUE;
}


/*************************************************
函数名称:
函数描述:获取文件md5加密
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
BOOL GetMd5(LPVOID lparam, char* srcfile, unsigned char* md5)
{
	if (srcfile == NULL || md5 == NULL)
	{
		ShowLstMsg2(lparam, _T("GetMd5|参数传入不正确"));
		return FALSE;
	}	

	//二进制方式打开文件
	FILE * fp;
	errno_t err;
	if ((err = fopen_s(&fp, srcfile, "rb")) != 0)
	{
		ShowLstMsg2(lparam, _T("GetMd5|无法打开此文件xml文件..."));
		return FALSE;
	}

	/* 获取文件大小 */
	fseek(fp, 0, SEEK_END);
	long len = ftell(fp);
	rewind(fp);

	//申请内存
	char *fpbuf;
	fpbuf = (char*)malloc(sizeof(char)*len);
	if (fpbuf == NULL)
	{
		ShowLstMsg2(lparam, _T("GetMd5|内存申请出错..."));
		return FALSE;
	}

	/* 将文件拷贝到buffer中 */
	size_t val = fread(fpbuf, 1, len, fp);
	if (val != len)
	{
		ShowLstMsg2(lparam, _T("GetMd5|读文件出错..."));
		return FALSE;
	}

	unsigned char dest[16];
	MD5_CTX context;
	MD5Init(&context);
	MD5Update(&context, (unsigned char*)fpbuf, len);
	MD5Final(dest, &context);
	
	for (int i = 0; i < 16; ++i)
	{
		sprintf_s((char*)md5 + i * 2, 3, "%02X" ,dest[i]);
		//printf("%02X", dest[i]);
	}

	fclose(fp);
	free(fpbuf);

	return TRUE;
}

/*************************************************
函数名称:
函数描述:压缩文件，压缩算法:LZ77
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
BOOL CompressFile(LPVOID lparam, char* infile, char* outfile)
{
	if (infile == NULL || outfile == NULL)
	{
		ShowLstMsg2(lparam, _T("CompressFile|参数传入不正确..."));
		return FALSE;
	}

	CCompressLZ77 cc;
	if (cc.CompressFile(infile, outfile) == 0)
		return TRUE;
	else
		return FALSE;
}
/*
BOOL CompressFile(LPVOID lparam, char* infile, char* outfile)
{
	if (infile == NULL || outfile == NULL)
	{
		ShowLstMsg2(lparam, _T("CompressFile|参数传入不正确..."));
		return FALSE;
	}	

	BYTE soubuf[65536] = { 0 };
	BYTE destbuf[65536 + 16] = { 0 };

	FILE* in;
	FILE* out;
	errno_t err;

	//打开源xml文件
	if ((err = fopen_s(&in, infile, "rb")) != 0)
	{
		ShowLstMsg2(lparam, _T("CompressFile|无法打开xml源文件..."));
		return FALSE;
	}

	//创建目标压缩文件
	if ((err = fopen_s(&out, outfile, "wb")) != 0)
	{
		ShowLstMsg2(lparam, _T("CompressFile|无法创建目标压缩文件..."));
		fclose(in);
		return FALSE;
	}

	fseek(in, 0, SEEK_END);
	long inlen = ftell(in);
	fseek(in, 0, SEEK_SET);

	CCompressLZ77 cc;
	WORD flag1, flag2;
	int last = inlen, act, outlen;
	while (last > 0)
	{
		act = min(65536, last);
		fread(soubuf, act, 1, in);
		last -= act;
		if (act == 65536)           // out 65536 bytes                 
			flag1 = 0;
		else                    // out last blocks   
			flag1 = act;
		fwrite(&flag1, sizeof(WORD), 1, out);

		outlen = cc.Compress((BYTE*)soubuf, act, (BYTE*)destbuf);
		if (outlen == 0)       // can't compress the block   
		{
			flag2 = flag1;
			fwrite(&flag2, sizeof(WORD), 1, out);
			fwrite(soubuf, act, 1, out);
		}
		else
		{
			flag2 = (WORD)outlen;
			fwrite(&flag2, sizeof(WORD), 1, out);
			fwrite(destbuf, outlen, 1, out);
		}
	}

	fclose(in);
	fclose(out);

	return TRUE;
}*/

/*************************************************
函数名称:
函数描述:解压缩文件，压缩算法:LZ77
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
BOOL UnCompressFile(LPVOID lparam, char* infile, char* outfile)
{
	if (infile == NULL || outfile == NULL)
	{
		ShowLstMsg2(lparam, _T("CompressFile|参数传入不正确..."));
		return FALSE;
	}

	CCompressLZ77 cc;
	if (cc.DecompressFile(infile, outfile) == 0)
		return TRUE;
	else
		return FALSE;
}
/*
BOOL UnCompressFile(LPVOID lparam, char* infile, char* outfile)
{
	if (infile == NULL || outfile == NULL)
	{
		ShowLstMsg2(lparam, _T("UnCompressFile|参数传入不正确..."));
		return FALSE;
	}

	BYTE soubuf[65536] = { 0 };
	BYTE destbuf[65536 + 16] = { 0 };

	FILE* in;
	FILE* out;
	errno_t err;

	//打开源文件
	if ((err = fopen_s(&in, infile, "rb")) != 0)
	{
		ShowLstMsg2(lparam, _T("UnCompressFile|无法打开源文件..."));
		return FALSE;
	}

	//创建目标压缩文件
	if ((err = fopen_s(&out, outfile, "wb")) != 0)
	{
		ShowLstMsg2(lparam, _T("UnCompressFile|无法创建目标解压缩文件..."));
		fclose(in);
		return FALSE;
	}

	fseek(in, 0, SEEK_END);
	long inlen = ftell(in);
	fseek(in, 0, SEEK_SET);

	CCompressLZ77 cc;
	WORD flag1, flag2;
	int last = inlen, act;
	while (last > 0)
	{
		fread(&flag1, sizeof(WORD), 1, in);
		fread(&flag2, sizeof(WORD), 1, in);
		last -= 2 * sizeof(WORD);
		if (flag1 == 0)
			act = 65536;
		else
			act = flag1;
		last -= flag2 ? (flag2) : act;

		if (flag2 == flag1)
		{
			fread(soubuf, act, 1, in);
		}
		else
		{
			fread(destbuf, flag2, 1, in);
			if (!cc.Decompress((BYTE*)soubuf, act, (BYTE*)destbuf))
			{
				ShowLstMsg2(lparam, _T("UnCompressFile|解压缩文件失败..."));
				fclose(in);
				fclose(out);
				return FALSE;
			}
		}
		fwrite((BYTE*)soubuf, act, 1, out);
	}

	fclose(in);
	fclose(out);

	return TRUE;
}*/

/*************************************************
函数名称:
函数描述:发送xml文件
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
BOOL SendCompressFile(LPVOID lparam, char* messageid, char* zipfile, unsigned char* md5)
{
	if (messageid == NULL || zipfile == NULL || md5 == NULL)
	{
		ShowLstMsg2(lparam, _T("SendCompressFile|参数传入不正确..."));
		return FALSE;
	}

	//////////////////////////////////////////////////////////////////////////
	// 连接网络
	//////////////////////////////////////////////////////////////////////////
	WORD request;
	WSADATA sadata;

	//socket通讯设置
	request = MAKEWORD(2, 2);

	if (WSAStartup(request, &sadata) != 0)
	{
		ShowLstMsg2(lparam, _T("SendCompressFile|socket通讯异常..."));
		return FALSE;
	}		

	if (LOBYTE(sadata.wVersion) != 2 || HIBYTE(sadata.wVersion) != 2)
	{
		WSACleanup();
		ShowLstMsg2(lparam, _T("SendCompressFile|socket通讯异常..."));
		return FALSE;
	}

	//捆绑本地地址

	SOCKET sockclient = socket(AF_INET, SOCK_STREAM, 0);
	SOCKADDR_IN addrsrv;

	//设置接收超时时间
	int nNetTimeout = 10000;//10秒，
	setsockopt(sockclient, SOL_SOCKET, SO_SNDTIMEO, (char *)&nNetTimeout, sizeof(int));

	inet_pton(AF_INET, g_config.serverip, (void*)&addrsrv.sin_addr.S_un.S_addr);
	addrsrv.sin_family = AF_INET;
	addrsrv.sin_port = htons(CLIENT_CONNECT_PORT);

	//设置socket接收时长限制参数，超过十秒，则继续往下执行
	int timeout = 10*1000;
	int ret = setsockopt(sockclient, SOL_SOCKET, SO_RCVTIMEO, (char *)&timeout, sizeof(timeout));

	//连接服务端
	int val = -1;
	while (TRUE)
	{
		val = connect(sockclient, (SOCKADDR*)&addrsrv, sizeof(SOCKADDR));
		if (val != 0)
		{
			ShowLstMsg1(lparam, _T("SendCompressFile|socket连接服务端失败..."));
			Sleep(1000);
			continue;
		}
		else
		{
			break;
		}
	}	

	SEND_STRU sendmsg = { 0 };
	memset((void*)&sendmsg, 0, sizeof(SEND_STRU));

	//////////////////////////////////////////////////////////////////////////
	// 读入压缩文件至内存
	//////////////////////////////////////////////////////////////////////////
	FILE* in;
	errno_t err;
	DWORD fiellen = 0;
	BYTE* content;

	//打开压缩文件
	if ((err = fopen_s(&in, zipfile, "rb")) != 0)
	{
		ShowLstMsg2(lparam, _T("SendCompressFile|无法打开xml压缩文件...|"), messageid);
		return FALSE;
	}

	fseek(in, 0, SEEK_END);
	fiellen = ftell(in);
	fseek(in, 0, SEEK_SET);

	int last = fiellen, act = 0;
	content = (BYTE*)sendmsg.xml_msg;
	while (last > 0)
	{
		//每次读入65536字节
		act = min(65536, last);		
		fread(content, act, 1, in);
		content += act;
		last -= act;
	}

	fclose(in);

	//////////////////////////////////////////////////////////////////////////
	// 发送数据
	//////////////////////////////////////////////////////////////////////////	
	REPLY_STRU replymsg = { 0 };
	DWORD sizelen = 0;
	DWORD result = 0;

	char recvbuf[22] = { 0 };

	//消息报文顺序号 （SSSSSSS） 20字节 + 压缩后的XML消息长度(CCCCCC) 6字节 + 32字节MD5(16进制显示字符串) + 二进制压缩后的XML消息包(AAAAAAAAAAAAAAA…)	
	strcpy_s(sendmsg.messageid, messageid);
	sprintf_s(sendmsg.xmllen, sizeof(sendmsg.xmllen), "%06lu", fiellen);
	strcpy_s(sendmsg.md5str, (const char*)md5);

	ShowLstMsg2(lparam, _T("sendmsg.messageid|"), sendmsg.messageid);
	ShowLstMsg2(lparam, _T("sendmsg.xmllen|"), sendmsg.xmllen);
	ShowLstMsg2(lparam, _T("sendmsg.md5str|"), sendmsg.md5str);

	//发送长度
	sizelen = 20 + 6 + 32 + fiellen + 1;
	ShowLstMsg2(lparam, _T("发送文件长度|"), sizelen);

	byte *p = NULL;
	p = new byte[sizelen];
	memcpy(p, (char*)sendmsg.messageid, 20);
	memcpy(p+20, (char*)sendmsg.xmllen, 6);
	memcpy(p+20+6, (char*)sendmsg.md5str, 32);
	memcpy(p+20+6+32, (char*)sendmsg.xml_msg, fiellen);

	//发送，3次
	for (int i = 0; i < 3; i++)
	{
		result = send(sockclient, (char*)p, sizelen, 0);
		//CreatFile((char*)p, sizelen);
		if (result == sizelen)
		{
			ShowLstMsg2(lparam, _T("发送文件成功，文件名|"), messageid);
			break;
		}
		else
		{
			ShowLstMsg2(lparam, _T("发送文件失败，再次发送...|"), messageid);
		}

		if (i >= 2)
		{
			ShowLstMsg2(lparam, _T("发送超过三次，本次发送文件失败，文件名|"), messageid);
			closesocket(sockclient);
			WSACleanup();
			delete (p);
			return FALSE;
		}
	}

	delete (p);

	//////////////////////////////////////////////////////////////////////////
	// 接收应答
	//////////////////////////////////////////////////////////////////////////
	//MessageId + tcp接收结果，21个字节
	memset((void*)&replymsg, 0, sizeof(REPLY_STRU));
	sizelen = recv(sockclient, recvbuf, 21, 0);
	if (sizelen == 21)
	{
		ShowLstMsg2(lparam, _T("SendCompressFile|收到文件应答"));

		//格式化指针
		REPLY_STRU* replymsg = (REPLY_STRU*)recvbuf;
		strcat_s(sendmsg.messageid, 22, "1");
		//&& replymsg->result == 1
		if (strcmp(sendmsg.messageid, replymsg->messageid) == 0 )
		{
			ShowLstMsg2(lparam, _T("文件应答接收成功"));
		}
		else
		{
			ShowLstMsg2(lparam, _T("文件应答接收失败...messageid:"), replymsg->messageid);
			ShowLstMsg2(lparam, _T("文件应答接收失败...result:"), replymsg->result);
			closesocket(sockclient);
			WSACleanup();
			return FALSE;
		}
	}
	//本意是想辨别recv超时进行处理
	else
	{
		ShowLstMsg2(lparam, _T("文件应答接收失败..."));
		closesocket(sockclient);
		WSACleanup();
		return FALSE;
	}

	//关闭通讯
	closesocket(sockclient);
	WSACleanup();

	return TRUE;
}


/*************************************************
函数名称:
函数描述:提取数据并打包(外省)
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
BOOL OtherProv_GetOriData(LPVOID lparam, ETC_XML_HEADER *header, ETC_XML_SEND_BODY_DESC *body, ETC_XML_SEND_TRANSACTION *trans, char *dirbuf, int NetNo, int CardType)
{
	if (header == NULL || body == NULL || trans == NULL || dirbuf == NULL)
	{
		ShowLstMsg1(lparam, _T("OtherProv_GetOriData: 参数传入不正确..."));
		WriteLog("结算线程-> OtherProv_GetOriData: 参数传入不正确...");
		return FALSE;
	}

	// 读取配置文件缓存
	char inBuf[80];

	long int nMessageId;//本次包号

	//sql语句
	char sql[1024] = { 0 };
	//xml文件名
	char xmlfilename[300] = { 0 };
	BOOL result = FALSE;

	//记录提取到的总条数
	int num = 0;

	OCI_Connection* cn;
	OCI_Statement* st;
	OCI_Resultset* rs;

	/*result = OCI_Initialize(NULL, NULL, OCI_ENV_THREADED);
	if (!result)
	{
		return FALSE;
	}*/

	//连接 第一个参数格式:【IP:端口/服务名】,第二个参数：登录用户名，第三个参数：密码
	cn = OCI_ConnectionCreate(g_config.dbservice, g_config.dbuser, g_config.dbpwd, OCI_SESSION_DEFAULT);
	if (cn == NULL)
	{
		return FALSE;
	}

	st = OCI_StatementCreate(cn);
	if (st == NULL)
	{
		return FALSE;
	}

	//消息头////////////////////////////////////////
	header->version = 1;
	header->messageclass = 5;
	header->messagetype = 7;
	memcpy(header->senderid, "00000000000000FD", sizeof(header->senderid));
	memcpy(header->receiverid, "0000000000000020", sizeof(header->senderid));

	//获取配置文件中的最新包号
	memset(inBuf, 0, sizeof(inBuf));
	WCHAR wszMessageId[256];
	memset(wszMessageId, 0, sizeof(wszMessageId));
	//char* 转LPCWSTR
	MultiByteToWideChar(CP_ACP, 0, inBuf, strlen(inBuf) + 1, wszMessageId, sizeof(wszMessageId) / sizeof(wszMessageId[0]));
	//读配置文件包号
	GetPrivateProfileString(L"config", L"num_of_deal", NULL, wszMessageId, sizeof(wszMessageId), L".\\db.ini");
	// LPCWSTR转char*
	DWORD dwMinSize;
	dwMinSize = WideCharToMultiByte(CP_ACP, NULL, wszMessageId, -1, NULL, 0, NULL, FALSE); //计算长度
	WideCharToMultiByte(CP_OEMCP, NULL, wszMessageId, -1, inBuf, dwMinSize, NULL, FALSE);
	//获取长整型包号
	nMessageId = atol(inBuf);

	header->messageid = nMessageId;//本省外省包号分奇数偶数，保证不重复

	//消息体////////////////////////////////////////
	body->contenttype = 1;
	GetDate(body->cleartargetdate);
	memcpy(body->serviceproviderid, "00000000000000FD", sizeof(header->senderid));
	memcpy(body->issuerid, "0000000000000020", sizeof(header->senderid));
	body->messageid = header->messageid;
	//查询记录明细，暂定为每次查询100条
	memset(sql, 0, sizeof(sql));
	sprintf_s(sql, sizeof(sql), "select count(*) as cnt,sum(F_NB_JIAOYJE) from B_TXF_JIESSJ \
					where F_NB_JIESJG = 98 and F_NB_XMLDBBJ = 0 and F_NB_KAWLH = %d and \
					F_NB_KALX = %d and F_NB_FUFFS = 2 and F_NB_JIAOYJE <>0 and rownum <= %d", NetNo, CardType,TRANS_MAX_COUNT);
	OCI_ExecuteStmt(st, sql);
	rs = OCI_GetResultset(st);
	while (OCI_FetchNext(rs))
	{
		body->count = OCI_GetInt(rs, 1);
		body->amount = OCI_GetInt(rs, 2) * 0.01;
	}

	OCI_ReleaseResultsets(st);

	//如果条数大于0，需打包，包号写回配置文件
	if (body->count > 0)
	{
		nMessageId = nMessageId + 2;

		memset(inBuf, 0, sizeof(inBuf));
		sprintf_s(inBuf, "%ld", nMessageId);
		memset(wszMessageId, 0, sizeof(wszMessageId));
		//char* 转LPCWSTR
		MultiByteToWideChar(CP_ACP, 0, inBuf, strlen(inBuf) + 1, wszMessageId, sizeof(wszMessageId) / sizeof(wszMessageId[0]));
		WritePrivateProfileString(L"config", L"num_of_deal", wszMessageId, L".\\db.ini");//写回配置
	}

	//交易明细//////////////////////////////////////
	memset(sql, 0, sizeof(sql));
	sprintf_s(sql, sizeof(sql), "select ROWNUM, to_char(F_DT_JIAOYSJ,'yyyy-mm-dd HH24:mi:ss') as F_DT_JIAOYSJ,F_NB_SHOUFZBH,F_NB_CHEDH, F_NB_JIAOYJE, \
					nvl(F_VC_JIAOYMS,'无') as F_VC_JIAOYMS, nvl(F_VC_JIAOYXXXX,'无') as F_VC_JIAOYXXXX,\
					F_NB_KALX, F_NB_KAWLH, F_VC_KABH, nvl(F_VC_KACP,'无') as F_VC_KACP, F_NB_KAQYE, F_NB_KAHYE,\
					F_VC_TACM, nvl(F_NB_JIAOYLX,0) as F_NB_JIAOYLX, F_VC_ZHONGDJBM, F_VC_ZHONGDJYXLH,\
					nvl(F_NB_OBUWLX,'') as F_NB_OBUWLX, F_VC_OBUBH, F_NB_OBUZT, nvl(F_VC_OBUCP,'无')as F_VC_OBUCP,\
					nvl(F_VC_YUANSJYXX,'0') as F_VC_YUANSJYXX \
					from B_TXF_JIESSJ where F_NB_JIESJG = 98 and F_NB_XMLDBBJ = 0 and F_NB_KAWLH = %d and F_NB_FUFFS = 2 \
					and F_NB_KALX = %d and F_NB_JIAOYJE <> 0 and rownum <= %d", NetNo, CardType, TRANS_MAX_COUNT);
	OCI_ExecuteStmt(st, sql);
	rs = OCI_GetResultset(st);
	while (OCI_FetchNext(rs))
	{
		trans[num].transid = OCI_GetInt2(rs, "ROWNUM");
		sprintf_s(trans[num].time,sizeof(trans[num].time), "%s", OCI_GetString2(rs, "F_DT_JIAOYSJ"));
		//memcpy(trans[num].time, OCI_GetString2(rs, "F_DT_JIAOYSJ"), sizeof(trans[num].time));
		trans[num].fee = OCI_GetInt2(rs, "F_NB_JIAOYJE") * 0.01;

		//20171117 begin
		trans[num].stationId = OCI_GetBigInt2(rs, "F_NB_SHOUFZBH");
		trans[num].lane = OCI_GetInt2(rs, "F_NB_CHEDH");
		//20171117 end

		//服务信息
		trans[num].service.servicetype = 2;
		sprintf_s(trans[num].service.description,sizeof(trans[num].service.description), "%s", OCI_GetString2(rs, "F_VC_JIAOYMS"));		//格式：XXXXXX停车场|3天1小时12分
		sprintf_s(trans[num].service.detail, sizeof(trans[num].service.detail),"%s", OCI_GetString2(rs, "F_VC_JIAOYXXXX"));

		//IC卡信息
		trans[num].iccard.cardtype = OCI_GetInt2(rs, "F_NB_KALX");
		sprintf_s(trans[num].iccard.netno, sizeof(trans[num].iccard.netno),"%s", OCI_GetString2(rs, "F_NB_KAWLH"));
		sprintf_s(trans[num].iccard.cardid, sizeof(trans[num].iccard.cardid),"%s", OCI_GetString2(rs, "F_VC_KABH"));
		sprintf_s(trans[num].iccard.license, sizeof(trans[num].iccard.license),"%s", OCI_GetString2(rs, "F_VC_KACP"));
		trans[num].iccard.prebalance = OCI_GetBigInt2(rs, "F_NB_KAQYE");
		trans[num].iccard.postbalance = OCI_GetBigInt2(rs, "F_NB_KAHYE");

		//TAC信息
		sprintf_s(trans[num].validation.tac, sizeof(trans[num].validation.tac),"%s", OCI_GetString2(rs, "F_VC_TACM"));
		sprintf_s(trans[num].validation.transtype, sizeof(trans[num].validation.transtype), "%02s", OCI_GetString2(rs, "F_NB_JIAOYLX"));
		sprintf_s(trans[num].validation.terminalno, sizeof(trans[num].validation.terminalno), "%s", OCI_GetString2(rs, "F_VC_ZHONGDJBM"));
		sprintf_s(trans[num].validation.terminaltransno, sizeof(trans[num].validation.terminaltransno),"%s", OCI_GetString2(rs, "F_VC_ZHONGDJYXLH"));

		//OBU信息
		sprintf_s(trans[num].obu.netno, sizeof(trans[num].obu.netno),"%s", OCI_GetString2(rs, "F_NB_KAWLH"));
		sprintf_s(trans[num].obu.obuid, sizeof(trans[num].obu.obuid),"%s", OCI_GetString2(rs, "F_VC_OBUBH"));
		sprintf_s(trans[num].obu.obestate, sizeof(trans[num].obu.obestate), "%04s", OCI_GetString2(rs, "F_NB_OBUZT"));
		sprintf_s(trans[num].obu.license, sizeof(trans[num].obu.license),"%s", OCI_GetString2(rs, "F_VC_OBUCP"));

		//原始信息
		sprintf_s(trans[num].customizeddata, sizeof(trans[num].customizeddata), "%s", OCI_GetString2(rs, "F_VC_YUANSJYXX"));

		num += 1;
	}

	//////////////////////////////////////////////////////////////////////////
	// 封装文件，打包xml
	//////////////////////////////////////////////////////////////////////////
	if (body->count > 0)
	{
		memset(xmlfilename, 0, sizeof(xmlfilename));

		if (CardType == 22)
		{
			//拼接文件名
			sprintf_s(xmlfilename, sizeof(xmlfilename),"%s\\CZ_%d_%020lu.xml", dirbuf, NetNo, header->messageid);
		}
		else if (CardType == 23)
		{
			//拼接文件名
			sprintf_s(xmlfilename, sizeof(xmlfilename),"%s\\JZ_%d_%020lu.xml", dirbuf, NetNo, header->messageid);
		}

		//判断文件是否存在，存在则先删除
		if (ACCESS(xmlfilename, 0) == 0)
		{
			DeleteFile((CString)xmlfilename);
		}

		//根据数据库数据生成xml文件
		result = CreateXmlFile(lparam, header, body, trans, xmlfilename);
		if (result == TRUE)
		{
			//ShowLstMsg1(lparam, _T("外省卡记账结算源数据打包xml成功"));
		}
		else
		{
			//ShowLstMsg1(lparam, _T("外省卡记账结算源数据打包xml失败..."));
			if (rs)
				OCI_ReleaseResultsets(st);

			//清除声明
			if (st)
				OCI_StatementFree(st);

			//清除连接
			if (cn)
				OCI_ConnectionFree(cn);

			//OCI_Cleanup();
			return FALSE;
		}
	}
	else
	{
		ShowLstMsg1(lparam, _T("(数据源记录为空,无数据处理..."));
		WriteLog("结算线程-> 省份 %d 暂无数据处理", NetNo);

	}

	//释放返回集
	if (rs)
		OCI_ReleaseResultsets(st);

	//清除声明
	if (st)
		OCI_StatementFree(st);

	//清除连接
	if (cn)
		OCI_ConnectionFree(cn);

	//OCI_Cleanup();

	return TRUE;
}

/*************************************************
函数名称:
函数描述:更新数据(外省卡)
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
BOOL OtherProv_SetOriData(LPVOID lparam, ETC_XML_SEND_BODY_DESC *body, ETC_XML_SEND_TRANSACTION *trans, int NetNo, int CardType)
{
	if (body == NULL || trans == NULL)
	{
		ShowLstMsg1(lparam, _T("OtherProv_SetOriData: 参数传入不正确..."));
		WriteLog("结算线程-> OtherProv_SetOriData: 参数传入不正确...");
		return FALSE;
	}

	BOOL result = FALSE;
	//sql语句
	char sql[1024] = { 0 };

	OCI_Connection* cn;

	/*result = OCI_Initialize(NULL, NULL, OCI_ENV_THREADED);
	if (!result)
	{
		return FALSE;
	}*/

	//连接 第一个参数格式:【IP:端口/服务名】,第二个参数：登录用户名，第三个参数：密码
	cn = OCI_ConnectionCreate(g_config.dbservice, g_config.dbuser, g_config.dbpwd, OCI_SESSION_DEFAULT);
	if (cn == NULL)
	{
		return FALSE;
	}

	char nowtime[20] = { 0 };
	//获取当前系统时间
	SYSTEMTIME st;
	GetLocalTime(&st);
	sprintf_s(nowtime, sizeof(nowtime),"%04d-%02d-%02d %02d:%02d:%02d", st.wYear, st.wMonth, st.wDay, st.wHour, st.wMinute, st.wSecond);

	//更新交易金额为0的交易，直接置为已清分
	memset(sql, 0, sizeof(sql));
	sprintf_s(sql, sizeof(sql), "update B_TXF_JIESSJ set F_NB_XMLDBBJ = 1,F_DT_XMLDBSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss'), \
            F_NB_CHULZT = 2,F_NB_JIESJG = 1,F_NB_QINGFZT = 1,F_DT_QINGFRQ = to_date('%s','yyyy-mm-dd') where F_NB_JIESJG = 98 and F_NB_XMLDBBJ = 0 \
            and F_NB_KAWLH = %d and F_NB_KALX = %d and F_NB_JIAOYJE = 0", nowtime, nowtime, NetNo, CardType);

	OCI_Immediate(cn, sql);
	OCI_Commit(cn);  //提交数据

    //更新结算结果，打包标记，打包时间，原始包包号以及原始包交易Id
	for (int i = 0; i < body->count; i++)
	{
		memset(sql, 0, sizeof(sql));
		sprintf_s(sql, sizeof(sql), "update B_TXF_JIESSJ set F_NB_JIESJG = 1,F_NB_XMLDBBJ = 1,\
             F_DT_XMLDBSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss'),F_NB_YUANSJYBBH = %lu,F_NB_YUANSJYBJYID = %d \
             where F_DT_JIAOYSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss') and F_NB_SHOUFZBH = %lld and F_NB_CHEDH = %d \
             and F_NB_KAWLH = %d and F_NB_KALX = %d", nowtime, body->messageid, trans[i].transid, trans[i].time, trans[i].stationId, trans[i].lane, NetNo, CardType);
		int ret = OCI_Immediate(cn, sql);
		if (!ret)
		{
			WriteLog("结算线程-> 更新外省卡OCI_Immediate 错误...");
			ShowLstMsg1(lparam, _T("OCI_Immediate: 不正确..."));
		}
		ret = OCI_Commit(cn);  //提交数据
		if (!ret)
		{
			WriteLog("结算线程-> 更新外省卡 OCI_Commit 错误...");
			ShowLstMsg1(lparam, _T("OCI_Commit: 不正确..."));
		}
		Sleep(1);

	}

	//清除连接
	if (cn)
		OCI_ConnectionFree(cn);

	//OCI_Cleanup();

	return TRUE;
}



/*************************************************
函数名称:
函数描述:结算工作线程(外省储值卡)
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
BOOL OtherProv_PreCard(LPVOID lparam, ETC_XML_HEADER* header, ETC_XML_SEND_BODY_DESC* body, ETC_XML_SEND_TRANSACTION* trans, char* dirbuf, int NetWork, int CardType)
{
	BOOL result = FALSE;

	if (header == NULL || body == NULL || trans == NULL || dirbuf == NULL)
	{
		ShowLstMsg1(lparam, _T("OtherProv_PreCard: 参数传入不正确..."));
		WriteLog("结算线程->OtherProv_PreCard 参数传入不正确...");
		return FALSE;
	}

	memset((void*)header, 0, sizeof(ETC_XML_HEADER));
	memset((void*)body, 0, sizeof(ETC_XML_SEND_BODY_DESC));
	memset((void*)trans, 0, sizeof(ETC_XML_SEND_TRANSACTION) * TRANS_MAX_COUNT);

	//////////////////////////////////////////////////////////////////////////
	// 根据条件，提取oracle数据，写入缓存结构中,将creatxml集成到此接口一并完成
	//////////////////////////////////////////////////////////////////////////
	result = OtherProv_GetOriData(lparam, header, body, trans, dirbuf, NetWork, CardType);
	if (result == TRUE)
	{
		WriteLog("结算线程-> 省份 %d 储值卡结算记账源数据提取打包成功", NetWork);
		ShowLstMsg1(lparam, _T("(储值卡)结算记账源数据提取打包成功"));
	}
	else
	{
		WriteLog("结算线程-> 省份 %d 储值卡结算记账源数据提取打包失败...", NetWork);
		ShowLstMsg1(lparam, _T("(储值卡)结算记账源数据提取打包失败..."));
		return result;
	}


	//////////////////////////////////////////////////////////////////////////
	// 修改数据库
	//////////////////////////////////////////////////////////////////////////
	result = OtherProv_SetOriData(lparam, body, trans, NetWork, CardType);
	if (result == TRUE)
	{
		WriteLog("结算线程-> 省份 %d 结算记账源数据更新成功", NetWork);
		ShowLstMsg1(lparam, _T("(储值卡)记账结算源数据更新成功"));
	}
	else
	{
		WriteLog("结算线程-> 省份 %d 结算记账源数据更新失败...", NetWork);
		ShowLstMsg1(lparam, _T("(储值卡)记账结算源数据更新失败..."));
		return result;
	}

	return TRUE;
}


/*************************************************
函数名称:
函数描述:结算工作线程(外省记账卡)
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
BOOL OtherProv_CreditCard(LPVOID lparam, ETC_XML_HEADER* header, ETC_XML_SEND_BODY_DESC* body, ETC_XML_SEND_TRANSACTION* trans, char* dirbuf, int NetWork, int CardType)
{
	BOOL result = FALSE;

	if (header == NULL || body == NULL || trans == NULL || dirbuf == NULL)
	{
		ShowLstMsg1(lparam, _T("OtherProv_CreditCard: 参数传入不正确..."));
		return FALSE;
	}

	memset((void*)header, 0, sizeof(ETC_XML_HEADER));
	memset((void*)body, 0, sizeof(ETC_XML_SEND_BODY_DESC));
	memset((void*)trans, 0, sizeof(ETC_XML_SEND_TRANSACTION) * TRANS_MAX_COUNT);

	//////////////////////////////////////////////////////////////////////////
	// 根据条件，提取oracle数据，写入缓存结构中,将creatxml集成到此接口一并完成
	//////////////////////////////////////////////////////////////////////////
	result = OtherProv_GetOriData(lparam, header, body, trans, dirbuf, NetWork, CardType);
	if (result == TRUE)
	{
		WriteLog("结算线程-> 省份 %d 记账卡结算记账源数据提取打包成功", NetWork);
		ShowLstMsg1(lparam, _T("(记账卡)结算记账源数据提取打包成功"));
	}
	else
	{
		WriteLog("结算线程-> 省份 %d 记账卡结算记账源数据提取打包失败...", NetWork);
		ShowLstMsg1(lparam, _T("(记账卡)结算记账源数据提取打包失败..."));
		return result;
	}


	//////////////////////////////////////////////////////////////////////////
	// 修改数据库
	//////////////////////////////////////////////////////////////////////////
	result = OtherProv_SetOriData(lparam, body, trans, NetWork, CardType);
	if (result == TRUE)
	{
		WriteLog("结算线程-> 省份 %d 记账卡结算记账源数据更新成功", NetWork);
		ShowLstMsg1(lparam, _T("(记账卡)记账结算源数据更新成功"));
	}
	else
	{
		WriteLog("结算线程-> 省份 %d 记账卡结算记账源数据更新失败...", NetWork);
		ShowLstMsg1(lparam, _T("(记账卡)记账结算源数据更新失败..."));
		return result;
	}

	return TRUE;
}


/*************************************************
函数名称:
函数描述:结算工作线程
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
UINT SettleWorkerThread(LPVOID lparam)
{
	//此处休眠是为了给出窗口创建时间
	Sleep(ONE_SEC * 3);
	ShowLstMsg1(lparam, _T("版本-20170713,开发--东南智能,描述--完成中心数据结算"));
	WriteLog("结算线程->版本-20170713,开发--东南智能,描述--完成中心数据结算!");
	//Sleep(ONE_SEC * 2);
	ShowLstMsg1(lparam, _T("结算线程启动中..."));
	WriteLog("结算线程->结算线程启动中...");
	//Sleep(ONE_SEC * 1);
	ShowLstMsg1(lparam, _T("结算开始..."));
	WriteLog("结算线程->结算开始 ...");

	long int n;
	while (TRUE)
	{
		//文件存储路径
		char dirbuf[150] = { 0 };

		char lastfile[300] = { 0 };
		char lasttime[30] = { 0 };

		SYSTEMTIME st;
		char nowdate[12] = { 0 };
		GetLocalTime(&st);
		sprintf_s(nowdate, sizeof(nowdate), "%04d-%02d-%02d", st.wYear, st.wMonth, st.wDay);

		// 创建目录
		sprintf_s(dirbuf, sizeof(dirbuf),"%s\\Send", g_config.dirpath);
		if (ACCESS(dirbuf, 0) != 0)
		{
			//路径不存在，就创建
			MKDIR(dirbuf);
		}
		//xml文件结构内容定义
		ETC_XML_HEADER header = { 0 };
		ETC_XML_SEND_BODY_DESC bodydesc = { 0 };
		ETC_XML_SEND_TRANSACTION trans[TRANS_MAX_COUNT] = { 0 };


		int i = 0;
		for (; i < NETWORK_CODE_CNT; i++)
		{
			WriteLog("结算线程-> %d 省份储值卡结算开始", gl_network[i]);
			OtherProv_PreCard(lparam, &header, &bodydesc, trans, dirbuf, gl_network[i], PRECARD);
			WriteLog("结算线程-> %d 省份储值卡结算完成", gl_network[i]);


			WriteLog("结算线程-> %d 省份记账卡结算开始", gl_network[i]);
			OtherProv_CreditCard(lparam, &header, &bodydesc, trans, dirbuf, gl_network[i], CREDITCARD);
			WriteLog("结算线程-> %d 省份记账卡结算完成", gl_network[i]);
		}

		//每次打包完写入上次打包时间（监测程序是否运行）
		time_t now;
		time(&now);// 等同于now = time(NULL)
		n = (long int)now;

		memset(lastfile, 0, sizeof(lastfile));
		sprintf_s(lastfile, sizeof(lastfile), "%s\\lastruntime.txt", g_config.dirpath);
		FILE *fp;
		if ((fopen_s(&fp, lastfile, "w+")) != 0)
		{
			CreateLogFile(_T("打开lastruntime.txt失败"), 1);
		}
		else
		{
			memset(lasttime, 0, sizeof(lasttime));
			_itoa_s(n, lasttime, 10);
			fwrite(lasttime, 1, strlen(lasttime), fp);

			fclose(fp);
		}


		ShowLstMsg1(lparam, _T("休眠30分钟..."));
		WriteLog("结算线程->本次扫描已完成，休眠30分钟!");
		Sleep(HALF_HOUR);
	}

	return EXIT_SUCCESS;
}


/*************************************************
函数名称:
函数描述:发送文件线程
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
UINT SendFileWorkerThread(LPVOID lparam)
{
	//此处休眠是为了给出窗口创建时间
	Sleep(ONE_SEC * 3);
	ShowLstMsg2(lparam, _T("版本-20170713,开发--东南智能,描述--完成中心数据交互"));
	CreateLogFile(_T("发送线程-> 版本-20170713,开发--东南智能,描述--完成中心数据交互!"), 2);
	//Sleep(ONE_SEC * 2);
	ShowLstMsg2(lparam, _T("发送线程启动中..."));
	CreateLogFile(_T("发送线程-> 发送线程启动中!"), 2);
	//Sleep(ONE_SEC * 1);
	ShowLstMsg2(lparam, _T("开始发送..."));
	CreateLogFile(_T("发送线程-> 开始发送!"), 2);

	//文件存储路径
	char dirbuf[150] = { 0 };
	char bakdirbuf[150] = { 0 };

	BOOL result = FALSE;

	FILE *fp;
	errno_t err;

	//xml源文件定义
	char xmlfile[150] = { 0 };

	//压缩后文件定义
	char zipfile[150] = { 0 };

	//md5校验码
	unsigned char md5[33] = { 0 };

	char messageid[21] = { 0 };

	tinyxml2::XMLDocument doc;

	vector<string> files, filesname;

	char mvfile[150] = { 0 };

	while (TRUE)
	{
		//////////////////////////////////////////////////////////////////////////
		// 提取文件
		//////////////////////////////////////////////////////////////////////////
		sprintf_s(dirbuf, sizeof(dirbuf),"%s\\Send", g_config.dirpath);
		sprintf_s(bakdirbuf, sizeof(bakdirbuf),"%s\\Sendbak", g_config.dirpath);
		if (ACCESS(dirbuf, 0) != 0)
		{
			//路径不存在，就创建
			MKDIR(dirbuf);
		}

		if (ACCESS(bakdirbuf, 0) != 0)
		{
			//路径不存在，就创建
			MKDIR(bakdirbuf);
		}

		//获取Send目录下所有文件名
		files.clear();
		filesname.clear();
		GetFiles(dirbuf, files, filesname);
		
		if (files.size() <= 0)
		{
			Sleep(ONE_SEC);
			ShowLstMsg2(lparam, _T("检索发送文件为空..."));
			CreateLogFile(_T("发送线程-> 检索发送路径下文件为空..."), 2);
			goto WORK_SLEEP;
		}

		for (UINT i = 0; i < files.size(); i++)
		{
			sprintf_s(xmlfile, sizeof(xmlfile), "%s", files[i].c_str());

			if ((err = fopen_s(&fp, xmlfile, "r")) != 0)
			{
				ShowLstMsg2(lparam, _T("无法打开此xml文件..."));
				CreateLogFile(_T("发送线程-> 无法打开此xml文件!"), 2);
				continue;
			}

			//////////////////////////////////////////////////////////////////////////
			// 提取messageid
			//////////////////////////////////////////////////////////////////////////
			memset(messageid, 0, sizeof(messageid));
			doc.LoadFile(xmlfile);
			XMLElement *root = doc.RootElement();
			sprintf_s(messageid, sizeof(messageid),"%020s", root->FirstChildElement("Header")->FirstChildElement("MessageId")->GetText());
			doc.Clear();

			//////////////////////////////////////////////////////////////////////////
			// 生成md5，MD5值由<Message> …</Message> XML报文压缩前明文的二进制流进行计算获得
			//////////////////////////////////////////////////////////////////////////
			memset(md5, 0, sizeof(md5));
			result = GetMd5(lparam, xmlfile, md5);
			if (result == TRUE)
			{
				ShowLstMsg2(lparam, _T("记账结算源数据生成md5码成功|"), (char*)md5);
				CreateLogFile(_T("发送线程-> 记账结算源数据生成md5码成功!"), 2);
			}
			else
			{
				ShowLstMsg2(lparam, _T("记账结算源数据生成md5码失败...|"), messageid);
				CreateLogFile(_T("发送线程-> 记账结算源数据生成md5码失败!"), 2);
				goto WORK_SLEEP;
			}

			//////////////////////////////////////////////////////////////////////////
			// 压缩文件，LZ77算法压缩
			//////////////////////////////////////////////////////////////////////////
			memset(zipfile, 0, sizeof(zipfile));
			sprintf_s(zipfile, sizeof(zipfile), "%s.lz77", xmlfile);
			//if (ACCESS(zipfile, 0) != 0)
			//{
			//	//路径不存在，就创建
			//	MKDIR(zipfile);
			//}
			result = CompressFile(lparam, xmlfile, zipfile);
			if (result == TRUE)
			{
				ShowLstMsg2(lparam, _T("记账结算源数据xml包压缩成功|"), messageid);
				CreateLogFile(_T("发送线程-> 记账结算源数据xml包压缩成功!"), 2);
			}
			else
			{
				ShowLstMsg2(lparam, _T("记账结算源数据xml包压缩失败...|"), messageid);
				CreateLogFile(_T("发送线程-> 记账结算源数据xml包压缩失败!"), 2);
				goto WORK_SLEEP;
			}

			//////////////////////////////////////////////////////////////////////////
			// 发送文件
			//////////////////////////////////////////////////////////////////////////
			ShowLstMsg2(lparam, _T("本次发送开始"));
			CreateLogFile(_T("发送线程-> 本次发送开始!"), 2);
			result = SendCompressFile(lparam, messageid, zipfile, md5);
			if (result == TRUE)
			{
				ShowLstMsg2(lparam, _T("记账结算源数据xml包发送联网中心成功|"), messageid);
				CreateLogFile(_T("发送线程-> 记账结算源数据xml包发送联网中心成功!"), 2);
				Sleep(2000);

				//迁移文件
				sprintf_s(mvfile, sizeof(mvfile),"%s\\%s", bakdirbuf, filesname[i].c_str());
				::CopyFile((CString)xmlfile, (CString)mvfile, FALSE);
				fclose(fp);
				DeleteFile((CString)xmlfile);
				DeleteFile((CString)zipfile);
			}
			else
			{
				ShowLstMsg2(lparam, _T("记账结算源数据xml包发送联网中心失败...|"), messageid);
				CreateLogFile(_T("发送线程-> 记账结算源数据xml包发送联网中心失败!"), 2);
				fclose(fp);
				DeleteFile((CString)zipfile);
				goto WORK_SLEEP;
			}
		}

	WORK_SLEEP:

		//////////////////////////////////////////////////////////////////////////
		// 循环
		//////////////////////////////////////////////////////////////////////////
		//休眠3分钟
		ShowLstMsg2(lparam, _T("休眠3分钟..."));
		CreateLogFile(_T("发送线程-> 休眠3分钟!"), 2);
		Sleep(THREE_MIN);
	}

	return EXIT_SUCCESS;
}


//UINT DebitCardSettleThread(LPVOID lparam)
//{
//	//此处休眠是为了给出窗口创建时间
//	Sleep(ONE_SEC * 3);
//	ShowLstMsg2(lparam, _T("版本-20180117,开发--东南智能,描述--完成本省记账卡结算结果更新"));
//	CreateLogFile(_T("记账卡结算线程-> 版本-20180117,开发--东南智能,描述--完成本省记账卡结算结果更新!"), 6);
//	//Sleep(ONE_SEC * 2);
//	ShowLstMsg2(lparam, _T("记账卡结算线程启动中..."));
//	CreateLogFile(_T("结算线程-> 记账卡结算线程启动中..."), 6);
//	//Sleep(ONE_SEC * 1);
//	ShowLstMsg2(lparam, _T("记账卡结算开始..."));
//	CreateLogFile(_T("结算线程-> 记账卡结算开始 ..."), 6);
//
//	while (1)
//	{
//		//sql语句
//		char sql[1024 * 2] = { 0 };
//		char tmpsql[1024 * 2] = { 0 };
//		int nStationId;
//		int nLane;
//		
//		char transtime[20];
//
//
//		OCI_Connection* cn1;
//		OCI_Statement* st1;
//		OCI_Resultset* rs1;
//
//		OCI_Connection* cn2;
//		OCI_Statement* st2;
//		OCI_Resultset* rs2;
//
//		rs1 = NULL;
//		rs2 = NULL;//初始化，不然过不了编译
//
//
//		//连接通行宝数据库 第一个参数格式:【IP:端口/服务名】,第二个参数：登录用户名，第三个参数：密码
//		cn1 = OCI_ConnectionCreate(g_config.dbservice, g_config.dbuser, g_config.dbpwd, OCI_SESSION_DEFAULT);
//		if (cn1 == NULL)
//		{
//			CreateLogFile(_T("本省记账卡结算线程-> 通行宝数据库连接 错误..."), 6);
//			return FALSE;
//		}
//
//		st1 = OCI_StatementCreate(cn1);
//		if (st1 == NULL)
//		{
//			CreateLogFile(_T("本省记账卡结算线程-> 通行宝数据库声明 错误..."), 6);
//			return FALSE;
//		}
//
//		//连接联网中心数据库 第一个参数格式:【IP:端口/服务名】,第二个参数：登录用户名，第三个参数：密码
//		cn2 = OCI_ConnectionCreate(g_config.lwzxservice, g_config.lwzxuser, g_config.lwzxpwd, OCI_SESSION_DEFAULT);
//		if (cn2 == NULL)
//		{
//			CreateLogFile(_T("本省记账卡结算线程-> 联网中心数据库连接 错误..."), 6);
//			return FALSE;
//		}
//
//		st2 = OCI_StatementCreate(cn2);
//		if (st2 == NULL)
//		{
//			CreateLogFile(_T("本省记账卡结算线程-> 联网中心数据库声明 错误..."), 6);
//			return FALSE;
//		}
//
//		memset(sql, 0, sizeof(sql));
//		//select * from b_txf_jiessj where F_NB_SHOUFZBH = 3208260001 and F_NB_CHEDH = 1101 and F_DT_JIAOYSJ = to_date('2017/10/4 14:43:39', 'yyyy/mm/dd hh24:mi:ss')
//		sprintf_s(sql, sizeof(sql), "select F_NB_SHOUFZBH,F_NB_CHEDH,F_DT_JIAOYSJ from B_TXF_JIESSJ where F_NB_JIESJG = 98 and F_NB_XMLDBBJ = 1");
//
//		OCI_ExecuteStmt(st1, sql);
//		rs1 = OCI_GetResultset(st1);
//		while (OCI_FetchNext(rs1))
//		{
//			int result = -1;
//			nStationId = OCI_GetInt2(rs1, "F_NB_SHOUFZBH");
//			nLane = OCI_GetInt2(rs1, "F_NB_CHEDH");
//			memset(transtime, 0, sizeof(transtime));
//			sprintf_s(transtime, sizeof(transtime), "%s", OCI_GetString2(rs1, "F_DT_JIAOYSJ"));
//
//			memset(tmpsql, 0, sizeof(tmpsql));
//			sprintf_s(tmpsql, sizeof(tmpsql), "select F_NB_JIESJG from B_TXF_JIESSJ where F_NB_SHOUFZBH = %d and F_NB_CHEDH = %d and F_DT_JIAOYSJ = to_date('%s','yyyy/mm/dd hh24:mi:ss')", nStationId, nLane, transtime);
//			OCI_ExecuteStmt(st2, tmpsql);
//			rs2 = OCI_GetResultset(st2);
//			if (OCI_FetchNext(rs2))
//			{
//				result = OCI_GetInt2(rs2, "F_NB_JIESJG");
//			}
//
//			SYSTEMTIME st;
//			char nowdate[12] = { 0 };
//			GetLocalTime(&st);
//			sprintf_s(nowdate, sizeof(nowdate), "%04d-%02d-%02d", st.wYear, st.wMonth, st.wDay);
//			if (result == 1)
//			{
//				memset(sql, 0, sizeof(sql));
//				sprintf_s(sql, sizeof(sql), "update B_TXF_JIESSJ set F_NB_JIESJG = 1,F_NB_QINGFZT = 1,F_DT_QINGFRQ = to_date('%s','yyyy-mm-dd') where F_NB_SHOUFZBH = %d and F_NB_CHEDH = %d and F_DT_JIAOYSJ = to_date('%s','yyyy/mm/dd hh24:mi:ss')", nowdate, nStationId, nLane, transtime);
//				int ret = OCI_Immediate(cn1, sql);
//				if (!ret)
//				{
//					CreateLogFile(_T("更新本省记账卡记账状态 失败.."), 6);
//					ShowLstMsg1(lparam, _T("更新本省记账卡记账状态 失败..."));
//				}
//				ret = OCI_Commit(cn1);  //提交数据
//				if (!ret)
//				{
//					CreateLogFile(_T("更新本省记账卡记账状态 失败.."), 6);
//					ShowLstMsg1(lparam, _T("更新本省记账卡记账状态: 不正确..."));
//				}
//			}
//		}
//
//
//		//休眠3分钟
//		ShowLstMsg2(lparam, _T("休眠3分钟..."));
//		CreateLogFile(_T("本省记账卡结算线程-> 休眠3分钟!"), 6);
//		//释放返回集
//		if (rs1)
//			OCI_ReleaseResultsets(st1);
//
//		//清除声明
//		if (st1)
//			OCI_StatementFree(st1);
//
//		//清除连接
//		if (cn1)
//			OCI_ConnectionFree(cn1);
//
//		if (rs2)
//			OCI_ReleaseResultsets(st2);
//
//		//清除声明
//		if (st2)
//			OCI_StatementFree(st2);
//
//		//清除连接
//		if (cn2)
//			OCI_ConnectionFree(cn2);
//		Sleep(THREE_MIN);
//	}
//
//	return EXIT_SUCCESS;
//
//}

/*************************************************
函数名称:
函数描述:接收xml文件
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
UINT RecvCompressFileThread(LPVOID lparam)
{
	ShowLstMsg2(lparam, _T("版本20170713 研发：东南智能  描述：接收xml文件包..."));
	ShowLstMsg2(lparam, _T("接收线程启动..."));
	CreateLogFile(_T("接收线程->版本20170713 研发：东南智能  描述：接收xml文件包"), 3);
	CreateLogFile(_T("接收线程->接收线程启动..."), 3);
	char recvdir[300] = { 0 };
	char readyanylyze[300] = { 0 };
	char recvbakdir[300] = { 0 };
	WORD request;
	WSADATA sadata;

	//socket通讯设置
	request = MAKEWORD(1, 1);

	if (WSAStartup(request, &sadata) != 0)
	{
		ShowLstMsg2(lparam, _T("RecvCompressFile|socket通讯异常..."));
		CreateLogFile(_T("接收线程->RecvCompressFile|socket 通讯异常 WSAStartup..."), 3);
		return FALSE;
	}

	if (LOBYTE(sadata.wVersion) != 1 || HIBYTE(sadata.wVersion) != 1)
	{
		WSACleanup();
		ShowLstMsg2(lparam, _T("RecvCompressFile|socket通讯异常..."));
		CreateLogFile(_T("接收线程->RecvCompressFile|socket 通讯异常 LOBYTE..."), 3);
		return FALSE;
	}

	//捆绑本地地址
	SOCKET socksrv = socket(AF_INET, SOCK_STREAM, 0);
	SOCKADDR_IN addrsrv;
	//addrsrv.sin_addr.S_un.S_addr = htonl(INADDR_ANY);
	addrsrv.sin_addr.s_addr = htonl(INADDR_ANY);
	addrsrv.sin_family = AF_INET;
	addrsrv.sin_port = htons(SERVER_LISTEN_PORT);

	int val = -1;
	while (TRUE)
	{
		val = bind(socksrv, (LPSOCKADDR)&addrsrv, sizeof(SOCKADDR));
		if (val != 0)
		{
			ShowLstMsg2(lparam, _T("服务端 -> socket通讯捆绑失败..."));
			CreateLogFile(_T("接收线程-> socket通讯捆绑失败"), 3);
			continue;
		}
		else
		{
			break;
		}
	}

	//监听客户端
	while (TRUE)
	{
		val = listen(socksrv, 1);
		if (val != 0)
		{
			ShowLstMsg2(lparam, _T("服务端 -> socket通讯监听失败..."));
			CreateLogFile(_T("接收线程-> socket通讯监听失败"), 3);
			continue;
		}
		else
		{
			break;
		}
	}

	SOCKET sockclient;

	sprintf_s(recvdir, sizeof(recvdir),"%s\\Recv", g_config.dirpath);
	sprintf_s(recvbakdir, sizeof(recvbakdir),"%s\\Recvbak", g_config.dirpath);
	sprintf_s(readyanylyze, sizeof(readyanylyze), "%s\\Ready", g_config.dirpath);
	if (ACCESS(recvdir, 0) != 0)
	{
		//路径不存在，就创建
		MKDIR(recvdir);
	}
	if (ACCESS(readyanylyze, 0) != 0)
	{
		//路径不存在，就创建
		MKDIR(readyanylyze);
	}

	if (ACCESS(recvbakdir, 0) != 0)
	{
		//路径不存在，就创建
		MKDIR(recvbakdir);
	}

	while (TRUE)
	{
	RECONNECT_CLIENT:
		//接收客户端
		ShowLstMsg2(lparam, _T("服务端 -> 服务监听中,等待客户端连接..."));
		CreateLogFile(_T("接收线程-> 服务端：服务监听中,等待客户端连接...."), 3);
		sockclient = accept(socksrv, NULL, NULL);

	RECEIVE_FILE:
		BOOL flag = FALSE;
		FILE *fp = NULL;
		errno_t err;
		int sum = 0;
		int i;
		int xml_length;
		int result;

		//接收报文缓存
		char recvbuf[65536] = { 0 };
		//应答报文
		REPLY_STRU replymsg = { 0 };
		
		char messageid[21] = { 0 };
		char xmllength[7] = { 0 };
		char readyfile[300] = { 0 };
		char sourcename[100] = { 0 };
		char destname[300] = { 0 };
		char strcmd[310] = { 0 };

		//int nTimeout = 10000;//10秒，
		//					 //设置发送超时
		//setsockopt(socksrv, SOL_SOCKET, SO_SNDTIMEO, (char *)&nTimeout, sizeof(int));
		////设置接收超时
		//setsockopt(socksrv, SOL_SOCKET, SO_RCVTIMEO, (char *)&nTimeout, sizeof(int));

		while (1)
		{
			memset(recvbuf, 0, sizeof(recvbuf));
			result = recv(sockclient, recvbuf, sizeof(recvbuf), 0);
			if (0 > result)
			{
				//关闭文件
				if (fp == NULL)
				{
					//不处理
				}
				else
				{
					//关闭并丢弃文件
					fclose(fp);
					fp = NULL;
					if (strlen(destname) > 0)
					{
						DeleteFile((CString)destname);
					}
				}
				
				//关闭套接字
				closesocket(sockclient);
				goto RECONNECT_CLIENT;
			}
			else if (result == 0)
			{
				ShowLstMsg2(lparam, _T("服务端 -> 客户端已关闭..."));
				Sleep(2000);
				CreateLogFile(_T("接收线程-> 客户端已关闭..."), 3);
				if (fp == NULL)
				{
					;
				}
				else
				{
					//关闭并丢弃文件
					fclose(fp);
					fp = NULL;
					if (strlen(destname) > 0)
					{
						DeleteFile((CString)destname);
					}
				}
				goto RECONNECT_CLIENT;
			}
			else if (result > 0)
			{
				//文件过大时分包接收时的第一次接收
				if (!flag)
				{
					//提取messageid和xml报文长度
					strncpy_s(messageid, 21, recvbuf, 20);
					strncpy_s(xmllength, 7, recvbuf + 20, 6);
					xml_length = atoi(xmllength);

					SYSTEMTIME st;
					GetLocalTime(&st);
					memset(sourcename, 0, sizeof(sourcename));
					sprintf_s(sourcename, sizeof(sourcename),"%04d%02d%02d%02d%02d%02d%03d", st.wYear, st.wMonth, st.wDay, st.wHour, st.wMinute, st.wSecond, st.wMilliseconds);
					//sprintf_s(sourcename,sizeof(sourcename), "%s", messageid);
					memset(destname, 0, sizeof(destname));
					sprintf_s(destname, sizeof(destname), "%s\\Recv\\%s.ori", g_config.dirpath, sourcename);

					while (1)
					{
						if ((err = fopen_s(&fp, destname, "ab+")) != 0)
						{
							ShowLstMsg2(lparam, _T("创建并打开目标文件失败"));
							Sleep(1000);
							continue;
						}
						else
						{
							flag = TRUE;
							break;
						}
					}

				}

				//将本次接收到的数据写入文件
				i = fwrite(recvbuf, sizeof(char), result, fp);
				if (i>0)
				{
					ShowLstMsg2(lparam, _T("本次写入文件长度:"), i);
					//累加文件写入长度
					sum = sum + i - 58;
					//写入文件总长度除去头长度等于xml_length时认为文件写完
					if (sum == xml_length)
					{
						fclose(fp);
						fp = NULL;
						CreateLogFile(_T("接收线程-> 写入文件总长度减去报文头长度 == xml文件压缩长度"), 3);

						//拼接接收成功的应答
						char snd_buf[22] = { 0 };
						strncpy_s(snd_buf, 21, messageid, 21);
						strcat_s(snd_buf, 22, "1");
						i = send(sockclient, (char*)&snd_buf, 21, 0);
						if (i == 21)
						{
							ShowLstMsg2(lparam, _T("服务端 -> 发送应答成功: "));
							CreateLogFile(_T("接收线程-> 服务端发送接收成功的应答成功"), 3);
						}
						memset(readyfile, 0, sizeof(readyfile));
						sprintf_s(readyfile, sizeof(readyfile), "%s\\Ready\\%s.ori", g_config.dirpath, sourcename);
						::CopyFile((CString)destname, (CString)readyfile, FALSE);
						DeleteFile((CString)destname);
						goto RECEIVE_FILE;
					}
					//写入文件总长度大于xml_length时删除该文件，重新开始recv
					else if (sum > xml_length)
					{
						fclose(fp);
						fp = NULL;
						DeleteFile((CString)destname);
						//拼接接收失败的应答报文
						char snd_buf[22] = { 0 };
						strncpy_s(snd_buf, 21, messageid, 21);
						strcat_s(snd_buf, 22, "0");
						i = send(sockclient, (char*)&snd_buf, 22, 0);
						if (i == 21)
						{
							ShowLstMsg2(lparam, _T("服务端 -> 发送应答成功: "));
							CreateLogFile(_T("接收线程-> 服务端发送接收失败的应答成功 2"), 3);
						}

						goto RECEIVE_FILE;
					}
					//写入文件总长度小于xml_length时认为文件没有写完
					else if (sum < xml_length)
					{
						//什么也不做等待循环下一次recv
					}
				}
			}
		}
		//关闭文件
		fclose(fp);
		fp = NULL;
	}

	//关闭通讯
	closesocket(socksrv);
	closesocket(sockclient);
	WSACleanup();

	return TRUE;
}


/*************************************************
函数名称:
函数描述:解析xml文件
输入参数:
输出参数:
函数返回:
其它说明:
*************************************************/
UINT AnalyzeXmlFileThread(LPVOID lparam)
{
	ShowLstMsg2(lparam, _T("版本-20170713,开发--东南智能,描述--完成中心数据交互"));
	ShowLstMsg2(lparam, _T("解析线程启动"));
	CreateLogFile(_T("解析线程-> 版本-20170713,开发--东南智能,描述--解析接收到的xml文件"), 4);
	CreateLogFile(_T("解析线程-> 解析线程启动"), 4);

	//文件存储路径
	char dirbuf[300] = { 0 };
	char xmldir[300] = { 0 };
	char zipdir[300] = { 0 };
	char xmlbakdir[300] = { 0 };
	vector<string> files, filesname;
	sprintf_s(dirbuf, sizeof(dirbuf),"%s\\Ready", g_config.dirpath);
	sprintf_s(xmldir, sizeof(xmldir),"%s\\Xmlfile", g_config.dirpath);
	sprintf_s(xmlbakdir,sizeof(xmlbakdir), "%s\\Xmlfilebak", g_config.dirpath);
	sprintf_s(zipdir,sizeof(zipdir), "%s\\Zipfile", g_config.dirpath);

	if (ACCESS(dirbuf, 0) != 0)
	{
		//路径不存在，就创建
		MKDIR(dirbuf);
	}
	if (ACCESS(xmldir, 0) != 0)
	{
		//路径不存在，就创建
		MKDIR(xmldir);
	}
	if (ACCESS(xmlbakdir, 0) != 0)
	{
		//路径不存在，就创建
		MKDIR(xmlbakdir);
	}
	if (ACCESS(zipdir, 0) != 0)
	{
		//路径不存在，就创建
		MKDIR(zipdir);
	}

	while (TRUE)
	{
		//连接数据库
		OCI_Connection* cn;

		BOOL result = FALSE;
		int act, act1;

		//连接 第一个参数格式:【IP:端口/服务名】,第二个参数：登录用户名，第三个参数：密码
		cn = OCI_ConnectionCreate(g_config.dbservice, g_config.dbuser, g_config.dbpwd, OCI_SESSION_DEFAULT);
		if (cn == NULL)
		{
			CreateLogFile(_T("解析线程-> OCI_ConnectionCreate 错误!"), 4);
			return FALSE;
		}

		//获取Recv目录下所有文件名
		files.clear();
		filesname.clear();
		GetFiles(dirbuf, files, filesname);
		if (files.size() <= 0)
		{
			Sleep(ONE_SEC);
			ShowLstMsg2(lparam, _T("检索接收文件为空..."));
			CreateLogFile(_T("解析线程-> 检索接收路径下文件为空"), 4);
			
			goto WORK_SLEEP;
		}

	CREAT_ZIPFILE:
		for (UINT i = 0; i < files.size(); i++)
		{
			int xml_length;
			int sum = 0;
			char xmllength[7] = { 0 };
			char messageid[21] = { 0 };
			//存储sql语句
			char sql[1024] = { 0 };
			//相关文件名
			char orifilename[50] = { 0 };
			char sourcefile[300] = { 0 };
			char zipfile[300] = { 0 };
			char destfile[300] = { 0 };
			char mvfile[300] = { 0 };
			char mvorifile[300] = { 0 };
			char xmlbuf[65536];

			FILE *fp_s;//源文件指针
			FILE *fp_d;//目的文件指针
			errno_t err;
			
			tinyxml2::XMLDocument doc;
			sprintf_s(sourcefile, sizeof(sourcefile),"%s", files[i].c_str());

			if ((err = fopen_s(&fp_s, sourcefile, "rb")) != 0)
			{
				ShowLstMsg2(lparam, _T("无法打开此源文件..."));
				CreateLogFile(_T("解析线程-> 无法打开此原始文件!"), 4);
				memset(mvorifile, 0, sizeof(mvorifile));
				sprintf_s(mvorifile, sizeof(mvorifile), "%s\\Recvbak\\%s", g_config.dirpath, filesname[i].c_str());
				::CopyFile((CString)files[i].c_str(), (CString)mvorifile, FALSE);
				DeleteFile((CString)files[i].c_str());
				continue;
			}
			//获取xml文件的压缩长度和messageid
			fseek(fp_s, 20, SEEK_SET);
			act = fread(xmllength, sizeof(char), 6, fp_s);
			if (act > 0)
			{
				ShowLstMsg2(lparam, _T("获取报文长度成功..."));
				CreateLogFile(_T("解析线程-> 获取报文长度成功!"), 4);
			}

			xml_length = atoi(xmllength);

			fseek(fp_s, 0, SEEK_SET);
			act = fread(messageid, sizeof(char), 20, fp_s);
			if (act > 0)
			{
				ShowLstMsg2(lparam, _T("获取messageid成功..."));
				CreateLogFile(_T("解析线程-> 获取messageid 成功!"), 4);
			}

			//拼压缩xml文件名,创建并打开
			snprintf(orifilename, 21, "%s",messageid);
			//snprintf(orifilename, 15, "%s", filesname[i].c_str());
			sprintf_s(zipfile, sizeof(zipfile),"%s\\Zipfile\\%s.lz77", g_config.dirpath, orifilename);
			if ((err = fopen_s(&fp_d, zipfile, "ab+")) != 0)
			{
				ShowLstMsg2(lparam, _T("无法创建并打开此xml压缩文件..."));
				CreateLogFile(_T("解析线程-> 无法打开解析出的xml压缩文件"), 4);
				Sleep(3000);
				continue;
			}

			//提取xml报文正文写入xmlbuf
			fseek(fp_s, 58, SEEK_SET);
			while (1)
			{
				memset(xmlbuf, 0, sizeof(xmlbuf));
				xml_length = min(xml_length, 65536);
				act1 = fread(xmlbuf, sizeof(char), xml_length, fp_s);
				//result = fread(xmlbuf, sizeof(char), sizeof(xmlbuf), fp_s);
				if (act1 > 0)
				{
					ShowLstMsg2(lparam, _T("读入xmlbuf的文件长度:"), act1);
				}

				act = fwrite(xmlbuf, sizeof(char), act1, fp_d);
				if (act > 0)
				{
					ShowLstMsg2(lparam, _T("写入压缩xml文件的长度:"), act);
					sum = sum + act;
					//写入文件长度等于压缩长度，说明已写完
					if (sum == xml_length)
					{
						fclose(fp_s);
						fclose(fp_d);
						CreateLogFile(_T("解析线程-> 写入文件长度等于压缩长度信息"), 4);
						break;
					}
					//写入文件长度小于压缩长度，继续写入
					else if (sum < xml_length)
					{
						CreateLogFile(_T("解析线程-> 写入文件长度小于压缩长度信息，继续写入"), 4);
						continue;
					}
					//写入文件长度大于压缩长度，出错
					else if (sum > xml_length)
					{
						CreateLogFile(_T("解析线程-> 写入文件长度大于压缩长度信息，出错"), 4);
						fclose(fp_s);
						fclose(fp_d);
						DeleteFile((CString)zipfile);
						goto CREAT_ZIPFILE;
					}
				}
				else if (act <= 0)
				{
					fclose(fp_s);
					fclose(fp_d);
					goto CREAT_ZIPFILE;
				}
			}

			//拼接xml文件名
			sprintf_s(destfile, sizeof(destfile),"%s\\Xmlfile\\%s.xml", g_config.dirpath, orifilename);
			
			//解压缩文件
			result = UnCompressFile(lparam, zipfile, destfile);
			if (result != TRUE)
			{
				ShowLstMsg2(lparam, _T("解压xml文件出错"));
				CreateLogFile(_T("解析线程-> 对解析出的压缩文件进行解压出错 !"), 4);
				continue;
			}

			//if (doc.LoadFile(sourcefile) != XML_NO_ERROR)
			if (doc.LoadFile(destfile) != XML_NO_ERROR)
			{
				ShowLstMsg2(lparam, _T("加载xml文件出错"));
				CreateLogFile(_T("解析线程-> 加载xml文件出错"), 4);
				continue;
			}
			CreateLogFile(_T("解析线程-> 加载xml文件成功"), 4);

			int iMessageType;
			int iMessageClass;
			char sMessageType[2] = { 0 };
			char sMessageClass[2] = { 0 };
			XMLElement *root = doc.RootElement();
			XMLElement *body = root->FirstChildElement("Body");
			XMLElement *header = root->FirstChildElement("Header");
			const XMLAttribute *typeAttr = body->FirstAttribute();
			string contenttype = typeAttr->Value();

			snprintf(sMessageType, 2, "%s", root->FirstChildElement("Header")->FirstChildElement("MessageType")->GetText());
			snprintf(sMessageClass, 2, "%s", root->FirstChildElement("Header")->FirstChildElement("MessageClass")->GetText());
			iMessageType = atoi(sMessageType);
			iMessageClass = atoi(sMessageClass);

			if ((iMessageClass == 6) && (iMessageType == 7) && (contenttype == "1"))
			{
				//外省原始应答包要处理，本省不处理
				CreateLogFile(_T("解析线程->  原始应答包处理流程"), 4);
				char ResMessageid[21] = { 0 };
				unsigned long int nResMessageid = 0;
				unsigned long int nMessageid = 0;
				char Processtime[20] = { 0 };
				char ChargeTime[20] = { 0 };
				char sResult[2] = { 0 };
				int  nResult;

				memset(messageid, 0, sizeof(messageid));
				snprintf(messageid, 21, "%s", body->FirstChildElement("MessageId")->GetText());
				nMessageid = atol(messageid);
				snprintf(ResMessageid, 21, "%s", header->FirstChildElement("MessageId")->GetText());
				nResMessageid = atol(ResMessageid);
				snprintf(ChargeTime, 20, "%s", body->FirstChildElement("ProcessTime")->GetText());
				snprintf(sResult, 2, "%s", body->FirstChildElement("Result")->GetText());
				nResult = atoi(sResult);

				//原始应答为1更新应答包号
				if (nResult == 1)
				{
					//处理时间2017-07-07T11:33:53转为 2017-07-07 11:33:53
					snprintf(Processtime, 11, "%s", ChargeTime);
					strcat_s(Processtime, 20, " ");
					strcat_s(Processtime, 20, ChargeTime + 11);
					//外省只更新外省的应答包，我们自己做一层保护，实际是不应该发过来本省应答包
					sprintf_s(sql, "update B_TXF_JIESSJ set F_NB_YUANSYDBBH = %lu,F_DT_YUANSYDBCLSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss') \
								 where F_NB_YUANSJYBBH = %lu and F_NB_KAWLH <> 3201", nResMessageid, Processtime, nMessageid);

					int ret = OCI_Immediate(cn, sql);
					if (!ret)
					{
						CreateLogFile(_T("外省卡原始应答包OCI_Immediate 错误.."), 4);
						ShowLstMsg1(lparam, _T("外省卡原始应答包OCI_Immediate: 不正确..."));
					}
					ret = OCI_Commit(cn);  //提交数据
					if (!ret)
					{
						CreateLogFile(_T("外省卡原始应答包OCI_Commit 错误.."), 4);
						ShowLstMsg1(lparam, _T("外省卡原始应答包OCI_Commit: 不正确..."));
					}
				}
				else
				{
					//外省只更新外省的应答包，我们自己做一层保护，实际是不应该发过来本省应答包
					memset(sql, 0, sizeof(sql));
					sprintf_s(sql, "update B_TXF_JIESSJ set F_NB_YUANSYDBBH = %lu,F_NB_CHULZT = 100,F_NB_BEIY1 = 1\
								 where F_NB_YUANSJYBBH = %lu and F_NB_KAWLH <> 3201", nResMessageid, nMessageid);

					int ret = OCI_Immediate(cn, sql);
					if (!ret)
					{
						CreateLogFile(_T("外省卡原始应答包OCI_Immediate 错误.."), 4);
						ShowLstMsg1(lparam, _T("外省卡原始应答包OCI_Immediate: 不正确..."));
					}
					ret = OCI_Commit(cn);  //提交数据
					if (!ret)
					{
						CreateLogFile(_T("外省卡原始应答包OCI_Commit 错误.."), 4);
						ShowLstMsg1(lparam, _T("外省卡原始应答包OCI_Commit: 不正确..."));
					}
				}

			}
			else if ((iMessageClass == 5) && (iMessageType == 5) && (contenttype == "1"))
			{
				//记账包处理
				CreateLogFile(_T("解析线程-> 进入记账包处理流程"), 4);
				char ChargeTime[20] = { 0 };
				char JzMessageid[21] = { 0 };
				unsigned long int nJzMessageid = 0;
				unsigned long int nMessageId = 0;
				char nowtime[20] = { 0 };
				char processtime[20] = { 0 };

				memset(JzMessageid,0,sizeof(JzMessageid));
				snprintf(JzMessageid, 21, "%s", header->FirstChildElement("MessageId")->GetText());
				nJzMessageid = atol(JzMessageid);
				memset(messageid, 0, sizeof(messageid));
				snprintf(messageid, 21, "%s", body->FirstChildElement("MessageId")->GetText());
				nMessageId = atol(messageid);
				memset(ChargeTime, 0, sizeof(ChargeTime));
				snprintf(ChargeTime, 20, "%s", body->FirstChildElement("ProcessTime")->GetText());
				//处理时间2017-07-07T11:33:53转为 2017-07-07 11:33:53
				memset(processtime,0,sizeof(processtime));
				snprintf(processtime, 11, "%s", ChargeTime);
				strcat_s(processtime, 20, " ");
				strcat_s(processtime, 20, ChargeTime + 11);
				SYSTEMTIME st;
				GetLocalTime(&st);
				sprintf_s(nowtime, "%04d-%02d-%02d %02d:%02d:%02d", st.wYear, st.wMonth, st.wDay, st.wHour, st.wMinute, st.wSecond);

				memset(sql, 0, sizeof(sql));
				sprintf_s(sql, sizeof(sql), "update B_TXF_JIESSJ set F_NB_JIZBBH = %lu,F_DT_JIZBCLSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss'),\
								F_NB_CHULZT = 2,F_NB_JIESJG = 1,F_DT_QUERZFSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss') \
                                where F_NB_YUANSJYBBH = %lu", nJzMessageid, processtime, nowtime, nMessageId);

				int ret = OCI_Immediate(cn, sql);
				if (!ret)
				{
					CreateLogFile(_T("外省卡更新记账状态OCI_Immediate 错误.."), 4);
					ShowLstMsg1(lparam, _T("外省卡更新记账状态OCI_Immediate: 不正确..."));
				}

				ret = OCI_Commit(cn);  //提交数据
				if (!ret)
				{
					CreateLogFile(_T("外省卡更新记账状态OCI_Commit 错误.."), 4);
					ShowLstMsg1(lparam, _T("外省卡更新记账状态OCI_Commit: 不正确..."));
				}
				XMLElement *disputedrecord = body->FirstChildElement("DisputedRecord");
				while (disputedrecord)
				{
					char transid[11] = { 0 };
					char DisputedType[4] = { 0 };
					snprintf(transid, 11, "%s", disputedrecord->FirstChildElement("TransId")->GetText());
					snprintf(DisputedType, 4, "%s", disputedrecord->FirstChildElement("Result")->GetText());

					memset(sql, 0, sizeof(sql));
					sprintf_s(sql, sizeof(sql), "update B_TXF_JIESSJ set F_NB_JIZBBH = '%s',F_DT_JIZBCLSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss'),\
									F_NB_CHULZT = 0,F_NB_JIESJG = 2,F_NB_ZHENGYLX = '%s',F_DT_QUERZFSJ = null \
									where F_NB_YUANSJYBBH = '%s' and F_NB_YUANSJYBJYID = '%s'", JzMessageid, processtime, DisputedType, messageid, transid);

					ret = OCI_Immediate(cn, sql);
					if (!ret)
					{
						CreateLogFile(_T("外省卡更新争议状态OCI_Immediate 错误.."), 4);
						ShowLstMsg1(lparam, _T("外省卡更新争议状态OCI_Immediate: 不正确..."));
					}
				    ret = OCI_Commit(cn);  //提交数据
					if (!ret)
					{
						CreateLogFile(_T("外省卡更新争议状态OCI_Commit 错误.."), 4);
						ShowLstMsg1(lparam, _T("外省卡更新争议状态OCI_Commit: 不正确..."));
					}

					disputedrecord = disputedrecord->NextSiblingElement();
				}
			}
			else if ((iMessageClass == 5) && (iMessageType == 7) && (contenttype == "2"))
			{
				//争议包处理
				CreateLogFile(_T("解析线程-> 进入争议包处理流程"), 4);
				char ProcessTime[20] = { 0 };
				char nowtime[20] = { 0 };
				char nowdate[11] = { 0 };
				char ZyMessageid[21] = { 0 };
				unsigned long int nZyMessageid = 0;
				char Handletime[20] = { 0 };

				SYSTEMTIME st;
				GetLocalTime(&st);
				sprintf_s(nowtime, sizeof(nowtime),"%04d-%02d-%02d %02d:%02d:%02d", st.wYear, st.wMonth, st.wDay, st.wHour, st.wMinute, st.wSecond);
				sprintf_s(nowdate, "%04d-%02d-%02d", st.wYear, st.wMonth, st.wDay);
				snprintf(ZyMessageid, 21, "%s", header->FirstChildElement("MessageId")->GetText());
				nZyMessageid = atol(ZyMessageid);
				snprintf(ProcessTime, 20, "%s", body->FirstChildElement("ProcessTime")->GetText());

				//处理时间2017-07-07T11:33:53转为 2017-07-07 11:33:53
				snprintf(Handletime, 11, "%s", ProcessTime);
				strcat_s(Handletime, 20, " ");
				strcat_s(Handletime, 20, ProcessTime + 11);

				XMLElement *messagelist = body->FirstChildElement("MessageList");
				while (messagelist)
				{
					char transid[9] = { 0 };
					char Result[2];
					int nResult;
					int nSettleResult;
					unsigned long int nMessageId = 0;
					memset(messageid, 0, sizeof(messageid));
					snprintf(messageid, 21, "%s", messagelist->FirstChildElement("MessageId")->GetText());
					nMessageId = atol(messageid);

					XMLElement *transaction = messagelist->FirstChildElement("Transaction");
					while (transaction)
					{
						int nTransId = 0;

						memset(transid,0,sizeof(transid));
						snprintf(transid, 9, "%s", transaction->FirstChildElement("TransId")->GetText());
						nTransId = atoi(transid);

						memset(Result,0,sizeof(Result));
						snprintf(Result, 2, "%s", transaction->FirstChildElement("Result")->GetText());

						//报文中result为0表示正常支付；为1表示此交易作坏账处理
						nResult = atoi(Result);
						if (nResult == 0)
						{
							nSettleResult = 1;
							memset(sql, 0, sizeof(sql));
							sprintf_s(sql, sizeof(sql), "update B_TXF_JIESSJ set F_NB_ZHENGYCLJGBBH = %lu,F_DT_ZHENGYBCLSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss'),\
										F_NB_JIESJG = %d,F_NB_QINGFZT = 1,F_DT_QINGFRQ = to_date('%s','yyyy-mm-dd hh24:mi:ss') ,F_DT_QUERZFSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss') where F_NB_YUANSJYBBH = %lu \
										and F_NB_YUANSJYBJYID = %d", nZyMessageid, Handletime, nSettleResult, nowtime, nowtime, nMessageId, nTransId);

						}
						else if (nResult == 1)
						{
							nSettleResult = 3;
							memset(sql, 0, sizeof(sql));
							sprintf_s(sql, sizeof(sql), "update B_TXF_JIESSJ set F_NB_ZHENGYCLJGBBH = %lu,F_DT_ZHENGYBCLSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss'),\
										F_NB_JIESJG = %d,F_DT_QUERZFSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss') where F_NB_YUANSJYBBH = %lu \
										and F_NB_YUANSJYBJYID = %d", nZyMessageid, Handletime, nSettleResult, nowtime, nMessageId, nTransId);
						}

						//memset(sql, 0, sizeof(sql));
						////外省卡原流程
						//sprintf_s(sql, sizeof(sql),"update B_TXF_JIESSJ set F_NB_ZHENGYCLJGBBH = %lu,F_DT_ZHENGYBCLSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss'),\
						//				F_NB_JIESJG = %d,F_DT_QUERZFSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss') where F_NB_YUANSJYBBH = %lu \
						//				and F_NB_YUANSJYBJYID = %d", nZyMessageid, Handletime, nSettleResult, nowtime, nMessageId, nTransId);

						int ret = OCI_Immediate(cn, sql);
						if (!ret)
						{
							CreateLogFile(_T("外省卡更新争议结果OCI_Immediate 错误.."), 4);
							ShowLstMsg1(lparam, _T("外省卡更新争议结果OCI_Immediate: 不正确..."));
						}
						ret = OCI_Commit(cn);  //提交数据
						if (!ret)
						{
							CreateLogFile(_T("外省卡更新争议结果OCI_Commit 错误.."), 4);
							ShowLstMsg1(lparam, _T("外省卡更新争议结果OCI_Commit: 不正确..."));
						}
						transaction = transaction->NextSiblingElement();
					}

					messagelist = messagelist->NextSiblingElement();
				}
			}
			else if ((iMessageClass == 5) && (iMessageType == 5) && (contenttype == "2"))
			{
				//清分包处理
				CreateLogFile(_T("解析线程-> 进入清分包处理流程"), 4);
				char QfMessageid[21] = { 0 };
				unsigned long int nQfMessageid = 0;
				char ProcessTime[20] = { 0 };
				char nowtime[20] = { 0 };
				char Handletime[20] = { 0 };

				//SYSTEMTIME st;
				//GetLocalTime(&st);
				//sprintf_s(nowtime, sizeof(nowtime),"%04d-%02d-%02d %02d:%02d:%02d", st.wYear, st.wMonth, st.wDay, st.wHour, st.wMinute, st.wSecond);

				snprintf(QfMessageid, 21, "%s", header->FirstChildElement("MessageId")->GetText());
				nQfMessageid = atol(QfMessageid);

				snprintf(ProcessTime, 20,"%s", body->FirstChildElement("ProcessTime")->GetText());

				snprintf(nowtime, 20, "%s", body->FirstChildElement("ClearTargetDate")->GetText());

				//处理时间报文中是 2017-07-07T11:33:53转为 2017-07-07 11:33:53方便to_date()转换123456
				snprintf(Handletime, 11, "%s", ProcessTime);
				strcat_s(Handletime, 20, " ");
				strcat_s(Handletime, 20, ProcessTime + 11);

				XMLElement *list = body->FirstChildElement("List");
				XMLElement *Messageid = list->FirstChildElement("MessageId");
				while (Messageid)
				{
					unsigned long int nMessageId = 0;
					
					if (strncmp(Messageid->Value(), "MessageId", 9))
					{
						break;
					}
					memset(messageid, 0, sizeof(messageid));
					sprintf_s(messageid, sizeof(messageid),"%s", Messageid->GetText());
					nMessageId = atol(messageid);

					memset(sql, 0, sizeof(sql));
					sprintf_s(sql,sizeof(sql), "update B_TXF_JIESSJ set F_NB_QINGFBBH = %lu,F_DT_QINGFBCLSJ = to_date('%s','yyyy-mm-dd hh24:mi:ss'),\
									F_NB_QINGFZT = 1,F_DT_QINGFRQ = to_date('%s','yyyy-mm-dd') where F_NB_YUANSJYBBH = %lu \
									and F_NB_JIESJG = 1", nQfMessageid, Handletime, nowtime, nMessageId);

					int ret = OCI_Immediate(cn, sql);
					if (!ret)
					{
						CreateLogFile(_T("外省卡更新清分状态OCI_Immediate 错误.."), 4);
						ShowLstMsg1(lparam, _T("外省卡更新清分状态OCI_Immediate: 不正确..."));
					}
					ret = OCI_Commit(cn);  //提交数据
					if (!ret)
					{
						CreateLogFile(_T("外省卡更新清分状态OCI_Commit 错误.."), 4);
						ShowLstMsg1(lparam, _T("外省卡更新清分状态OCI_Commit: 不正确..."));
					}

					Messageid = Messageid->NextSiblingElement();
				}
			}
			else
			{
				ShowLstMsg2(lparam, _T("该xml文件不是原始应答包、记账包、清分包、争议包中的一种"));
				CreateLogFile(_T("解析线程-> 该xml文件不是原始应答包、记账包、清分包、争议包中的一种 "), 4);
			}

			//清除加载的文件
			doc.Clear();
			//关闭文件
			fclose(fp_s);
			fclose(fp_d);
			fp_s = NULL;
			fp_d = NULL;

			//转移已经处理好的xml文件
			memset(mvfile,0,sizeof(mvfile));
			sprintf_s(mvfile, sizeof(mvfile),"%s\\Xmlfilebak\\%s.xml", g_config.dirpath, orifilename);
			::CopyFile((CString)destfile, (CString)mvfile, FALSE);
			DeleteFile((CString)destfile);

			//转移recv文件夹下已处理的原始文件
			memset(mvorifile, 0, sizeof(mvorifile));
			sprintf_s(mvorifile, sizeof(mvorifile),"%s\\Recvbak\\%s", g_config.dirpath, filesname[i].c_str());
			::CopyFile((CString)files[i].c_str(), (CString)mvorifile, FALSE);
			DeleteFile((CString)files[i].c_str());
			
		}


	WORK_SLEEP:

		//////////////////////////////////////////////////////////////////////////
		// 循环
		//////////////////////////////////////////////////////////////////////////
		//休眠3分钟

		//清除连接
		if (cn)
			OCI_ConnectionFree(cn);

		//OCI_Cleanup();

		ShowLstMsg2(lparam, _T("休眠3分钟..."));
		CreateLogFile(_T("解析线程-> 休眠三分钟 !"), 4);
		Sleep(THREE_MIN);
	}

	return EXIT_SUCCESS;
}









