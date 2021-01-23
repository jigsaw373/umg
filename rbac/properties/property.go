package properties

type Property struct {
	ID         int64  `xorm:"pk not null autoincr 'id'" json:"id"`
	MeteringID int64  `xorm:"metering_id" json:"metering_id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
}
