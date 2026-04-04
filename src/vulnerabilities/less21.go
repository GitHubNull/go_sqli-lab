package vulnerabilities

import (
	"database/sql"
	"encoding/base64"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less21 POST - Cookie Injection - User-Agent - Base64
type Less21 struct {
	BaseLesson
}

func (l *Less21) ID() string {
	return "less-21"
}

func (l *Less21) Name() string {
	return "Less-21: POST - Cookie Injection - User-Agent - Base64"
}

func (l *Less21) Description() string {
	return "POST请求 - Base64编码的Cookie注入"
}

func (l *Less21) Category() string {
	return "Cookie Injection"
}

func (l *Less21) Route(r *gin.RouterGroup) {
	group := r.Group("/less-21")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
	}
}

func (l *Less21) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-21")

	// 检查Cookie
	cookie, err := c.Cookie("uname")
	if err == nil && cookie != "" {
		// Base64解码
		decoded, err := base64.StdEncoding.DecodeString(cookie)
		if err != nil {
			l.renderCookieForm(c, "less-21", "Cookie解码失败")
			return
		}

		// 从解码后的Cookie获取用户名并查询
		query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' LIMIT 1 OFFSET 0", string(decoded))
		var user models.User
		err = l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				l.renderCookieForm(c, "less-21", "Cookie用户不存在")
			} else {
				l.renderCookieForm(c, "less-21", err.Error())
			}
			return
		}
		result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
		l.renderCookieForm(c, "less-21", result)
		return
	}

	uname := c.PostForm("uname")
	passwd := c.PostForm("passwd")

	if uname == "" && passwd == "" {
		l.renderLoginForm(c, "less-21", "")
		return
	}

	// 验证用户
	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' and password='%s' LIMIT 1 OFFSET 0", uname, passwd)
	var user models.User
	err = l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			l.renderLoginForm(c, "less-21", "登录失败")
		} else {
			l.renderLoginForm(c, "less-21", err.Error())
		}
		return
	}

	// 设置Base64编码的Cookie
	encoded := base64.StdEncoding.EncodeToString([]byte(uname))
	c.SetCookie("uname", encoded, 3600, "/", "", false, true)
	result := fmt.Sprintf("登录成功！Base64 Cookie已设置<br>Your Login name: %s", user.Username)
	l.renderLoginForm(c, "less-21", result)
}

func (l *Less21) renderLoginForm(c *gin.Context, lesson, result string) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<div class="container">
		<h1>%s - Base64 Cookie注入</h1>
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

func (l *Less21) renderCookieForm(c *gin.Context, lesson, result string) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<div class="container">
		<h1>%s - Base64 Cookie注入</h1>
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

func (r *Registry) registerLess21() {
	r.Register(&Less21{BaseLesson{db: r.db, logger: r.logger}})
}
