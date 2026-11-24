// Package assets holds static files that are compiled into the docops binary.
package assets

import _ "embed"

// CSS is the stylesheet inlined into generated HTML pages.
//
//go:embed style.css
var CSS string
