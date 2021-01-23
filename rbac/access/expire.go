package access

import (
	"time"
)

type Expire struct {
	ID        int64     `xorm:"pk not null autoincr 'id'"`
	UserID    int64     `xorm:"not null 'user_id'"`
	ExpireAt  time.Time `xorm:"expire_at"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

type ExpireReq struct {
	UserID int64  `json:"user_id"`
	Date   string `json:"date"`
}
