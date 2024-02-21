package message_buffer_test

import (
	"fmt"
	"math/rand"
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
var sPM_TEST_COUNT int64 = 10000000
var sw *chrono.StopWatch = chrono.NeoStopWatch()

func _addMessage() {
	sLock.Lock()
	defer sLock.Unlock()

	hdr := memory.NeoO1L31C16Header(0, 0)
	msg := messages.NeoProcTestMessage(false)
	if msg == nil {
		panic("create msg failed")
	}
	_, rc := msg.Serialize(hdr, sBuffer)
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

	hdr, rc := memory.O1L31C16HeaderFromBuffer(sBuffer)
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return 0, false
		}
		panic("IsMessageComplete Failed")
	}
	sBuffer.ReaderSeek(memory.BUFFER_SEEK_CUR, 6)
	msg, dataLength := messages.ProcTestMessageDeserialize(sBuffer, hdr)
	if msg == nil {
		panic("Deser to create msg failed")
	}

	if core.Err(msg.(*messages.ProcTestMessage).Validate()) {
		panic("Data validation Failed")
	}

	return dataLength, true
}

var bsCount int64 = 0

func debugBuffer() {
	var bytesFlushed int64 = 0
	for {
		buf := sBuffer.(*memory.LinkedListByteBuffer).InternalDataForReading()
		if buf != nil {
			fmt.Printf("Get Buffer Len:%d\n", len(buf))
			bytesFlushed += int64(len(buf))
			if !sBuffer.ReaderSeek(memory.BUFFER_SEEK_CUR, int64(len(buf))) {
				panic("error")
			}
		} else {
			fmt.Printf("end buffer_len:%d\n", sBuffer.ReadAvailable())
			return
		}
	}
}

func Test_Serialization_Functional_Basic(t *testing.T) {
	rand.Seed(0)

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
