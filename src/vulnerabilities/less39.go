package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less39 GET - Stacked Query - Integer
type Less39 struct {
	BaseLesson
}

func (l *Less39) ID() string {
	return "less-39"
}

func (l *Less39) Name() string {
	return "Less-39: GET - Stacked Query - Integer"
}

func (l *Less39) Description() string {
	return "GET请求 - 整数类型堆叠查询注入"
}

func (l *Less39) Category() string {
	return "Stacked Injection"
}

func (l *Less39) Route(r *gin.RouterGroup) {
	group := r.Group("/less-39")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less39) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-39")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-39", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id=%s LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-39", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-39", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-39", result, "")
}

func (r *Registry) registerLess39() {
	r.Register(&Less39{BaseLesson{db: r.db, logger: r.logger}})
}
