package storage

import (
	"CenterSettlement-go/database"
	"CenterSettlement-go/types"
	log "github.com/sirupsen/logrus"
)

//根据停车场编号 查询 停车场名称
func GetTingcc(tingccbh string) string {
	database.DBInit()
	xorm := database.XormClient
	//停车场信息
	tingcc := new(types.BTccTingcc)
	is, err := xorm.Where("F_VC_TINGCCBH=?", tingccbh).Get(tingcc)
	if is == false {
		log.Println("没有该停车场")
	}
	if err != nil {
		log.Fatal("查询停车场名称 error ")
	}
	if is == true {
		return tingcc.F_VC_MINGC
	}
	return ""
}
