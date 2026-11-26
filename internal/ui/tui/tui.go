package tui

import (
	"context"
	"io"
	"log"
	"os"

	"charm.land/bubbletea/v2"

	"github.com/olle/docops/internal/app"
	"github.com/olle/docops/internal/ui/text"
)

type Options struct {
	In               io.Reader
	Out              io.Writer
	Width            int
	Height           int
	ForceFallback    bool
	ForceInteractive bool
}

type screen struct {
	source app.SnapshotSource
	opts   Options
}

func New(source app.SnapshotSource, opts Options) app.Screen {
	if opts.In == nil {
		opts.In = os.Stdin
	}
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	return &screen{source: source, opts: opts}
}

func (s *screen) Run(ctx context.Context) error {
	snapshot, err := s.source.Snapshot(ctx)
	if err != nil {
		log.Printf("reading snapshot: %v", err)
	}
	if s.opts.ForceFallback || (!s.opts.ForceInteractive && !isTerminal(s.opts.In, s.opts.Out)) {
		return text.Fallback(s.opts.Out, snapshot)
	}

	programOpts := []tea.ProgramOption{
		tea.WithInput(s.opts.In),
		tea.WithOutput(s.opts.Out),
		tea.WithContext(ctx),
	}
	if s.opts.Width > 0 && s.opts.Height > 0 {
		programOpts = append(programOpts, tea.WithWindowSize(s.opts.Width, s.opts.Height))
	}
	_, err = tea.NewProgram(newModel(ctx, s.source, snapshot, err), programOpts...).Run()
	return err
}

func isTerminal(in io.Reader, out io.Writer) bool {
	stdin, ok := in.(*os.File)
	if !ok {
		return false
	}
	stdout, ok := out.(*os.File)
	if !ok {
		return false
	}
	inInfo, err := stdin.Stat()
	if err != nil {
		return false
	}
	outInfo, err := stdout.Stat()
	if err != nil {
		return false
	}
	return inInfo.Mode()&os.ModeCharDevice != 0 && outInfo.Mode()&os.ModeCharDevice != 0
}
