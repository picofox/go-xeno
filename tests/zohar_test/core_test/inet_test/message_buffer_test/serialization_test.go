package message_buffer_test

import (
	"sync"
	"testing"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/inet/transcomm"
	"xeno/zohar/core/memory"
)

var gBufList *memory.ByteBufferList = memory.NeoByteBufferList()
var gLock sync.Mutex

func _addMessage(t *testing.T, m message_buffer.INetMessage) int64 {
	gLock.Lock()
	defer gLock.Unlock()
	sw := chrono.NeoStopWatch()
	sw.Begin()
	totalLen, checkLen, rc := m.PiecewiseSerialize(gBufList)
	if core.Err(rc) {
		t.Errorf("[P] : PiecewiseSerialize Failed")
	}
	if checkLen != m.BodyLength() {
		t.Errorf("[P] : body len Check Failed")
	}
	if totalLen < checkLen {
		t.Errorf("[P] : total len Check Failed")
	}
	cost := sw.Stop()
	return cost
}

func _getMessage(t *testing.T, handler *transcomm.O1L15COT15CodecClientHandler) (message_buffer.INetMessage, int64, int32) {

	sw := chrono.NeoStopWatch()
	gLock.Lock()
	defer gLock.Unlock()

	sw.Begin()
	bodyLen, cmd, rc := handler.CheckCompletion(gBufList.Front())
	if core.Err(rc) {
		return nil, -1, core.MkSuccess(0)
	}
	rMsg := messages.GetDefaultMessageBufferDeserializationMapper().Deserialize(cmd, gBufList, bodyLen)
	if rMsg == nil {
		t.Errorf("[C] : Derialization Failed")
		return nil, -1, core.MkErr(core.EC_DESERIALIZE_FIELD_FAIELD, 1)
	}
	cost := sw.Stop()

	return rMsg, cost, core.MkSuccess(0)
}

func messageConsumer(t *testing.T) {
	r := &transcomm.HandlerRegistration{}
	handler := r.NeoO1L15COT15DecodeClientHandler(nil)

	for {
		msg, cost, rc := _getMessage(t, handler)
		if core.Err(rc) {
			t.Errorf("[C] : Derialization Failed")
		} else {
			if msg != nil {
				t.Logf("[C] - message (%s) cost (%d)", msg.IdentifierString(), cost)
			}

		}
	}
}

func messageProducer1(t *testing.T) {

	for {
		m := messages.NeoProcTestMessage(false)
		if m == nil {
			t.Errorf("[P] : Create Message Failed")
		}
		rc := _addMessage(t, m)
		t.Logf("[P] - message (%s) cost (%d)", m.IdentifierString(), rc)
		time.Sleep(10 * time.Millisecond)
	}
}

func messageProducer2(t *testing.T) {
	for {
		m := messages.NeoKeepAliveMessage(false)
		if m == nil {
			t.Errorf("[P] : Create Message Failed")
		}
		_addMessage(t, m)
		t.Log("[P] + message")
		time.Sleep(140 * time.Millisecond)
	}
}

func Test_Serialization_Functional_Basic(t *testing.T) {
	go messageConsumer(t)
	go messageProducer1(t)
	go messageProducer1(t)
	for {
		t.Logf("BL=%d", gBufList.Count())
		time.Sleep(1000 * time.Millisecond)
	}
}
