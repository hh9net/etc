/*****************************
** LZ77.CPP
******************************/

#include "StdAfx.h"
#include "lz77.h"

/*******************************************************
** ȡlog2(n)��upper_bound*/
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
** ȡlog2(n)��lower_bound*/
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
** ��λָ��*piByte(�ֽ�ƫ��), *piBit(�ֽ���λƫ��)����numλ*/
void CCompressLZ77::MovePos(int* piByte, int* piBit, int num)
{
	num += (*piBit);
	(*piByte) += num / 8;
	(*piBit) = num % 8;
}
/* MovePos
************************************************************/

/***********************************************************
** �õ��ֽ�byte��posλ��ֵ
**		pos˳��Ϊ��λ���0����������*/
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
** ����byte�ĵ�iBitλΪaBit
**		iBit˳��Ϊ��λ���0����������*/
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
** ��DWORDֵ�Ӹ�λ�ֽڵ���λ�ֽ�����*/
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
** CopyBits : �����ڴ��е�λ��
**		memDest - Ŀ��������
**		nDestPos - Ŀ����������һ���ֽ��е���ʼλ
**		memSrc - Դ������
**		nSrcPos - Դ��������һ���ֽڵ�����ʼλ
**		nBits - Ҫ���Ƶ�λ��
**	˵����
**		��ʼλ�ı�ʾԼ��Ϊ���ֽڵĸ�λ����λ���������ң�
**		����Ϊ 0��1��... , 7
**		Ҫ���Ƶ������������������غ�*/
void CCompressLZ77::CopyBits(BYTE* memDest, int nDestPos, BYTE* memSrc, int nSrcPos, int nBits)
{
	int iByteDest = 0, iBitDest;
	int iByteSrc = 0, iBitSrc = nSrcPos;

	int nBitsToFill, nBitsCanFill;

	while (nBits > 0)
	{
		/* ����Ҫ��Ŀ������ǰ�ֽ�����λ��*/
		nBitsToFill = min(nBits, iByteDest ? 8 : 8 - nDestPos);
		/* Ŀ������ǰ�ֽ�Ҫ������ʼλ*/
		iBitDest = iByteDest ? 0 : nDestPos;
		/* �������һ�δ�Դ�������и��Ƶ�λ��*/
		nBitsCanFill = min(nBitsToFill, 8 - iBitSrc);
		/* �ֽ��ڸ���*/
		CopyBitsInAByte(memDest + iByteDest, iBitDest, 
			memSrc + iByteSrc, iBitSrc, nBitsCanFill);		
		/* �����û�и����� nBitsToFill ��*/
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

		nBits -= nBitsToFill;	/* �Ѿ������nBitsToFillλ*/
		iByteDest++;
	}	
}
/* CopyBits
********************************************************/

/*******************************************************
** CopyBitsInAByte : ��һ���ֽڷ�Χ�ڸ���λ��
** ��������ͬ CopyBits �Ĳ���
** ˵����
**		�˺����� CopyBits ���ã����������飬��
**		�ٶ�Ҫ���Ƶ�λ����һ���ֽڷ�Χ��*/
void CCompressLZ77::CopyBitsInAByte(BYTE* memDest, int nDestPos, BYTE* memSrc, int nSrcPos, int nBits)
{
	BYTE b1, b2;
	b1 = *memSrc;
	b1 <<= nSrcPos; b1 >>= 8 - nBits;	/* �����ø��Ƶ�λ��0*/
	b1 <<= 8 - nBits - nDestPos;		/* ��Դ��Ŀ���ֽڶ���*/
	*memDest |= b1;		/* ����ֵΪ1��λ*/
	b2 = 0xff; b2 <<= 8 - nDestPos;		/* �����ø��Ƶ�λ��1*/
	b1 |= b2;
	b2 = 0xff; b2 >>= nDestPos + nBits;
	b1 |= b2;
	*memDest &= b1;		/* ����ֵΪ0��λ*/
}
/* CopyBitsInAByte
********************************************************/


/*------------------------------------------------------------------*/


/* ��ʼ���������ͷ��ϴ�ѹ���õĿռ�*/
void CCompressLZ77::_InitSortTable()
{
	memset(SortTable, 0, sizeof(WORD) * 65536);
	nWndSize = 0;
	HeapPos = 1;
}

/* �����������һ��2�ֽڴ�*/
void CCompressLZ77::_InsertIndexItem(int off)
{
	WORD q;
	BYTE ch1, ch2;
	ch1 = pWnd[off]; ch2 = pWnd[off + 1];	
	
	if (ch1 != ch2)
	{
		/* �½��ڵ�*/
		q = HeapPos;
		HeapPos++;
		SortHeap[q].off = off;
		SortHeap[q].next = SortTable[ch1 * 256 + ch2];
		SortTable[ch1 * 256 + ch2] = q;
	}
	else
	{
		/* ���ظ�2�ֽڴ�
		** ��Ϊû������ƫ��Ҳû��ɾ��������ֻҪ�Ƚϵ�һ���ڵ�
		** �Ƿ�� off �����Ӽ���*/
		q = SortTable[ch1 * 256 + ch2];
		if (q != 0 && off == SortHeap[q].off2 + 1)
		{		
			/* �ڵ�ϲ�*/
			SortHeap[q].off2 = off;
		}		
		else
		{
			/* �½��ڵ�*/
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
** ���������һ���n���ֽ�*/
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
** �õ��Ѿ�ƥ����2���ֽڵĴ���λ��offset
** ����ƥ����ٸ��ֽ�*/
int CCompressLZ77::_GetSameLen(BYTE* src, int srclen, int nSeekStart, int offset)
{
	int i = 2; /* �Ѿ�ƥ����2���ֽ�*/
	int maxsame = min(srclen - nSeekStart, nWndSize - offset);
	while (i < maxsame
			&& src[nSeekStart + i] == pWnd[offset + i])
		i++;
	/*_ASSERT(nSeekStart + i <= srclen && offset + i <= nWndSize);*/
	return i;
}

/*********************************************************
** �ڻ��������в�������
** nSeekStart - �Ӻδ���ʼƥ��
** offset, len - ���ڽ��ս������ʾ�ڻ��������ڵ�ƫ�ƺͳ���
** ����ֵ- �Ƿ�鵽����Ϊ2��2���ϵ�ƥ���ֽڴ�*/
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
** ���ѹ����
** code - Ҫ�������
** bits - Ҫ�����λ��(��isGamma=TRUEʱ��Ч)
** isGamma - �Ƿ����Ϊ�ñ���*/
void CCompressLZ77::_OutCode(BYTE* dest, DWORD code, int bits, BOOL isGamma)
{	
	if ( isGamma )
	{
		BYTE* pb;
		DWORD out;
		/* �������λ��*/
		int GammaCode = (int)code - 1;
		int q = LowerLog2(GammaCode);
		if (q > 0)
		{
			out = 0xffffffff;
			pb = (BYTE*)&out;
			/* ���q��1*/
			CopyBits(dest + CurByte, CurBit, 
				pb, 0, q);
			MovePos(&CurByte, &CurBit, q);
		}
		/* ���һ��0*/
		out = 0;
		pb = (BYTE*)&out;		
		CopyBits(dest + CurByte, CurBit, pb + 3, 7, 1);
		MovePos(&CurByte, &CurBit, 1);
		if (q > 0)
		{
			/* �������, qλ*/
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
** ѹ��һ���ֽ���
** src - Դ������
** srclen - Դ�������ֽڳ���
** dest - ѹ��������������ǰ����srclen+5�ֽ��ڴ�
** ����ֵ > 0 ѹ�����ݳ���
** ����ֵ = 0 �����޷�ѹ��
** ����ֵ < 0 ѹ�����쳣����*/
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
			/* ���ƥ������ flag(1bit) + len(�ñ���) + offset(���16bit)*/
			_OutCode(dest, 1, 1, FALSE);
			_OutCode(dest, len, 0, TRUE);

			/* �ڴ��ڲ���64k��Сʱ������Ҫ16λ�洢ƫ��*/
			_OutCode(dest, off, UpperLog2(nWndSize), FALSE);
						
			_ScrollWindow(len);
			i += len - 1;
		}
		else
		{
			/* ���������ƥ���ַ� 0(1bit) + char(8bit)*/
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
** ��ѹ��һ���ֽ���
** src - ����ԭʼ���ݵ��ڴ���
** srclen - Դ�������ֽڳ���
** dest - ѹ��������
** ����ֵ - �ɹ����*/
BOOL CCompressLZ77::Decompress(BYTE* src, int srclen, BYTE* dest)
{
	int i,j;
	int len, off,bits;
	DWORD dw = 0;
	BYTE* pb;
	BYTE b;
	int q ;


	CurByte = 0; CurBit = 0;
	pWnd = src;		/* ��ʼ������*/
	nWndSize = 0;
	

	if (srclen > 65536) 
		return FALSE;
	
	for (i = 0; i < srclen; i++)
	{		
		b = GetBit(dest[CurByte], CurBit);
		MovePos(&CurByte, &CurBit, 1);
		if (b == 0) /* �����ַ�*/
		{
			CopyBits(src + i, 0, dest + CurByte, CurBit, 8);
			MovePos(&CurByte, &CurBit, 8);
			nWndSize++;
		}
		else		/* �����ڵ�����*/
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

			/* �ڴ��ڲ���64k��Сʱ������Ҫ16λ�洢ƫ��*/
			dw = 0;
			pb = (BYTE*)&dw;
			bits = UpperLog2(nWndSize);
			CopyBits(pb + (32 - bits) / 8, (32 - bits) % 8, dest + CurByte, CurBit, bits);
			MovePos(&CurByte, &CurBit, bits);
			InvertDWord(&dw);
			off = (int)dw;
			/* �������*/
			for (j = 0; j < len; j++)
			{
				/*_ASSERT(i + j <  srclen);
				_ASSERT(off + j <  _MAX_WINDOW_SIZE);*/
				
				src[i + j] = pWnd[off + j];
			}
			nWndSize += len;
			i += len - 1;
		}
		/* ��������*/
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

	if ((err = fopen_s(&in, inputfile, "rb")) != 0)
	{
		puts("Can't open source file");
		return -1;
	}

	if ((err = fopen_s(&out, outputfile, "wb")) != 0)
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

	if ((err = fopen_s(&in, inputfile, "rb")) != 0)
	{
		puts("Can't open source file");
		return -1;
	}

	if ((err = fopen_s(&out, outputfile, "wb")) != 0)
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
void main(int argc, char* argv[])
{	
	BYTE soubuf[65536];
	BYTE destbuf[65536 + 16];

	FILE* in;
	FILE* out;
	long soulen,destlen;
	int last,act;
	WORD flag1, flag2,netflag;

	if (argc != 4)
	{
		puts("Usage: ");
		printf("    Compress : %s c sourcefile destfile\n", argv[0]);
		printf("  Decompress : %s d sourcefile destfile\n", argv[0]);
		return;
	}

	in = fopen(argv[2], "rb");
	if (in == NULL)
	{
		puts("Can't open source file");
		return;
	}
	out = fopen(argv[3], "wb");
	if (out == NULL)
	{
		puts("Can't open dest file");
		fclose(in);
		return;
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

			destlen = Compress((BYTE*)soubuf, act, (BYTE*)destbuf);
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
				if (!Decompress((BYTE*)soubuf, act, (BYTE*)destbuf))
				{
					puts("Decompress error");
					fclose(in);
					fclose(out);
					return;
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
}
*/