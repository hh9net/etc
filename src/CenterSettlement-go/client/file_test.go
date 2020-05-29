package client

import "testing"

func TestMoveFile(t *testing.T) {
	s := "../sendzipxml/CZ_3201_00000000000000100027.xml.lz77"
	d := "../sendxmlsucceed/CZ_3201_00000000000000100027.xml.lz77"
	MoveFile(s, d)
}
