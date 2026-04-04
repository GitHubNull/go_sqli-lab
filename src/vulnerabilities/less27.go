package vulnerabilities

import (
	"database/sql"
	"fmt"
	"strings"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less27 GET - Error Based - UNION & SELECT Bypass
type Less27 struct {
	BaseLesson
}

func (l *Less27) ID() string {
	return "less-27"
}

func (l *Less27) Name() string {
	return "Less-27: GET - Error Based - UNION & SELECT Bypass"
}

func (l *Less27) Description() string {
	return "GET请求 - 过滤UNION和SELECT的SQL注入"
}

func (l *Less27) Category() string {
	return "WAF Bypass"
}

func (l *Less27) Route(r *gin.RouterGroup) {
	group := r.Group("/less-27")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less27) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-27")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-27", "请输入ID参数", "")
		return
	}

	// 过滤UNION和SELECT（不区分大小写）
	id = strings.ReplaceAll(id, "UNION", "")
	id = strings.ReplaceAll(id, "union", "")
	id = strings.ReplaceAll(id, "SELECT", "")
	id = strings.ReplaceAll(id, "select", "")

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-27", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-27", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-27", result, "")
}

func (r *Registry) registerLess27() {
	r.Register(&Less27{BaseLesson{db: r.db, logger: r.logger}})
}
