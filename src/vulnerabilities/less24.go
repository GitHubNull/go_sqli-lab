package vulnerabilities

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Less24 POST - Second Order Injection
type Less24 struct {
	BaseLesson
}

func (l *Less24) ID() string {
	return "less-24"
}

func (l *Less24) Name() string {
	return "Less-24: POST - Second Order Injection"
}

func (l *Less24) Description() string {
	return "POST请求 - 二次注入"
}

func (l *Less24) Category() string {
	return "Second Order"
}

func (l *Less24) Route(r *gin.RouterGroup) {
	group := r.Group("/less-24")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
		group.GET("/login", l.handleLogin)
		group.GET("/new_user", l.handleNewUser)
		group.POST("/new_user", l.handleCreateUser)
	}
}

func (l *Less24) handleIndex(c *gin.Context) {
	l.renderMainPage(c)
}

func (l *Less24) handleLogin(c *gin.Context) {
	l.renderLoginForm(c)
}

func (l *Less24) handleNewUser(c *gin.Context) {
	l.renderNewUserForm(c, "")
}

func (l *Less24) handleCreateUser(c *gin.Context) {
	l.logRequest(c, "less-24-create")

	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		l.renderNewUserForm(c, "请输入用户名和密码")
		return
	}

	// 有漏洞的INSERT语句
	query := fmt.Sprintf("INSERT INTO users (username, password) VALUES ('%s', '%s')", username, password)
	_, err := l.db.Exec(query)
	if err != nil {
		l.renderNewUserForm(c, err.Error())
		return
	}

	l.renderNewUserForm(c, "用户创建成功")
}

func (l *Less24) renderMainPage(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Less-24: Second Order Injection</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<div class="container">
		<h1>Less-24: 二次注入</h1>
		<div class="nav-links">
			<a href="/less-24/login">登录</a><br>
			<a href="/less-24/new_user">创建新用户</a>
		</div>
		<div class="nav">
			<a href="/">返回首页</a>
		</div>
	</div>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}

func (l *Less24) renderLoginForm(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Less-24: Login</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<div class="container">
		<h1>登录</h1>
		<form method="POST" action="/less-24/">
			<div>
				<label>Username:</label>
				<input type="text" name="uname" value=""/>
			</div>
			<div>
				<label>Password:</label>
				<input type="password" name="passwd" value=""/>
			</div>
			<div>
				<input type="submit" value="Login"/>
			</div>
		</form>
		<div class="nav">
			<a href="/less-24">返回</a>
		</div>
	</div>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}

func (l *Less24) renderNewUserForm(c *gin.Context, msg string) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Less-24: Create User</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<div class="container">
		<h1>创建新用户</h1>
		<form method="POST" action="/less-24/new_user">
			<div>
				<label>Username:</label>
				<input type="text" name="username" value=""/>
			</div>
			<div>
				<label>Password:</label>
				<input type="password" name="password" value=""/>
			</div>
			<div>
				<input type="submit" value="Create"/>
			</div>
		</form>
		%s
		<div class="nav">
			<a href="/less-24">返回</a>
		</div>
	</div>
</body>
</html>`, renderResult(msg))
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}

func (r *Registry) registerLess24() {
	r.Register(&Less24{BaseLesson{db: r.db, logger: r.logger}})
}
