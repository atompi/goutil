package bar

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestGetPercent(t *testing.T) {
	tests := []struct {
		name     string
		current  int64
		total    int64
		expected int
	}{
		{"50%", 50, 100, 50},
		{"0%", 0, 100, 0},
		{"100%", 100, 100, 100},
		{"zero total", 50, 0, 0},
		{"over 100%", 150, 100, 150},
		{"fraction", 1, 3, 33},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bar{current: tt.current, total: tt.total}
			if got := b.getPercent(); got != tt.expected {
				t.Errorf("getPercent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetTime(t *testing.T) {
	now := time.Now().Add(-70 * time.Minute)
	b := &Bar{start: now, current: 50, total: 100}

	elapsed, remaining := b.getTime()

	if elapsed == "" {
		t.Error("elapsed time should not be empty")
	}
	if remaining == "" {
		t.Error("remaining time should not be empty when current > 0")
	}

	b.current = 0
	_, remaining = b.getTime()
	if remaining != "" {
		t.Errorf("remaining should be empty when current is 0, got %s", remaining)
	}
}

func TestNewBar(t *testing.T) {
	b := NewBar(20, 100)

	if b.current != 20 {
		t.Errorf("current = %d, want 20", b.current)
	}
	if b.total != 100 {
		t.Errorf("total = %d, want 100", b.total)
	}
	if b.graph != "█" {
		t.Errorf("graph = %s, want █", b.graph)
	}
	if b.percent != 20 {
		t.Errorf("percent = %d, want 20", b.percent)
	}
}

func TestNewBarWithGraph(t *testing.T) {
	b := NewBarWithGraph(20, 100, "#")

	if b.graph != "#" {
		t.Errorf("graph = %s, want #", b.graph)
	}
	if b.current != 20 {
		t.Errorf("current = %d, want 20", b.current)
	}
}

func TestBarAdd(t *testing.T) {
	var buf bytes.Buffer
	b := newBar(0, 100, "█", &buf)

	b.Add(30)

	if b.current != 30 {
		t.Errorf("current = %d, want 30", b.current)
	}
	if b.percent != 30 {
		t.Errorf("percent = %d, want 30", b.percent)
	}
}

func TestBarReset(t *testing.T) {
	var buf bytes.Buffer
	b := newBar(0, 100, "█", &buf)

	b.Reset(50)

	if b.current != 50 {
		t.Errorf("current = %d, want 50", b.current)
	}
	if b.percent != 50 {
		t.Errorf("percent = %d, want 50", b.percent)
	}
}

func TestConcurrentAccess(t *testing.T) {
	b := NewBar(0, 1000)
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10 {
				b.Add(1)
			}
		}()
	}

	wg.Wait()

	if b.current != 100 {
		t.Errorf("current = %d, want 100", b.current)
	}
}

func TestBarRender(t *testing.T) {
	var buf bytes.Buffer
	b := newBar(10, 100, "█", &buf)

	b.render()

	output := buf.String()
	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Error("output should contain progress bar brackets")
	}
	if !strings.Contains(output, "%") {
		t.Error("output should contain percentage")
	}
}

func TestBarEdgeCases(t *testing.T) {
	b := NewBar(0, 0)
	if b.percent != 0 {
		t.Errorf("zero total should give 0 percent, got %d", b.percent)
	}

	b = NewBar(150, 100)
	if b.percent != 150 {
		t.Errorf("over 100 should give percent > 100, got %d", b.percent)
	}
}
