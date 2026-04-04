package vulnerabilities

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Less58 Less-58: Challenge 5 - Stacked Query
type Less58 struct {
	BaseLesson
}

func (l *Less58) ID() string {
	return "less-58"
}

func (l *Less58) Name() string {
	return "Less-58: Challenge 5 - Stacked Query"
}

func (l *Less58) Description() string {
	return "挑战关卡5：堆叠查询注入"
}

func (l *Less58) Category() string {
	return "Challenge"
}

func (l *Less58) Route(r *gin.RouterGroup) {
	group := r.Group("/less-58")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less58) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-58")

	id := c.Query("id")
	if id == "" {
		id = "1"
	}

	query := fmt.Sprintf("SELECT * FROM challenge5 WHERE id='%s' LIMIT 1 OFFSET 0;", id)

	rows, err := l.db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-58", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-58", "", err.Error())
		}
		return
	}
	defer rows.Close()

	result := "Query executed successfully"
	l.renderHTML(c, "less-58", result, "")
}

func (r *Registry) registerLess58() {
	r.Register(&Less58{BaseLesson{db: r.db, logger: r.logger}})
}
