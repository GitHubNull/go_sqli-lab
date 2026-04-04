package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less40 GET - Stacked Query - String - Blind
type Less40 struct {
	BaseLesson
}

func (l *Less40) ID() string {
	return "less-40"
}

func (l *Less40) Name() string {
	return "Less-40: GET - Stacked Query - String - Blind"
}

func (l *Less40) Description() string {
	return "GET请求 - 字符串类型盲注堆叠查询"
}

func (l *Less40) Category() string {
	return "Stacked Injection"
}

func (l *Less40) Route(r *gin.RouterGroup) {
	group := r.Group("/less-40")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less40) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-40")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-40", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-40", "", "")
		} else {
			l.renderHTML(c, "less-40", "", "")
		}
		return
	}

	l.renderHTML(c, "less-40", "You are in...........", "")
}

func (r *Registry) registerLess40() {
	r.Register(&Less40{BaseLesson{db: r.db, logger: r.logger}})
}
