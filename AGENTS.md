# 项目知识库

**生成时间：** 2026-05-19
**提交：** cdc9187
**分支：** main

## 概述

Go 工具库 (github.com/atompi/goutil)，提供进度条（`bar`）和日志（`log`）工具。代码注释使用中文。

## 项目结构

```
goutil/
├── bar/           # 进度条工具
│   ├── bar.go
│   └── bar_test.go
├── log/           # 日志工具（slog + zap 双后端支持）
│   ├── common.go  # 文件写入器、路径验证
│   ├── slog.go    # slog.Handler 实现
│   ├── zap.go     # zap.Logger 封装
│   └── *_test.go  # 测试
└── README.md
```

## 速查表

| 需求 | 文件位置 | 关键函数 |
|------|----------|----------|
| 进度条 | `bar/bar.go` | `NewBar()`, `NewBarWithGraph()` |
| 日志配置 | `log/common.go` | `NewLoggerOptions()`, `WithLevel()`, `WithFormat()`, `WithPath()`, `WithMultiFiles()` |
| slog 后端 | `log/slog.go` | `NewSlogLogger()` |
| zap 后端 | `log/zap.go` | `NewZapLogger()` |
| 文件写入器 | `log/common.go` | `NewLogFile()` |

## 代码符号表

| 符号 | 类型 | 位置 | 说明 |
|------|------|------|------|
| Bar | 结构体 | bar/bar.go:11 | 进度条状态 |
| NewBar | 函数 | bar/bar.go:92 | 构造函数 |
| Logger | 结构体 | log/common.go:23 | 配置容器 |
| Options | 类型 | log/common.go:47 | 函数式选项 |
| NewSlogLogger | 函数 | log/slog.go:68 | slog 日志后端 |
| NewZapLogger | 函数 | log/zap.go:11 | zap 日志后端 |
| NewLogFile | 函数 | log/common.go:86 | 文件写入器工厂 |

## 约定规范

- **多模块结构**: 每个子目录是独立的 Go 模块，有自己的 `go.mod`
- **函数式选项**: 日志配置使用函数式选项模式（`WithLevel`、`WithFormat` 等）
- **依赖注入**: `MkdirAll`、`OpenFile`、`PathSeparator` 在 `log/common.go` 中定义为包变量，便于测试

## 本项目反模式

- 无 CI/CD 配置（无 workflows、无 Makefile）
- 无 linter 配置（无 `.golangci.yaml`）
- 无 golangci-lint 强制执行

## 独特风格

- 代码注释使用中文
- 进度条默认使用 `█` 作为图形字符
- 日志路径验证拒绝跨平台不兼容字符

## 常用命令

```bash
# 测试各模块
cd bar && go test ./...
cd log && go test ./...

# 构建各模块
cd bar && go build ./...
cd log && go build ./...
```

## 注意事项

- `bar/` 模块：无外部依赖
- `log/` 模块：依赖 `testify` (v1.8.1) 和 `zap` (v1.27.1)
- 每个模块需要单独执行 `go mod tidy` / `go get`
