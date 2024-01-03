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

	fmt.Printf("check len %d\n", checkLen)

	r := &transcomm.HandlerRegistration{}
	handler := r.NeoO1L15COT15DecodeClientHandler(nil)
	var bodyLen int64 = 0
	bodyLen, rc = handler.CheckCompletion(bufList.Front())
	fmt.Printf("bodylen = %d rc = %d\n", bodyLen, rc)

	messages.ReadHeaderContent(hcache, bufList)

	rMsg := messages.ProcTestMessagePiecewiseDeserialize(bufList, bodyLen)

	fmt.Printf("%s\n", rMsg.String())

	//var idx = 0
	//var lb *memory.LinearBuffer = memory.NeoLinearBuffer(1024 * 2048)
	//for n := bufList.Front(); n != nil; n = n.Next() {
	//	lb.WriteRawBytes(*n.InternalData(), 0, n.ReadAvailable())
	//	idx++
	//}
	//
	//
	//
	//for {
	//	lenAndOpt, _ := lb.ReadInt16()
	//	cmdAndOpt, _ := lb.ReadInt16()
	//	h := message_buffer.NeoMessageHeader()
	//	h.SetRaw2(lenAndOpt, cmdAndOpt)
	//	fmt.Printf("header <%s>\n", h.String())
	//	fmt.Printf("totalLen=%d rc=%d seeklen=%d\n", totalLen, rc, h.Length())
	//
	//	if h.Length() == 0 {
	//		break
	//	}
	//	b := lb.ReaderSeek(memory.BUFFER_SEEK_CUR, int64(h.Length()))
	//	if !b {
	//		break
	//	}
	//}

}
