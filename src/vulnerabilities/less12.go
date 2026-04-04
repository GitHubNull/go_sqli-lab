package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less12 POST - Error Based - Double Quotes
type Less12 struct {
	BaseLesson
}

func (l *Less12) ID() string {
	return "less-12"
}

func (l *Less12) Name() string {
	return "Less-12: POST - Error Based - Double Quotes"
}

func (l *Less12) Description() string {
	return "POST请求 - 基于错误信息的SQL注入 (双引号)"
}

func (l *Less12) Category() string {
	return "POST Injection"
}

func (l *Less12) Route(r *gin.RouterGroup) {
	group := r.Group("/less-12")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
	}
}

func (l *Less12) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-12")

	uname := c.PostForm("uname")
	passwd := c.PostForm("passwd")

	if uname == "" && passwd == "" {
		l.renderLoginForm(c)
		return
	}

	query := fmt.Sprintf("SELECT username, password FROM users WHERE username=\"%s\" and password=\"%s\" LIMIT 1 OFFSET 0", uname, passwd)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderLoginFormWithError(c, "less-12", "登录失败")
		} else {
			l.renderLoginFormWithError(c, "less-12", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderLoginFormWithResult(c, "less-12", result)
}

func (l *Less12) renderLoginForm(c *gin.Context) {
	l.renderLoginFormWithError(c, "less-12", "")
}

func (l *Less12) renderLoginFormWithError(c *gin.Context, lesson, err string) {
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

func (l *Less12) renderLoginFormWithResult(c *gin.Context, lesson, result string) {
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

func (r *Registry) registerLess12() {
	r.Register(&Less12{BaseLesson{db: r.db, logger: r.logger}})
}
