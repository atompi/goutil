package bar

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestGetPercent 测试 getPercent 方法
func TestGetPercent(t *testing.T) {
	tests := []struct {
		name     string
		current  int64
		total    int64
		expected int
	}{
		{"正常情况50%", 50, 100, 50},
		{"开始状态0%", 0, 100, 0},
		{"完成状态100%", 100, 100, 100},
		{"除零保护", 50, 0, 0}, // 避免除零错误
		{"超出范围150%", 150, 100, 150},
		{"小数处理", 1, 3, 33}, // 0.333... -> 33%
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := &Bar{
				current: tt.current,
				total:   tt.total,
			}
			result := bar.getPercent()
			if result != tt.expected {
				t.Errorf("getPercent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetTime 测试 getTime 方法
func TestGetTime(t *testing.T) {
	// 为了稳定测试，我们需要控制时间，模拟经过 70m
	now := time.Now().Add(-70 * time.Minute)
	bar := &Bar{
		start:   now,
		current: 50,
		total:   100,
	}
	pasted, left := bar.getTime()

	// 基本检查，实际时间值难以精确预测
	if pasted == "" || left == "" {
		t.Error("getTime() should return non-empty strings")
	}

	// 测试刚开始状态
	barStart := &Bar{
		start:   now,
		current: 0,
		total:   100,
	}
	pastedStart, leftStart := barStart.getTime()
	if pastedStart == "" || leftStart == "" {
		t.Error("getTime() should return non-empty strings for start state")
	}

	// 测试已完成状态
	barFinished := &Bar{
		start:   now,
		current: 100,
		total:   100,
	}
	_, leftFinished := barFinished.getTime()
	if !strings.Contains(leftFinished, "0s") {
		t.Errorf("Finished bar should have 0s left, got: %s", leftFinished)
	}
}

// TestLoad 测试 load 方法的输出格式
func TestLoad(t *testing.T) {
	// 保存原始 stderr
	originalStderr := os.Stderr

	// 创建管道来捕获输出
	r, w, _ := os.Pipe()
	os.Stderr = w

	bar := &Bar{
		graph:   "█",
		current: 10,
		total:   100,
		start:   time.Now(),
	}
	bar.percent = bar.getPercent()

	// 初始化 rate
	for i := 0; i < bar.percent; i += 2 {
		bar.rate += bar.graph
	}

	// 执行 load
	bar.load()

	// 关闭写入端并恢复 stderr
	w.Close()
	os.Stderr = originalStderr

	// 读取捕获的输出
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// 验证输出包含关键元素
	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Error("load() output should contain progress bar brackets")
	}
	if !strings.Contains(output, "%") {
		t.Error("load() output should contain percentage")
	}
}

// TestReset 测试 Reset 方法
func TestReset(t *testing.T) {
	bar := NewBar(0, 100)

	// 保存原始 stderr
	originalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	// 重置进度
	bar.Reset(50)

	// 关闭写入端并恢复 stderr
	w.Close()
	os.Stderr = originalStderr

	// 检查内部状态
	if bar.current != 50 {
		t.Errorf("Reset() failed to set current, got %d, want 50", bar.current)
	}
	if bar.percent != 50 {
		t.Errorf("Reset() failed to update percent, got %d, want 50", bar.percent)
	}
}

// TestAdd 测试 Add 方法
func TestAdd(t *testing.T) {
	bar := NewBar(0, 100)

	// 保存原始 stderr
	originalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	// 增加进度
	bar.Add(30)

	// 关闭写入端并恢复 stderr
	w.Close()
	os.Stderr = originalStderr

	// 检查内部状态
	if bar.current != 30 {
		t.Errorf("Add() failed to increase current, got %d, want 30", bar.current)
	}
	if bar.percent != 30 {
		t.Errorf("Add() failed to update percent, got %d, want 30", bar.percent)
	}

	// 继续增加
	bar.Add(20)
	if bar.current != 50 {
		t.Errorf("Add() failed on second call, got %d, want 50", bar.current)
	}
}

// TestNewBar 测试 NewBar 函数
func TestNewBar(t *testing.T) {
	bar := NewBar(20, 100)

	if bar.current != 20 {
		t.Errorf("NewBar() failed to set current, got %d, want 20", bar.current)
	}
	if bar.total != 100 {
		t.Errorf("NewBar() failed to set total, got %d, want 100", bar.total)
	}
	if bar.graph != "█" {
		t.Errorf("NewBar() failed to set default graph, got %s, want █", bar.graph)
	}
	if bar.percent != 20 {
		t.Errorf("NewBar() failed to calculate initial percent, got %d, want 20", bar.percent)
	}

	// 验证 rate 被正确初始化（每2%一个字符）
	expectedRateLength := 20 / 2 * 3
	if len(bar.rate) != expectedRateLength {
		t.Errorf("NewBar() failed to initialize rate, got length %d, want %d", len(bar.rate), expectedRateLength)
	}
}

// TestNewBarWithGraph 测试 NewBarWithGraph 函数
func TestNewBarWithGraph(t *testing.T) {
	bar := NewBarWithGraph(20, 100, "#")

	if bar.graph != "#" {
		t.Errorf("NewBarWithGraph() failed to set custom graph, got %s, want #", bar.graph)
	}

	// 确保其他属性正常初始化
	if bar.current != 20 {
		t.Errorf("NewBarWithGraph() failed to set current, got %d, want 20", bar.current)
	}
	if bar.total != 100 {
		t.Errorf("NewBarWithGraph() failed to set total, got %d, want 100", bar.total)
	}
}

// TestConcurrentAccess 测试并发访问安全性
func TestConcurrentAccess(t *testing.T) {
	bar := NewBar(0, 1000)
	var wg sync.WaitGroup

	// 启动多个 goroutine 并发调用 Add
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				bar.Add(1)
			}
		}()
	}

	wg.Wait()

	// 最终值应该是 100
	if bar.current != 100 {
		t.Errorf("Concurrent Add() calls resulted in incorrect value, got %d, want 100", bar.current)
	}
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	// 测试 total 为 0 的情况
	bar := NewBar(0, 0)
	if bar.percent != 0 {
		t.Errorf("Bar with zero total should have 0 percent, got %d", bar.percent)
	}

	// 测试 current 大于 total 的情况
	bar = NewBar(150, 100)
	if bar.percent != 150 {
		t.Errorf("Bar with current > total should have percent > 100, got %d", bar.percent)
	}

	// 测试负数输入
	bar = NewBar(-10, 100)
	if bar.percent != -10 {
		t.Errorf("Bar with negative current should have negative percent, got %d", bar.percent)
	}
}
