package services

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

func IsValidSort(sortBy ...string) bool {
	if len(sortBy) == 0 {
		return false
	}

	for _, by := range sortBy {
		if by != Name && by != Username && by != Email && by != LastLogin && by != ID && by != Created {
			return false
		}
	}

	return true
}
