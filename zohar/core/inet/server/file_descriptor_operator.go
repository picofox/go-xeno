package server

import (
	"runtime"
	"sync/atomic"
)

// FileDescriptorOperator is a collection of operations on file descriptors.
type FileDescriptorOperator struct {
	// FD is file descriptor, poll will bind when register.
	FD int

	// The FileDescriptorOperator provides three operations of reading, writing, and hanging.
	// The poll actively fire the FileDescriptorOperator when fd changes, no check the return value of FileDescriptorOperator.
	OnRead  func(p IPoll) error
	OnWrite func(p IPoll) error
	OnHup   func(p IPoll) error

	// The following is the required fn, which must exist when used, or directly panic.
	// Fns are only called by the poll when handles connection events.
	Inputs   func(vs [][]byte) (rs [][]byte)
	InputAck func(n int) (err error)

	// Outputs will locked if len(rs) > 0, which need unlocked by OutputAck.
	Outputs   func(vs [][]byte) (rs [][]byte, supportZeroCopy bool)
	OutputAck func(n int) (err error)

	// poll is the registered location of the file descriptor.
	poll IPoll

	// protect only detach once
	detached int32

	// private, used by operatorCache
	next  *FileDescriptorOperator
	state int32 // CAS: 0(unused) 1(inuse) 2(do-done)
	index int32 // index in operatorCache
}

func (op *FileDescriptorOperator) Control(event PollEvent) error {
	if event == PollDetach && atomic.AddInt32(&op.detached, 1) > 1 {
		return nil
	}
	return op.poll.Control(op, event)
}

func (op *FileDescriptorOperator) Free() {
	op.poll.Free(op)
}

func (op *FileDescriptorOperator) do() (can bool) {
	return atomic.CompareAndSwapInt32(&op.state, 1, 2)
}

func (op *FileDescriptorOperator) done() {
	atomic.StoreInt32(&op.state, 1)
}

func (op *FileDescriptorOperator) inuse() {
	for !atomic.CompareAndSwapInt32(&op.state, 0, 1) {
		if atomic.LoadInt32(&op.state) == 1 {
			return
		}
		runtime.Gosched()
	}
}

func (op *FileDescriptorOperator) unused() {
	for !atomic.CompareAndSwapInt32(&op.state, 1, 0) {
		if atomic.LoadInt32(&op.state) == 0 {
			return
		}
		runtime.Gosched()
	}
}

func (op *FileDescriptorOperator) isUnused() bool {
	return atomic.LoadInt32(&op.state) == 0
}

func (op *FileDescriptorOperator) reset() {
	op.FD = 0
	op.OnRead, op.OnWrite, op.OnHup = nil, nil, nil
	op.Inputs, op.InputAck = nil, nil
	op.Outputs, op.OutputAck = nil, nil
	op.poll = nil
	op.detached = 0
}
