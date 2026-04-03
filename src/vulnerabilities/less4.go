package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less4 Less-4漏洞: Error Based - Double Quotes
type Less4 struct {
	BaseLesson
}

func (l *Less4) ID() string {
	return "less-4"
}

func (l *Less4) Name() string {
	return "Less-4: Error Based - Double Quotes"
}

func (l *Less4) Description() string {
	return "基于错误信息的SQL注入 - 双引号字符串类型"
}

func (l *Less4) Category() string {
	return "Error Based"
}

func (l *Less4) Route(r *gin.RouterGroup) {
	lessonPath := "/less-4"
	group := r.Group(lessonPath)
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less4) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-4")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-4", "请输入ID参数", "")
		return
	}

	// 有漏洞的SQL查询 - 使用双引号
	query := fmt.Sprintf("SELECT * FROM users WHERE id=\"%s\" LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-4", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-4", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-4", result, "")
}

func (r *Registry) registerLess4() {
	r.Register(&Less4{BaseLesson{db: r.db, logger: r.logger}})
}
