package domains

// Domain used for saving Domains
// each Product in the system can be subset of a domain
type Domain struct {
	ID   int64  `xorm:"pk not null autoincr 'id'" json:"id"`
	Name string `xorm:"varchar(64) not null unique" json:"name"`
}
