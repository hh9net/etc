/*****************************
** LZ77.h
*****************************/

/* 使用在自己的堆中分配索引节点，不滑动窗口
** 每次最多压缩 65536 字节数据
** 的优化版本*/


#ifndef _WIX_LZ77_COMPRESS_HEADER_001_
#define _WIX_LZ77_COMPRESS_HEADER_001_

/* 滑动窗口的字节大小*/
#define _MAX_WINDOW_SIZE	65536
#define min(a,b)            (((a) < (b)) ? (a) : (b))
#define FALSE               0
#define TRUE                1



typedef unsigned long       DWORD;
typedef int                 BOOL;
typedef unsigned char       BYTE;
typedef unsigned short      WORD;

class CCompressLZ77
{
public:


	typedef union _UDWORD { DWORD dw; BYTE b[4]; } UDWORD;

	/*******************************************************
	** CopyBitsInAByte : 在一个字节范围内复制位流
	** 参数含义同 CopyBits 的参数
	** 说明：
	**		此函数由 CopyBits 调用，不做错误检查，即
	**		假定要复制的位都在一个字节范围内*/
	void CopyBitsInAByte(BYTE* memDest, int nDestPos, BYTE* memSrc, int nSrcPos, int nBits);

	/*******************************************************
	** CopyBits : 复制内存中的位流
	**		memDest - 目标数据区
	**		nDestPos - 目标数据区第一个字节中的起始位
	**		memSrc - 源数据区
	**		nSrcPos - 源数据区第一个字节的中起始位
	**		nBits - 要复制的位数
	**	说明：
	**		起始位的表示约定为从字节的高位至低位（由左至右）
	**		依次为 0，1，... , 7
	**		要复制的两块数据区不能有重合*/
	void CopyBits(BYTE* memDest, int nDestPos, BYTE* memSrc, int nSrcPos, int nBits);

	/*************************************************************
	** 将DWORD值从高位字节到低位字节排列*/
	void InvertDWord(DWORD* pDW);

	/***********************************************************
	** 设置byte的第iBit位为aBit
	**		iBit顺序为高位起从0记数（左起）*/
	void SetBit(BYTE* byte, int iBit, BYTE aBit);

	/***********************************************************
	** 得到字节byte第pos位的值
	**		pos顺序为高位起从0记数（左起）*/
	BYTE GetBit(BYTE byte, int pos);

	/***********************************************************
	** 将位指针*piByte(字节偏移), *piBit(字节内位偏移)后移num位*/
	void MovePos(int* piByte, int* piBit, int num);

	/*******************************************************
	** 取log2(n)的upper_bound*/
	int UpperLog2(int n);

	/*******************************************************
	** 取log2(n)的lower_bound*/
	int LowerLog2(int n);



	/*******************************************
	** 压缩一段字节流
	** src - 源数据区
	** srclen - 源数据区字节长度, srclen <= 65536
	** dest - 压缩数据区，调用前分配srclen字节内存
	** 返回值 > 0 压缩数据长度
	** 返回值 = 0 数据无法压缩
	** 返回值 < 0 压缩中异常错误*/
	int Compress(BYTE* src, int srclen, BYTE* dest);

	/*******************************************
	** 解压缩一段字节流
	** src - 接收原始数据的内存区, srclen <= 65536
	** srclen - 源数据区字节长度
	** dest - 压缩数据区
	** 返回值 - 成功与否*/
	BOOL Decompress(BYTE* src, int srclen, BYTE* dest);


	BYTE* pWnd;
	/* 窗口大小最大为 64k ，并且不做滑动
	** 每次最多只压缩 64k 数据，这样可以方便从文件中间开始解压
	** 当前窗口的长度*/
	int nWndSize;

	/* 对滑动窗口中每一个2字节串排序
	** 排序是为了进行快速术语匹配
	** 排序的方法是用一个64k大小的指针数组
	** 数组下标依次对应每一个2字节串：(00 00) (00 01) ... (01 00) (01 01) ...
	** 每一个指针指向一个链表，链表中的节点为该2字节串的每一个出现位置*/
	typedef struct _STIDXNODE
	{
		WORD off;		/* 在src中的偏移*/
		WORD off2;		/* 用于对应的2字节串为重复字节的节点*/
						/* 指从 off 到 off2 都对应了该2字节串*/
		WORD next;		/* 在SortHeap中的指针*/
	} STIDXNODE;

	WORD SortTable[65536];  /* 256 * 256 指向SortHeap中下标的指针*/

	/* 因为窗口不滑动，没有删除节点的操作，所以
	** 节点可以在SortHeap 中连续分配*/

	STIDXNODE SortHeap[_MAX_WINDOW_SIZE];
	int HeapPos;	/* 当前分配位置*/

	/* 当前输出位置(字节偏移及位偏移)*/
	int CurByte, CurBit;

	/***************************************
	** 输出压缩码
	** code - 要输出的数
	** bits - 要输出的位数(对isGamma=TRUE时无效)
	** isGamma - 是否输出为γ编码*/
	void _OutCode(BYTE* dest, DWORD code, int bits, BOOL isGamma);

	/*********************************************************
	** 在滑动窗口中查找术语
	** nSeekStart - 从何处开始匹配
	** offset, len - 用于接收结果，表示在滑动窗口内的偏移和长度
	** 返回值- 是否查到长度为3或3以上的匹配字节串*/
	BOOL _SeekPhase(BYTE* src, int srclen, int nSeekStart, int* offset, int* len);

	/*********************************************************
	** 得到已经匹配了3个字节的窗口位置offset
	** 共能匹配多少个字节*/
	int _GetSameLen(BYTE* src, int srclen, int nSeekStart, int offset);

	/*****************************************
	** 将窗口向右滑动n个字节*/
	void _ScrollWindow(int n);

	/* 向索引中添加一个2字节串*/
	void _InsertIndexItem(int off);

	/* 初始化索引表，释放上次压缩用的空间*/
	void _InitSortTable();

	int GetDecompressBufLen(unsigned char *input, int inputlen);
	int CompressFile(char *inputfile, char *outputfile);
	int CompressBuf(unsigned char *input, int inputlen, unsigned char *output);
	int DecompressFile(char *inputfile, char *outputfile);
	int DecompressBuf(unsigned char *input, int inputlen, unsigned char *output);
};


extern "C" int Compressfile(char *inputfile, char *outputfile);
extern "C" int Decompressfile(char *inputfile, char *outputfile);
extern "C" int Testcode();

#endif /* _WIX_LZW_COMPRESS_HEADER_001_*/