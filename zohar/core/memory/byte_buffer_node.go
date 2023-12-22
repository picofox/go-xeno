package memory

type ByteBufferNode struct {
	_buffer *LinearBuffer
	_next   *ByteBufferNode
}

func (ego *ByteBufferNode) Clear() {
	ego._buffer.Clear()
}

func (ego *ByteBufferNode) Buffer() *LinearBuffer {
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
		_buffer: NeoLinearBuffer(capa),
		_next:   nil,
	}
}
