# Agent Instructions

## Critical Rules for AI Agents

### READ-ONLY Dependencies

#### ref/sqli-labs (Git Submodule)

**⚠️ CRITICAL: This directory is a READ-ONLY dependency**

- **Path**: `ref/sqli-labs/`
- **Type**: Git submodule (external dependency)
- **Source**: https://github.com/Audi-1/sqli-labs

**FORBIDDEN ACTIONS**:
- ❌ NEVER modify any files in `ref/sqli-labs/`
- ❌ NEVER delete files in `ref/sqli-labs/`
- ❌ NEVER add new files to `ref/sqli-labs/`
- ❌ NEVER rename files in `ref/sqli-labs/`

**ALLOWED ACTIONS**:
- ✅ Read and reference files for implementation
- ✅ Study the PHP code to understand vulnerability patterns
- ✅ Use as reference when porting to Go

**HOW TO UPDATE**:
Only sync updates when the upstream repository changes:

```bash
# Fetch latest changes from upstream
git submodule update --init --remote ref/sqli-labs

# Or update all submodules
git submodule update --init --remote --recursive
```

**WHY THIS MATTERS**:
This is the reference PHP implementation of sqli-labs. We are porting it to Go in the `src/` directory. Keeping it as a git submodule ensures:
1. We can always reference the original implementation
2. We can track upstream updates easily
3. We maintain a clean separation between our code and third-party code

---

## Project Structure Rules

| 路径 | 职责 |
|------|------|
| `src/` | 全部 Go 源码 |
| `src/main.go` | 程序入口：flag 解析、项目根目录自动定位、配置/日志/数据库初始化、Gin 路由注册、优雅关闭 |
| `src/config/` | 配置加载：内置默认值 < `config/app.yaml` < `SQLI_LAB_*` 环境变量 |
| `src/db/` | 数据库层：`Database` 接口（GetDB/Migrate/Seed/Reset/Query/QueryRow/Exec），SQLite/MySQL/PostgreSQL 方言适配 |
| `src/handlers/` | 通用 HTTP 处理器：健康检查、db-status、setup-db 页面、reset-db AJAX、结果页 |
| `src/logger/` | 结构化日志（lumberjack 滚动）+ Gin HTTP 请求中间件 |
| `src/models/` | 数据模型（User/Email/UAgent/Referer） |
| `src/vulnerabilities/` | **所有漏洞关卡**：`less1.go` ~ `less71.go` + 共享逻辑文件 + `registry.go` 注册表 |
| `src/web/` | 静态资源：`index.html`（关卡入口）、css、assets、`less-N/` 关卡静态资源目录 |
| `config/` | 配置文件（`app.yaml`） |
| `scripts/` | 跨平台一键启动脚本（start-windows.bat / start-macos.command / start-linux.sh） |
| `doc/` | 文档：installation.md / usage.md / development.md |
| `data/` `logs/` `tmp/` | 运行时目录：数据库文件、日志、关卡请求记录（`tmp/less-N/result.txt`） |
| `tools/genresume/` | 简历 XLSX 模板生成工具（配合 Less-68/69 的 `assets/resume_template.xlsx`） |
| `ref/sqli-labs/` | 只读参考（Git 子模块），禁止任何写入 |

## Vulnerability Module Development Guide

### 接口契约（src/vulnerabilities/vulnerability.go）

每个关卡实现 `Vulnerability` 接口（全部方法必须实现）：

```go
type Vulnerability interface {
    ID() string          // 唯一 ID，如 "less-71"
    Name() string        // 显示名，如 "Less-71: Upload XLSX Injection"
    Description() string // 中文漏洞描述
    Category() string    // 分类（Error Based / Blind Injection / POST Injection / WAF Bypass 等）
    Route(r *gin.RouterGroup) // 在 /less 组下挂载本关路由
}
```

### 新增关卡的标准文件链路

新增/复制一个漏洞关卡时必须同步完成以下全部步骤，缺任一步关卡都不会在首页出现或无法访问：

1. **`src/vulnerabilities/<lesson>.go`**：定义 struct（通常嵌入 `BaseLesson` 获得 `db db.Database` 与 `logger`），实现接口方法；`Route()` 内挂载页面入口与操作接口（参考 Less-70 的 GET `/template` + POST `/upload` 多端点模式）。
2. **共享逻辑抽到同级独立文件**：多个关卡复用的解析/渲染/查询函数放入 `export.go` / `upload.go` / `resume_export.go` 这类公共文件，不要在 each lesson 中重复。
3. **`src/vulnerabilities/registry.go`**：在 `RegisterAll()` 中追加对应的 `r.registerLessN()` 调用（并定义 `registerLessN()`），否则路由不会被注册。
4. **`src/web/index.html`**：按分类区块为关卡添加 `<a href="/less/less-N">` 链接。
5. **`src/web/css/style.css`**：为分类区块新增 `.category-<name>` 样式规则（首页分类配色依赖该规则）。

### 关卡实现约定

- **调用 `l.logRequest(c, "less-N")`**：请求会以追加方式记录到 `tmp/less-N/result.txt`（BaseLesson 自动建目录）。
- **统一用 `l.renderHTML(...)` / 各关自带的表单渲染函数** 输出页面（Go 代码直接输出 HTML，不是 gohtml 模板），复用 `/static/css/style.css`。
- **漏洞点写法**：用 `fmt.Sprintf` 直接拼接用户输入构造 SQL（这正是靶场要保留的漏洞），再通过 `l.db.QueryRow/Query/Exec` 执行。恶意拼接到参数名或 SQL 结构（如表名/排序字段）时同样如此。
- **业务合理性**：POST 类关卡先渲染登录/注册表单，登录验证通过后才进入有漏洞的 INSERT/UPDATE（参考 Less-17 改密、Less-18/19 头注入）。
- **模板等静态二进制资源**：放入 `src/vulnerabilities/assets/`，用 `//go:embed` 嵌入（参考 `resume_template.xlsx` 与 `upload_template.xls`），运行时经 handler 以 `Content-Disposition: attachment` 下发。
- **导入新依赖**：改动 `go.mod` 后运行 `go mod tidy`；注意 `.xls` 解析使用 `github.com/extrame/xls`（仅接受真正的二进制 .xls，不接受 HTML 改扩展名的伪 .xls）。
- **文档同步**：新增/变更关卡后同步更新 `src/web/index.html`、README 关卡地图、`doc/usage.md`（如需）与 CHANGELOG.md。

## Database Layer (src/db/db.go)

### 多数据库适配约定

- 数据库类型由 `config.app.yaml` 的 `database.type` 决定：`sqlite`（默认）/ `mysql` / `postgres`，驱动分别为 modernc.org/sqlite、go-sql-driver/mysql、lib/pq。
- **SQL 方言分支**：`Migrate()` 中所有建表语句按数据库类型三路分支（AUTOINCREMENT vs SERIAL vs AUTO_INCREMENT、utf8mb4 引擎、TRUNCATE vs TRUNCATE RESTART IDENTITY vs DELETE）。新增表时必须为三种方言各写一份。
- **参数占位符**：SQLite/MySQL 用 `?`，PostgreSQL 用 `$1, $2...`。写公共代码时按 `cfg.Type == "postgres"` 分支生成占位符（参考 `seedResumes` 的动态占位符生成）。
- **跨库 SQL 兼容**：查询统一使用 `LIMIT 1 OFFSET 0` 写法（**禁止 `LIMIT 0,1`**，PostgreSQL 不兼容）；`db_test.go` 中 `TestDatabase_Query_LIMIT` 专门守护该语法。
- **关卡 SQL 不得参数化**：漏洞代码刻意使用 `fmt.Sprintf` 拼接，此与常规开发相反，是靶场特性。

### Migrate / Seed / Reset 语义

- `Migrate()`：幂等建表 —— `users`、`emails`、`uagents`、`referers`、`challenge1`~`challenge12`（挑战关）、`resumes`（Less-68/69 简历，20 个字段）。
- `Seed()`：幂等种数据 —— `resumes` 独立判断空表才插入；`users`/`emails` 仅在 users 表为空时插入（8 个用户 + 8 个邮箱，账号含 admin/admin）。
- `Reset()`：清空全部表（含 challenge 表与 resumes）→ `forceSeed()` 强制重种 → `seedChallenges()` 按每表 `3 + 表号%3` 条插入挑战数据（共 3~5 条，用户名形如 challenge_user1）。**注意：挑战表数据只在 Reset 流程注入，`--setup-db`（Migrate+Seed）不插入挑战数据。**
- **修改数据模型时**：新增/变更表必须同步 `Migrate()`、`Seed()`、`Reset()`（及其内部清表清单与挑战种子）三处逻辑，并更新 `handlers.go` 中的重置统计（ClearedTables/InsertedRecords）与 `db_test.go` 断言（users 8 条、表清单）。

### 测试与 CI

- 单元/集成测试在 `src/db/db_test.go`：默认跑 SQLite（`./test_data/test_sqli_lab.db`）；设置 `TEST_MYSQL_HOST` / `TEST_POSTGRES_HOST` 环境变量后自动追加 MySQL/PostgreSQL 用例。
- CI（`.github/workflows/ci.yml`）：`go test ./src/db/...`（SQLite job 固定跑，MySQL/PostgreSQL 用 services 容器），外加 `golangci-lint` 与 `/health` 冒烟测试。

## Logger Best Practices

- 通过 `logger.New(cfg.Log)` 创建（`output: file` 时自动 lumberjack 滚动；`stdout` 输出到控制台），实现注入到 `db`/各 handler；**禁止**在业务代码中直接 `fmt.Println` 打日志。
- 级别与格式在 `config/app.yaml` 的 `log` 段配置（`%timestamp% %level% %goroutine% %module% %function% %file% %line% %message%` 占位符）。
- HTTP 请求统一走 `logger.GinLogger(log)` 中间件；漏洞关卡内用结构化 KV 调用 `l.logger.Info/Debug/Error`。
- 关卡请求明细（URL/参数）写入 `tmp/<lesson>/result.txt`（`BaseLesson.logRequest`），这是靶场"结果记录"能力，勿随意改动路径约定（相对项目根目录）。

## Config Management Best Practices

- 配置结构见 `src/config/config.go`（Server/Database/Log 三组），默认值在 `DefaultConfig()`。
- 新增配置项需同步四处：`config.go` 结构体与默认值 → `applyEnvOverrides()`（若需环境变量）→ `config/app.yaml` 示例 → 相关文档。
- 应用启动时会自动向上查找 `go.mod` + `config/app.yaml` 定位项目根目录并 `chdir`，因此代码中相对路径（`./data/...`、`./logs/...`、`./tmp/...`、`./src/web`）一律相对项目根目录书写；勿假设当前工作目录。

## Code Style & Workflow

- 代码格式：`go fmt`；静态检查：`go vet`、golangci-lint（`.golangci.yml` 已在 CI 启用）。
- 注释与页面文案使用中文；代码标识符保持英文（Go 风格）。
- 测试：`make test`（`go test -v ./...`），数据库相关改动至少保证 `go test ./src/db/...` 通过。
- 保持与既有关卡一致的实现模式（参考同类型最近的关卡文件），不要引入新的页面渲染方式或重复造 helper。
- **提交代码必须使用 smart-commit 技能**，禁止手动 `git commit`；版本号变更时必须校验 JSON 语法与 CHANGELOG 使用真实时间戳（详见仓库约定）。
- 维护 README/doc 与代码事实一致：文档描述的每个功能点都应在代码中有对应实现；废弃功能需从文档移除。
