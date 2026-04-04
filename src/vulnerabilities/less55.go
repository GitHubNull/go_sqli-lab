package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less55 Less-55: Challenge 2 - Error Based Blind
type Less55 struct {
	BaseLesson
}

func (l *Less55) ID() string {
	return "less-55"
}

func (l *Less55) Name() string {
	return "Less-55: Challenge 2 - Error Blind"
}

func (l *Less55) Description() string {
	return "挑战关卡2：错误回显盲注"
}

func (l *Less55) Category() string {
	return "Challenge"
}

func (l *Less55) Route(r *gin.RouterGroup) {
	group := r.Group("/less-55")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less55) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-55")

	id := c.Query("id")
	if id == "" {
		id = "1"
	}

	query := fmt.Sprintf("SELECT * FROM challenge2 WHERE id=(%s) LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-55", "", "")
		} else {
			l.renderHTML(c, "less-55", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-55", result, "")
}

func (r *Registry) registerLess55() {
	r.Register(&Less55{BaseLesson{db: r.db, logger: r.logger}})
}
