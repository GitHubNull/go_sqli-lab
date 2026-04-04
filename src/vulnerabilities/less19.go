package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less19 POST - Header Injection - Referer
type Less19 struct {
	BaseLesson
}

func (l *Less19) ID() string {
	return "less-19"
}

func (l *Less19) Name() string {
	return "Less-19: POST - Header Injection - Referer"
}

func (l *Less19) Description() string {
	return "POST请求 - Referer头注入"
}

func (l *Less19) Category() string {
	return "Header Injection"
}

func (l *Less19) Route(r *gin.RouterGroup) {
	group := r.Group("/less-19")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
	}
}

func (l *Less19) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-19")

	uname := c.PostForm("uname")
	passwd := c.PostForm("passwd")

	if uname == "" && passwd == "" {
		l.renderLoginForm(c, "less-19", "")
		return
	}

	// 验证用户
	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' and password='%s' LIMIT 1 OFFSET 0", uname, passwd)
	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			l.renderLoginForm(c, "less-19", "登录失败")
		} else {
			l.renderLoginForm(c, "less-19", err.Error())
		}
		return
	}

	// 获取Referer并插入数据库
	referer := c.GetHeader("Referer")
	ip := c.ClientIP()

	// 有漏洞的INSERT语句
	insertQuery := fmt.Sprintf("INSERT INTO referers (referer, ip_address) VALUES ('%s', '%s')", referer, ip)
	_, err = l.db.Exec(insertQuery)
	if err != nil {
		l.renderLoginForm(c, "less-19", err.Error())
		return
	}

	result := fmt.Sprintf("Your Referer is: %s<br>Your IP Address is: %s", referer, ip)
	l.renderLoginForm(c, "less-19", result)
}

func (l *Less19) renderLoginForm(c *gin.Context, lesson, result string) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<div class="container">
		<h1>%s - Referer注入</h1>
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

func (r *Registry) registerLess19() {
	r.Register(&Less19{BaseLesson{db: r.db, logger: r.logger}})
}
