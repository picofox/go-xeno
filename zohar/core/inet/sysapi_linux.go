package inet

import (
	"os"
	"syscall"
	"unsafe"
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/logging"
)

func SetDefaultSockopts(s, family, sotype int, ipv6only bool) int32 {
	if family == syscall.AF_INET6 && sotype != syscall.SOCK_RAW {
		// Allow both IP versions even if the OS default
		// is otherwise. Note that some operating systems
		// never admit this option.
		syscall.SetsockoptInt(s, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, datatype.BoolToInt(ipv6only))
	}

	// Allow broadcast.
	if os.NewSyscallError("setsockopt", syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)) != nil {
		return core.MkErr(core.EC_SET_NONBLOCK_ERROR, 1)
	}
	return core.MkSuccess(0)
}

func SysSocket(family, sotype, proto int) (int, int32) {
	// See ../syscall/exec_unix.go for description of ForkLock.
	syscall.ForkLock.RLock()
	s, err := syscall.Socket(family, sotype, proto)
	if err == nil {
		syscall.CloseOnExec(s)
	}
	syscall.ForkLock.RUnlock()
	if err != nil {
		return -1, core.MkErr(core.EC_CREATE_SOCKET_ERROR, 0)
	}
	if err = syscall.SetNonblock(s, true); err != nil {
		syscall.Close(s)
		return -1, core.MkErr(core.EC_SET_NONBLOCK_ERROR, 0)
	}
	return s, core.MkSuccess(0)
}

func SysRead(fd int, ba []byte) (int64, int32) {
	n, err := syscall.Read(fd, ba)
	if err != nil {
		logging.GetLoggerManager().GetDefaultLogger().Log(core.LL_SYS, "SysRead: read error %s", err.Error())
		if err == syscall.EAGAIN || err == syscall.EINTR {
			return int64(n), core.MkErr(core.EC_TRY_AGAIN, 1)
		}
		return int64(n), core.MkErr(core.EC_FILE_READ_FAILED, 1)
	}
	return int64(n), core.MkSuccess(0)
}

func SysWriteN(fd int, p []byte) (int64, int32) {
	var nDone int64 = 0
	var nToBeWrite = int64(len(p))
	for nDone < nToBeWrite {
		n, err := syscall.Write(fd, p[nDone:])
		if err != nil {
			if err == syscall.EAGAIN {
				return nDone + int64(n), core.MkErr(core.EC_TRY_AGAIN, 1)
			}
			return nDone + int64(n), core.MkErr(core.EC_FILE_WRITE_FAILED, 1)
		}
		nDone += int64(n)
	}
	return nDone, core.MkSuccess(0)
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

func Socket(family, sotype, proto int) (int, int32) {
	fd, rc := SysSocket(family, sotype, proto)
	if core.Err(rc) {
		return -1, rc
	}

	rc = SetDefaultSockopts(fd, family, sotype, false)
	if core.Err(rc) {
		syscall.Close(fd)
		return -1, rc
	}

	return fd, core.MkSuccess(0)
}
