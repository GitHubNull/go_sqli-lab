package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less15 POST - Blind - Boolean/Time Based - Single Quotes
type Less15 struct {
	BaseLesson
}

func (l *Less15) ID() string {
	return "less-15"
}

func (l *Less15) Name() string {
	return "Less-15: POST - Blind - Boolean/Time Based - Single Quotes"
}

func (l *Less15) Description() string {
	return "POST请求 - 基于布尔/时间的盲注 (单引号)"
}

func (l *Less15) Category() string {
	return "POST Blind"
}

func (l *Less15) Route(r *gin.RouterGroup) {
	group := r.Group("/less-15")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
	}
}

func (l *Less15) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-15")

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
			l.renderLoginFormWithError(c, "less-15", "")
		} else {
			l.renderLoginFormWithError(c, "less-15", "")
		}
		return
	}

	l.renderLoginFormWithResult(c, "less-15", "You are in...........")
}

func (l *Less15) renderLoginForm(c *gin.Context) {
	l.renderLoginFormWithError(c, "less-15", "")
}

func (l *Less15) renderLoginFormWithError(c *gin.Context, lesson, err string) {
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

func (l *Less15) renderLoginFormWithResult(c *gin.Context, lesson, result string) {
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

func (r *Registry) registerLess15() {
	r.Register(&Less15{BaseLesson{db: r.db, logger: r.logger}})
}
