package message_buffer

import "xeno/zohar/core/memory"

type CheckBufferCompletionState struct {
	_lastComplete  bool
	_lastEndBuffer *memory.ByteBufferNode
	_lastReadIndex int64
}

func (ego *CheckBufferCompletionState) CouldTry(list *memory.ByteBufferList) bool {
	if list == nil || list.Back() == nil {
		return false
	} else {
		if ego._lastComplete {
			return true
		} else {
			if ego._lastEndBuffer != list.Back() || ego._lastReadIndex != list.Back().WritePos() {
				return true
			}
		}
	}
	return false
}

func (ego *CheckBufferCompletionState) Update(bComplete bool, buffer *memory.ByteBufferNode, readIndex int64) {
	ego._lastComplete = true
	ego._lastEndBuffer = buffer
	ego._lastReadIndex = readIndex
}

func NeoCheckBufferCompletionState() *CheckBufferCompletionState {
	return &CheckBufferCompletionState{
		_lastComplete:  true,
		_lastEndBuffer: nil,
		_lastReadIndex: 0,
	}
}
