package storage

import (
	//"CenterSettlement-go/sysinit"
	//"CenterSettlement-go/types"
	"github.com/go-xorm/xorm"
	"log"
)

func TransactionCommit(session *xorm.Session) {
	defer session.Close()
	// add Commit() after all actions
	err := session.Commit()
	if err != nil {
		return
	}
}

func TransactionBegin(engine *xorm.Engine) *xorm.Session {
	session := engine.NewSession()

	// add Begin() before any action
	serr := session.Begin()
	if serr != nil {
		log.Fatalln("Transaction session.Begin error ")
	}
	return session

	//数据库表操作
	//user1 := types.BJsJiessj{}
	//_, err = session.Insert(&user1)
	//if err != nil {
	//	session.Rollback()
	//	return
	//}
	//user2 := types.BJsJiessj{}
	//_, err = session.Where("id = ?", 2).Update(&user2)
	//if err != nil {
	//	session.Rollback()
	//	return
	//}
	//
	//_, err = session.Exec("delete from userinfo where username = ?", user2.FNbDabzt)
	//if err != nil {
	//	session.Rollback()
	//	return
	//}
}
