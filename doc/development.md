# 开发文档

本文档面向 go-sqli-lab 的贡献者与二次开发者，描述项目实际代码结构、核心设计约定（配置 / 日志 / 数据库 / 漏洞模块）以及测试与贡献流程。

- 新手安装与运行见 [installation.md](./installation.md)
- 漏洞利用练习指南见 [usage.md](./usage.md)
- 仓库级说明见 [../README.md](../README.md)

---

## 1. 环境要求

| 依赖 | 版本要求 | 说明 |
| --- | --- | --- |
| Go | 1.25+（`go.mod` 声明 `go 1.25.0`） | 数据库驱动 `modernc.org/sqlite` 为纯 Go 实现，默认构建**不需要 CGO** |
| Git | 任意较新版本 | 需拉取 `ref/sqli-labs` 子模块 |
| GNU Make | 可选 | 仅 `Makefile` 目标需要；Windows 可直接使用 `go` 命令 |
| MySQL / PostgreSQL | 可选 | 仅当开发与多数据库兼容性相关功能时需要，见 [第 9 节](#9-测试与ci) |

> 仓库没有内置 Dockerfile。`Makefile` 中 `docker-build` / `docker-run` 目标需要自行在仓库根目录提供 Dockerfile 后方可使用（参考 [installation.md](./installation.md) 的示例）。

## 2. 仓库结构

```
go-sqli-lab/
├── src/                          # Go 源码（module: go-sqli-lab，入口为 src/main.go）
│   ├── main.go                   # 程序入口：flag 解析、项目根定位、路由注册、优雅关闭
│   ├── config/
│   │   └── config.go             # 配置加载：默认值 < app.yaml < SQLI_LAB_* 环境变量
│   ├── db/
│   │   ├── db.go                 # 数据库层：Dialect 接口 + 多库适配（迁移/种子/重置）
│   │   ├── db_test.go            # 数据库集成测试（默认 SQLite，可选 MySQL/PostgreSQL）
│   │   └── test_data/            # SQLite 测试库文件（被 .gitignore 忽略的运行时文件）
│   ├── handlers/
│   │   └── handlers.go           # 健康检查、DB 状态、重置数据库、setup 结果页
│   ├── logger/
│   │   ├── logger.go             # 日志接口 + lumberjack 轮转默认实现
│   │   └── gin.go                # GinLogger HTTP 请求日志中间件
│   ├── models/
│   │   └── user.go               # User / Email / UAgent / Referer 数据模型
│   ├── vulnerabilities/          # 漏洞实现（71 个关卡 + 共享工具）
│   │   ├── vulnerability.go      # Vulnerability 接口定义
│   │   ├── registry.go           # Registry 注册表 + /api/vulnerabilities API
│   │   ├── base.go               # BaseLesson：db/logger 注入、renderHTML、logRequest
│   │   ├── less1.go ~ less71.go  # 各关卡实现
│   │   ├── export.go             # Less-66/67 导出（关键字查询 + xlsx/xls 生成）
│   │   ├── upload.go             # Less-70/71 上传（xlsx/xls 解析与入库）
│   │   ├── resume_export.go      # Less-68/69 简历导出（模板填充 xlsx/xls）
│   │   └── assets/               # go:embed 嵌入的模板文件（resume_template.xlsx 等）
│   └── web/                      # 前端资源（go:embed 之外，经 /static 与模板提供）
│       ├── index.html            # 首页（71 个关卡入口，html/template 渲染）
│       ├── loading.html / setup_success.html / setup_error.html
│       ├── css/                  # style.css / common.css
│       ├── assets/               # 横幅 SVG 等
│       └── less-1/ ~ less-65/    # 各关卡静态资源目录（index.html/style.css/script.js）
├── tools/
│   └── genresume/                # 简历模板生成器（输出 src/vulnerabilities/assets/resume_template.xlsx）
├── config/app.yaml               # 默认配置文件
├── data/                         # SQLite 数据库文件目录（运行时生成 sqli_lab.db）
├── logs/                         # 应用日志目录（运行时生成 app.log）
├── tmp/                          # logRequest 记录的练习请求痕迹（tmp/<lesson>/result.txt）
├── bin/                          # 构建产物目录（go-sqli-lab / go-sqli-lab.exe）
├── scripts/                      # 三平台一键启动脚本（构建 + 初始化 + 打开浏览器）
├── ref/sqli-labs/                # ⚠️ 只读 Git 子模块（原始 PHP 参考实现）
└── doc/                          # 文档（installation / usage / development）
```

## 3. 启动流程与代码路径（main.go）

启动链路：`main()` → flag 解析 → `ensureProjectRoot()` → `config.Load()` → `logger.New()` → `db.New()` → 按模式执行。

1. **flag 解析**：`--config`（配置文件路径，默认 `config/app.yaml`，相对路径会先转绝对路径）；`--setup-db`（仅初始化数据库后退出）。
2. **项目根定位**（`ensureProjectRoot`）：依次尝试 —— 当前目录直接是项目根（同时存在 `go.mod` 与 `config/app.yaml`）；可执行文件位于 `bin/` 时取其父目录；向上逐级查找。因此**预编译二进制可以在仓库内任意子目录运行**。
3. **数据库初始化**：`--setup-db` 模式下执行 `database.Migrate()` + `database.Seed()` 后退出。
4. **Web 模式**：`gin.New()` + `gin.Recovery()` + `logger.GinLogger` 中间件，随后 `registerRoutes`：
   - `r.LoadHTMLGlob("./src/web/*.html")`：首页与 setup 页模板；
   - `r.Static("/static", "./src/web")` + `r.StaticFile("/", ...)`：静态资源；
   - `/health`、`/api/db-status`、`/api/reset-db`（AJAX）、`/setup-db`、`/loading`、`/setup-success`、`/setup-error`；
   - `vulnerabilities.NewRegistry(...).RegisterAll(r)`：注册全部关卡路由（`/less/*`）与漏洞列表 API。
5. **优雅关闭**：监听 `SIGINT` / `SIGTERM` 后退出。

## 4. 配置管理（config）

配置优先级（后者覆盖前者）：

```
代码内默认值 < config/app.yaml < SQLI_LAB_* 环境变量
```

核心配置项（`config/app.yaml`，注释已说明取值范围）：

| 分组 | 键 | 默认值 | 环境变量 |
| --- | --- | --- | --- |
| server | host / port / mode | `0.0.0.0` / `8080` / `debug` | `SQLI_LAB_HOST` / `SQLI_LAB_PORT` / `SQLI_LAB_MODE` |
| database | type | `sqlite`（可选 `mysql` / `postgres`） | `SQLI_LAB_DB_TYPE` |
| database | dsn | `./data/sqli_lab.db` | `SQLI_LAB_DB_DSN` |
| database | host/port/user/password/dbname/sslmode | 见 app.yaml | `SQLI_LAB_DB_HOST` / `DB_PORT` / `DB_USER` / `DB_PASSWORD` / `DB_NAME` |
| log | level / output / filepath 等 | `info` / `file` / `./logs/app.log` | `SQLI_LAB_LOG_LEVEL` / `SQLI_LAB_LOG_FILE` |

**新增配置项的落地步骤**：

1. `src/config/config.go`：在对应结构体（`ServerConfig` / `DatabaseConfig` / `LogConfig`）中添加字段；
2. 在 `Load()` 的默认值与 env 覆盖段中同步（名称保持 `SQLI_LAB_*` 前缀规范）；
3. `config/app.yaml` 补充注释化的默认值；
4. 使用方（如 `db.New` 的 DSN 组装）读取新字段。

**多数据库切换示例**（本地开发时常用）：

```bash
# 使用 MySQL
set SQLI_LAB_DB_TYPE=mysql
set SQLI_LAB_DB_HOST=127.0.0.1
set SQLI_LAB_DB_PORT=3306
set SQLI_LAB_DB_USER=root
set SQLI_LAB_DB_PASSWORD=yourpass
set SQLI_LAB_DB_NAME=security
go run src/main.go --setup-db
```

## 5. 日志系统（logger）

- `logger.New(cfg LogConfig)` 返回 `logger.Logger` 接口（`Debug/Info/Warn/Error/Fatal`），支持 **KV 结构化字段**，例如 `log.Info("注册漏洞", "id", v.ID(), "name", v.Name())`。
- 默认实现写文件并轮转：由 app.yaml 的 `log` 节控制（`filepath`、`maxsize`(MB)、`maxbackups`、`maxage`(天)、`compress`）；`output` 可为 `file` 或 `stdout`。日志行格式由 `format` 定义，支持 `%timestamp% %level% %goroutine% %module% %function% %file% %line% %message%` 令牌。
- `logger.GinLogger` 是 Gin 中间件，自动记录每个 HTTP 请求（方法、路径、状态等），在 `main.go` 中全局挂载。
- 测试中常用配置：`logger.New(config.LogConfig{Level: "error", Output: "stdout"})`（见 `db_test.go`）。

**最佳实践**：

- 结构化字段优先于拼接字符串：`log.Error("查询失败", "lesson", "less-8", "error", err)`；
- 业务错误按级别区分：正常练习流量无需记录到应用日志（由 `logRequest` 承接，见下），真正异常才用 `Error`；
- **区分两套“记录”**：
  - `BaseLesson.logRequest(c, lesson)` → 追加写入 `tmp/<lesson>/result.txt`，记录请求方法与查询参数，供练习复盘（如 sqlmap 扫描痕迹）；`make clean` 会清空 `tmp/`；
  - 应用日志 `logs/app.log` → 系统运行日志。

## 6. 数据库层（db）

### 6.1 接口

`db.go` 中 `Database` 接口统一对上层（关卡、handlers）暴露：

| 方法 | 用途 |
| --- | --- |
| `GetDB() *sql.DB` | 原始连接（如事务） |
| `Migrate() error` | 建表（幂等，方言分支） |
| `Seed() error` | 插入种子数据（幂等） |
| `Reset() error` | 清空全部表并强制重灌 |
| `Query / QueryRow / Exec` | 便捷执行（占位符参数或裸 SQL，供关卡拼接注入） |
| `Close() error` | 关闭连接池 |

`db.New(cfg, log)` 依据 `cfg.Type` 选择实现并组装 DSN：sqlite（`modernc.org/sqlite`，纯 Go）、mysql（`go-sql-driver/mysql`）、postgres（`lib/pq`）。

### 6.2 表结构总览（Migrate 创建 17 张表）

| 表 | 用途 |
| --- | --- |
| `users` / `emails` / `uagents` / `referers` | 练习基础表（Less-1~22、46~53 等） |
| `challenge1` ~ `challenge12` | 挑战关卡专用（Less-54~65 固定查询各自挑战表） |
| `resumes` | 简历导出专用（Less-68/69），含 20 列中文求职数据 |

种子账号固定 8 条（与 PHP 原版一致）：`Dumb/Dumb`、`Angelina/I-kill-you`、…、`admin/admin`。

### 6.3 迁移 / 种子 / 重置语义（重要）

| 方法 | 行为 |
| --- | --- |
| `Migrate()` | 按方言分支 `CREATE TABLE IF NOT EXISTS`，可重复执行 |
| `Seed()` | **幂等**：先保证 `resumes` 已种入（独立判断）；再检查 `users` 是否非空，非空则跳过。首次 `--setup-db` 后仅有基础表数据 |
| `Reset()` | 清空 **17 张表**（`users/emails/uagents/referers/resumes` + `challenge1~12`）→ `forceSeed()` 强制重灌 → `seedChallenges()` 注入挑战表数据 |

清空语句按方言区分：sqlite `DELETE FROM`、mysql `TRUNCATE TABLE`、postgres `TRUNCATE TABLE ... RESTART IDENTITY`。

> 挑战关卡（Less-54~65）的 `challengeN` 数据只经 `Reset()` 注入。首次安装（`--setup-db`）后请先在首页执行一次“重置数据库”，或访问 `http://localhost:8080/setup-db?reset=true` 触发，挑战关卡才有数据可注入。

### 6.4 跨数据库兼容约定（新增 SQL 时必读）

漏洞 SQL 以**字符串拼接**刻意实现，因此每一条语句都必须兼容三种数据库：

1. **分页写法统一 `LIMIT 1 OFFSET 0`**（禁止 `LIMIT 0,1`，PostgreSQL 不支持）；
2. **占位符**：sqlite/mysql 用 `?`，postgres 用 `$1`；同一逻辑尽量写成方言分支（参考 `Seed()`/`forceSeed()` 中对 postgres 的分支写法）；
3. 自增主键显式插入 id 时注意 postgres 需要 `RESTART IDENTITY` 才能重置自增序列；
4. 新增表请保持小写表名与既有风格（单数）；字符串值默认 `'...'` 包裹，与各关卡利用方式保持一致。

### 6.5 修改数据模型时的同步清单

当新增业务表或改动字段时，请同步：

1. `db.go Migrate()` 中对应方言建表语句；
2. `Seed()` / `forceSeed()` / `seedResumes()` / `seedChallenges()` 的种子数据；
3. `Reset()` 的 `tables` 清空列表；
4. `src/db/db_test.go` 的 `tables` 断言与行数断言；
5. 文档：`README.md`（功能描述）、`doc/usage.md`（如影响关卡）与 `CHANGELOG.md`。

## 7. 漏洞关卡开发（vulnerabilities）

### 7.1 接口与基类

```go
type Vulnerability interface {
    ID() string          // 唯一 ID，如 "less-71"
    Name() string        // 展示名
    Description() string // 中文描述
    Category() string    // 分类（与首页分类对应，如 "Error Based"/"Blind"/"WAF Bypass" 等）
    Route(r *gin.RouterGroup) // 在 /less 组下注册本关卡路由
}
```

`BaseLesson` 被所有关卡内嵌，提供：

- `db db.Database`、`logger logger.Logger`（构造时由 `registerLessN` 注入）；
- `logRequest(c, lesson)`：请求记录到 `tmp/<lesson>/result.txt`（lesson 可带动作后缀，如 `less-24-create`）；
- `renderHTML(c, lesson, result, errMsg)`：内联输出统一风格 HTML 页面（含 `/static/css/style.css` 与“返回首页”链接），结果与错误分别经 `renderResult` / `renderError` 包成成功/失败样式块。

### 7.2 典型关卡骨架（以 GET 单参数为例）

```go
package vulnerabilities

import (
    "database/sql"
    "fmt"

    "go-sqli-lab/src/models"

    "github.com/gin-gonic/gin"
)

type LessX struct {
    BaseLesson
}

func (l *LessX) ID() string          { return "less-x" }
func (l *LessX) Name() string        { return "Less-X: 示例关卡" }
func (l *LessX) Description() string { return "GET请求 - 字符型注入示例" }
func (l *LessX) Category() string    { return "Error Based" }

func (l *LessX) Route(r *gin.RouterGroup) {
    group := r.Group("/less-x")
    {
        group.GET("", l.handleIndex)
        group.GET("/", l.handleIndex)
        // 按需补充 POST /login 等路由（参考 less11.go / less24.go）
    }
}

func (l *LessX) handleIndex(c *gin.Context) {
    l.logRequest(c, "less-x")

    id := c.Query("id")
    if id == "" {
        id = "1"
    }

    // 刻意使用字符串拼接构造漏洞（仅用于教学靶场）
    query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

    var user models.User
    err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            l.renderHTML(c, "less-x", "", "没有找到记录")
        } else {
            l.renderHTML(c, "less-x", "", err.Error()) // 报错注入：透出 SQL 错误
        }
        return
    }

    result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
    l.renderHTML(c, "less-x", result, "")
}

// registry.go: RegisterAll 中按编号加入一行
func (r *Registry) registerLessX() {
    r.Register(&LessX{BaseLesson{db: r.db, logger: r.logger}})
}
```

### 7.3 新增一个关卡的标准步骤

1. **实现**：在 `src/vulnerabilities/lessX.go` 中按 7.2 骨架实现 `Vulnerability` 接口；有漏洞点的地方必须**刻意用 `fmt.Sprintf` 拼接**，并注释说明漏洞类型；
2. **注册**：在 `registry.go` 的 `RegisterAll()` 中按编号顺序添加 `r.registerLessX()`（`Register` 会自动记录日志并挂到 `/less` 组）；
3. **入口与资源**：简单的关卡由 handler 内联渲染即可；如需独立页面/脚本，在 `src/web/less-X/` 下建资源目录（通过 `/static/less-X/...` 引用，样式统一引 `/static/css/style.css`）；
4. **首页入口**：在 `src/web/index.html` 对应分类区块添加 `<a href="/less/less-x?id=1">Less-X: 描述</a>`；若引入新分类，需在 `src/web/css/style.css` / `common.css` 中为该分类补充样式类并保持首页区块结构一致；
5. **数据**：若查询新表，按 6.5 清单同步 `db.go` 与测试；
6. **文档**：更新 `README.md` 关卡地图、`doc/usage.md`（如提供新利用思路）并追加 `CHANGELOG.md`。

### 7.4 特殊类型关卡共享实现

- **导出（Less-66/67，export.go）**：页面提供关键字搜索框，后端拼 `LIKE '%keyword%'` 查询并生成 `xlsx/xls` 文件流下载；xlsx 用标准库 `archive/zip` 手写 OOXML，xls 用 HTML + office 命名空间伪表格 —— 关键字处即注入点。
- **简历导出（Less-68/69，resume_export.go）**：查询 `resumes` 表（20 列），xlsx 走 `excelize` + `go:embed` 模板（`assets/resume_template.xlsx`，字段按 `RESUME:<字段名>` 占位符替换）；模板本身由 `tools/genresume` 生成 —— 改模板时改 `tools/genresume/main.go` 并重新生成，不要手改二进制模板。
- **上传（Less-70/71，upload.go）**：先下载模板（xlsx 动态生成 / xls 嵌入资源），用户以单元格内容作 keyword 上传；`processUploadData` 跳过表头行、取每行第一列入库；上传大小限制 10MB；xls 解析依赖 `github.com/extrame/xls`。

## 8. 前端说明

项目没有通用 HTML 模板引擎（旧版文档所述的 `{{.Title}}` 模板体系**不存在**），页面分三种形态：

1. **html/template 页面**：仅 `src/web/*.html`（`index.html`、`loading.html`、`setup_success.html`、`setup_error.html`），由 `r.LoadHTMLGlob` 加载、handlers 渲染 —— 首页即在此种；
2. **关卡静态资源**：`src/web/less-1/ ~ less-65/` 的 `index.html` / `style.css` / `script.js`，经 `/static` 前缀访问；这些页面负责表单 UI，实际数据渲染仍由后端完成；
3. **后端内联渲染**：多数关卡（含 Less-66~71，`renderExportPage` / `renderUploadPage` 等）直接在 Go 代码里以字符串拼 HTML 输出，`renderHTML` 提供统一壳层。

前端样式集中在 `src/web/css/style.css`（关卡通用）与 `common.css`（布局基元）；修改首页结构时注意与 `handlers` 渲染的 setup 流程页面（`/setup-db` 等）保持一致的视觉风格。

## 9. 测试与 CI

### 9.1 数据库集成测试（src/db/db_test.go）

五个测试用例覆盖数据库层全部契约：

| 用例 | 验证点 |
| --- | --- |
| `TestDatabase_Connect` | 连接 + Ping |
| `TestDatabase_Migrate` | 16 张业务表（4 基础 + challenge1~12）按方言存在（resumes 表由其他用例间接覆盖） |
| `TestDatabase_Seed` | `users`/`emails` 各 8 条；`id=1` 为 `Dumb/Dumb` |
| `TestDatabase_Reset` | 重置后数据恢复为种子状态 |
| `TestDatabase_Query_LIMIT` | `LIMIT 1 OFFSET 0` 语法三库可用（回归保护） |

**测试矩阵开关**：默认只跑 SQLite（测试库 `src/db/test_data/test_sqli_lab.db`）；设置以下环境变量即自动加入 MySQL / PostgreSQL：

```
TEST_MYSQL_HOST / TEST_MYSQL_PORT(3306) / TEST_MYSQL_USER(root) /
TEST_MYSQL_PASSWORD / TEST_MYSQL_DB(test_sqli_lab)

TEST_POSTGRES_HOST / TEST_POSTGRES_PORT(5432) / TEST_POSTGRES_USER(postgres) /
TEST_POSTGRES_PASSWORD / TEST_POSTGRES_DB(test_sqli_lab)
```

本地执行：

```bash
go test -v ./...                                  # SQLite 全量
go test -v ./src/db/...                           # 仅数据库层
go test -cover ./...                              # 覆盖率
# 加入 MySQL/PostgreSQL（PowerShell 同理，用 $env:TEST_MYSQL_HOST=...）
TEST_MYSQL_HOST=127.0.0.1 TEST_MYSQL_PASSWORD=secret go test -v ./src/db/...
```

### 9.2 CI（.github/workflows/ci.yml）

每次 push / PR（`master`、`main` 分支）自动执行：

- `test-sqlite`：`go mod tidy` → `go test ./src/db/...` → `go build` → `--setup-db` 后启动并探测 `/health`；
- `test-mysql` / `test-postgres`：以 GitHub Actions `services` 起 MySQL 8.0 / PostgreSQL 15 容器，注入 `TEST_*` 环境变量后跑 `go test ./src/db/...`；
- `lint`：golangci-lint（v2.13.2，规则见 `.golangci.yml`）。

提交前请本地跑齐：`go fmt ./...`、`go vet ./...`、`go test ./...`。

## 10. 常用命令

### Makefile（类 Unix 环境）

| 目标 | 等价手动命令 | 说明 |
| --- | --- | --- |
| `make build` | `go build -o bin/go-sqli-lab src/main.go` | 构建二进制 |
| `make run` | `go run src/main.go` | 直接运行 |
| `make run-setup` | `go run src/main.go --setup-db` | 初始化数据库并运行（注意：`--setup-db` 后即退出，不进入 Web 模式） |
| `make setup-db` | `go run src/main.go --setup-db` | 仅初始化数据库 |
| `make deps` | `go mod tidy` | 整理依赖 |
| `make test` / `make test-coverage` | `go test -v ./...` / `go test -cover ./...` | 测试 |
| `make fmt` / `make vet` | `go fmt ./...` / `go vet ./...` | 静态检查 |
| `make logs` | `tail -f logs/app.log` | 跟踪运行日志 |
| `make clean` | — | 删除 `bin/`、`logs/`、`tmp/`（谨慎：清空练习记录） |
| `make help` | — | 目标列表 |

Windows 无 Make 时直接使用左侧 `go` 命令即可（PowerShell）。

### 启动脚本（scripts/）

`start-linux.sh` / `start-macos.command` / `start-windows.bat` 的行为一致：检测 Go → `go build` 到 `bin/` → 执行 `--setup-db` 初始化 → 启动服务并自动打开浏览器（设置 `SQLI_LAB_NO_BROWSER=1` 可禁用自动打开）。脚本内路径均相对仓库根目录，请从根目录调用。

## 11. 代码规范与贡献指南

### 规范要点

- 全部代码遵循 `gofmt` 风格，并通过 `go vet` 与 golangci-lint；
- 关卡内**刻意保留 SQL 拼接漏洞**，但必须有注释标明漏洞成因与类型；公共/基础设施代码（config、db、logger、handlers）不允许出现注入写法；
- 中文注释与用户可见文案保持中文；提交信息建议遵循仓库既有风格；
- ⚠️ **`ref/sqli-labs/` 为只读子模块：禁止修改 / 新增 / 删除 / 重命名其中任何文件**；仅可阅读参考，移植实现一律写入 `src/vulnerabilities/`。

### 提交流程

1. `git clone --recursive`（或 clone 后 `git submodule update --init --recursive`）；
2. 新建分支：`git checkout -b feature/less-new` 或 `fix/xxx`；
3. 本地验证：`make fmt vet test`（涉及多库改动请按 9.1 补跑 MySQL/PG）；
4. 文档同步：功能/关卡变化必须同步 `README.md` 关卡地图、相关 `doc/*.md`，并在 `CHANGELOG.md` 按语义化版本追加条目；
5. 提交并推送，创建 Pull Request（描述中说明：漏洞类型、对应表、测试范围）。

### 发布（版本号同步约定）

版本号同步维护在 `CHANGELOG.md`（如 `## v1.5.0 (2026-09-07)`）与 `src/main.go` 启动日志的 `"version"` 字段中，发版时两处需一致；若以 tag 发布，tag 名与二者保持一致。
