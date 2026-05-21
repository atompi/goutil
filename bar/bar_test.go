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
	if !strings.Contains(output, "[") {
		t.Error("output should contain opening bracket")
	}
	if !strings.Contains(output, "]") {
		t.Error("output should contain closing bracket")
	}
	if !strings.Contains(output, "%") {
		t.Error("output should contain percentage")
	}
	if !strings.Contains(output, "10/100") {
		t.Errorf("output should contain current/total, got %s", output)
	}

	buf.Reset()
	b.current = 50
	b.percent = 50
	b.rate = strings.Repeat("█", 25)
	b.render()
	output = buf.String()
	if !strings.Contains(output, "50%") {
		t.Errorf("output should contain 50%%, got %s", output)
	}

	if strings.Contains(output, "100%") {
		t.Error("output should not contain 100% when percent is 50")
	}

	bracketEnd := strings.Index(output, "]")
	bracketStart := strings.Index(output, "[")
	if bracketStart >= 0 && bracketEnd > bracketStart {
		rateContent := output[bracketStart+1 : bracketEnd]
		runeCount := []rune(rateContent)
		if len(runeCount) != 50 {
			t.Errorf("rate field should be 50 characters, got %d", len(runeCount))
		}
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

func TestBarUpdateRate(t *testing.T) {
	tests := []struct {
		name        string
		percent     int
		graph       string
		expectedLen int
	}{
		{"50 percent", 50, "█", 75},
		{"100 percent", 100, "█", 150},
		{"25 percent", 25, "█", 39},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bar{graph: tt.graph, rate: "", percent: tt.percent}
			b.updateRate()
			if len(b.rate) != tt.expectedLen {
				t.Errorf("expected rate length %d, got %d", tt.expectedLen, len(b.rate))
			}
		})
	}
}

func TestBarLargeNumbers(t *testing.T) {
	b := NewBar(0, 1e9)
	if b.percent != 0 {
		t.Errorf("expected 0 percent for 0/1e9, got %d", b.percent)
	}

	b.current = 5e8
	b.percent = b.getPercent()
	if b.percent != 50 {
		t.Errorf("expected 50 percent for 5e8/1e9, got %d", b.percent)
	}
}

func TestBarTimeFormatting(t *testing.T) {
	tests := []struct {
		name     string
		elapsed  time.Duration
		expected string
	}{
		{"seconds only", 30 * time.Second, "30s"},
		{"minutes and seconds", 90 * time.Second, "1m 30s"},
		{"hours minutes seconds", 3661 * time.Second, "1h 1m 1s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bar{
				start:   time.Now().Add(-tt.elapsed),
				current: 50,
				total:   100,
			}
			elapsed, _ := b.getTime()
			if elapsed != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, elapsed)
			}
		})
	}
}

func TestBarAddEdgeCases(t *testing.T) {
	var buf bytes.Buffer
	b := newBar(0, 100, "█", &buf)

	b.Add(0)
	if b.current != 0 {
		t.Errorf("Add(0) should not change current, got %d", b.current)
	}

	b.Add(100)
	if b.current != 100 {
		t.Errorf("expected current 100, got %d", b.current)
	}

	b.Add(-50)
	if b.current != 50 {
		t.Errorf("Add(-50) should result in current 50, got %d", b.current)
	}

	b.Add(-100)
	if b.current != -50 {
		t.Errorf("Add(-100) when current is 50 should result in -50, got %d", b.current)
	}
}

func TestBarResetEdgeCases(t *testing.T) {
	var buf bytes.Buffer
	b := newBar(50, 100, "█", &buf)

	// Reset to same value
	b.Reset(50)
	if b.current != 50 || b.percent != 50 {
		t.Errorf("reset to same value failed: current=%d, percent=%d", b.current, b.percent)
	}

	// Reset to 0
	b.Reset(0)
	if b.current != 0 || b.percent != 0 {
		t.Errorf("reset to 0 failed: current=%d, percent=%d", b.current, b.percent)
	}
}
