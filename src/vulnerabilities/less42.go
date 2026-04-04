package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less42 POST - Stacked Query - String
type Less42 struct {
	BaseLesson
}

func (l *Less42) ID() string {
	return "less-42"
}

func (l *Less42) Name() string {
	return "Less-42: POST - Stacked Query - String"
}

func (l *Less42) Description() string {
	return "POST请求 - 字符串类型堆叠查询注入"
}

func (l *Less42) Category() string {
	return "Stacked Injection"
}

func (l *Less42) Route(r *gin.RouterGroup) {
	group := r.Group("/less-42")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
	}
}

func (l *Less42) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-42")

	login_user := c.PostForm("login_user")
	login_password := c.PostForm("login_password")

	if login_user == "" && login_password == "" {
		l.renderLoginForm(c, "less-42", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' and password='%s' LIMIT 1 OFFSET 0", login_user, login_password)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderLoginForm(c, "less-42", "登录失败")
		} else {
			l.renderLoginForm(c, "less-42", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderLoginForm(c, "less-42", result)
}

func (l *Less42) renderLoginForm(c *gin.Context, lesson, result string) {
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
				<input type="text" name="login_user" value=""/>
			</div>
			<div>
				<label>Password:</label>
				<input type="password" name="login_password" value=""/>
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

func (r *Registry) registerLess42() {
	r.Register(&Less42{BaseLesson{db: r.db, logger: r.logger}})
}
