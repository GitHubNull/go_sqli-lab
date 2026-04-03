package models

// User 用户模型
type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

// Email 邮箱模型
type Email struct {
	ID      int    `json:"id" db:"id"`
	EmailID string `json:"email_id" db:"email_id"`
}

// UAgent 用户代理模型
type UAgent struct {
	ID        int    `json:"id" db:"id"`
	UAgent    string `json:"uagent" db:"uagent"`
	IPAddress string `json:"ip_address" db:"ip_address"`
	Username  string `json:"username" db:"username"`
}

// Referer 来源页面模型
type Referer struct {
	ID        int    `json:"id" db:"id"`
	Referer   string `json:"referer" db:"referer"`
	IPAddress string `json:"ip_address" db:"ip_address"`
}
