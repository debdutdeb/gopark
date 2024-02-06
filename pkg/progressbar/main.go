package progressbar

import (
	"fmt"
	"io"
	"math"

	"golang.org/x/sys/unix"
)

const extraCharacters = "[] [] xxx.xx%"

type progressBar struct {
	order int

	writer   io.Writer
	label    string
	maxBytes int64
	perc     float64

	barWidth int

	//state
	bytesWritten   int64
	columnsWritten int
	bar            []byte
}

type ProgressBarOptions struct {
	BarWidth float64
}

func (p *progressBar) Write(data []byte) (int, error) {
	n, err := p.writer.Write(data)
	if err != nil {
		return 0, err
	}

	if n == 0 {
		return 0, nil
	}

	fmt.Printf("\r[%s] [%s] %.2f%%", p.label, p.buildprogressbar(n), p.perc)

	return n, err
}

func (p *progressBar) buildprogressbar(bytes int) []byte {
	p.bytesWritten += int64(bytes)

	p.perc = float64(p.bytesWritten) / float64(p.maxBytes) * 100

	fill := int(math.Floor(float64(p.barWidth) * p.perc / 100.00))

	for i := p.columnsWritten; i < fill; i += 1 {
		p.bar[i] = '='
	}

	p.columnsWritten = fill

	return p.bar
}

func NewWriteProgressBar(label string, maxBytes int64, writer io.Writer, opts *ProgressBarOptions) (io.Writer, error) {
	size, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err != nil {
		return nil, err
	}

	p := &progressBar{writer: writer, label: label, maxBytes: maxBytes}

	columns := int(size.Col)

	if opts != nil {
		columns = int(math.Floor(float64(columns) * float64(opts.BarWidth) / 100.0))
	}

	p.barWidth = int(columns) - len(label) - len(extraCharacters)

	p.bar = make([]byte, p.barWidth)
	for i := range p.bar {
		p.bar[i] = ' '
	}

	return p, nil
}
