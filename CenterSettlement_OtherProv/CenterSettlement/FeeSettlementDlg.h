
// FeeSettlementDlg.h : 头文件
//

#pragma once
#include "afxwin.h"


// CFeeSettlementDlg 对话框
class CFeeSettlementDlg : public CDialogEx
{
// 构造
public:
	CFeeSettlementDlg(CWnd* pParent = NULL);	// 标准构造函数
	CString GetNow();		
	BOOL GetConfig();
	void ToTray();
	void DeleteTray();

// 对话框数据
#ifdef AFX_DESIGN_TIME
	enum { IDD = IDD_FEESETTLEMENT_DIALOG };
#endif

	protected:
	virtual void DoDataExchange(CDataExchange* pDX);	// DDX/DDV 支持


// 实现
protected:
	HICON m_hIcon;

	// 生成的消息映射函数
	virtual BOOL OnInitDialog();
	afx_msg void OnPaint();
	afx_msg HCURSOR OnQueryDragIcon();
	afx_msg LRESULT OnMsg1Back(WPARAM wParam, LPARAM lParam);
	afx_msg LRESULT OnMsg2Back(WPARAM wParam, LPARAM lParam);
	afx_msg LRESULT OnShowTask(WPARAM wParam, LPARAM lParam);
	DECLARE_MESSAGE_MAP()
public:
	CListBox m_lstMsg1;
	CListBox m_lstMsg2;
	afx_msg void OnBnClickedBtnMin();
	afx_msg void OnBnClickedCancel();
};
