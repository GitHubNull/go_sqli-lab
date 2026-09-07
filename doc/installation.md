# 安装指南

## 系统要求

- **操作系统**: Windows 10+ / Linux / macOS
- **Go 版本**: 1.25 或更高（由 [go.mod](../go.mod) 的 `go 1.25.0` 指令指定，旧版本无法编译）
- **内存**: 至少 512MB
- **磁盘**: 至少 200MB 可用空间（含 Go 工具链与依赖缓存）
- **数据库**: SQLite 为默认嵌入式数据库（**无需额外安装**）；MySQL / PostgreSQL 可选

## 快速安装

### 1. 获取源码

```bash
# 克隆仓库（建议带上 --recursive 同时拉取只读参考子模块）
git clone --recursive https://github.com/yourusername/go-sqli-lab.git   # 请替换为实际仓库地址
cd go-sqli-lab

# 若克隆时未带子模块，执行：
git submodule update --init --recursive
```

> `ref/sqli-labs/` 为只读参考子模块（原始 PHP 实现），缺省不影响程序运行，但会影响对照学习，建议拉取。

### 2. 安装依赖并构建

```bash
go mod tidy
make build          # 产物输出到 bin/go-sqli-lab
```

> Windows 无 make 时可直接执行：`go build -o bin\go-sqli-lab.exe src\main.go`

### 3. 初始化数据库（幂等，可重复执行）

```bash
make setup-db
# 或直接执行：
go run src/main.go --setup-db
```

初始化会创建全部业务表（`users`、`emails`、`uagents`、`referers`、`challenge1~12`、`resumes`），并插入基础种子数据（`users`/`emails` 各 8 条、`resumes` 简历数据）；已有数据时自动跳过，不会重复插入。挑战表（`challenge1~12`）的数据在“重置数据库”时注入，见下方 FAQ。

### 4. 启动服务

```bash
make run
# 或直接执行：
go run src/main.go
```

访问 <http://localhost:8080> 即可进入首页。服务支持 Ctrl+C 优雅关闭。

> 程序启动时会自动向上查找 `go.mod` 与 `config/app.yaml` 定位项目根目录并切换工作目录，因此**从任意目录直接运行已编译的二进制**（如下方的 `bin/go-sqli-lab`）也能正常工作。

---

## 各平台安装说明

三种平台均提供一键启动脚本：**自动检测 Go → 编译到 `bin/` → 初始化数据库（幂等）→ 启动服务 → 自动打开浏览器**。设置环境变量 `SQLI_LAB_NO_BROWSER=1` 可禁止自动打开浏览器。

### Windows

方式一：双击启动脚本

1. 安装 Go 1.25+ 并确保 `go` 已加入 PATH（安装包默认勾选）。
2. 双击 [`scripts/start-windows.bat`](../scripts/start-windows.bat)。
3. 脚本自动编译、初始化数据库并启动服务，浏览器自动打开 <http://localhost:8080>。

方式二：使用预编译程序（无需 Go 环境）

1. 将 `bin/go-sqli-lab.exe` 与项目目录（含 `config/`、`src/web/` 等）一起分发到目标机器。
2. 在项目根目录执行：

```bat
bin\go-sqli-lab.exe --setup-db
bin\go-sqli-lab.exe
```

方式三：源码运行（PowerShell / CMD）

```powershell
go mod tidy
go run src/main.go --setup-db
go run src/main.go
```

### macOS

1. 安装 Go 1.25+：`brew install go` 或从 <https://go.dev/dl/> 下载 pkg 安装包。
2. 首次使用需给启动脚本执行权限：

```bash
chmod +x scripts/start-macos.command
```

3. 双击 `scripts/start-macos.command`（或在终端执行 `./scripts/start-macos.command`），自动完成编译、初始化并打开浏览器。

### Linux

1. 安装 Go 1.25+（发行版软件源版本过旧时请从 <https://go.dev/dl/> 安装官方二进制包）。
2. 执行启动脚本：

```bash
chmod +x scripts/start-linux.sh
./scripts/start-linux.sh
```

3. 服务启动于 <http://localhost:8080>，脚本依赖 `xdg-open` 自动打开浏览器（无桌面环境时设置 `SQLI_LAB_NO_BROWSER=1`）。

---

## 数据库配置

### SQLite（默认，零配置）

无需任何额外配置，数据库文件位于 `./data/sqli_lab.db`（首次初始化自动创建目录）。

### MySQL

1. 创建数据库（需已存在，程序不会自动建库）：

```sql
CREATE DATABASE IF NOT EXISTS security DEFAULT CHARSET utf8mb4;
```

2. 修改 [config/app.yaml](../config/app.yaml)：

```yaml
database:
  type: mysql
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: security
```

或通过环境变量覆盖（优先级高于 YAML）：

```bash
export SQLI_LAB_DB_TYPE=mysql
export SQLI_LAB_DB_HOST=localhost
export SQLI_LAB_DB_PORT=3306
export SQLI_LAB_DB_USER=root
export SQLI_LAB_DB_PASSWORD=your_password
export SQLI_LAB_DB_NAME=security
```

3. 重新执行 `--setup-db` 完成建表与种子数据。

### PostgreSQL

1. 创建数据库（程序不会自动建库）：

```sql
CREATE DATABASE security;
```

2. 修改配置：

```yaml
database:
  type: postgres
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: security
  sslmode: disable   # PostgreSQL 专用，按需改为 require
```

或设置环境变量 `SQLI_LAB_DB_TYPE=postgres` 及上述 `SQLI_LAB_DB_*` 系列。

3. 重新执行 `--setup-db`。

> 说明：程序对三种数据库的建表 SQL、清表语句（SQLite `DELETE` / MySQL `TRUNCATE` / PostgreSQL `TRUNCATE ... RESTART IDENTITY`）与参数占位符（`?` / `$N`）已做自动适配，关卡 SQL 均使用三种数据库通用的语法（如 `LIMIT 1 OFFSET 0`）。个别利用手法（如时间盲注的 `SLEEP()`）依赖数据库函数，详见 [usage.md](usage.md)。

---

## 容器化部署

仓库**当前未提供 Dockerfile**，`make docker-build` 目标需要先自行补充：

```dockerfile
# 以多阶段构建示例
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go build -o /go-sqli-lab src/main.go

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /go-sqli-lab .
COPY config ./config
COPY src/web ./src/web
EXPOSE 8080
CMD ["./go-sqli-lab"]
```

```bash
docker build -t go-sqli-lab .
docker run -p 8080:8080 go-sqli-lab
```

> 注意：镜像内必须包含 `config/app.yaml` 与 `src/web/`（HTML 模板与静态资源依赖相对项目根目录的路径）。SQLite 模式下建议挂载 `data/` 卷持久化数据库。

---

## 常见问题（FAQ）

### 1. 提示 "Go not found" / "未找到 Go 环境"

- 确认已安装 Go 1.25+ 且 `go version` 可正常输出；
- Windows 安装后需**重新打开**终端/资源管理器使 PATH 生效；
- 也可以直接使用 `bin/` 下已编译的二进制，免去 Go 环境依赖。

### 2. 端口被占用

修改 `config/app.yaml` 中 `server.port`，或临时覆盖：

```bash
export SQLI_LAB_PORT=8081
```

### 3. 数据库连接失败（MySQL/PostgreSQL）

- 确认数据库服务已启动、库已创建（`dbname`）；
- 核对账号密码与 `sslmode`（PostgreSQL）；
- 查看日志定位原因：`logs/app.log`（或 `tail -f logs/app.log`）。

### 4. 启动后页面样式/关卡 404

- 程序必须在**项目根目录结构完整**的环境运行（含 `src/web/` 与 `config/`）；程序自动定位根目录的依据是 `go.mod` + `config/app.yaml` 同目录，请勿拆分这两个文件。
- 刷新浏览器缓存后再试。

### 5. 挑战关卡/导出/上传相关报错

- 使用前请先执行 `--setup-db` 初始化（简历导出依赖 `resumes` 表，该表会随初始化自动填充）。
- 挑战关卡（Less-54~65）的 `challenge1~12` 数据只在“重置数据库”时注入：`--setup-db` 只建表不插挑战数据，首次练习前请先重置一次（首页「🔄 重置数据库」或访问 `http://localhost:8080/setup-db?reset=true`）。
- 关卡打到一半数据"打脏"了？用首页右上角「🔄 重置数据库」一键还原全部表数据。

### 6. Windows 下 `make` 不可用

使用 `scripts/start-windows.bat`，或直接使用 `go build` / `go run` 命令（见上文 Windows 方式三）。

## 卸载

- **脚本/源码运行**：直接删除项目目录即可。
- **残留数据**（如需彻底清理）：

```bash
# Linux / macOS
rm -rf data/ logs/ tmp/

# Windows（PowerShell）
Remove-Item -Recurse -Force data, logs, tmp
```

- **自建数据库**（MySQL/PostgreSQL）：按需 `DROP DATABASE security;`。
