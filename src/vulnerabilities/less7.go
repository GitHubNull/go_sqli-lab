package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less7 Less-7漏洞: Dump into outfile
type Less7 struct {
	BaseLesson
}

func (l *Less7) ID() string {
	return "less-7"
}

func (l *Less7) Name() string {
	return "Less-7: Dump into outfile"
}

func (l *Less7) Description() string {
	return "导出到文件 - 使用INTO OUTFILE的SQL注入"
}

func (l *Less7) Category() string {
	return "File Operations"
}

func (l *Less7) Route(r *gin.RouterGroup) {
	group := r.Group("/less-7")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less7) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-7")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-7", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id=(('%s')) LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-7", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-7", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-7", result, "")
}

func (r *Registry) registerLess7() {
	r.Register(&Less7{BaseLesson{db: r.db, logger: r.logger}})
}
