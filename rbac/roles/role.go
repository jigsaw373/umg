package roles

import "time"

const (
	// table fields
	ID   = "id"
	Name = "name"

	// order
	Desc = "desc"
	Asc  = "asc"
)

// Role used for save roles in database
type Role struct {
	ID        int64     `xorm:"pk not null autoincr 'id'" json:"id"`
	Name      string    `xorm:"varchar(64) not null unique" json:"name"`
	CreatedAt time.Time `xorm:"created" json:"-"`
}
