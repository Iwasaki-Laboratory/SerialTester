package main

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output winapi_generate.go winapi.go

//sys CreateFile(filename *byte, access uint32, mode uint32, sa *syscall.SecurityAttributes, createmode uint32, flags uint32, templatefile syscall.Handle) (handle syscall.Handle, err error) = kernel32.CreateFileA
//sys GetCommState(handle syscall.Handle, dcb *DCB) (err error) = kernel32.GetCommState
//sys SetCommState(handle syscall.Handle, dcb *DCB) (err error) = kernel32.SetCommState
//sys WriteFile(handle syscall.Handle, buf *byte, n uint32, written *uint32, overlapped *syscall.Overlapped) (err error) = kernel32.WriteFile
//sys ReadFile(handle syscall.Handle, buf *byte, n uint32, read *uint32, overlapped *syscall.Overlapped) (err error) = kernel32.ReadFile
//sys CloseHandle(handle syscall.Handle) (err error) = kernel32.CloseHandle
//sys FlushFileBuffers(handle syscall.Handle) (err error) = kernel32.FlushFileBuffers
//sys SetCommTimeouts(handle syscall.Handle, timeouts *CommTimeouts) (err error) = kernel32.SetCommTimeouts
//sys SetCommMask(handle syscall.Handle, mask uint32) (err error) = kernel32.SetCommMask
//sys SetupComm(handle syscall.Handle, inQueue	uint32, outQueue uint32) (err error) = kernel32.SetupComm
//sys ClearCommError(handle syscall.Handle, errors *uint32, stat *COMMSTAT) (err error) = kernel32.ClearCommError
//sys GetCommPorts(lpPortNumbers *uint32, uPortNumbersCount uint32, puPortNumbersFound *uint32) (err error) = KernelBase.GetCommPorts

type DCB struct {
	DCBlength  uint32
	BaudRate   uint32
	Flags      uint32 // ビットフィールドを組み合わせる
	WReserved  uint16
	XonLim     uint16
	XoffLim    uint16
	ByteSize   byte
	Parity     byte
	StopBits   byte
	XonChar    byte
	XoffChar   byte
	ErrorChar  byte
	EofChar    byte
	EvtChar    byte
	WReserved1 uint16
}

const (
	NOPARITY   = 0
	EVENPARITY = 2
	ODDPARITY  = 1
)

const (
	ONESTOPBITS = 0
	TWOSTOPBITS = 2
)

// Flagsフィールド内のビットを表す定数
const (
	fBinary uint32 = 1 << iota
	fParity
	fOutxCtsFlow
	fOutxDsrFlow
	fDtrControl
	fDtrControlHigh = 3 << (iota - 1)
	fDsrSensitivity
	fTXContinueOnXoff
	fOutX
	fInX
	fErrorChar
	fNull
	fRtsControl
	fRtsControlHigh = 3 << (iota - 1)
	fAbortOnError
	fDummy2
)

type COMMSTAT struct {
	Flags    uint32 // ビットフィールドを組み合わせる
	CbInQue  uint32
	CbOutQue uint32
}

// Flags フィールドのビットを表す定数
const (
	fCtsHold uint32 = 1 << iota
	fDsrHold
	fRlsdHold
	fXoffHold
	fXoffSent
	fEof
	fTxim
)

type CommTimeouts struct {
	ReadIntervalTimeout         uint32
	ReadTotalTimeoutMultiplier  uint32
	ReadTotalTimeoutConstant    uint32
	WriteTotalTimeoutMultiplier uint32
	WriteTotalTimeoutConstant   uint32
}

const (
	GENERIC_READ          = 0x80000000
	GENERIC_WRITE         = 0x40000000
	OPEN_EXISTING         = 3
	NULL                  = 0
	FILE_ATTRIBUTE_NORMAL = 0x80
)
