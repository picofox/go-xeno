package message_buffer_test

import (
	"sync"
	"sync/atomic"
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
var pCount atomic.Int64
var cCount int64
var TEST_COUNT int = 500000
var sCompleteState *message_buffer.CheckBufferCompletionState = message_buffer.NeoCheckBufferCompletionState()

var totalCount int64 = 0
var okCount int64 = 0

func _addMessage(t *testing.T, m message_buffer.INetMessage) int64 {
	sw := chrono.NeoStopWatch()
	gLock.Lock()
	defer gLock.Unlock()
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
	pCount.Add(1)

	return cost
}

func _getMessage(t *testing.T, handler *transcomm.O1L15COT15CodecServerHandler) (message_buffer.INetMessage, int64, int32) {

	sw := chrono.NeoStopWatch()

	gLock.Lock()
	defer gLock.Unlock()

	//gCond.Wait()

	if !sCompleteState.CouldTry(gBufList) {
		return nil, -1, core.MkSuccess(0)
	}

	sw.Begin()
	totalCount++
	bodyLen, cmd, rc := handler.CheckCompletion(gBufList.Front())
	if core.Err(rc) {
		sCompleteState.Update(false, gBufList.Back(), gBufList.Back().WritePos())
		return nil, -1, core.MkSuccess(0)
	}
	sCompleteState.Update(true, gBufList.Back(), gBufList.Back().WritePos())
	okCount++
	rMsg := messages.GetDefaultMessageBufferDeserializationMapper().Deserialize(cmd, gBufList, bodyLen)
	if rMsg == nil {
		t.Errorf("[C] : Derialization Failed")
		return nil, -1, core.MkErr(core.EC_DESERIALIZE_FIELD_FAIELD, 1)
	}
	cost := sw.Stop()

	cCount++
	return rMsg, cost, core.MkSuccess(0)
}

func messageConsumer(t *testing.T) {
	r := &transcomm.HandlerRegistration{}
	handler := r.NeoO1L15COT15DecodeServerHandler(nil)

	for {
		msg, _, rc := _getMessage(t, handler)
		if core.Err(rc) {
			t.Errorf("[C] : Derialization Failed")
		} else {
			if msg != nil {
				//t.Logf("[C] - message (%s) cost (%d)", msg.IdentifierString(), cost)
			}
		}
	}
}

func messageProducer1(t *testing.T) {
	for i := 0; i < TEST_COUNT; i++ {
		m := messages.NeoProcTestMessage(false)
		if m == nil {
			t.Errorf("[P] : Create Message Failed")
		}
		_ = _addMessage(t, m)
		//t.Logf("[P] - message (%s) cost (%d)", m.IdentifierString(), rc)
		time.Sleep(1 * time.Millisecond)
	}
}

var done bool = false

func Test_Serialization_Functional_Basic(t *testing.T) {
	go messageConsumer(t)
	go messageProducer1(t)
	go messageProducer1(t)
	go messageProducer1(t)

	for {
		gLock.Lock()
		t.Logf("BL=%d pCount=%d cCount=%d, tCount=%d ocount=%d", gBufList.Count(), pCount.Load(), cCount, totalCount, okCount)
		gLock.Unlock()
		time.Sleep(1000 * time.Millisecond)
		//if !done {
		//	if pCount.Load() >= int64(TEST_COUNT) {
		//		fmt.Printf("list:\n")
		//		done = true
		//
		//		str := gBufList.String()
		//		fmt.Printf(str)
		//
		//	}
		//}

	}
}
