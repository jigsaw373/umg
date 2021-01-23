package policies

// Policy used for saving policies
type Policy struct {
	ID         int64    `xorm:"pk not null autoincr 'id'" json:"id"`
	RoleID     int64    `xorm:"not null 'role_id'" json:"role_id"`
	Type       string   `xorm:"varchar(1) not null" json:"type"`
	Actions    []string `xorm:"not null" json:"actions"`
	Properties []int64  `json:"properties"`
	ProductID  int64    `xorm:"'product_id'" json:"product_id"`
	DomainID   int64    `xorm:"'domain_id'" json:"domain_id"`
}
