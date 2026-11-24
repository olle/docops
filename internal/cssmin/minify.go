// Package cssmin provides a tiny, dependency-free CSS minifier.
//
// It strips comments, collapses whitespace, tightens punctuation, drops
// redundant trailing semicolons, and shortens leading zeros in numbers.
// String contents are preserved byte-for-byte.
package cssmin

// Minify returns a minified version of css.
func Minify(css string) string {
	return string(tighten(stripCommentsAndSpaces([]byte(css))))
}

func stripCommentsAndSpaces(src []byte) []byte {
	out := make([]byte, 0, len(src))
	var quote byte
	var pendingSpace bool

	for i := 0; i < len(src); i++ {
		c := src[i]

		if quote != 0 {
			out = append(out, c)
			if c == '\\' && i+1 < len(src) {
				i++
				out = append(out, src[i])
			} else if c == quote {
				quote = 0
			}
			continue
		}

		if c == '/' && i+1 < len(src) && src[i+1] == '*' {
			i += 2
			for i+1 < len(src) && !(src[i] == '*' && src[i+1] == '/') {
				i++
			}
			i++ // closing slash
			continue
		}

		if c == '"' || c == '\'' {
			flushSpace(&out, pendingSpace)
			pendingSpace = false
			quote = c
			out = append(out, c)
			continue
		}

		if isSpace(c) {
			pendingSpace = true
			continue
		}

		flushSpace(&out, pendingSpace)
		pendingSpace = false
		out = append(out, c)
	}
	return out
}

func tighten(src []byte) []byte {
	out := make([]byte, 0, len(src))
	var quote byte
	var pendingSpace, skipSpace bool

	for i := 0; i < len(src); i++ {
		c := src[i]

		if quote != 0 {
			out = append(out, c)
			if c == '\\' && i+1 < len(src) {
				i++
				out = append(out, src[i])
			} else if c == quote {
				quote = 0
			}
			continue
		}

		if c == '"' || c == '\'' {
			flushSpace(&out, pendingSpace)
			pendingSpace, skipSpace = false, false
			quote = c
			out = append(out, c)
			continue
		}

		if isSpace(c) {
			if !skipSpace {
				pendingSpace = true
			}
			continue
		}

		// A redundant semicolon just before a closing brace.
		if c == '}' && len(out) > 0 && out[len(out)-1] == ';' {
			out = out[:len(out)-1]
		}

		// Drop the leading zero in numbers like 0.5rem or 0.07.
		if c == '0' && i+2 < len(src) && src[i+1] == '.' && isDigit(src[i+2]) &&
			(len(out) == 0 || !isIdentByte(out[len(out)-1])) {
			continue
		}

		// Shorten six-digit colors whose channels repeat (e.g. #ffffff).
		if c == '#' && valueContext(out) {
			if replacement, consumed, ok := compactHex(src[i:]); ok {
				out = append(out, replacement...)
				i += consumed - 1
				pendingSpace, skipSpace = false, false
				continue
			}
		}

		switch c {
		case '(':
			// Keep a preceding space (e.g. `@media (`), but drop what follows.
			flushSpace(&out, pendingSpace)
			pendingSpace, skipSpace = false, true
			out = append(out, c)
			continue
		case ':':
			// Keep any preceding space: `form :where(...)` is a descendant
			// selector, not a pseudo-class attached to `form`.
			flushSpace(&out, pendingSpace)
			pendingSpace, skipSpace = false, true
			out = append(out, c)
			continue
		case ')':
			// Drop space before the parenthesis, but keep space after it:
			// `var(--x) 1fr` needs that separator.
			pendingSpace, skipSpace = false, false
			out = append(out, c)
			continue
		case '{', '}', ';', ',', '>', '+', '~', '/':
			pendingSpace = false
			out = append(out, c)
			skipSpace = true
			continue
		}

		flushSpace(&out, pendingSpace)
		pendingSpace, skipSpace = false, false
		out = append(out, c)
	}
	return out
}

func flushSpace(out *[]byte, pending bool) {
	if !pending || len(*out) == 0 {
		return
	}
	last := (*out)[len(*out)-1]
	if isSpace(last) || last == '(' {
		return
	}
	*out = append(*out, ' ')
}

func compactHex(s []byte) ([]byte, int, bool) {
	run := 0
	for run+1 < len(s) && isHex(s[run+1]) {
		run++
	}
	if run != 6 && run != 8 {
		return nil, 0, false
	}
	if s[1] != s[3] || s[2] != s[4] || s[3] != s[5] {
		return nil, 0, false
	}
	if run == 8 && (s[7] != 'f' || s[8] != 'f') {
		return nil, 0, false
	}
	return []byte{'#', s[1], s[2], s[3]}, run + 1, true
}

func valueContext(out []byte) bool {
	if len(out) == 0 {
		return false
	}
	switch out[len(out)-1] {
	case ':', ',', '(':
		return true
	}
	return false
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isHex(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isIdentByte(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-' || c == '_' || c == '.'
}
