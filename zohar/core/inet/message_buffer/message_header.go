package message_buffer

import (
	"fmt"
	"xeno/zohar/core/memory"
)

type MessageHeader struct {
	_data []byte
}

func (ego *MessageHeader) Data() []byte {
	return ego._data
}

func (ego *MessageHeader) String() string {
	u0 := memory.BytesToInt16BE(&ego._data, 0)
	u1 := memory.BytesToInt16BE(&ego._data, 2)
	o1 := u0>>15&0x1 == 1
	o2 := u1>>15&0x1 == 1
	cmd := int16(u1 & 0x7FFF)
	l := int16(u0 & 0x7FFF)

	return fmt.Sprintf("%d:%d:%t:%t", l, cmd, o1, o2)
}

func (ego *MessageHeader) Clear() {
	ego._data[0] = 0
	ego._data[1] = 0
	ego._data[2] = 0
	ego._data[3] = 0
}

func (ego *MessageHeader) SetRaw2(lenAndO0 int16, cmdAndO1 int16) {
	memory.Int16IntoBytesBE(lenAndO0, &ego._data, 0)
	memory.Int16IntoBytesBE(cmdAndO1, &ego._data, 2)
}

func (ego *MessageHeader) Length() int16 {
	var lenAndO0 int16
	lenAndO0 = memory.BytesToInt16BE(&ego._data, 0)
	return lenAndO0 & 0x7FFF
}

func (ego *MessageHeader) Set(o0 bool, o1 bool, length int16, cmd int16) {
	var lenAndO0 int16 = length
	var cmdAndO1 int16 = cmd
	if o0 {
		iv := 1 << 15
		lenAndO0 = length | int16(iv)
	}
	if o1 {
		iv := 1 << 15
		cmdAndO1 = cmd | int16(iv)
	}

	memory.Int16IntoBytesBE(lenAndO0, &ego._data, 0)
	memory.Int16IntoBytesBE(cmdAndO1, &ego._data, 2)
}

func NeoMessageHeader() MessageHeader {
	return MessageHeader{
		_data: make([]byte, 4),
	}
}
