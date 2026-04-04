package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less18 POST - Header Injection - User-Agent
type Less18 struct {
	BaseLesson
}

func (l *Less18) ID() string {
	return "less-18"
}

func (l *Less18) Name() string {
	return "Less-18: POST - Header Injection - User-Agent"
}

func (l *Less18) Description() string {
	return "POST请求 - User-Agent头注入"
}

func (l *Less18) Category() string {
	return "Header Injection"
}

func (l *Less18) Route(r *gin.RouterGroup) {
	group := r.Group("/less-18")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
	}
}

func (l *Less18) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-18")

	uname := c.PostForm("uname")
	passwd := c.PostForm("passwd")

	if uname == "" && passwd == "" {
		l.renderLoginForm(c, "less-18", "")
		return
	}

	// 验证用户
	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' and password='%s' LIMIT 1 OFFSET 0", uname, passwd)
	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			l.renderLoginForm(c, "less-18", "登录失败")
		} else {
			l.renderLoginForm(c, "less-18", err.Error())
		}
		return
	}

	// 获取User-Agent并插入数据库
	uagent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	// 有漏洞的INSERT语句
	insertQuery := fmt.Sprintf("INSERT INTO uagents (uagent, ip_address, username) VALUES ('%s', '%s', '%s')", uagent, ip, uname)
	_, err = l.db.Exec(insertQuery)
	if err != nil {
		l.renderLoginForm(c, "less-18", err.Error())
		return
	}

	result := fmt.Sprintf("Your User Agent is: %s<br>Your IP Address is: %s", uagent, ip)
	l.renderLoginForm(c, "less-18", result)
}

func (l *Less18) renderLoginForm(c *gin.Context, lesson, result string) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<div class="container">
		<h1>%s - User-Agent注入</h1>
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

func (r *Registry) registerLess18() {
	r.Register(&Less18{BaseLesson{db: r.db, logger: r.logger}})
}
