package types

//
type BTccTingcc struct {
	F_VC_TINGCCBH      string //`F_VC_TINGCCBH` varchar(32) NOT NULL COMMENT '停车场编号',
	F_VC_GONGSBH       string //`F_VC_GONGSBH` varchar(32) DEFAULT NULL COMMENT '公司/集团编号 idx-',
	F_VC_QUDBH         string //`F_VC_QUDBH` varchar(32) DEFAULT NULL COMMENT '渠道编号',
	F_VC_TINGCCWLBH    string //`F_VC_TINGCCWLBH` varchar(32) DEFAULT NULL COMMENT '停车场网络编号 由于前期要与旧平台同步，改字段请用数字表示',
	F_NB_TINGCCLX      int64  //`F_NB_TINGCCLX` int(11) NOT NULL DEFAULT '1' COMMENT '停车场类型 1：单点；2：总对总；',
	F_VC_MINGC         string //`F_VC_MINGC` varchar(32) DEFAULT NULL COMMENT '名称-NEW',
	F_VC_DIZ           string //`F_VC_DIZ` varchar(512) DEFAULT NULL COMMENT '地址',
	F_VC_JINGD         string //`F_VC_JINGD` decimal(32,10) DEFAULT NULL COMMENT '经度',
	F_VC_WEID          string //`F_VC_WEID` decimal(32,10) DEFAULT NULL COMMENT '维度',
	F_VC_GUANLYID      string //`F_VC_GUANLYID` varchar(32) NOT NULL COMMENT '管理员ID-NEW',
	F_DT_CHUANGJSJ     string //`F_DT_CHUANGJSJ` datetime DEFAULT NULL COMMENT '创建时间',
	F_VC_CHUANGJR      string //`F_VC_CHUANGJR` varchar(32) DEFAULT NULL COMMENT '创建人',
	F_NB_ZHUANGT       int    //`F_NB_ZHUANGT` int(11) DEFAULT '1' COMMENT '状态-U 1：正常；2：待审核；3：停用；',
	F_VC_VERIFY_STATUS int    //`F_VC_VERIFY_STATUS` int(11) DEFAULT NULL COMMENT '审核结果-NEW 1：审核通过；2：待审核；3：审核驳回，需修改信息；4：审核拒绝；',
	F_VC_FUZRDH        string //`F_VC_FUZRDH` varchar(32) DEFAULT NULL COMMENT '负责人电话-D',
	F_VC_FUZRXM        string //`F_VC_FUZRXM` varchar(32) DEFAULT NULL COMMENT '负责人姓名-D',
	F_VC_SHENGDM       string //`F_VC_SHENGDM` varchar(32) DEFAULT NULL COMMENT '省代码',
	F_VC_SHENGMC       string //`F_VC_SHENGMC` varchar(32) DEFAULT NULL COMMENT '省名称',
	F_VC_SHIDM         string //`F_VC_SHIDM` varchar(32) DEFAULT NULL COMMENT '市代码',
	F_VC_SHIMC         string //`F_VC_SHIMC` varchar(32) DEFAULT NULL COMMENT '市名称',
	F_VC_QUDM          string //`F_VC_QUDM` varchar(32) DEFAULT NULL COMMENT '区代码',
	F_VC_QUMC          string //`F_VC_QUMC` varchar(32) DEFAULT NULL COMMENT '区名称',
	F_NB_FEIL          int    //`F_NB_FEIL` int(11) DEFAULT NULL COMMENT '费率 万分比',
}
