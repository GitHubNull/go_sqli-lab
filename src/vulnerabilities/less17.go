package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less17 POST - Update Query - Error Based - String
type Less17 struct {
	BaseLesson
}

func (l *Less17) ID() string {
	return "less-17"
}

func (l *Less17) Name() string {
	return "Less-17: POST - Update Query - Error Based - String"
}

func (l *Less17) Description() string {
	return "POST请求 - UPDATE语句注入 (基于错误)"
}

func (l *Less17) Category() string {
	return "Update Query"
}

func (l *Less17) Route(r *gin.RouterGroup) {
	group := r.Group("/less-17")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
	}
}

func (l *Less17) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-17")

	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" && password == "" {
		l.renderPasswordResetForm(c, "less-17", "")
		return
	}

	// 先检查用户是否存在
	checkQuery := fmt.Sprintf("SELECT * FROM users WHERE username='%s' LIMIT 0,1", username)
	var user models.User
	err := l.db.QueryRow(checkQuery).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			l.renderPasswordResetForm(c, "less-17", "用户不存在")
		} else {
			l.renderPasswordResetForm(c, "less-17", err.Error())
		}
		return
	}

	// 更新密码 - 有漏洞的UPDATE语句
	updateQuery := fmt.Sprintf("UPDATE users SET password='%s' WHERE username='%s'", password, username)
	_, err = l.db.Exec(updateQuery)
	if err != nil {
		l.renderPasswordResetForm(c, "less-17", err.Error())
		return
	}

	l.renderPasswordResetForm(c, "less-17", "密码更新成功")
}

func (l *Less17) renderPasswordResetForm(c *gin.Context, lesson, msg string) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<div class="container">
		<h1>%s - 密码重置</h1>
		<form method="POST" action="">
			<div>
				<label>Username:</label>
				<input type="text" name="username" value=""/>
			</div>
			<div>
				<label>New Password:</label>
				<input type="text" name="password" value=""/>
			</div>
			<div>
				<input type="submit" value="Reset"/>
			</div>
		</form>
		%s
		<div class="nav">
			<a href="/">返回首页</a>
		</div>
	</div>
</body>
</html>`, lesson, lesson, renderResult(msg))
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}

func (r *Registry) registerLess17() {
	r.Register(&Less17{BaseLesson{db: r.db, logger: r.logger}})
}
