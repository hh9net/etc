package storage

import (
	"CenterSettlement-go/database"
	"CenterSettlement-go/types"
	log "github.com/sirupsen/logrus"
)

//根据停车场编号 查询 停车场名称
func GetTingcc(tingccbh string) string {
	xorm := database.XormClient
	//停车场信息
	tingcc := &types.BTccTingcc{FVcTingccbh: tingccbh}
	is, err := xorm.Get(tingcc)
	if err != nil {
		log.Println("查询停车场名称 error ", err)
		return ""
	}
	if is == false {
		log.Println("没有该停车场", err)
		return ""
	}

	if is == true {
		log.Println("查询停车场名称  ", tingcc.FVcMingc)
		return tingcc.FVcMingc
	}
	return tingcc.FVcMingc
}
