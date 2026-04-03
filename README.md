# Go SQLi Lab

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

一个基于 Go 语言实现的 SQL 注入学习平台，参考经典的 [sqli-labs](https://github.com/Audi-1/sqli-labs) 项目构建。

## 🎯 项目简介

Go SQLi Lab 是一个用于学习和测试 SQL 注入漏洞的 Web 应用程序。它包含了 65 个不同难度和类型的 SQL 注入关卡，涵盖了从基础的 Error-based 注入到高级的 WAF 绕过技术。

### 主要特性

- 🔒 **65 个 SQL 注入关卡** - 从基础到高级的完整学习路径
- 🗄️ **多数据库支持** - 支持 SQLite (默认)、MySQL 和 PostgreSQL
- 📝 **多种注入类型** - Error-based、Blind、Time-based、Stacked Query 等
- 🛡️ **WAF 绕过技术** - 注释过滤、关键字绕过、编码绕过等
- 📊 **模块化设计** - 易于扩展和添加新的漏洞类型
- 🎨 **现代化前端** - 原生 HTML/CSS/JS，响应式设计
- 📝 **完整日志系统** - 类 log4j 格式的滚动日志

## 🚀 快速开始

### 环境要求

- Go 1.21 或更高版本
- SQLite (默认，无需额外安装)
- 或 MySQL/PostgreSQL (可选)

### 安装步骤

1. 克隆项目

```bash
git clone https://github.com/yourusername/go-sqli-lab.git
cd go-sqli-lab
```

2. 安装依赖

```bash
go mod tidy
```

3. 初始化数据库

```bash
go run src/main.go --setup-db
```

4. 启动服务

```bash
go run src/main.go
```

5. 访问应用

打开浏览器访问 http://localhost:8080

### 使用 Makefile

```bash
# 构建项目
make build

# 运行测试
make test

# 启动服务
make run

# 初始化数据库
make setup-db
```

## 📁 项目结构

```
go-sqli-lab/
├── src/                        # Go 源代码
│   ├── main.go                # 程序入口
│   ├── config/                # 配置管理
│   ├── db/                    # 数据库层
│   ├── handlers/              # HTTP 处理器
│   ├── logger/                # 日志系统
│   ├── models/                # 数据模型
│   ├── vulnerabilities/       # 漏洞实现 (Less-1 ~ Less-65)
│   └── web/                   # 静态资源
│       ├── css/               # 样式文件
│       ├── less-1/            # 关卡 1 的独立资源
│       ├── less-2/            # 关卡 2 的独立资源
│       ├── ...                # 其他关卡资源
│       ├── less-65/           # 关卡 65 的独立资源
│       └── index.html         # 主页
├── config/                     # 配置文件
│   └── app.yaml               # 应用配置
├── data/                       # SQLite 数据库文件
├── logs/                       # 日志文件
├── tmp/                        # 临时文件
├── doc/                        # 文档
├── bin/                        # 编译后的二进制文件
├── ref/                        # 参考项目（Git Submodule）
│   └── sqli-labs/             # 原始 PHP 实现（只读）
├── README.md                   # 项目说明
├── LICENSE                     # 开源协议
└── DISCLAIMER.md              # 免责声明
```

### 📦 Git 子模块说明

本项目使用 Git 子模块来管理参考项目 `ref/sqli-labs/`：

- **子模块路径**: `ref/sqli-labs/`
- **来源**: https://github.com/Audi-1/sqli-labs
- **用途**: 提供原始 PHP 实现作为参考
- **权限**: **只读**，永远不要修改其中的文件

**克隆包含子模块的项目**:
```bash
git clone --recursive https://github.com/yourusername/go-sqli-lab.git
```

**如果已克隆但没有子模块内容**:
```bash
git submodule update --init --recursive
```

**同步上游更新（仅当原始项目更新时）**:
```bash
git submodule update --init --remote ref/sqli-labs
```

## ⚙️ 配置说明

配置文件位于 `config/app.yaml`，支持以下配置项：

```yaml
server:
  host: 0.0.0.0
  port: 8080
  mode: debug  # debug 或 release

database:
  type: sqlite     # sqlite, mysql, postgres
  dsn: ./data/sqli_lab.db  # SQLite 数据库路径
  # MySQL/PostgreSQL 配置
  host: localhost
  port: 3306
  user: root
  password: ""
  dbname: security

log:
  level: info
  format: "[%timestamp%] [%level%] [GID:%goroutine%] [%module%] [%function%] [%file%:%line%] - %message%"
  output: file
  filepath: ./logs/app.log
  maxsize: 100      # MB
  maxbackups: 10
  maxage: 30        # 天
  compress: true
```

### 环境变量

也可以通过环境变量覆盖配置：

```bash
export SQLI_LAB_PORT=8080
export SQLI_LAB_DB_TYPE=mysql
export SQLI_LAB_DB_HOST=localhost
export SQLI_LAB_DB_USER=root
export SQLI_LAB_DB_PASSWORD=password
export SQLI_LAB_LOG_LEVEL=debug
```

## 📚 关卡分类

| 分类 | 关卡 | 说明 |
|------|------|------|
| **基础注入** | Less-1 ~ Less-10 | Error-based、Double Injection、Blind |
| **POST 注入** | Less-11 ~ Less-22 | POST 请求、Header、Cookie 注入 |
| **WAF 绕过** | Less-23 ~ Less-37 | 注释过滤、关键字绕过、编码绕过 |
| **堆叠注入** | Less-38 ~ Less-53 | Stacked Query、ORDER BY 注入 |
| **挑战关卡** | Less-54 ~ Less-65 | 综合挑战 |

## 📝 使用说明

### 基础注入示例

访问 Less-1 (Error Based - String):
```
http://localhost:8080/less/less-1?id=1'
```

### POST 注入示例

访问 Less-11 (POST - Error Based):
```
http://localhost:8080/less/less-11
```

在表单中输入:
- Username: `admin' OR '1'='1`
- Password: `anything`

### 数据库初始化

首次使用或重置数据库:
```bash
go run src/main.go --setup-db
```

或使用 API:
```
http://localhost:8080/setup-db
```

## 🔧 开发文档

详见 [doc/development.md](doc/development.md)

### 添加新的漏洞关卡

1. 在 `src/vulnerabilities/` 目录下创建新的 Go 文件
2. 实现 `Vulnerability` 接口
3. 在 `registry.go` 中注册新关卡

示例:
```go
type LessX struct {
    BaseLesson
}

func (l *LessX) ID() string {
    return "less-x"
}

func (l *LessX) Name() string {
    return "Less-X: Description"
}

func (r *Registry) registerLessX() {
    r.Register(&LessX{BaseLesson{db: r.db, logger: r.logger}})
}
```

## ⚠️ 免责声明

**本平台仅供学习和研究使用！**

1. 请勿在未经授权的系统上使用本工具
2. 请勿用于非法用途或恶意攻击
3. 使用本平台造成的任何后果由使用者自行承担
4. 请在合法的环境和授权下进行测试

详见 [DISCLAIMER.md](DISCLAIMER.md)

## 📄 开源协议

本项目采用 [MIT 协议](LICENSE) 开源。

## 🙏 致谢

- 感谢 [sqli-labs](https://github.com/Audi-1/sqli-labs) 项目提供的灵感
- 感谢所有贡献者和安全研究人员

## 📞 联系方式

- GitHub Issues: [提交问题](https://github.com/yourusername/go-sqli-lab/issues)
- Email: your.email@example.com

---

**⚠️ 重要提示: 请合法使用，切勿用于非法用途！**
