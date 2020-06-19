package centerServer

import (
	"CenterSettlement-go/client"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	//处理文件
	Handle()

	//模拟联网中心，处理结算数据
	//Server()
	//处理文件的解析
	//HandleFile()
	//处理文件的发送
	//CenterClient()

}

func TestHandleFile(t *testing.T) {
	//模拟联网中心，处理结算数据
	//Server()
	//处理文件的解析
	HandleFile()
	//处理文件的发送
	//CenterClient()
}

func TestCenterClient(t *testing.T) {
	//模拟联网中心，发送数据
	CenterClient()
}

func TestNewDatabase(t *testing.T) {
	for {
		//扫描原始数据包
		pwd := "../compressedXml/"
		fileList, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("该centerYuanshi文件夹下有xml文件的数量 ：", len(fileList))
		if len(fileList) == 0 {
			log.Println("该centerYuanshi 文件夹下没有需要解析的文件")
			return
		}

		for i := range fileList {
			log.Println("该centerYuanshi文件夹需要数据入库的 xml 名字为:", fileList[i].Name()) //打印当前文件或目录下的文件或目录名
			//判断文件的结尾名
			if strings.HasSuffix(fileList[i].Name(), ".xml") {

				s1 := strings.Split(fileList[i].Name(), "_")

				s2 := strings.Split(s1[2], ".")
				s, _ := strconv.Atoi(s2[0])

				log.Println(fileList[i].Name(), s1, s2, s)
				for i1 := 100300; i1 < 100347; i1++ {
					if i1 == s {
						log.Println(i, s)
						//删除文件
						f1 := "../compressedXml/" + fileList[i].Name()
						log.Println(f1)
						client.DelFile(f1)
						log.Println("删除文件成功")
						break
					}
				}

			}
		}
	}
}

func TestDB_NewTable(t *testing.T) {
	db := new(DB)

	db.NewTable()
}

func TestGetFileMd5(t *testing.T) {
	f := "00000000000000010612.xml"
	s := GetFileMd5(f)
	log.Println(s)
}

func TestUnzip(t *testing.T) {
	f := "00000000000000010612.xml"
	s := UnZipLz77(f)
	log.Println(s)
}

//bytes.Equa
func TestToABCD(t *testing.T) {
	s := "12asd456fdghfsh"
	fmt.Println("s =", s)
	str := strings.ToUpper(s)
	fmt.Println("str =", str)
}

func TestGenerateKeepAccountXml(t *testing.T) {
	GenerateKeepAccountXml()

}

func TestQueryKeepAccountdata(t *testing.T) {
	QueryKeepAccountdata()
}

func TestUpdatedata(t *testing.T) {
	//
	Updatedata()

}

func TestUpdateJZdata(t *testing.T) {
	//
	UpdateJZdata()
}

func TestXorminitTest(t *testing.T) {
	//XorminitTest()
	XormInsert()
}

func TestQueryKeepAccountMsgdata(t *testing.T) {
	QueryKeepAccountMsgdata()
}
