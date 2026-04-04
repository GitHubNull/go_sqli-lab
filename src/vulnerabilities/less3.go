package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less3 Less-3漏洞: Error Based - String with Twist
type Less3 struct {
	BaseLesson
}

func (l *Less3) ID() string {
	return "less-3"
}

func (l *Less3) Name() string {
	return "Less-3: Error Based - String (with twist)"
}

func (l *Less3) Description() string {
	return "基于错误信息的SQL注入 - 字符串类型 (带括号 twist)"
}

func (l *Less3) Category() string {
	return "Error Based"
}

func (l *Less3) Route(r *gin.RouterGroup) {
	lessonPath := "/less-3"
	group := r.Group(lessonPath)
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less3) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-3")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-3", "请输入ID参数", "")
		return
	}

	// 有漏洞的SQL查询 - 使用('')格式
	query := fmt.Sprintf("SELECT * FROM users WHERE id=('%s') LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-3", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-3", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-3", result, "")
}

func (r *Registry) registerLess3() {
	r.Register(&Less3{BaseLesson{db: r.db, logger: r.logger}})
}
