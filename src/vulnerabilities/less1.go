package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less1 Less-1漏洞: Error Based - String
type Less1 struct {
	BaseLesson
}

func (l *Less1) ID() string {
	return "less-1"
}

func (l *Less1) Name() string {
	return "Less-1: Error Based - String"
}

func (l *Less1) Description() string {
	return "基于错误信息的SQL注入 - 字符串类型 (单引号)"
}

func (l *Less1) Category() string {
	return "Error Based"
}

func (l *Less1) Route(r *gin.RouterGroup) {
	lessonPath := "/less-1"
	group := r.Group(lessonPath)
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less1) handleIndex(c *gin.Context) {
	// 记录请求到result.txt
	l.logRequest(c, "less-1")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-1", "请输入ID参数", "")
		return
	}

	// 有漏洞的SQL查询 - 直接拼接
	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-1", "", "没有找到记录")
		} else {
			// 显示SQL错误信息（漏洞点）
			l.renderHTML(c, "less-1", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-1", result, "")
}

func (r *Registry) registerLess1() {
	r.Register(&Less1{BaseLesson{db: r.db, logger: r.logger}})
}
