package message_buffer_test

import (
	"fmt"
	"sync"
	"testing"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
)

var sBuffer memory.IByteBuffer = memory.NeoLinkedListByteBuffer(datatype.SIZE_4K)
var sLock sync.Mutex
var sCount int64
var sPM_TEST_COUNT int64 = 1000000
var sw *chrono.StopWatch = chrono.NeoStopWatch()

func _addMessage() {
	sLock.Lock()
	defer sLock.Unlock()
	msg := messages.NeoProcTestMessage(false)
	if msg == nil {
		panic("create msg failed")
	}
	_, rc := msg.O1L15O1T15Serialize(sBuffer)
	if core.Err(rc) {
		panic("Serialize msg Failed")
	}
}

func messageProducer(t *testing.T) {
	for i := 0; i < int(sPM_TEST_COUNT); i++ {
		_addMessage()
	}
	fmt.Printf("Producer Done")
}

func _getMessage() (int64, bool) {
	sLock.Lock()
	defer sLock.Unlock()

	_, _, ll, el, rc := messages.IsMessageComplete(sBuffer)
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return 0, false
		}
		panic("IsMessageComplete Failed")
	}
	msg, dataLength := messages.ProcTestMessageDeserialize(sBuffer, ll, el)
	if msg == nil {
		panic("Deser to create msg failed")
	}

	if !msg.(*messages.ProcTestMessage).Validate() {
		panic("Data validation Failed")
	}

	return dataLength, true
}

var bsCount int64 = 0

func Test_Serialization_Functional_Basic(t *testing.T) {
	//r := &transcomm.HandlerRegistration{}
	//handler := r.NeoO1L15COT15DecodeServerHandler(nil)
	go messageProducer(t)
	var bsPrintNext int64 = 1
	sw.Begin("x")
	for {
		l, r := _getMessage()
		if r {
			bsCount += l
			sCount++
			if sCount == sPM_TEST_COUNT {
				fmt.Printf("Done\n")
				break
			}

			if sCount%1000000 == 0 {
				fmt.Printf("Total: %d count\n", sCount)
			}

			if bsCount > bsPrintNext*1048576*1000 {
				sw.Mark(".")
				t := sw.GetRecordRecent()
				avgspd := float64(bsCount) / 1048576.0 / (float64(t) / 1000000000.0)
				fmt.Printf("Total: %.2f MB MeanSPD = %.2f MBPS\n", float64(bsCount)/1048576.0, avgspd)
				bsPrintNext++
			}

		}
	}
}
