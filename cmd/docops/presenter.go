package main

import "io"

type presenter interface {
	Render(w io.Writer, msg message) error
}

type webPresenter interface {
	presenter
	ContentType() string
}

type screenPresenter interface {
	Present(msg message) error
}
