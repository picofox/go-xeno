package message_buffer

import "xeno/zohar/core/memory"

type MessageHeader struct {
	_data []byte
}

func (ego *MessageHeader) Data() []byte {
	return ego._data
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
