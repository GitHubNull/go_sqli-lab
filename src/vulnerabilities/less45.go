package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less45 POST - Stacked Query - Blind
type Less45 struct {
	BaseLesson
}

func (l *Less45) ID() string {
	return "less-45"
}

func (l *Less45) Name() string {
	return "Less-45: POST - Stacked Query - Blind"
}

func (l *Less45) Description() string {
	return "POST请求 - 盲注堆叠查询注入"
}

func (l *Less45) Category() string {
	return "Stacked Injection"
}

func (l *Less45) Route(r *gin.RouterGroup) {
	group := r.Group("/less-45")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
	}
}

func (l *Less45) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-45")

	login_user := c.PostForm("login_user")
	login_password := c.PostForm("login_password")

	if login_user == "" && login_password == "" {
		l.renderLoginForm(c, "less-45", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' and password='%s' LIMIT 0,1", login_user, login_password)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderLoginForm(c, "less-45", "")
		} else {
			l.renderLoginForm(c, "less-45", "")
		}
		return
	}

	l.renderLoginForm(c, "less-45", "You are in...........")
}

func (l *Less45) renderLoginForm(c *gin.Context, lesson, result string) {
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

func (r *Registry) registerLess45() {
	r.Register(&Less45{BaseLesson{db: r.db, logger: r.logger}})
}
