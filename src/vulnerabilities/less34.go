package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less34 POST - Bypass addslashes()
type Less34 struct {
	BaseLesson
}

func (l *Less34) ID() string {
	return "less-34"
}

func (l *Less34) Name() string {
	return "Less-34: POST - Bypass addslashes()"
}

func (l *Less34) Description() string {
	return "POST请求 - 绕过addslashes()"
}

func (l *Less34) Category() string {
	return "WAF Bypass"
}

func (l *Less34) Route(r *gin.RouterGroup) {
	group := r.Group("/less-34")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
	}
}

func (l *Less34) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-34")

	uname := c.PostForm("uname")
	passwd := c.PostForm("passwd")

	if uname == "" && passwd == "" {
		l.renderLoginForm(c, "less-34", "")
		return
	}

	query := fmt.Sprintf("SELECT username, password FROM users WHERE username='%s' and password='%s' LIMIT 0,1", uname, passwd)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderLoginForm(c, "less-34", "登录失败")
		} else {
			l.renderLoginForm(c, "less-34", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderLoginForm(c, "less-34", result)
}

func (l *Less34) renderLoginForm(c *gin.Context, lesson, result string) {
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
</html>`, lesson, lesson, renderResult(result))
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}

func (r *Registry) registerLess34() {
	r.Register(&Less34{BaseLesson{db: r.db, logger: r.logger}})
}
