// Package assets holds static files that are compiled into the docops binary.
package assets

import _ "embed"

//go:generate go run ../../cmd/minifycss -in style.css -out style.min.css

// CSS is the minified stylesheet inlined into generated HTML pages.
// It is generated from style.css; run `go generate ./...` after editing it.
//
//go:embed style.min.css
var CSS string
