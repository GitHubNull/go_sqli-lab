package vulnerabilities

import (
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less57 Less-57: Challenge 4 - Boolean Based Blind
type Less57 struct {
	BaseLesson
}

func (l *Less57) ID() string {
	return "less-57"
}

func (l *Less57) Name() string {
	return "Less-57: Challenge 4 - Boolean Blind"
}

func (l *Less57) Description() string {
	return "挑战关卡4：布尔盲注"
}

func (l *Less57) Category() string {
	return "Challenge"
}

func (l *Less57) Route(r *gin.RouterGroup) {
	group := r.Group("/less-57")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less57) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-57")

	id := c.Query("id")
	if id == "" {
		id = "1"
	}

	query := fmt.Sprintf("SELECT * FROM challenge4 WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		l.renderHTML(c, "less-57", "", "")
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-57", result, "")
}

func (r *Registry) registerLess57() {
	r.Register(&Less57{BaseLesson{db: r.db, logger: r.logger}})
}
