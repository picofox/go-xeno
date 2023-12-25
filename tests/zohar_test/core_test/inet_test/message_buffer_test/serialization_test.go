package message_buffer_test

import (
	"fmt"
	"testing"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
)

func Test_Serialization_Functional_Basic(t *testing.T) {
	m := messages.NeoProcTestMessage(false)
	bufList := memory.NeoByteBufferList()
	totalLen, rc := m.SerializeToList(bufList)
	fmt.Printf("totalLen=%d rc=%d\n", totalLen, rc)
}
