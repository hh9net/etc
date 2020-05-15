package common

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
)

//this file is made to finish it that change a number from base 10 to hex
func ToHex(ten int) string {
	s := fmt.Sprintf("%020x", ten)
	return s
}

func ToHex1(ten int) (hex []int, length int) {
	m := 0
	hex = make([]int, 0)
	length = 0
	for {
		m = ten / 16
		ten = ten % 16
		if m == 0 {
			hex = append(hex, ten)
			length++
			break
		}
		hex = append(hex, m)
		length++
	}
	return
}

func HexToTen(Hex string) {
	bs, _ := hex.DecodeString(Hex)
	num := binary.BigEndian.Uint16(bs[:2])
	fmt.Println(num)
}

func ToHexFormat(ten int, numb int) string {
	var s string
	switch numb {
	case 2:
		s = fmt.Sprintf("%02x", ten)
		return s
	case 6:
		s = fmt.Sprintf("%02x", ten)
		return s
	case 8:
		s = fmt.Sprintf("%08x", ten)
		return s
	case 12:
		s = fmt.Sprintf("%08x", ten)
		return s
	default:
		log.Println("请检查函数参数")
	}
	return ""
}
