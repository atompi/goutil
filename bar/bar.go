package bar

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Bar struct {
	mu      sync.Mutex
	graph   string
	rate    string
	percent int
	current int64
	total   int64
	start   time.Time
	output  io.Writer
}

func (b *Bar) getPercent() int {
	if b.total == 0 {
		return 0
	}
	return int(float64(b.current) / float64(b.total) * 100)
}

func (b *Bar) getTime() (elapsed, remaining string) {
	elapsedSec := time.Since(b.start).Seconds()
	h := int(elapsedSec) / 3600
	m := (int(elapsedSec) % 3600) / 60
	s := int(elapsedSec) % 60

	if h > 0 {
		elapsed = fmt.Sprintf("%dh %dm %ds", h, m, s)
	} else if m > 0 {
		elapsed = fmt.Sprintf("%dm %ds", m, s)
	} else {
		elapsed = fmt.Sprintf("%ds", s)
	}

	if b.current > 0 {
		remainingSec := (float64(b.total) / float64(b.current)) * elapsedSec
		rh := int(remainingSec) / 3600
		rm := (int(remainingSec) % 3600) / 60
		rs := int(remainingSec) % 60

		if rh > 0 {
			remaining = fmt.Sprintf("%dh %dm %ds", rh, rm, rs)
		} else if rm > 0 {
			remaining = fmt.Sprintf("%dm %ds", rm, rs)
		} else {
			remaining = fmt.Sprintf("%ds", rs)
		}
	}

	return
}

func (b *Bar) render() {
	elapsed, remaining := b.getTime()
	fmt.Fprintf(b.output, "\r[%-50s]% 3d%% %12s/%-12s %d/%d", b.rate, b.percent, elapsed, remaining, b.current, b.total)
}

func (b *Bar) Reset(current int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current = current
	b.percent = b.getPercent()
	b.updateRate()
	b.render()
}

func (b *Bar) Add(delta int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current += delta
	b.percent = b.getPercent()
	if b.percent%2 == 0 {
		b.rate += b.graph
	}
	b.render()
}

func (b *Bar) updateRate() {
	for i := 0; i < b.percent; i += 2 {
		b.rate += b.graph
	}
}

func NewBar(current, total int64) *Bar {
	return newBar(current, total, "█", os.Stderr)
}

func NewBarWithGraph(current, total int64, graph string) *Bar {
	return newBar(current, total, graph, os.Stderr)
}

func newBar(current, total int64, graph string, output io.Writer) *Bar {
	b := &Bar{
		current: current,
		total:   total,
		graph:   graph,
		output:  output,
		start:   time.Now(),
	}
	b.percent = b.getPercent()
	b.updateRate()
	return b
}
