/*****************************
** LZ77.CPP
******************************/

#include "Lz77.h"
#include <stdio.h>
#include <arpa/inet.h>
#include <string.h>

/*******************************************************
** 取log2(n)的upper_bound*/
int CCompressLZ77::UpperLog2(int n)
{
	int i = 0;
	if (n > 0)
	{
		int m = 1;
		while(1)
		{
			if (m >= n)
				return i;
			m <<= 1;
			i++;
		}
	}
	else 
		return -1;
}
/* UpperLog2
********************************************************/

/*******************************************************
** 取log2(n)的lower_bound*/
int CCompressLZ77::LowerLog2(int n)
{
	int i = 0;
	if (n > 0)
	{
		int m = 1;
		while(1)
		{
			if (m == n)
				return i;
			if (m > n)
				return i - 1;
			m <<= 1;
			i++;
		}
	}
	else 
		return -1;
}
/* LowerLog2
********************************************************/

/***********************************************************
** 将位指针*piByte(字节偏移), *piBit(字节内位偏移)后移num位*/
void CCompressLZ77::MovePos(int* piByte, int* piBit, int num)
{
	num += (*piBit);
	(*piByte) += num / 8;
	(*piBit) = num % 8;
}
/* MovePos
************************************************************/

/***********************************************************
** 得到字节byte第pos位的值
**		pos顺序为高位起从0记数（左起）*/
BYTE CCompressLZ77::GetBit(BYTE byte, int pos)
{
	int j = 1;
	j <<= 7 - pos;
	if (byte & j)
		return 1;
	else 
		return 0;
}
/* GetBit
************************************************************/

/***********************************************************
** 设置byte的第iBit位为aBit
**		iBit顺序为高位起从0记数（左起）*/
void CCompressLZ77::SetBit(BYTE* byte, int iBit, BYTE aBit)
{
	if (aBit)
		(*byte) |= (1 << (7 - iBit));
	else
		(*byte) &= ~(1 << (7 - iBit));
}
/* SetBit
**************************************************************/

/*************************************************************
** 将DWORD值从高位字节到低位字节排列*/
void CCompressLZ77::InvertDWord(DWORD* pDW)
{
	*pDW = htonl(*pDW);
	return;
	/*
	UDWORD* pUDW = (UDWORD*)pDW;
	BYTE b;
	b = pUDW->b[0];	pUDW->b[0] = pUDW->b[3]; pUDW->b[3] = b;
	b = pUDW->b[1];	pUDW->b[1] = pUDW->b[2]; pUDW->b[2] = b;
	*/
}
/* InvertDWord
**************************************************************/

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
void CCompressLZ77::CopyBits(BYTE* memDest, int nDestPos, BYTE* memSrc, int nSrcPos, int nBits)
{
	int iByteDest = 0, iBitDest;
	int iByteSrc = 0, iBitSrc = nSrcPos;

	int nBitsToFill, nBitsCanFill;

	while (nBits > 0)
	{
		/* 计算要在目标区当前字节填充的位数*/
		nBitsToFill = min(nBits, iByteDest ? 8 : 8 - nDestPos);
		/* 目标区当前字节要填充的起始位*/
		iBitDest = iByteDest ? 0 : nDestPos;
		/* 计算可以一次从源数据区中复制的位数*/
		nBitsCanFill = min(nBitsToFill, 8 - iBitSrc);
		/* 字节内复制*/
		CopyBitsInAByte(memDest + iByteDest, iBitDest, 
			memSrc + iByteSrc, iBitSrc, nBitsCanFill);		
		/* 如果还没有复制完 nBitsToFill 个*/
		if (nBitsToFill > nBitsCanFill)
		{
			iByteSrc++; iBitSrc = 0; iBitDest += nBitsCanFill;
			CopyBitsInAByte(memDest + iByteDest, iBitDest, 
					memSrc + iByteSrc, iBitSrc, 
					nBitsToFill - nBitsCanFill);
			iBitSrc += nBitsToFill - nBitsCanFill;
		}
		else 
		{
			iBitSrc += nBitsCanFill;
			if (iBitSrc >= 8)
			{
				iByteSrc++; iBitSrc = 0;
			}
		}

		nBits -= nBitsToFill;	/* 已经填充了nBitsToFill位*/
		iByteDest++;
	}	
}
/* CopyBits
********************************************************/

/*******************************************************
** CopyBitsInAByte : 在一个字节范围内复制位流
** 参数含义同 CopyBits 的参数
** 说明：
**		此函数由 CopyBits 调用，不做错误检查，即
**		假定要复制的位都在一个字节范围内*/
void CCompressLZ77::CopyBitsInAByte(BYTE* memDest, int nDestPos, BYTE* memSrc, int nSrcPos, int nBits)
{
	BYTE b1, b2;
	b1 = *memSrc;
	b1 <<= nSrcPos; b1 >>= 8 - nBits;	/* 将不用复制的位清0*/
	b1 <<= 8 - nBits - nDestPos;		/* 将源和目的字节对齐*/
	*memDest |= b1;		/* 复制值为1的位*/
	b2 = 0xff; b2 <<= 8 - nDestPos;		/* 将不用复制的位置1*/
	b1 |= b2;
	b2 = 0xff; b2 >>= nDestPos + nBits;
	b1 |= b2;
	*memDest &= b1;		/* 复制值为0的位*/
}
/* CopyBitsInAByte
********************************************************/


/*------------------------------------------------------------------*/


/* 初始化索引表，释放上次压缩用的空间*/
void CCompressLZ77::_InitSortTable()
{
	memset(SortTable, 0, sizeof(WORD) * 65536);
	nWndSize = 0;
	HeapPos = 1;
}

/* 向索引中添加一个2字节串*/
void CCompressLZ77::_InsertIndexItem(int off)
{
	WORD q;
	BYTE ch1, ch2;
	ch1 = pWnd[off]; ch2 = pWnd[off + 1];	
	
	if (ch1 != ch2)
	{
		/* 新建节点*/
		q = HeapPos;
		HeapPos++;
		SortHeap[q].off = off;
		SortHeap[q].next = SortTable[ch1 * 256 + ch2];
		SortTable[ch1 * 256 + ch2] = q;
	}
	else
	{
		/* 对重复2字节串
		** 因为没有虚拟偏移也没有删除操作，只要比较第一个节点
		** 是否和 off 相连接即可*/
		q = SortTable[ch1 * 256 + ch2];
		if (q != 0 && off == SortHeap[q].off2 + 1)
		{		
			/* 节点合并*/
			SortHeap[q].off2 = off;
		}		
		else
		{
			/* 新建节点*/
			q = HeapPos;
			HeapPos++;
			SortHeap[q].off = off;
			SortHeap[q].off2 = off;
			SortHeap[q].next = SortTable[ch1 * 256 + ch2];
			SortTable[ch1 * 256 + ch2] = q;
		}
	}
}

/*****************************************
** 将窗口向右滑动n个字节*/
void CCompressLZ77::_ScrollWindow(int n)
{	
	int i;
	for (i = 0; i < n; i++)
	{		
		nWndSize++;		
		if (nWndSize > 1)			
			_InsertIndexItem(nWndSize - 2);
	}
}

/*********************************************************
** 得到已经匹配了2个字节的窗口位置offset
** 共能匹配多少个字节*/
int CCompressLZ77::_GetSameLen(BYTE* src, int srclen, int nSeekStart, int offset)
{
	int i = 2; /* 已经匹配了2个字节*/
	int maxsame = min(srclen - nSeekStart, nWndSize - offset);
	while (i < maxsame
			&& src[nSeekStart + i] == pWnd[offset + i])
		i++;
	/*_ASSERT(nSeekStart + i <= srclen && offset + i <= nWndSize);*/
	return i;
}

/*********************************************************
** 在滑动窗口中查找术语
** nSeekStart - 从何处开始匹配
** offset, len - 用于接收结果，表示在滑动窗口内的偏移和长度
** 返回值- 是否查到长度为2或2以上的匹配字节串*/
BOOL CCompressLZ77::_SeekPhase(BYTE* src, int srclen, int nSeekStart, int* offset, int* len)
{	
	int j, m, n;
	WORD p;
	if (nSeekStart < srclen - 1)
	{
		BYTE ch1, ch2;
		ch1 = src[nSeekStart]; ch2 = src[nSeekStart + 1];
		
		p = SortTable[ch1 * 256 + ch2];
		if (p != 0)
		{
			m = 2; n = SortHeap[p].off;
			while (p != 0)
			{
				j = _GetSameLen(src, srclen, 
					nSeekStart, SortHeap[p].off);
				if ( j > m )
				{ 
					m = j; 
					n = SortHeap[p].off;
				}			
				p = SortHeap[p].next;
			}	
			(*offset) = n; 
			(*len) = m;
			return TRUE;		
		}	
	}
	return FALSE;
}

/***************************************
** 输出压缩码
** code - 要输出的数
** bits - 要输出的位数(对isGamma=TRUE时无效)
** isGamma - 是否输出为γ编码*/
void CCompressLZ77::_OutCode(BYTE* dest, DWORD code, int bits, BOOL isGamma)
{	
	if ( isGamma )
	{
		BYTE* pb;
		DWORD out;
		/* 计算输出位数*/
		int GammaCode = (int)code - 1;
		int q = LowerLog2(GammaCode);
		if (q > 0)
		{
			out = 0xffffffff;
			pb = (BYTE*)&out;
			/* 输出q个1*/
			CopyBits(dest + CurByte, CurBit, 
				pb, 0, q);
			MovePos(&CurByte, &CurBit, q);
		}
		/* 输出一个0*/
		out = 0;
		pb = (BYTE*)&out;		
		CopyBits(dest + CurByte, CurBit, pb + 3, 7, 1);
		MovePos(&CurByte, &CurBit, 1);
		if (q > 0)
		{
			/* 输出余数, q位*/
			int sh = 1;
			sh <<= q;
			out = GammaCode - sh;
			pb = (BYTE*)&out;
			InvertDWord(&out);
			CopyBits(dest + CurByte, CurBit, 
				pb + (32 - q) / 8, (32 - q) % 8, q);
			MovePos(&CurByte, &CurBit, q);
		}
	}
	else 
	{
		DWORD dw = (DWORD)code;
		BYTE* pb = (BYTE*)&dw;
		InvertDWord(&dw);
		CopyBits(dest + CurByte, CurBit, 
				pb + (32 - bits) / 8, (32 - bits) % 8, bits);
		MovePos(&CurByte, &CurBit, bits);
	}
}

/*******************************************
** 压缩一段字节流
** src - 源数据区
** srclen - 源数据区字节长度
** dest - 压缩数据区，调用前分配srclen+5字节内存
** 返回值 > 0 压缩数据长度
** 返回值 = 0 数据无法压缩
** 返回值 < 0 压缩中异常错误*/
int CCompressLZ77::Compress(BYTE* src, int srclen, BYTE* dest)
{
	int i;
	int off, len,destlen ;
	CurByte = 0; CurBit = 0;	

	if (srclen > 65536) 
		return -1;

	pWnd = src;
	_InitSortTable();
	for (i = 0; i < srclen; i++)
	{		
		if (CurByte >= srclen)
			return 0;
		if (_SeekPhase(src, srclen, i, &off, &len))
		{			
			/* 输出匹配术语 flag(1bit) + len(γ编码) + offset(最大16bit)*/
			_OutCode(dest, 1, 1, FALSE);
			_OutCode(dest, len, 0, TRUE);

			/* 在窗口不满64k大小时，不需要16位存储偏移*/
			_OutCode(dest, off, UpperLog2(nWndSize), FALSE);
						
			_ScrollWindow(len);
			i += len - 1;
		}
		else
		{
			/* 输出单个非匹配字符 0(1bit) + char(8bit)*/
			_OutCode(dest, 0, 1, FALSE);
			_OutCode(dest, (DWORD)(src[i]), 8, FALSE);
			_ScrollWindow(1);
		}
	}
	destlen = CurByte + ((CurBit) ? 1 : 0);
	if (destlen >= srclen)
		return 0;
	return destlen;
}


/*******************************************
** 解压缩一段字节流
** src - 接收原始数据的内存区
** srclen - 源数据区字节长度
** dest - 压缩数据区
** 返回值 - 成功与否*/
BOOL CCompressLZ77::Decompress(BYTE* src, int srclen, BYTE* dest)
{
	int i,j;
	int len, off,bits;
	DWORD dw = 0;
	BYTE* pb;
	BYTE b;
	int q ;


	CurByte = 0; CurBit = 0;
	pWnd = src;		/* 初始化窗口*/
	nWndSize = 0;
	

	if (srclen > 65536) 
		return FALSE;
	
	for (i = 0; i < srclen; i++)
	{		
		b = GetBit(dest[CurByte], CurBit);
		MovePos(&CurByte, &CurBit, 1);
		if (b == 0) /* 单个字符*/
		{
			CopyBits(src + i, 0, dest + CurByte, CurBit, 8);
			MovePos(&CurByte, &CurBit, 8);
			nWndSize++;
		}
		else		/* 窗口内的术语*/
		{
			q = -1;
			while (b != 0)
			{
				q++;
				b = GetBit(dest[CurByte], CurBit);
				MovePos(&CurByte, &CurBit, 1);				
			}
			dw = 0;
			if (q > 0)
			{				
				pb = (BYTE*)&dw;
				CopyBits(pb + (32 - q) / 8, (32 - q) % 8, dest + CurByte, CurBit, q);
				MovePos(&CurByte, &CurBit, q);
				InvertDWord(&dw);
				len = 1;
				len <<= q;
				len += dw;
				len += 1;
			}
			else
				len = 2;

			/* 在窗口不满64k大小时，不需要16位存储偏移*/
			dw = 0;
			pb = (BYTE*)&dw;
			bits = UpperLog2(nWndSize);
			CopyBits(pb + (32 - bits) / 8, (32 - bits) % 8, dest + CurByte, CurBit, bits);
			MovePos(&CurByte, &CurBit, bits);
			InvertDWord(&dw);
			off = (int)dw;
			/* 输出术语*/
			for (j = 0; j < len; j++)
			{
				/*_ASSERT(i + j <  srclen);
				_ASSERT(off + j <  _MAX_WINDOW_SIZE);*/
				
				src[i + j] = pWnd[off + j];
			}
			nWndSize += len;
			i += len - 1;
		}
		/* 滑动窗口*/
		if (nWndSize > _MAX_WINDOW_SIZE)
		{
			pWnd += nWndSize - _MAX_WINDOW_SIZE;
			nWndSize = _MAX_WINDOW_SIZE;			
		}
	}

	return TRUE;
}

int CCompressLZ77::CompressFile(char *inputfile, char *outputfile )
{	
	BYTE soubuf[65536];
	BYTE destbuf[65536 + 16];

	FILE* in;
	FILE* out;
	errno_t err;
	long soulen,destlen;
	int last,act;
	WORD flag1, flag2,netflag;

	printf("compress file1\n");
	in = fopen(inputfile,"rb");
	if (in == NULL)
	{
		puts("Can't open source file");
		return -1;
	}

	printf("compress file2\n");
	out = fopen(outputfile, "wb");
	if (out == NULL)
	{
		puts("Can't open dest file");
		fclose(in);
		return -1;
	}

	fseek(in, 0, SEEK_END);
	soulen = ftell(in);
	fseek(in, 0, SEEK_SET);

	last = soulen;
	while ( last > 0 )
	{
		act = min(65536, last);
		fread(soubuf, act, 1, in);
		last -= act;
		if (act == 65536)			/* out 65536 bytes*/
			flag1 = 0;		
		else					/* out last blocks*/
			flag1 = act;
		netflag = htons(flag1);
		fwrite(&netflag, sizeof(WORD), 1, out);

		destlen = Compress((BYTE*)soubuf, act, (BYTE*)destbuf);
		if (destlen == 0)		/* can't compress the block*/
		{
			flag2 = flag1;
			netflag = htons(flag2);
			fwrite(&netflag, sizeof(WORD), 1, out);
			fwrite(soubuf, act, 1, out);
		}
		else
		{
			flag2 = (WORD)destlen;
			netflag = htons(flag2);
			fwrite(&netflag, sizeof(WORD), 1, out);
			fwrite(destbuf, destlen, 1, out);
		}
	}
	fclose(in);
	fclose(out);
	return 0;
}

int CCompressLZ77::DecompressFile(char *inputfile, char *outputfile )
{	
	BYTE soubuf[65536];
	BYTE destbuf[65536 + 16];

	FILE* in;
	FILE* out;
	errno_t err;
	long soulen;
	int last,act;
	WORD flag1, flag2;

	in = fopen(inputfile , "rb");

	if (NULL == in)
	{
		puts("Can't open source file");
		return -1;
	}

	out = fopen(outputfile,"wb");
	if (NULL == out)
	{
		puts("Can't open dest file");
		fclose(in);
		return -1;
	}

	fseek(in, 0, SEEK_END);
	soulen = ftell(in);
	fseek(in, 0, SEEK_SET);

	last = soulen;
	while (last > 0)
	{
		fread(&flag1, sizeof(WORD), 1, in);
		flag1 = ntohs(flag1);
		fread(&flag2, sizeof(WORD), 1, in);
		flag2 = ntohs(flag2);
		last -= 2 * sizeof(WORD);
		if (flag1 == 0)
			act = 65536;
		else
			act = flag1;
		last-= flag2 ? (flag2) : act;

		if (flag2 == flag1)
		{
			fread(soubuf, act, 1, in);				
		}
		else
		{
			fread(destbuf, flag2, 1, in);
			if (!Decompress((BYTE*)soubuf, act, (BYTE*)destbuf))
			{
				puts("Decompress error");
				fclose(in);
				fclose(out);
				return -1;
			}
		}
		fwrite((BYTE*)soubuf, act, 1, out);				
	}
	fclose(in);
	fclose(out);
	return 0;
}

int CCompressLZ77::CompressBuf(unsigned char *input, int inputlen,unsigned char *output )
{	
	BYTE soubuf[65536];
	BYTE destbuf[65536 + 16];

	long soulen,destlen,outlen=0;
	int last,act;
	WORD flag1, flag2,netflag;
	
	soulen = inputlen;

	last = soulen;
	while ( last > 0 )
	{
		act = min(65536, last);
		/*fread(soubuf, act, 1, in);*/
		memcpy(soubuf,input+soulen-last,act);
		last -= act;
		if(act == 65536)			/* out 65536 bytes*/
			flag1 = 0;
		else					/* out last blocks*/
			flag1 = act;
		netflag = htons(flag1);
		/*fwrite(&netflag, sizeof(WORD), 1, out);*/
		memcpy(output+outlen,&netflag, sizeof(WORD));
		outlen+= sizeof(WORD);
		destlen = Compress((BYTE*)soubuf, act, (BYTE*)destbuf);
		if (destlen == 0)		/* can't compress the block*/
		{
			flag2 = flag1;
			netflag = htons(flag2);
			/*fwrite(&netflag, sizeof(WORD), 1, out);*/
			memcpy(output+outlen,&netflag,sizeof(WORD));
			outlen+= sizeof(WORD);
			/*fwrite(soubuf, act, 1, out);*/
			memcpy(output+outlen,soubuf,act);
			outlen+=act;
		}
		else
		{
			flag2 = (WORD)destlen;
			netflag = htons(flag2);
			/*fwrite(&netflag, sizeof(WORD), 1, out);*/
			memcpy(output+outlen,&netflag, sizeof(WORD));
			outlen+= sizeof(WORD);
			/*fwrite(destbuf, destlen, 1, out);*/
			memcpy(output+outlen,destbuf,destlen);
			outlen+=destlen;
		}
	}
	return outlen;
}

int CCompressLZ77::DecompressBuf(unsigned char *input, int inputlen,unsigned char *output )
{
	BYTE soubuf[65536];
	BYTE destbuf[65536 + 16];

	long soulen,outlen=0;
	int last,act;
	WORD flag1, flag2;

	soulen = inputlen;
	last = soulen;
	while (last > 0)
	{
		/*fread(&flag1, sizeof(WORD), 1, in);*/
		memcpy(&flag1,input+soulen-last,sizeof(WORD));
		last -= sizeof(WORD);
		flag1 = ntohs(flag1);
		
		/*fread(&flag2, sizeof(WORD), 1, in);*/
		memcpy(&flag2,input+soulen-last,sizeof(WORD));
		last -= sizeof(WORD);
		flag2 = ntohs(flag2);
		
		if (flag1 == 0)
			act = 65536;
		else
			act = flag1;
		
		if (flag2 == flag1)
		{
			/*fread(soubuf, act, 1, in);*/
			memcpy(soubuf,input+soulen-last,act);
		}
		else
		{
			/*fread(destbuf, flag2, 1, in);*/
			memcpy(destbuf,input+soulen-last,flag2);
			if (!Decompress((BYTE*)soubuf, act, (BYTE*)destbuf))
			{
				puts("Decompress error");
				return -1;
			}
		}
		last-= flag2 ? (flag2) : act;
		/*fwrite((BYTE*)soubuf, act, 1, out);*/
		memcpy(output+outlen,soubuf,act);
		outlen +=act;
	}
	return outlen;
}

int CCompressLZ77::GetDecompressBufLen(unsigned char *input, int inputlen)
{
	long soulen,outlen=0;
	int last,act;
	WORD flag1, flag2;

	soulen = inputlen;
	last = soulen;
	while (last > 0)
	{
		memcpy(&flag1,input+soulen-last,sizeof(WORD));
		last -= sizeof(WORD);
		flag1 = ntohs(flag1);

		memcpy(&flag2,input+soulen-last,sizeof(WORD));
		last -= sizeof(WORD);
		flag2 = ntohs(flag2);
		if (flag1 == 0)
			act = 65536;
		else
			act = flag1;
		last-= flag2 ? (flag2) : act;
		outlen +=act;
	}
	return outlen;
}

/*
int  main(int argc, char* argv[])
{	
	BYTE soubuf[65536];
	BYTE destbuf[65536 + 16];

	FILE* in;
	FILE* out;
	long soulen,destlen;
	int last,act;
	WORD flag1, flag2,netflag;

	CCompressLZ77 cc;

	if (argc != 4)
	{
		puts("Usage: ");
		printf("    Compress : %s c sourcefile destfile\n", argv[0]);
		printf("  Decompress : %s d sourcefile destfile\n", argv[0]);
		return 0;
	}

	in = fopen(argv[2], "rb");
	if (in == NULL)
	{
		puts("Can't open source file");
		return 0;
	}
	out = fopen(argv[3], "wb");
	if (out == NULL)
	{
		puts("Can't open dest file");
		fclose(in);
		return 0;
	}
	fseek(in, 0, SEEK_END);
	soulen = ftell(in);
	fseek(in, 0, SEEK_SET);

	if (argv[1][0] == 'c')
	{
		last = soulen;
		while ( last > 0 )
		{
			act = min(65536, last);
			fread(soubuf, act, 1, in);
			last -= act;
			if (act == 65536)		
				flag1 = 0;		
			else					
				flag1 = act;
			netflag = htons(flag1);
			fwrite(&netflag, sizeof(WORD), 1, out);

			destlen = cc.Compress((BYTE*)soubuf, act, (BYTE*)destbuf);
			if (destlen == 0)	
			{
				flag2 = flag1;
				netflag = htons(flag2);
				fwrite(&netflag, sizeof(WORD), 1, out);
				fwrite(soubuf, act, 1, out);
			}
			else
			{
				flag2 = (WORD)destlen;
				netflag = htons(flag2);
				fwrite(&netflag, sizeof(WORD), 1, out);
				fwrite(destbuf, destlen, 1, out);
			}
		}
	}
	else if (argv[1][0] == 'd') 
	{
		int last = soulen, act;
		while (last > 0)
		{
			fread(&flag1, sizeof(WORD), 1, in);
			flag1 = ntohs(flag1);
			fread(&flag2, sizeof(WORD), 1, in);
			flag2 = ntohs(flag2);
			last -= 2 * sizeof(WORD);
			if (flag1 == 0)
				act = 65536;
			else
				act = flag1;
			last-= flag2 ? (flag2) : act;

			if (flag2 == flag1)
			{
				fread(soubuf, act, 1, in);				
			}
			else
			{
				fread(destbuf, flag2, 1, in);
				if (!cc.Decompress((BYTE*)soubuf, act, (BYTE*)destbuf))
				{
					puts("Decompress error");
					fclose(in);
					fclose(out);
					return 0;
				}
			}
			fwrite((BYTE*)soubuf, act, 1, out);				
		}
	}
	else
	{
		puts("Usage: ");
		printf("    Compress : %s c sourcefile destfile\n", argv[0]);
		printf("  Decompress : %s d sourcefile destfile\n", argv[0]);		
	}

	fclose(in);
	fclose(out);

	return 0;
}
*/

int Compressfilefunc(char *inputfile, char *outputfile)
{
	CCompressLZ77 lz;
	return lz.CompressFile(inputfile, outputfile);
}

int Decompressfilefunc(char *inputfile, char *outputfile)
{
	CCompressLZ77 lz;
	return lz.DecompressFile(inputfile, outputfile);
}

int Testcode()
{
	printf("hello cgodll  test!\n");
	return 1;
}
