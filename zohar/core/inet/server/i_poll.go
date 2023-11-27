package server

type IPoll interface {
	// Wait will poll all registered fds, and schedule processing based on the triggered event.
	// The call will block, so the usage can be like:
	//
	//  go wait()
	//
	Wait() error

	// Close the poll and shutdown Wait().
	Close() error

	// Trigger can be used to actively refresh the loop where Wait is located when no event is triggered.
	// On linux systems, eventfd is used by default, and kevent by default on bsd systems.
	Trigger() error

	// Control the event of file descriptor and the operations is defined by PollEvent.
	Control(operator *FileDescriptorOperator, event PollEvent) error

	// Alloc the operator from cache.
	Alloc() (operator *FileDescriptorOperator)

	// Free the operator from cache.
	Free(operator *FileDescriptorOperator)
}

// PollEvent defines the operation of poll.Control.
type PollEvent int

const (
	// PollReadable is used to monitor whether the FDOperator registered by
	// listener and connection is readable or closed.
	PollReadable PollEvent = 0x1

	// PollWritable is used to monitor whether the FDOperator created by the dialer is writable or closed.
	// ET mode must be used (still need to poll hup after being writable)
	PollWritable PollEvent = 0x2

	// PollDetach is used to remove the FDOperator from poll.
	PollDetach PollEvent = 0x3

	// PollR2RW is used to monitor writable for FDOperator,
	// which is only called when the socket write buffer is full.
	PollR2RW PollEvent = 0x5

	// PollRW2R is used to remove the writable monitor of FDOperator, generally used with PollR2RW.
	PollRW2R PollEvent = 0x6
)
