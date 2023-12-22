package container

type ISinglyLinkedListNode interface {
	Next() ISinglyLinkedListNode
	SetNext(node ISinglyLinkedListNode)
}
