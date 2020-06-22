package centerServer

import (
	log "github.com/sirupsen/logrus"
	"os"
)

//联网中心 压缩 记账包 "jz"、争议包"zy"、清分包"qf"

func CenterZipxml(fname string, lx string) error {
	var pwd string
	switch lx {
	//记账
	case "jz":
		pwd = "../centerkeepaccount/" + fname
	//争议
	case "zy":
		pwd = "../centerdispute/" + fname
	//清分
	case "qf":
		pwd = "../centerClearing/" + fname
	}
	log.Println("zipxmlpwd:", pwd)

	origin, oerr := os.Open(pwd)
	if oerr != nil {
		log.Fatalln(oerr)
		return oerr
	}
	defer origin.Close()
	outpwd := "../centerSendxmlzip/" + fname + ".lz77"
	out, cerr := os.Create(outpwd)
	if cerr != nil {
		log.Fatalln(cerr)
		return cerr
	}
	defer out.Close()

	if zerr := Compress(origin, out); zerr != nil {
		log.Fatalln(zerr)
		return zerr
	}
	return nil
}
