package inet

import (
	"syscall"
	"unsafe"
	"xeno/zohar/core"
)

func SysRead(fd int, ba []byte) (int64, int32) {
	n, err := syscall.Read(fd, ba)
	if err != nil {
		if err == syscall.EAGAIN || err == syscall.EINTR {
			return int64(n), core.MkErr(core.EC_TRY_AGAIN, 1)
		}
		return int64(n), core.MkErr(core.EC_FILE_READ_FAILED, 1)
	}
	return int64(n), core.MkSuccess(0)
}

const EPOLLET = -syscall.EPOLLET

type EPollEvent struct {
	Events uint32
	Data   [8]byte // unaligned uintptr
}

func EpollCreate(flag int) (fd int, err error) {
	var r0 uintptr
	r0, _, err = syscall.RawSyscall(syscall.SYS_EPOLL_CREATE1, uintptr(flag), 0, 0)
	if err == syscall.Errno(0) {
		err = nil
	}
	return int(r0), err
}

// EpollCtl implements epoll_ctl.
func EpollCtl(epfd int, op int, fd int, event *EPollEvent) (err error) {
	_, _, err = syscall.RawSyscall6(syscall.SYS_EPOLL_CTL, uintptr(epfd), uintptr(op), uintptr(fd), uintptr(unsafe.Pointer(event)), 0, 0)
	if err == syscall.Errno(0) {
		err = nil
	}
	return err
}

// EpollWait implements epoll_wait.
func EpollWait(epfd int, events []EPollEvent, msec int) (n int, err error) {
	var r0 uintptr
	var _p0 = unsafe.Pointer(&events[0])
	if msec == 0 {
		r0, _, err = syscall.RawSyscall6(syscall.SYS_EPOLL_WAIT, uintptr(epfd), uintptr(_p0), uintptr(len(events)), 0, 0, 0)
	} else {
		r0, _, err = syscall.Syscall6(syscall.SYS_EPOLL_WAIT, uintptr(epfd), uintptr(_p0), uintptr(len(events)), uintptr(msec), 0, 0)
	}
	if err == syscall.Errno(0) {
		err = nil
	}
	return int(r0), err
}

func EPollWaitUntil(pollFD int, events []EPollEvent, msec int) (int, int32) {
WAIT:
	rc := core.MkSuccess(0)
	n, err := EpollWait(pollFD, events, msec)
	if err != nil {
		if err == syscall.EINTR {
			goto WAIT
		} else {
			rc = core.MkErr(core.EC_EPOLL_WAIT_ERROR, 1)
		}
	}
	return n, rc
}
