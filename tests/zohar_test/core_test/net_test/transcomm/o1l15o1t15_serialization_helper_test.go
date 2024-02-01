package transcomm

import (
	"fmt"
	"testing"
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/strs"
)

func Test_O1L15O1C15Serializer_Functional_Basic(t *testing.T) {
	var buffer *memory.LinkedListByteBuffer = memory.NeoLinkedListByteBuffer(datatype.SIZE_4K)

	ctx, rc := messages.InitializeSerialization(buffer, 1)
	if core.Err(rc) {
		panic("InitializeSerialization Failed")
	}

	s := strs.CreateSampleString(32760, "@", "$")
	rc = ctx.WriteString(s)
	if core.Err(rc) {
		panic("InitializeSerialization Failed")
	}
	rc = ctx.FinalizeSerialization()
	if core.Err(rc) {
		panic("InitializeSerialization Failed")
	}
	fmt.Printf("%s\n", ctx.String())

	v, rc := buffer.PeekInt16(0)
	fmt.Print(v)
}
