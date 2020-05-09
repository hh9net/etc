
// FeeSettlementDlg.h : ͷ�ļ�
//

#pragma once
#include "afxwin.h"


// CFeeSettlementDlg �Ի���
class CFeeSettlementDlg : public CDialogEx
{
// ����
public:
	CFeeSettlementDlg(CWnd* pParent = NULL);	// ��׼���캯��
	CString GetNow();		
	BOOL GetConfig();
	void ToTray();
	void DeleteTray();

// �Ի�������
#ifdef AFX_DESIGN_TIME
	enum { IDD = IDD_FEESETTLEMENT_DIALOG };
#endif

	protected:
	virtual void DoDataExchange(CDataExchange* pDX);	// DDX/DDV ֧��


// ʵ��
protected:
	HICON m_hIcon;

	// ���ɵ���Ϣӳ�亯��
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
