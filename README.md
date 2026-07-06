# 60s TUI

基于 [vikiboss/60s API](https://github.com/vikiboss/60s) 构建的终端信息工具。每天 60 秒看世界，也可以查看热榜、天气、实用数据和常用工具。

默认内置公共 API 地址，也支持填写自己部署的 API 地址。更推荐自行部署，稳定性和可控性更好。

## 首次使用

启动程序：

```bash
60s
```

首次打开会显示配置引导：

- 可以从内置的公共实例列表中选择服务器
- 也可以选择自部署，并手动输入 API 地址
- 可前往 [vikiboss/60s](https://github.com/vikiboss/60s) 查看部署方式

在应用内保存 API 地址后才会同步数据。输入地址的过程中不会触发请求；检测成功后会自动保存配置。

## 常用操作

| 按键 | 说明 |
| --- | --- |
| `↑` / `↓` 或 `k` / `j` | 移动选择 |
| `/` | 搜索接口，支持按名称、路径或序号搜索 |
| `Enter` | 确认 / 请求接口 |
| `s` | 设置服务器和默认格式 |
| `r` | 刷新或重新请求 |
| `Esc` / `b` | 返回 |
| `q` | 退出 |

## 配置

配置文件默认保存到：

```text
~/.config/endpoint-tui/config.json
```

示例：

```json
{
  "base_url": "https://your-api.example.com",
  "server_mode": "custom",
  "default_encoding": "json",
  "language": "zh"
}
```

## 发行版

可以直接点击 Release 下载最新版本

## 从源码运行

依赖：Go 1.21+、curl。

```bash
git clone https://github.com/olddoge/60s-tui.git
cd 60s-tui
go mod tidy
go build -o 60s .
./60s
```

Debian 安装 curl：

```bash
sudo apt update && sudo apt install -y curl
```

## 开发

```bash
go test ./...
go vet ./...
gofmt -w .
```

## 相关项目

- [vikiboss/60s](https://github.com/vikiboss/60s)