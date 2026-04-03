package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less41 GET - Stacked Query - Integer - Blind
type Less41 struct {
	BaseLesson
}

func (l *Less41) ID() string {
	return "less-41"
}

func (l *Less41) Name() string {
	return "Less-41: GET - Stacked Query - Integer - Blind"
}

func (l *Less41) Description() string {
	return "GET请求 - 整数类型盲注堆叠查询"
}

func (l *Less41) Category() string {
	return "Stacked Injection"
}

func (l *Less41) Route(r *gin.RouterGroup) {
	group := r.Group("/less-41")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less41) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-41")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-41", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id=%s LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-41", "", "")
		} else {
			l.renderHTML(c, "less-41", "", "")
		}
		return
	}

	l.renderHTML(c, "less-41", "You are in...........", "")
}

func (r *Registry) registerLess41() {
	r.Register(&Less41{BaseLesson{db: r.db, logger: r.logger}})
}
