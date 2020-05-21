package types

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_XMLUnMarshall(t *testing.T) {
	m := Message{}

	data := `<Message>
    <Header>
        <Version>00010000</Version>
        <MessageClass>5</MessageClass>
        <MessageType>7</MessageType>
        <SenderId>00000000000000FD</SenderId>
        <ReceiverId>0000000000000020</ReceiverId>
        <MessageId>97640</MessageId>
    </Header>
    <Body ContentType="1">
        <ClearTargetDate>2020-04-15</ClearTargetDate>
        <ServiceProviderId>00000000000000FD</ServiceProviderId>
        <IssuerId>0000000000000020</IssuerId>
        <MessageId>97640</MessageId>
        <Count>4</Count>
        <Amount>35.00</Amount>
        <Transaction>
            <TransId>1</TransId>
            <Time>2020-04-15 11:25:27</Time>
            <Fee>18.00</Fee>
            <Service>
                <ServiceType>2</ServiceType>
                <Description>姹</Description>
                <Detail>1|04|3201|3201000003|1104|20200415 112527|03|3201|3201000003|1003|20200415 081421</Detail>
            </Service>
            <ICCard>
                <CardType>23</CardType>
                <NetNo>3401</NetNo>
                <CardId>1030230212304558</CardId>
                <License></License>
                <PreBalance>19999216.82</PreBalance>
                <PostBalance>19999198.82</PostBalance>
            </ICCard>
            <Validation>
                <TAC>f05e6ba9</TAC>
                <TransType>09</TransType>
                <TerminalNo>01320002d9fd</TerminalNo>
                <TerminalTransNo>00010bc0</TerminalTransNo>
            </Validation>
            <OBU>
                <NetNo>3401</NetNo>
                <OBUId>34011703080a830a</OBUId>
                <OBEState>0288</OBEState>
                <License></License>
            </OBU>
            <CustomizedData>f05e6ba9000007080901320002d9fd00010bc020200415112527000077345B0A773462120070</CustomizedData>
        </Transaction>
        <Transaction>
            <TransId>2</TransId>
            <Time>2020-04-15 11:25:40</Time>
            <Fee>1.50</Fee>
            <Service>
                <ServiceType>2</ServiceType>
                <Description>姹</Description>
                <Detail>1|04|3201|3201000003|1105|20200415 112540|03|3201|3201000003|1003|20200415 110044</Detail>
            </Service>
            <ICCard>
                <CardType>23</CardType>
                <NetNo>3401</NetNo>
                <CardId>1030230213539454</CardId>
                <License></License>
                <PreBalance>19980136.58</PreBalance>
                <PostBalance>19980135.08</PostBalance>
            </ICCard>
            <Validation>
                <TAC>319249ce</TAC>
                <TransType>09</TransType>
                <TerminalNo>01320002da0e</TerminalNo>
                <TerminalTransNo>00010099</TerminalTransNo>
            </Validation>
            <OBU>
                <NetNo>3401</NetNo>
                <OBUId>3401150908f6f099</OBUId>
                <OBEState>0288</OBEState>
                <License></License>
            </OBU>
            <CustomizedData>319249ce000000960901320002da0e0001009920200415112540000077174444771744DA0502</CustomizedData>
        </Transaction>
    </Body>
</Message>`

	err := xml.Unmarshal([]byte(data), &m)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Printf("XMLName: %#v\n", m.XMLName)
	fmt.Printf("Header: %q\n", m.Header)
	fmt.Printf("m.Header.MessageClass: %q\n", m.Header.MessageClass)
	fmt.Printf("m.Header.MessageId: %v\n", m.Header.MessageId)
	fmt.Printf(" m.Header.SenderId: %v\n", m.Header.SenderId)
	fmt.Printf("Transaction: %v\n", m.Body.Transaction)
}

func Test_XMLMarshal(t *testing.T) {
	//把原始交易数据转化成xml文件
	jiaoyi := &Message{
		Header: Header{Version: "000100", MessageClass: 1, MessageType: 5, SenderId: "12345", ReceiverId: "42131", MessageId: 0000000001},
		Body:   Body{}}
	jiaoyi.Body.Transaction = []Transaction{{TransId: 123}, {TransId: 456}}

	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(jiaoyi, "  ", "    ")
	//使用Marshal函数，生成的XML格式无缩进
	//outputxml,err:=xml.Marshal(v)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	//os.Stdout.Write(outputxml)

	//存为xml文件
	//
	//num := strconv.Itoa(9999999)
	//NeedAdd := 20 - len(num)
	//var newByte string
	//if NeedAdd > 0 {
	//	newByte = strings.Repeat("0",NeedAdd)
	//	newByte = newByte + num
	//}
	//fname:=newByte
	num := 888888
	fname := fmt.Sprintf("%020d", num)
	//创建文件"
	fw, f_werr := os.Create("../sendzipxml/" + fname + ".xml")
	if f_werr != nil {
		log.Fatal("Read:", f_werr)
	}
	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, outputxml...)
	//写入文件
	ioutil.WriteFile("test.xml", xmlOutPutData, os.ModeAppend)
	defer fw.Close()
	_, ferr := fw.Write((xmlOutPutData))
	if ferr != nil {
		fmt.Printf("error: %v\n", err)
	}
}
