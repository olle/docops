//go:build darwin

package main

import (
	"syscall"
	"unsafe"
)

type terminalState syscall.Termios

func makeRaw(fd uintptr) (terminalState, error) {
	var old syscall.Termios
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&old))); errno != 0 {
		return terminalState{}, errno
	}

	raw := old
	raw.Iflag &^= syscall.IXON | syscall.ICRNL | syscall.INLCR
	raw.Oflag &^= syscall.OPOST
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&raw))); errno != 0 {
		return terminalState{}, errno
	}
	return terminalState(old), nil
}

func restoreTerminal(fd uintptr, state terminalState) error {
	old := syscall.Termios(state)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&old)))
	if errno != 0 {
		return errno
	}
	return nil
}
