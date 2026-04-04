package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less14 POST - Double Injection - Double Quotes
type Less14 struct {
	BaseLesson
}

func (l *Less14) ID() string {
	return "less-14"
}

func (l *Less14) Name() string {
	return "Less-14: POST - Double Injection - Double Quotes"
}

func (l *Less14) Description() string {
	return "POST请求 - 双查询注入 (双引号)"
}

func (l *Less14) Category() string {
	return "POST Injection"
}

func (l *Less14) Route(r *gin.RouterGroup) {
	group := r.Group("/less-14")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
	}
}

func (l *Less14) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-14")

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
			l.renderLoginFormWithError(c, "less-14", "登录失败")
		} else {
			l.renderLoginFormWithError(c, "less-14", err.Error())
		}
		return
	}

	l.renderLoginFormWithResult(c, "less-14", "You are in...........")
}

func (l *Less14) renderLoginForm(c *gin.Context) {
	l.renderLoginFormWithError(c, "less-14", "")
}

func (l *Less14) renderLoginFormWithError(c *gin.Context, lesson, err string) {
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

func (l *Less14) renderLoginFormWithResult(c *gin.Context, lesson, result string) {
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

func (r *Registry) registerLess14() {
	r.Register(&Less14{BaseLesson{db: r.db, logger: r.logger}})
}
