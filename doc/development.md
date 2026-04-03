# 开发文档

## 项目结构

```
go-sqli-lab/
├── src/                        # Go 源代码
│   ├── main.go                # 程序入口
│   ├── config/                # 配置管理
│   │   └── config.go         # 配置加载和解析
│   ├── db/                    # 数据库层
│   │   └── db.go             # 数据库连接和操作
│   ├── handlers/              # HTTP 处理器
│   │   └── handlers.go       # 路由处理器
│   ├── logger/                # 日志系统
│   │   ├── logger.go         # 日志实现
│   │   └── gin.go            # Gin 中间件
│   ├── models/                # 数据模型
│   │   └── user.go           # 用户模型
│   ├── vulnerabilities/       # 漏洞实现
│   │   ├── vulnerability.go   # 漏洞接口
│   │   ├── registry.go        # 漏洞注册表
│   │   ├── base.go            # 基础结构
│   │   └── less*.go           # 65个漏洞级别
│   └── web/                   # 静态资源
│       ├── css/               # 样式文件
│       ├── less-1/            # 关卡 1 的独立资源
│       ├── less-2/            # 关卡 2 的独立资源
│       ├── ...                # 其他关卡资源
│       ├── less-65/           # 关卡 65 的独立资源
│       └── index.html         # 主页
├── config/                     # 配置文件
│   └── app.yaml              # 配置文件
├── data/                       # SQLite 数据库
├── logs/                       # 日志目录
├── tmp/                        # 临时目录
├── ref/                        # 参考项目（Git 子模块）
│   └── sqli-labs/             # 原始 PHP 实现（只读）
└── doc/                        # 文档目录
```

## Git 子模块说明

### ref/sqli-labs（只读依赖）

`ref/sqli-labs/` 目录是一个 **Git 子模块**，指向原始 PHP 实现 [sqli-labs](https://github.com/Audi-1/sqli-labs)。

**⚠️ 重要：此目录为只读，永远不要修改！**

#### 子模块信息
- **路径**: `ref/sqli-labs/`
- **来源**: https://github.com/Audi-1/sqli-labs
- **用途**: 作为参考实现，用于理解漏洞模式并移植到 Go
- **更新方式**: 仅通过 Git 子模块命令同步上游更新

#### 克隆包含子模块的项目

```bash
# 方式一：克隆时包含子模块
git clone --recursive https://github.com/yourusername/go-sqli-lab.git

# 方式二：克隆后初始化子模块
git clone https://github.com/yourusername/go-sqli-lab.git
cd go-sqli-lab
git submodule update --init --recursive
```

#### 同步上游更新

仅当原始 sqli-labs 项目更新时才执行：

```bash
# 获取上游最新版本
git submodule update --init --remote ref/sqli-labs

# 或者更新所有子模块
git submodule update --init --remote --recursive
```

#### 禁止的操作

在 `ref/sqli-labs/` 目录中：
- ❌ **永远不要修改任何文件**
- ❌ **永远不要添加新文件**
- ❌ **永远不要删除文件**
- ❌ **永远不要重命名文件**

如需参考代码，请：
- ✅ 在 `src/vulnerabilities/` 中创建新的 Go 实现
- ✅ 在 `src/web/less-X/` 中创建独立的前端资源
- ✅ 阅读并理解参考代码的漏洞逻辑

---

## 添加新漏洞关卡

### 1. 创建漏洞文件

在 `src/vulnerabilities/` 下创建 `lessX.go`：

```go
package vulnerabilities

import (
	"database/sql"
	"fmt"
	"go-sqli-lab/src/models"
	"github.com/gin-gonic/gin"
)

// LessX 描述你的漏洞类型
type LessX struct {
	BaseLesson
}

func (l *LessX) ID() string {
	return "less-x"
}

func (l *LessX) Name() string {
	return "Less-X: 漏洞名称"
}

func (l *LessX) Description() string {
	return "漏洞描述"
}

func (l *LessX) Category() string {
	return "分类"  // 如: Error Based, Blind, WAF Bypass 等
}

func (l *LessX) Route(r *gin.RouterGroup) {
	group := r.Group("/less-x")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		// 添加其他路由
	}
}

func (l *LessX) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-x")

	// 获取参数
	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-x", "提示信息", "")
		return
	}

	// 构造有漏洞的SQL查询
	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-x", "", "未找到记录")
		} else {
			l.renderHTML(c, "less-x", "", err.Error())  // 显示SQL错误
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", 
		user.Username, user.Password)
	l.renderHTML(c, "less-x", result, "")
}

func (r *Registry) registerLessX() {
	r.Register(&LessX{BaseLesson{db: r.db, logger: r.logger}})
}
```

### 2. 注册漏洞

在 `registry.go` 中的 `RegisterAll` 方法中添加：

```go
func (r *Registry) RegisterAll(router *gin.Engine) {
	// ... 其他注册
	r.registerLessX()  // 添加新行
	
	// 注册路由
	vulnGroup := router.Group("/less")
	for _, v := range r.vulnerabilities {
		v.Route(vulnGroup)
	}
}
```

### 3. 创建前端资源

在 `src/web/less-X/` 目录下创建独立的静态资源：

```
src/web/less-X/
├── index.html    # 关卡页面
├── style.css     # 关卡样式
└── script.js     # 关卡脚本
```

### 4. 更新主页

在 `src/web/index.html` 中添加链接：

```html
<a href="/less/less-x?id=1" class="category-error">Less-X: 描述</a>
```

---

## 数据库操作

### 执行查询

```go
// 单行查询
var user models.User
err := l.db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&user.ID, &user.Username, &user.Password)

// 多行查询
rows, err := l.db.Query("SELECT * FROM users")
defer rows.Close()
for rows.Next() {
	var user models.User
	rows.Scan(&user.ID, &user.Username, &user.Password)
}

// 执行语句
result, err := l.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
```

### 事务支持

```go
tx, err := l.db.GetDB().Begin()
// ... 执行操作
tx.Commit()  // 或 tx.Rollback()
```

---

## 日志记录

### 在漏洞中使用日志

```go
func (l *LessX) handleIndex(c *gin.Context) {
	// 记录请求
	l.logRequest(c, "less-x")
	
	// 使用 logger
	l.logger.Info("处理请求", "ip", c.ClientIP())
	l.logger.Debug("SQL查询", "query", query)
	l.logger.Error("发生错误", "error", err)
}
```

---

## 前端开发

### 使用模板

基础HTML模板：

```html
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>{{.Title}}</title>
	<link rel="stylesheet" href="/static/css/common.css">
	<link rel="stylesheet" href="/static/less-X/style.css">
</head>
<body>
	<div class="container">
		<h1>{{.Title}}</h1>
		<!-- 内容 -->
		<div class="nav">
			<a href="/">返回首页</a>
		</div>
	</div>
	<script src="/static/less-X/script.js"></script>
</body>
</html>
```

### 添加CSS样式

每个关卡使用独立的 `style.css`，公共样式放在 `src/web/css/common.css`。

---

## 测试

### 运行测试

```bash
go test ./...
```

### 添加单元测试

创建 `src/vulnerabilities/lessX_test.go`：

```go
package vulnerabilities

import (
	"testing"
	"github.com/gin-gonic/gin"
)

func TestLessX(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// 测试代码
}
```

---

## 贡献指南

1. Fork 项目
2. 创建分支 (`git checkout -b feature/new-lesson`)
3. 提交更改 (`git commit -am 'Add new lesson'`)
4. 推送分支 (`git push origin feature/new-lesson`)
5. 创建 Pull Request

## 代码规范

- 使用 `go fmt` 格式化代码
- 添加必要的注释
- 遵循 Go 命名规范
- 编写清晰的提交信息
- **永远不要修改 `ref/sqli-labs/` 中的文件**
