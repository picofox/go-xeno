package transcomm

import "sync"

type HandlerRegistration struct {
}

var sHandlerRegistration *HandlerRegistration
var sHandlerRegistrationOnce sync.Once

func GetHandlerRegistration() *HandlerRegistration {
	sHandlerRegistrationOnce.Do(func() {
		sHandlerRegistration = &HandlerRegistration{}
	})
	return sHandlerRegistration
}
