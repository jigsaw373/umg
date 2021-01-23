package products

// Product is a type for saving products of the system
type Product struct {
	ID       int64  `xorm:"pk not null autoincr 'id'" json:"id"`
	DomainID int64  `xorm:"'domain_id'" json:"domain_id"`
	Name     string `json:"name"`
}
