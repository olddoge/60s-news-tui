# 60s TUI

基于 [60s API](https://github.com/vikiboss/60s) 的终端新闻浏览工具。

[60s](https://github.com/vikiboss/60s) 是一个汇聚每日精选新闻的开源项目，提供涵盖科技、社会、娱乐、财经等数十个类别的 API 接口。60s TUI 让你在终端中浏览和阅读这些内容，无需离开命令行。

## 功能

- 自动获取 60s API 的全部可用接口
- 终端中浏览 70+ 个新闻/资讯分类
- 支持 JSON / Text / Markdown 三种输出格式
- 终端内直接查看新闻内容，支持滚动翻页
- JSON 自动格式化，Markdown 保留原始排版
- 根路径可配置，配置持久化保存

## 快速开始

### 从 Release 下载

在 [Releases](https://github.com/your-username/endpoint-tui/releases) 页面下载对应平台的二进制文件。

```bash
# Linux AMD64
sudo install -m 755 endpoint-tui-linux-amd64 /usr/local/bin/60s

# 运行
60s
```

### 从源码编译

**依赖：** Go 1.21+、curl

```bash
# Debian
sudo apt update && sudo apt install -y curl

# 编译
go mod tidy
go build -o 60s .

# 运行
./60s
```

## 使用方式

```bash
# 基本使用
60s

# 指定配置文件
60s --config /path/to/config.json

# 自定义发现服务地址
60s --discovery-url http://your-server:13205
```

### 接口列表页

终端展示 60s 提供的所有接口，`↑` `↓` 选择感兴趣的新闻类别：

| 按键 | 功能 |
|------|------|
| `↑` / `k` | 选择上一项 |
| `↓` / `j` | 选择下一项 |
| `Enter` | 选择接口，进入格式选择 |
| `r` | 刷新接口列表 |
| `s` | 打开设置 |
| `q` | 退出 |

### 格式选择

选择输出格式（JSON / Text / Markdown）：

| 按键 | 功能 |
|------|------|
| `↑` / `↓` | 切换格式 |
| `Enter` | 发起请求 |
| `Esc` | 返回列表 |

### 结果查看

请求结果支持完整滚动阅读：

| 按键 | 功能 |
|------|------|
| `↑` `↓` / `PgUp` `PgDn` | 滚动 |
| `Home` / `End` | 跳到顶部/底部 |
| `r` | 重新请求 |
| `b` / `Esc` | 返回列表 |

### 设置

| 按键 | 功能 |
|------|------|
| `↑` / `↓` | 选择默认格式 |
| `Ctrl+S` | 保存 |
| `Esc` | 取消 |

## 配置

配置文件 `~/.config/endpoint-tui/config.json`：

```json
{
  "base_url": "http://127.0.0.1:8080",
  "default_encoding": "json"
}
```

### 环境变量

| 变量 | 说明 |
|------|------|
| `ENDPOINT_DISCOVERY_URL` | 接口发现服务地址 |
| `NO_COLOR` | 设为任意值禁用颜色输出 |

## 技术栈

- 语言：Go
- TUI 框架：[Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Bubbles](https://github.com/charmbracelet/bubbles) + [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- HTTP 请求：系统 `curl`（安全沙箱调用）

## 项目结构

```
60s-tui/
├── main.go              # 入口
├── cmd/root.go          # CLI 参数
├── internal/
│   ├── app/             # Bubble Tea 页面模型
│   ├── api/             # 接口发现、解析、curl 调用
│   ├── config/          # 配置管理
│   ├── ui/              # 终端样式
│   └── urlutil/         # URL 安全拼接
└── tests/               # 47 个单元测试
```

## 开发

```bash
# 运行测试
go test ./...

# 静态检查
go vet ./...

# 格式化
gofmt -w .
```

## 手动编译当前测试版：

```
go build -o endpoint-tui-test.exe .
```

## 相关项目

- [60s](https://github.com/vikiboss/60s) —— 每日精选新闻 API

## 许可证

MIT
