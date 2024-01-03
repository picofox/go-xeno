package message_buffer_test

import (
	"fmt"
	"testing"
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/inet/transcomm"
	"xeno/zohar/core/memory"
)

func Test_Serialization_Functional_Basic(t *testing.T) {
	var hcache []byte = make([]byte, 4)
	m := messages.NeoProcTestMessage(false)
	bufList := memory.NeoByteBufferList()
	totalLen, checkLen, rc := m.PiecewiseSerialize(bufList)
	if core.Err(rc) {
		t.Errorf("Validation failed")
	}
	if checkLen != m.BodyLength() {
		t.Errorf("Validation 2 failed")
	}
	if totalLen == 0 {
		t.Errorf("Validation 2 failed")
	}

	//m = messages.NeoProcTestMessage(false)
	//totalLen, checkLen, rc = m.PiecewiseSerialize(bufList)
	//if core.Err(rc) {
	//	t.Errorf("Validation failed")
	//}
	//if checkLen != m.BodyLength() {
	//	t.Errorf("Validation 2 failed")
	//}
	//
	//if totalLen == 0 {
	//	t.Errorf("Validation 2 failed")
	//}
	//
	//fmt.Printf("check len %d\n", checkLen)

	r := &transcomm.HandlerRegistration{}
	handler := r.NeoO1L15COT15DecodeClientHandler(nil)
	var bodyLen int64 = 0

	bodyLen, rc = handler.CheckCompletion(bufList.Front())

	fmt.Printf("bodylen = %d rc = %d\n", bodyLen, rc)

	messages.ReadHeaderContent(hcache, bufList)

	rMsg := messages.ProcTestMessagePiecewiseDeserialize(bufList, bodyLen)
	if rMsg == nil {
		t.Errorf("Serial failed")
	}

	//messages.ReadHeaderContent(hcache, bufList)

}
