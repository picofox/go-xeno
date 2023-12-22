package memory

import "xeno/zohar/core/container"

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

func (ego *ByteBufferNode) Next() container.ISinglyLinkedListNode {
	return ego._next
}

func (ego *ByteBufferNode) SetNext(node container.ISinglyLinkedListNode) {
	if node == nil {
		ego._next = nil
	} else {
		ego._next = node.(*ByteBufferNode)
	}
}

func NeoByteBufferNode(capa int64) *ByteBufferNode {
	return &ByteBufferNode{
		_buffer: NeoLinearBuffer(capa),
		_next:   nil,
	}
}

var _ container.ISinglyLinkedListNode = &ByteBufferNode{}
