// Code generated by 'go generate'; DO NOT EDIT.

package main

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modKernelBase = windows.NewLazySystemDLL("KernelBase.dll")
	modkernel32   = windows.NewLazySystemDLL("kernel32.dll")

	procGetCommPorts     = modKernelBase.NewProc("GetCommPorts")
	procClearCommError   = modkernel32.NewProc("ClearCommError")
	procCloseHandle      = modkernel32.NewProc("CloseHandle")
	procCreateFileA      = modkernel32.NewProc("CreateFileA")
	procFlushFileBuffers = modkernel32.NewProc("FlushFileBuffers")
	procGetCommState     = modkernel32.NewProc("GetCommState")
	procReadFile         = modkernel32.NewProc("ReadFile")
	procSetCommMask      = modkernel32.NewProc("SetCommMask")
	procSetCommState     = modkernel32.NewProc("SetCommState")
	procSetCommTimeouts  = modkernel32.NewProc("SetCommTimeouts")
	procSetupComm        = modkernel32.NewProc("SetupComm")
	procWriteFile        = modkernel32.NewProc("WriteFile")
)

func GetCommPorts(lpPortNumbers *uint32, uPortNumbersCount uint32, puPortNumbersFound *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetCommPorts.Addr(), 3, uintptr(unsafe.Pointer(lpPortNumbers)), uintptr(uPortNumbersCount), uintptr(unsafe.Pointer(puPortNumbersFound)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func ClearCommError(handle syscall.Handle, errors *uint32, stat *COMMSTAT) (err error) {
	r1, _, e1 := syscall.Syscall(procClearCommError.Addr(), 3, uintptr(handle), uintptr(unsafe.Pointer(errors)), uintptr(unsafe.Pointer(stat)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func CloseHandle(handle syscall.Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procCloseHandle.Addr(), 1, uintptr(handle), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func CreateFile(filename *byte, access uint32, mode uint32, sa *syscall.SecurityAttributes, createmode uint32, flags uint32, templatefile syscall.Handle) (handle syscall.Handle, err error) {
	r0, _, e1 := syscall.Syscall9(procCreateFileA.Addr(), 7, uintptr(unsafe.Pointer(filename)), uintptr(access), uintptr(mode), uintptr(unsafe.Pointer(sa)), uintptr(createmode), uintptr(flags), uintptr(templatefile), 0, 0)
	handle = syscall.Handle(r0)
	if handle == 0 {
		err = errnoErr(e1)
	}
	return
}

func FlushFileBuffers(handle syscall.Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procFlushFileBuffers.Addr(), 1, uintptr(handle), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func GetCommState(handle syscall.Handle, dcb *DCB) (err error) {
	r1, _, e1 := syscall.Syscall(procGetCommState.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(dcb)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func ReadFile(handle syscall.Handle, buf *byte, n uint32, read *uint32, overlapped *syscall.Overlapped) (err error) {
	r1, _, e1 := syscall.Syscall6(procReadFile.Addr(), 5, uintptr(handle), uintptr(unsafe.Pointer(buf)), uintptr(n), uintptr(unsafe.Pointer(read)), uintptr(unsafe.Pointer(overlapped)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetCommMask(handle syscall.Handle, mask uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procSetCommMask.Addr(), 2, uintptr(handle), uintptr(mask), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetCommState(handle syscall.Handle, dcb *DCB) (err error) {
	r1, _, e1 := syscall.Syscall(procSetCommState.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(dcb)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetCommTimeouts(handle syscall.Handle, timeouts *CommTimeouts) (err error) {
	r1, _, e1 := syscall.Syscall(procSetCommTimeouts.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(timeouts)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetupComm(handle syscall.Handle, inQueue uint32, outQueue uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procSetupComm.Addr(), 3, uintptr(handle), uintptr(inQueue), uintptr(outQueue))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func WriteFile(handle syscall.Handle, buf *byte, n uint32, written *uint32, overlapped *syscall.Overlapped) (err error) {
	r1, _, e1 := syscall.Syscall6(procWriteFile.Addr(), 5, uintptr(handle), uintptr(unsafe.Pointer(buf)), uintptr(n), uintptr(unsafe.Pointer(written)), uintptr(unsafe.Pointer(overlapped)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
