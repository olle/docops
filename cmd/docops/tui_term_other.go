//go:build !darwin && !linux

package main

import "errors"

type terminalState struct{}

func makeRaw(uintptr) (terminalState, error) {
	return terminalState{}, errors.New("raw terminal mode is not supported on this platform")
}

func restoreTerminal(uintptr, terminalState) error {
	return nil
}
