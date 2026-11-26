package app

import (
	"context"
	"io"

	"github.com/olle/docops/internal/view"
)

const DefaultName = "World"

type SnapshotSource interface {
	Snapshot(ctx context.Context) (view.Snapshot, error)
}

type Screen interface {
	Run(ctx context.Context) error
}

type Presenter interface {
	Render(w io.Writer, snapshot view.Snapshot) error
}

type WebPresenter interface {
	Presenter
	ContentType() string
}
