
// FeeSettlement.h : PROJECT_NAME Ӧ�ó������ͷ�ļ�
//

#pragma once

#ifndef __AFXWIN_H__
	#error "�ڰ������ļ�֮ǰ������stdafx.h�������� PCH �ļ�"
#endif

#include "resource.h"		// ������


// CFeeSettlementApp: 
// �йش����ʵ�֣������ FeeSettlement.cpp
//

class CFeeSettlementApp : public CWinApp
{
public:
	CFeeSettlementApp();

// ��д
public:
	virtual BOOL InitInstance();

// ʵ��

	DECLARE_MESSAGE_MAP()
};

extern CFeeSettlementApp theApp;