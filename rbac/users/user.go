package users

import (
	"time"

	"github.com/boof/umg/rbac/roles"
)

const (
	// table fields
	ID        = "id"
	Name      = "name"
	Username  = "username"
	Email     = "email"
	LastLogin = "last_login"
	Created   = "created_at"

	// order
	Desc = "desc"
	Asc  = "asc"
)

// User used for save users in database
type User struct {
	ID        int64     `xorm:"pk not null autoincr 'id'" json:"id"`
	Username  string    `xorm:"not null unique" json:"username"`
	Password  string    `xorm:"not null" json:"password"`
	Email     string    `xorm:"not null unique" json:"email"`
	Name      string    `xorm:"not null" json:"name"`
	Company   string    `json:"company"`
	Website   string    `json:"website"`
	Address1  string    `json:"address1"`
	Address2  string    `json:"address2"`
	Phone1    string    `json:"phone1"`
	Phone2    string    `json:"phone2"`
	Fax1      string    `json:"fax1"`
	Fax2      string    `json:"fax2"`
	RoleIDs   []int64   `xorm:"'role_ids'" json:"role_ids"`
	LastLogin time.Time `xorm:"last_login" json:"last_login"`
	CreatedAt time.Time `xorm:"created" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
}

type SimpleUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Online    bool   `json:"online"`
	LastLogin string `json:"last_login"`
	ExpireAt  string `json:"expire_at"`
}

type RoleUser struct {
	ID        int64         `json:"id"`
	Username  string        `json:"username"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Online    bool          `json:"online"`
	LastLogin string        `json:"last_login"`
	CreatedAt string        `json:"created_at"`
	ExpireAt  string        `json:"expire_at"`
	Roles     []*roles.Role `json:"roles"`
}
