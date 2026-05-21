# log/ AGENTS.md

**模块：** github.com/atompi/goutil/log

## 概述

日志工具库，支持双后端：`slog`（标准库）和 `zap`（uber）。支持文件输出、JSON/文本格式、级别过滤。

## 项目结构

```
log/
├── common.go    # Logger 配置 + 文件写入器工厂 + 路径验证
├── slog.go      # slog.Handler 实现
└── zap.go       # zap.Logger 封装
```

## 速查表

| 需求 | 文件位置 | 说明 |
|------|----------|------|
| 日志配置 | `log/common.go` | `NewLoggerOptions()`, `WithLevel()`, `WithFormat()`, `WithPath()`, `WithMultiFiles()` |
| 文件写入器 | `log/common.go` | `NewLogFile()` - 验证路径、创建目录 |
| slog 后端 | `log/slog.go` | `NewSlogLogger()` - 自定义 handler |
| zap 后端 | `log/zap.go` | `NewZapLogger()` - tee core |

## 约定规范 (log 模块)

- **函数式选项**: 通过 `Options` 函数类型配置
- **依赖注入**: `MkdirAll`、`OpenFile`、`PathSeparator` 是包变量（可测试）
- **路径验证**: Unix 拒绝 `[]!"#$%&'()*+,/:;<=>?@^{|}~`，Windows 额外拒绝 `,`
- **无自动日志轮转**: 文件无限增长

## 符号表

| 符号 | 类型 | 位置 | 说明 |
|------|------|------|------|
| Logger | 结构体 | common.go:23 | 配置容器 |
| Options | 类型 | common.go:47 | 函数式选项类型 |
| NewSlogLogger | 函数 | slog.go:68 | 构造函数 |
| NewZapLogger | 函数 | zap.go:11 | 构造函数 |
| NewLogFile | 函数 | common.go:86 | 文件写入器工厂 |
