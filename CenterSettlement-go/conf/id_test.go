package conf

import (
	"log"
	"testing"
)

func TestGenerateMessageId(t *testing.T) {
	messageid := GenerateMessageId()
	log.Println(messageid)
}
