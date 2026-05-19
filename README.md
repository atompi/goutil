# goutil - 常用 Go 工具包

Go 工具库，提供进度条（bar）和日志（log）工具。

## 安装

```bash
# 安装 bar 模块
go get github.com/atompi/goutil/bar

# 安装 log 模块
go get github.com/atompi/goutil/log
```

## bar - 进度条

无外部依赖的轻量级进度条工具。

### 快速开始

```go
package main

import (
    "github.com/atompi/goutil/bar"
)

func main() {
    b := bar.NewBar(0, 100)
    for i := int64(0); i <= 100; i++ {
        b.Add(1)
    }
}
```

### 自定义图形字符

```go
b := bar.NewBarWithGraph(0, 100, "▓")
```

### 主要函数

| 函数 | 说明 |
|------|------|
| `NewBar(current, total int64)` | 创建进度条 |
| `NewBarWithGraph(current, total int64, graph string)` | 创建自定义图形的进度条 |
| `Add(i int64)` | 增加进度 |
| `Reset(current int64)` | 重置进度 |

## log - 日志

支持 `slog`（标准库）和 `zap`（uber）双后端的日志工具。

### 快速开始

```go
package main

import (
    "github.com/atompi/goutil/log"
)

func main() {
    opts := log.NewLoggerOptions(
        log.WithLevel("info"),
        log.WithFormat("console"),
        log.WithPath("app.log"),
    )

    slogger := log.NewSlogLogger(opts)
    slogger.Info("hello world")
}
```

### 使用 zap 后端

```go
zlogger := log.NewZapLogger(opts)
zlogger.Info("hello world")
```

### 多文件输出

```go
opts := log.NewLoggerOptions(
    log.WithLevel("info"),
    log.WithFormat("json"),
    log.WithPath("app"),
    log.WithMultiFiles(true), // 生成 app.debug.log, app.info.log, app.warn.log, app.error.log
)
```

### 配置选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `WithLevel(level)` | 日志级别 (debug/info/warn/error) | "info" |
| `WithFormat(format)` | 格式 (console/json) | "console" |
| `WithPath(path)` | 日志文件路径 | "logger" |
| `WithMultiFiles(multiFiles)` | 按级别分文件输出 | false |

### 主要函数

| 函数 | 说明 |
|------|------|
| `NewLoggerOptions(opts...)` | 创建日志配置 |
| `NewSlogLogger(opts)` | 创建 slog 后端日志器 |
| `NewZapLogger(opts)` | 创建 zap 后端日志器 |
| `NewLogFile(path)` | 创建文件写入器 |

### YAML 配置示例

```yaml
# config.yaml
log:
  level: "info"
  format: "console"      # console 或 json
  path: "app.log"        # 日志文件路径
  multi_files: false     # 是否按级别分文件输出
```

**解析 YAML 配置：**

```go
package main

import (
    "github.com/atompi/goutil/log"
    "gopkg.in/yaml.v3"
    "os"
)

type Config struct {
    Log log.Config `yaml:"log"`
}

func main() {
    data, _ := os.ReadFile("config.yaml")
    var cfg Config
    yaml.Unmarshal(data, &cfg)

    opts := cfg.Log.ToOptions()
    logger := log.NewSlogLogger(log.NewLoggerOptions(opts...))
    logger.Info("hello world")
}
```

**多文件输出时生成的日志文件：**

| 级别 | 文件名 |
|------|--------|
| debug | app.debug.log |
| info | app.info.log |
| warn | app.warn.log |
| error | app.error.log |

## 测试

```bash
# 测试 bar 模块
cd bar && go test ./...

# 测试 log 模块
cd log && go test ./...
```

## 构建

```bash
# 构建 bar 模块
cd bar && go build ./...

# 构建 log 模块
cd log && go build ./...
```

## 项目结构

```
goutil/
├── bar/           # 进度条工具
│   ├── bar.go
│   └── bar_test.go
├── log/           # 日志工具
│   ├── common.go  # 文件写入器、路径验证、Logger 配置
│   ├── slog.go    # slog.Handler 实现
│   ├── zap.go     # zap.Logger 封装
│   └── *_test.go
└── README.md
```

## 注意事项

- `bar/` 模块无外部依赖
- `log/` 模块依赖 `testify` (v1.8.1) 和 `zap` (v1.27.1)
- 每个模块是独立的 Go 模块，需要单独执行 `go mod tidy`
