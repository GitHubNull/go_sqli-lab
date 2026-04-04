package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less2 Less-2漏洞: Error Based - Integer
type Less2 struct {
	BaseLesson
}

func (l *Less2) ID() string {
	return "less-2"
}

func (l *Less2) Name() string {
	return "Less-2: Error Based - Integer"
}

func (l *Less2) Description() string {
	return "基于错误信息的SQL注入 - 整数类型 (无引号)"
}

func (l *Less2) Category() string {
	return "Error Based"
}

func (l *Less2) Route(r *gin.RouterGroup) {
	lessonPath := "/less-2"
	group := r.Group(lessonPath)
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less2) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-2")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-2", "请输入ID参数", "")
		return
	}

	// 有漏洞的SQL查询 - 整数类型，无引号
	query := fmt.Sprintf("SELECT * FROM users WHERE id=%s LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-2", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-2", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-2", result, "")
}

func (r *Registry) registerLess2() {
	r.Register(&Less2{BaseLesson{db: r.db, logger: r.logger}})
}
