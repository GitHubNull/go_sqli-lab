package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less38 GET - Stacked Query
type Less38 struct {
	BaseLesson
}

func (l *Less38) ID() string {
	return "less-38"
}

func (l *Less38) Name() string {
	return "Less-38: GET - Stacked Query"
}

func (l *Less38) Description() string {
	return "GET请求 - 堆叠查询注入"
}

func (l *Less38) Category() string {
	return "Stacked Injection"
}

func (l *Less38) Route(r *gin.RouterGroup) {
	group := r.Group("/less-38")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less38) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-38")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-38", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-38", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-38", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-38", result, "")
}

func (r *Registry) registerLess38() {
	r.Register(&Less38{BaseLesson{db: r.db, logger: r.logger}})
}
