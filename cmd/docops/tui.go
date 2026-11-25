package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type terminal interface {
	Raw() (restore func(), ok bool)
}

type stdioTerminal struct{}

func (stdioTerminal) Raw() (func(), bool) {
	fd := os.Stdin.Fd()
	state, err := makeRaw(fd)
	if err != nil {
		return func() {}, false
	}
	return func() { _ = restoreTerminal(fd, state) }, true
}

type tuiScreen struct {
	out  io.Writer
	in   io.Reader
	term terminal
}

func (s tuiScreen) Present(msg message) error {
	restore, styled := s.term.Raw()
	defer func() {
		restore()
		if styled {
			_, _ = io.WriteString(s.out, "\x1b[?25h\x1b[?1049l")
		}
	}()

	if err := renderTUI(s.out, msg, styled); err != nil {
		return err
	}
	return waitForQuit(s.in, styled)
}

func renderTUI(w io.Writer, msg message, styled bool) error {
	if styled {
		if _, err := io.WriteString(w, "\x1b[?1049h\x1b[?25l\x1b[2J\x1b[H"); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "\x1b[1m%s\x1b[0m\n\n", msg.greeting); err != nil {
			return err
		}
	} else if _, err := fmt.Fprintf(w, "%s\n\n", msg.greeting); err != nil {
		return err
	}

	if len(msg.specifications) > 0 {
		if _, err := io.WriteString(w, "Specifications:\n"); err != nil {
			return err
		}
		for _, spec := range msg.specifications {
			entry := spec.kind.Label() + ": " + spec.title
			if styled {
				if _, err := fmt.Fprintf(w, "  \x1b[36m*\x1b[0m %s\n", entry); err != nil {
					return err
				}
				continue
			}
			if _, err := fmt.Fprintf(w, "  * %s\n", entry); err != nil {
				return err
			}
		}
	}

	if styled {
		_, err := io.WriteString(w, "\n\x1b[2mq: quit\x1b[0m")
		return err
	}
	_, err := io.WriteString(w, "\nq: quit\n")
	return err
}

func waitForQuit(in io.Reader, raw bool) error {
	reader := bufio.NewReader(in)
	if raw {
		for {
			b, err := reader.ReadByte()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return err
			}
			if b == 'q' || b == 'Q' || b == 3 || b == 27 {
				return nil
			}
		}
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if strings.ContainsAny(line, "qQ") || errors.Is(err, io.EOF) {
			return nil
		}
	}
}
