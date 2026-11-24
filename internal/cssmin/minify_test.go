package cssmin

import "testing"

func TestMinify(t *testing.T) {
	input := `
		/* A comment that should disappear. */
		body > main { color: #ffffff; margin: 0.5em 0; }
		@media (prefers-color-scheme: dark) {
			:root { --x: rgb(0 0 0 / 0.07); }
		}
	`
	want := `body>main{color:#fff;margin:.5em 0}@media (prefers-color-scheme:dark){:root{--x:rgb(0 0 0/.07)}}`
	if got := Minify(input); got != want {
		t.Errorf("Minify() = %q, want %q", got, want)
	}
}

func TestMinifyPreservesStrings(t *testing.T) {
	input := `a::before{content:"a  b; } /* x */";font-family:'A, B'"}`
	want := `a::before{content:"a  b; } /* x */";font-family:'A, B'"}`
	if got := Minify(input); got != want {
		t.Errorf("Minify() = %q, want %q", got, want)
	}
}

func TestMinifyKeepsRequiredSpaces(t *testing.T) {
	input := `form :where(p, small) { color: color-mix(in srgb, var(--accent) 76%, var(--text-muted)); gap: var(--a) var(--b) }`
	want := `form :where(p,small){color:color-mix(in srgb,var(--accent) 76%,var(--text-muted));gap:var(--a) var(--b)}`
	if got := Minify(input); got != want {
		t.Errorf("Minify() = %q, want %q", got, want)
	}
}

func TestMinifyKeepsIDSelectorsAndWholeNumbers(t *testing.T) {
	input := `#ffffff { --x: 10.5em; }`
	want := `#ffffff{--x:10.5em}`
	if got := Minify(input); got != want {
		t.Errorf("Minify() = %q, want %q", got, want)
	}
}
