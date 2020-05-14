package storage

import (
	"testing"
)

//测试查询本省的结算数据
func TestQueryJiessj(t *testing.T) {
	QueryJiessj()
}

//测试查询其他地区的结算数据
func TestQueryQiTaJiessj(t *testing.T) {
	QueryQiTaJiessj()
}
