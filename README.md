# Endpoint TUI

终端接口调试工具 —— 在 Debian Linux 服务器上通过 TUI 浏览和调试 API 接口。

## 功能

- 从发现服务自动获取接口列表
- 终端中浏览、选择接口
- 配置根路径，自动拼接完整请求地址
- 选择返回格式（JSON / Text / Markdown）
- 通过系统 `curl` 发送 GET 请求
- 请求结果支持滚动查看
- JSON 返回自动格式化
- 配置持久化保存

## 安装

### 依赖

- Go 1.21+
- curl

```bash
# Debian
sudo apt update && sudo apt install -y curl
```

### 编译

```bash
# 本地编译
go mod tidy
go build -o endpoint-tui .

# Linux AMD64 交叉编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -trimpath \
  -ldflags="-s -w" \
  -o endpoint-tui .

# Linux ARM64 交叉编译
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
  -trimpath \
  -ldflags="-s -w" \
  -o endpoint-tui-arm64 .
```

### 部署

```bash
sudo install -m 755 endpoint-tui /usr/local/bin/endpoint-tui
```

### 运行

```bash
endpoint-tui
```

## 命令行参数

```
用法:
  endpoint-tui [选项]

选项:
  --config FILE        配置文件路径 (默认: ~/.config/endpoint-tui/config.json)
  --discovery-url URL  接口发现服务地址 (可通过环境变量 ENDPOINT_DISCOVERY_URL 设置)
  --version            显示版本号
  --help               显示帮助信息
```

## 操作方式

### 接口列表页

| 按键 | 功能 |
|------|------|
| `↑` / `k` | 选择上一项 |
| `↓` / `j` | 选择下一项 |
| `Enter` | 打开请求配置页面 |
| `r` | 重新获取接口列表 |
| `s` | 打开全局设置 |
| `q` | 退出程序 |
| `Ctrl+C` | 强制退出程序 |

### encoding 选择页

| 按键 | 功能 |
|------|------|
| `↑` / `k` | 选择上一项 |
| `↓` / `j` | 选择下一项 |
| `Enter` | 执行请求 |
| `Esc` | 返回接口列表 |

### 请求结果页

| 按键 | 功能 |
|------|------|
| `↑` / `k` | 向上滚动 |
| `↓` / `j` | 向下滚动 |
| `PgUp` | 向上翻页 |
| `PgDn` | 向下翻页 |
| `Home` | 跳到顶部 |
| `End` | 跳到底部 |
| `r` | 重新请求 |
| `b` / `Esc` | 返回接口列表 |
| `q` | 退出程序 |

### 设置页

| 按键 | 功能 |
|------|------|
| `↑` / `↓` | 选择默认格式 |
| `Ctrl+S` | 保存设置 |
| `Esc` | 取消并返回 |

## 配置文件

配置文件位于 `~/.config/endpoint-tui/config.json`：

```json
{
  "base_url": "http://127.0.0.1:8080",
  "default_encoding": "json"
}
```

## 项目结构

```
endpoint-tui/
├── main.go                 # 程序入口
├── cmd/
│   └── root.go             # 命令行参数解析
├── internal/
│   ├── app/                # Bubble Tea Model/Update/View
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── messages.go
│   ├── api/                # 接口发现、解析、curl 调用
│   │   ├── discovery.go
│   │   ├── parser.go
│   │   └── request.go
│   ├── config/             # 配置文件读写
│   │   └── config.go
│   ├── ui/                 # 样式定义
│   │   └── styles.go
│   └── urlutil/            # URL 拼接工具
│       └── url.go
└── tests/                  # 单元测试
    ├── url_test.go
    ├── endpoint_parser_test.go
    ├── config_test.go
    └── request_test.go
```

## 测试

```bash
go test ./...
go vet ./...
```

## 安全说明

- 使用 `exec.CommandContext` 执行 curl，不通过 shell 拼接命令
- 根路径只允许 HTTP/HTTPS
- 配置文件权限为 0600
- 自动过滤返回内容中的 ANSI 转义序列

## 许可证

MIT
