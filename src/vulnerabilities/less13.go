package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less13 POST - Double Injection - Single Quotes
type Less13 struct {
	BaseLesson
}

func (l *Less13) ID() string {
	return "less-13"
}

func (l *Less13) Name() string {
	return "Less-13: POST - Double Injection - Single Quotes"
}

func (l *Less13) Description() string {
	return "POST请求 - 双查询注入 (单引号)"
}

func (l *Less13) Category() string {
	return "POST Injection"
}

func (l *Less13) Route(r *gin.RouterGroup) {
	group := r.Group("/less-13")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
	}
}

func (l *Less13) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-13")

	uname := c.PostForm("uname")
	passwd := c.PostForm("passwd")

	if uname == "" && passwd == "" {
		l.renderLoginForm(c)
		return
	}

	query := fmt.Sprintf("SELECT username, password FROM users WHERE username='%s' and password='%s' LIMIT 1 OFFSET 0", uname, passwd)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderLoginFormWithError(c, "less-13", "登录失败")
		} else {
			l.renderLoginFormWithError(c, "less-13", err.Error())
		}
		return
	}

	l.renderLoginFormWithResult(c, "less-13", "You are in...........")
}

func (l *Less13) renderLoginForm(c *gin.Context) {
	l.renderLoginFormWithError(c, "less-13", "")
}

func (l *Less13) renderLoginFormWithError(c *gin.Context, lesson, err string) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<div class="container">
		<h1>%s</h1>
		<form method="POST" action="">
			<div>
				<label>Username:</label>
				<input type="text" name="uname" value=""/>
			</div>
			<div>
				<label>Password:</label>
				<input type="password" name="passwd" value=""/>
			</div>
			<div>
				<input type="submit" value="Submit"/>
			</div>
		</form>
		%s
		<div class="nav">
			<a href="/">返回首页</a>
		</div>
	</div>
</body>
</html>`, lesson, lesson, renderError(err))
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}

func (l *Less13) renderLoginFormWithResult(c *gin.Context, lesson, result string) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<div class="container">
		<h1>%s</h1>
		%s
		<div class="nav">
			<a href="/">返回首页</a>
		</div>
	</div>
</body>
</html>`, lesson, lesson, renderResult(result))
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}

func (r *Registry) registerLess13() {
	r.Register(&Less13{BaseLesson{db: r.db, logger: r.logger}})
}
