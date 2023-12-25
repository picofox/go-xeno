package memory

type ByteBufferNode struct {
	_buffer *LinearBufferFixed
	_next   *ByteBufferNode
}

func (ego *ByteBufferNode) Clear() {
	ego._buffer.Clear()
}

func (ego *ByteBufferNode) Buffer() *LinearBufferFixed {
	return ego._buffer
}

func (ego *ByteBufferNode) Next() *ByteBufferNode {
	return ego._next
}

func (ego *ByteBufferNode) SetNext(node *ByteBufferNode) {
	if node == nil {
		ego._next = nil
	} else {
		ego._next = node
	}
}

func NeoByteBufferNode(capa int64) *ByteBufferNode {
	return &ByteBufferNode{
		_buffer: NeoLinearBufferFixed(capa),
		_next:   nil,
	}
}

func AdoptByteBufferNode(byteBuff *LinearBufferFixed) *ByteBufferNode {
	return &ByteBufferNode{
		_buffer: byteBuff,
		_next:   nil,
	}
}
