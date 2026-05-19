# bar/ AGENTS.md

**模块：** github.com/atompi/goutil/bar

## 概述

轻量级进度条工具，无外部依赖。

## 项目结构

```
bar/
├── bar.go
└── bar_test.go
```

## 速查表

| 需求 | 文件位置 | 说明 |
|------|----------|------|
| 进度条 | `bar/bar.go` | `NewBar()`, `NewBarWithGraph()` |

## 约定规范 (bar 模块)

- **默认图形**: `█`（实心方块）
- **线程安全**: 使用 `sync.Mutex` 保护 `Add()` 和 `Reset()`
- **输出**: 使用回车符 (`\r`) 写入 `os.Stderr`

## 符号表

| 符号 | 类型 | 位置 | 说明 |
|------|------|------|------|
| Bar | 结构体 | bar.go:11 | 进度条状态 |
| NewBar | 函数 | bar.go:77 | 构造函数 |
| NewBarWithGraph | 函数 | bar.go:92 | 支持自定义图形的构造函数 |
