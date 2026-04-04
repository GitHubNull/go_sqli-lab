package vulnerabilities

import (
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less63 Less-63: Challenge 10 - Blind Integer
type Less63 struct {
	BaseLesson
}

func (l *Less63) ID() string {
	return "less-63"
}

func (l *Less63) Name() string {
	return "Less-63: Challenge 10 - Blind Integer"
}

func (l *Less63) Description() string {
	return "挑战关卡10：整数型布尔盲注"
}

func (l *Less63) Category() string {
	return "Challenge"
}

func (l *Less63) Route(r *gin.RouterGroup) {
	group := r.Group("/less-63")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less63) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-63")

	id := c.Query("id")
	if id == "" {
		id = "1"
	}

	query := fmt.Sprintf("SELECT * FROM challenge10 WHERE id=%s LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		// Blind - no error shown
		l.renderHTML(c, "less-63", "", "")
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-63", result, "")
}

func (r *Registry) registerLess63() {
	r.Register(&Less63{BaseLesson{db: r.db, logger: r.logger}})
}
