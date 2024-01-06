package memory

type ByteBufferNode struct {
	LinearBufferFixed
	_next *ByteBufferNode
}

func (ego *ByteBufferNode) ReadAvailableByOffset(offset int64) int64 {
	if offset < ego._beginPos {
		panic("offset < begin")
	} else {
		return ego._length - (offset - ego._beginPos)
	}
}

func (ego *ByteBufferNode) Clear() {
	ego.LinearBufferFixed.Clear()
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
		LinearBufferFixed: LinearBufferFixed{
			_capacity: capa,
			_beginPos: 0,
			_length:   0,
			_data:     make([]byte, capa),
			_cache:    make([]byte, 8),
		},
		_next: nil,
	}
}
