# CLAUDE.md

## 项目概述

开发一个运行在 Debian Linux 服务器上的 TUI 终端工具。

该工具从指定服务获取接口列表，在终端中展示接口，用户可以通过方向键选择接口。选中接口后按回车键，程序需要调用系统中的 `curl` 命令访问该接口，并在 TUI 页面中展示请求结果。

项目需要编译为单个可执行文件，不依赖图形界面。

---

## 技术栈

使用以下技术开发：

* 开发语言：Go
* TUI 框架：Bubble Tea
* TUI 组件：Bubbles
* 样式库：Lip Gloss
* 配置文件格式：JSON
* HTTP 请求执行方式：调用系统中的 `curl`
* 目标系统：Debian Linux

依赖包建议：

```text
github.com/charmbracelet/bubbletea
github.com/charmbracelet/bubbles
github.com/charmbracelet/lipgloss
```

不要使用浏览器、Electron 或其他 GUI 框架。

---

## 数据源

程序启动后访问以下地址获取接口数据：

```text
http://localhost:13205/
```

从返回的 JSON 数据中读取 `endpoint` 字段。

`endpoint` 字段中包含可供用户选择的接口地址。

程序必须对返回结构做兼容处理，至少支持以下形式。

### endpoint 是字符串数组

```json
{
  "endpoint": [
    "/api/article",
    "/api/document",
    "/api/status"
  ]
}
```

### endpoint 是对象数组

```json
{
  "endpoint": [
    {
      "name": "获取文章",
      "path": "/api/article"
    },
    {
      "name": "获取文档",
      "path": "/api/document"
    }
  ]
}
```

### endpoint 是对象

```json
{
  "endpoint": {
    "article": "/api/article",
    "document": "/api/document"
  }
}
```

解析时应尽量兼容实际返回结构。

如果无法解析 `endpoint`，需要在页面中展示明确的错误信息，并允许用户重新加载。

---

## 核心功能

### 1. 获取接口列表

程序启动后自动访问：

```text
http://localhost:13205/
```

获取并解析返回数据中的 `endpoint` 字段。

接口列表请求可以使用 Go 标准库 `net/http`，不强制使用 `curl`。

接口列表加载期间显示加载状态：

```text
正在加载接口列表...
```

加载成功后进入接口选择页面。

加载失败时显示：

* HTTP 状态码
* 网络错误
* JSON 解析错误
* `endpoint` 字段不存在或格式不支持

用户可以按 `r` 重新加载。

---

### 2. 展示接口列表

接口列表需要在终端中展示。

每一项至少显示：

* 接口名称
* 接口路径

如果接口没有名称，则直接使用接口路径作为名称。

示例：

```text
接口列表

> 获取文章       /api/article
  获取文档       /api/document
  服务状态       /api/status
```

操作方式：

| 按键        | 功能       |
| --------- | -------- |
| `↑` / `k` | 选择上一项    |
| `↓` / `j` | 选择下一项    |
| `Enter`   | 打开请求配置页面 |
| `r`       | 重新获取接口列表 |
| `s`       | 打开全局设置   |
| `q`       | 退出程序     |
| `Ctrl+C`  | 强制退出程序   |

接口数量超过终端高度时必须支持滚动。

---

### 3. 根路径配置

接口列表中的地址可能只包含接口路径，没有协议、域名和端口。

例如：

```text
/api/article
```

程序需要将用户配置的根路径与接口路径组合：

```text
根路径：http://127.0.0.1:8080
接口路径：/api/article
完整地址：http://127.0.0.1:8080/api/article
```

根路径由用户设置，并持久化保存。

推荐配置文件位置：

```text
~/.config/endpoint-tui/config.json
```

配置文件示例：

```json
{
  "base_url": "http://127.0.0.1:8080",
  "default_encoding": "json"
}
```

如果目录不存在，程序需要自动创建。

配置文件权限建议为：

```text
0600
```

首次启动时，如果没有配置根路径，应自动打开设置页面，要求用户填写根路径。

根路径输入示例：

```text
http://127.0.0.1:8080
```

保存前进行以下处理：

1. 去除首尾空格。
2. 检查是否以 `http://` 或 `https://` 开头。
3. 去除根路径末尾多余的 `/`。
4. 确保接口路径与根路径之间只有一个 `/`。
5. 不要对接口地址进行重复拼接。
6. 如果 `endpoint` 已经是完整的 HTTP 或 HTTPS 地址，则直接使用该地址。

地址拼接必须使用可靠的 URL 处理方式，不要简单地无条件拼接字符串。

---

### 4. encoding 参数选择

所有业务接口都是 GET 请求。

业务接口有一个查询参数：

```text
encoding
```

允许的值只有：

```text
json
text
markdown
```

用户不能自由输入 encoding，只能从固定选项中选择。

接口选择完成后进入请求配置页面：

```text
接口：/api/article

请选择返回格式：

> json
  text
  markdown
```

操作方式：

| 按键        | 功能     |
| --------- | ------ |
| `↑` / `k` | 选择上一项  |
| `↓` / `j` | 选择下一项  |
| `Enter`   | 执行请求   |
| `Esc`     | 返回接口列表 |

最终请求地址示例：

```text
http://127.0.0.1:8080/api/article?encoding=json
```

如果原接口地址已经有查询参数，例如：

```text
/api/article?id=100
```

必须正确追加参数：

```text
http://127.0.0.1:8080/api/article?id=100&encoding=json
```

必须使用 URL 查询参数处理方法，不能通过判断字符串中是否存在 `?` 来粗暴拼接。

---

### 5. 使用 curl 请求接口

用户确认 encoding 后，程序在系统内调用 `curl`。

基础命令示例：

```bash
curl --silent --show-error --location \
  --connect-timeout 10 \
  --max-time 60 \
  "http://127.0.0.1:8080/api/article?encoding=json"
```

Go 中必须使用：

```go
exec.CommandContext
```

禁止通过以下形式执行用户输入：

```go
exec.Command("sh", "-c", command)
```

禁止使用 shell 拼接命令，避免命令注入。

正确方式示例：

```go
exec.CommandContext(
    ctx,
    "curl",
    "--silent",
    "--show-error",
    "--location",
    "--connect-timeout",
    "10",
    "--max-time",
    "60",
    requestURL,
)
```

必须捕获：

* 标准输出
* 标准错误
* curl 退出码
* 请求耗时
* 最终请求地址

请求执行期间显示：

```text
正在请求接口...
```

请求必须异步执行，不能阻塞 Bubble Tea 的主事件循环。

请求过程中应支持按键取消，例如：

```text
Ctrl+C
```

或者：

```text
Esc
```

取消时需要终止当前 `curl` 子进程。

---

## 页面设计

程序至少包含以下页面状态：

```text
PageLoading
PageEndpointList
PageEncodingSelect
PageResult
PageSettings
PageError
```

建议使用统一的 Model 管理页面状态，不要为每个页面启动独立的 Bubble Tea 程序。

---

## 接口列表页面

页面需要显示：

* 程序标题
* 当前配置的根路径
* 接口列表
* 当前选中项
* 底部快捷键说明
* 加载错误或状态提示

示例：

```text
Endpoint TUI

根路径：http://127.0.0.1:8080

接口列表：

> /api/article
  /api/document
  /api/status

↑/↓ 选择  Enter 请求  s 设置  r 刷新  q 退出
```

---

## 请求结果页面

请求完成后，在结果页面展示：

* 接口名称
* 完整请求地址
* encoding
* curl 退出状态
* 请求耗时
* 返回内容
* 错误信息

示例：

```text
请求完成

接口：/api/article
格式：json
地址：http://127.0.0.1:8080/api/article?encoding=json
耗时：325ms
状态：成功

返回内容：

{
  "code": 200,
  "message": "success"
}
```

操作方式：

| 按键          | 功能                    |
| ----------- | --------------------- |
| `↑` / `k`   | 向上滚动                  |
| `↓` / `j`   | 向下滚动                  |
| `PgUp`      | 向上翻页                  |
| `PgDn`      | 向下翻页                  |
| `Home`      | 跳到顶部                  |
| `End`       | 跳到底部                  |
| `r`         | 使用相同接口和 encoding 重新请求 |
| `b` / `Esc` | 返回接口列表                |
| `q`         | 退出程序                  |

返回内容可能很长，必须使用 Bubbles 的 `viewport` 组件展示，不能直接一次性输出导致页面失控。

---

## 不同 encoding 的显示规则

### json

当 encoding 为 `json` 时：

1. 尝试解析返回内容。
2. 如果是合法 JSON，使用缩进格式化后显示。
3. 如果不是合法 JSON，展示原始内容。
4. JSON 格式化失败不能导致程序崩溃。

格式化建议：

```go
json.Indent
```

### text

当 encoding 为 `text` 时：

* 原样显示文本。
* 保留换行。
* 过滤或安全处理可能影响终端显示的控制字符。

### markdown

当 encoding 为 `markdown` 时：

* 第一版可以直接显示 Markdown 原文。
* 可以使用 Glamour 对 Markdown 进行终端渲染，但不是强制要求。
* 如果引入 Glamour，不得影响纯文本模式和小尺寸终端。

可选依赖：

```text
github.com/charmbracelet/glamour
```

---

## 设置页面

设置页面至少包含：

```text
根路径
默认 encoding
```

其中：

* 根路径允许输入。
* 默认 encoding 只能选择 `json`、`text` 或 `markdown`。
* 保存后立即生效。
* 保存失败时显示明确错误。
* 用户可以按 `Esc` 放弃修改并返回。

示例：

```text
设置

根路径：
http://127.0.0.1:8080

默认格式：
> json
  text
  markdown

Ctrl+S 保存
Esc    取消
```

保存成功后显示短暂状态：

```text
配置已保存
```

---

## 配置结构

建议配置结构：

```go
type Config struct {
    BaseURL         string `json:"base_url"`
    DefaultEncoding string `json:"default_encoding"`
}
```

默认值：

```go
Config{
    BaseURL:         "",
    DefaultEncoding: "json",
}
```

配置加载规则：

1. 配置文件不存在时使用默认配置。
2. JSON 无法解析时不能 panic。
3. 配置错误时提示用户重新设置。
4. 非法的默认 encoding 自动回退到 `json`。
5. 根路径为空时进入设置页面。

---

## 接口数据结构

内部统一转换成以下结构：

```go
type Endpoint struct {
    Name string
    Path string
}
```

无论服务端的 `endpoint` 是字符串数组、对象数组还是对象，解析后都转换为 `[]Endpoint`。

接口列表应过滤：

* 空字符串
* 无法识别的值
* 完全重复的接口

可以保留原始顺序。

---

## 项目目录

推荐目录结构：

```text
endpoint-tui/
├── CLAUDE.md
├── README.md
├── go.mod
├── go.sum
├── main.go
├── cmd/
│   └── root.go
├── internal/
│   ├── app/
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── messages.go
│   ├── api/
│   │   ├── discovery.go
│   │   ├── parser.go
│   │   └── request.go
│   ├── config/
│   │   └── config.go
│   ├── ui/
│   │   ├── styles.go
│   │   ├── endpoint_list.go
│   │   ├── encoding_select.go
│   │   ├── result.go
│   │   └── settings.go
│   └── urlutil/
│       └── url.go
└── tests/
    ├── endpoint_parser_test.go
    ├── url_test.go
    └── config_test.go
```

避免过度拆分，但不要将全部代码写进 `main.go`。

---

## 命令行参数

程序建议命名为：

```text
endpoint-tui
```

支持以下参数：

```text
endpoint-tui
endpoint-tui --config /path/to/config.json
endpoint-tui --version
endpoint-tui --help
```

可选支持：

```text
endpoint-tui --discovery-url http://localhost:13205/
```

默认发现地址始终为：

```text
http://localhost:13205/
```

命令行参数优先级高于配置文件，但不要将临时命令行参数自动覆盖保存到配置文件。

---

## 错误处理

以下情况不能导致程序崩溃：

* 接口列表服务无法连接
* 接口列表请求超时
* 服务端返回非 2xx 状态码
* 返回内容不是 JSON
* `endpoint` 字段不存在
* `endpoint` 字段格式不支持
* 接口列表为空
* 根路径格式错误
* 系统没有安装 curl
* curl 请求超时
* curl 返回非零退出码
* 返回内容为空
* 配置文件无法创建
* 配置文件没有读写权限
* 终端尺寸太小

如果系统未安装 curl，显示：

```text
未找到 curl 命令。

请在 Debian 中执行：
sudo apt update && sudo apt install -y curl
```

检查方式：

```go
exec.LookPath("curl")
```

所有错误都应该作为页面状态展示，不要只写入日志后静默失败。

---

## 安全要求

1. 不允许使用 `sh -c`、`bash -c` 执行 curl。
2. 不允许将用户输入直接拼接成 shell 命令。
3. 根路径只允许 HTTP 或 HTTPS。
4. 对配置文件使用合理的文件权限。
5. 对接口返回中的 ANSI 转义序列进行过滤或转义，避免终端注入。
6. curl 必须设置连接超时和总超时。
7. 请求 URL 必须通过 `net/url` 构造。
8. 不记录敏感请求内容。
9. 错误日志中避免输出不必要的系统环境信息。
10. 禁止因为服务端返回异常数据而 panic。

---

## 终端兼容性

需要适配常见 SSH 终端。

最低要求：

* 支持终端窗口尺寸变化。
* 支持 80×24 的普通终端。
* 小尺寸终端下显示提示，而不是布局错乱。
* 不依赖鼠标。
* 不依赖 Unicode 图标才能正常使用。
* 配色不可影响无彩色终端。
* 设置 `NO_COLOR` 时禁用或减少颜色。

---

## 编译与部署

本地编译：

```bash
go mod tidy
go build -o endpoint-tui .
```

Linux AMD64 交叉编译：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -trimpath \
  -ldflags="-s -w" \
  -o endpoint-tui .
```

Linux ARM64 交叉编译：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
  -trimpath \
  -ldflags="-s -w" \
  -o endpoint-tui-arm64 .
```

部署示例：

```bash
sudo install -m 755 endpoint-tui /usr/local/bin/endpoint-tui
```

运行：

```bash
endpoint-tui
```

---

## 测试要求

至少编写以下单元测试。

### URL 拼接测试

覆盖：

```text
http://localhost:8080 + /api/test
http://localhost:8080/ + /api/test
http://localhost:8080 + api/test
http://localhost:8080/ + api/test
```

覆盖已有查询参数：

```text
/api/test?id=1
```

最终应正确追加：

```text
encoding=json
```

覆盖完整接口地址：

```text
https://example.com/api/test
```

完整地址不应再次拼接根路径。

### endpoint 解析测试

覆盖：

* 字符串数组
* 对象数组
* 键值对象
* 空 endpoint
* 缺失 endpoint
* endpoint 类型错误
* 重复接口
* 空接口路径

### 配置测试

覆盖：

* 配置文件不存在
* 正常读取
* 正常保存
* 非法 JSON
* 非法 encoding
* 无法创建配置目录

### curl 请求测试

不要在单元测试中依赖公网。

将命令执行逻辑抽象成接口，以便使用模拟执行器测试：

* 成功返回
* curl 不存在
* 请求超时
* 非零退出码
* 标准错误输出
* 用户取消请求

运行测试：

```bash
go test ./...
```

运行静态检查：

```bash
go vet ./...
```

---

## 验收标准

项目完成后必须满足以下条件：

1. 可以在 Debian Linux 终端正常启动。
2. 启动后可以从指定地址读取 `endpoint`。
3. 能够在终端中使用方向键选择接口。
4. 可以选择 `json`、`text` 或 `markdown`。
5. encoding 只能选择，不能自由输入。
6. 所有业务接口均通过 GET 请求访问。
7. 按回车后通过系统 curl 发起请求。
8. 能够将根路径与接口路径正确组合。
9. 根路径可以在 TUI 中修改。
10. 根路径配置可以持久化保存。
11. 重启程序后可以读取之前保存的配置。
12. 请求结果可以在终端页面中滚动查看。
13. JSON 返回结果可以自动格式化。
14. 请求失败时有明确错误信息。
15. 没有安装 curl 时有明确安装提示。
16. 用户可以重新加载接口列表。
17. 用户可以取消正在执行的请求。
18. 长返回内容不会破坏 TUI 页面。
19. 不使用 shell 拼接命令执行 curl。
20. `go test ./...` 可以正常通过。

---

## 开发约束

1. 优先实现稳定可运行的版本，不要加入无关功能。
2. 不要擅自改变接口列表来源地址。
3. 不要将根路径硬编码到程序中。
4. 不要让用户手动输入 encoding。
5. 不要将业务请求改为 POST。
6. 不要用 Go HTTP 客户端代替业务请求中的 curl。
7. 获取接口列表可以使用 Go HTTP 客户端。
8. 不要在请求时退出 TUI 界面。
9. 不要直接使用 `fmt.Println` 输出业务结果破坏 TUI。
10. 所有异步操作通过 Bubble Tea `Cmd` 和消息机制返回结果。
11. 代码必须通过 `gofmt`。
12. 公共函数和复杂逻辑需要有必要注释。
13. 不要为了抽象而抽象，保持代码可读。
14. 首个版本不实现接口的新增、删除和修改。
15. 首个版本不支持自定义 HTTP 请求方法。

---

## 实现顺序

按照以下顺序开发：

1. 初始化 Go 项目和依赖。
2. 实现配置文件读取、保存和默认值。
3. 实现根路径与接口路径组合。
4. 实现接口列表请求。
5. 实现 `endpoint` 多格式解析。
6. 实现接口列表页面。
7. 实现 encoding 选择页面。
8. 实现安全的 curl 调用。
9. 实现请求结果页面和滚动。
10. 实现设置页面。
11. 实现错误页面、重试和取消。
12. 编写单元测试。
13. 编写 README 和 Debian 部署说明。
14. 执行 `gofmt`、`go vet` 和 `go test ./...`。

