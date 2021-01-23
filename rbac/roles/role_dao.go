package roles

import (
	"encoding/json"
	"errors"
	"github.com/boof/umg/db"
	"github.com/boof/umg/util/validator"
	"strings"
)

const (
	adminRole = "admin"
)

func init() {
	db.Sync(new(Role))
}

func (r *Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(&map[string]interface{}{
		"id":         r.ID,
		"name":       r.Name,
		"created_at": r.CreatedAt.Format("Jan 02, 2006"),
	})
}

// Save inserts a new role into the database
func (r *Role) Save() error {
	if err := r.ValidateForInsert(); err != nil {
		return err
	}

	// remove id
	r.ID = 0

	_, err := db.Engine.Insert(r)
	return err
}

func (r *Role) Update() error {
	if err := r.validateFroUpdate(); err != nil {
		return err
	}

	_, err := db.Engine.Id(r.ID).AllCols().Update(r)
	return err
}

// GetByID returns a role with the given id
func (r *Role) GetByID() (*Role, error) {
	if r.ID < 1 {
		return r, errors.New("invalid id")
	}

	role := &Role{ID: r.ID}
	if has, err := db.Engine.Get(role); !has || err != nil {
		return role, errors.New("role not found")
	}

	return role, nil
}

// GetByName returns a Role with the given name
func (r *Role) GetByName() (*Role, error) {
	if r.Name == "" {
		return r, errors.New("invalid name")
	}

	role := &Role{Name: r.Name}
	if has, err := db.Engine.Get(role); !has || err != nil {
		return role, errors.New("role not found")
	}

	return role, nil
}

// HasUniqueNname indicates that role's name is unique or not
func (r *Role) HasUniqueName() (bool, error) {
	has, err := db.Engine.Get(&Role{Name: r.Name})
	return !has, err
}

// ValidateForInsert validate the role
func (r *Role) ValidateForInsert() error {
	if r.Name == "" {
		return errors.New("invalid role name")
	}

	errs := validator.Validate.Var(r.Name, "min=1,max=64")
	if errs != nil {
		return errors.New("role name should have a length between 1 and 64")
	}

	var roles []Role
	err := db.Engine.SQL("SELECT * FROM role WHERE LOWER(name) = ?", strings.ToLower(r.Name)).Find(&roles)
	if err != nil {
		return errors.New("database error")
	} else if roles != nil && len(roles) > 0 {
		return errors.New("duplicated role name")
	}

	return nil
}

func (r *Role) validateFroUpdate() error {
	old, err := r.GetByID()
	if err != nil {
		return errors.New("invalid role")
	}

	errs := validator.Validate.Var(r.Name, "min=1,max=64")
	if errs != nil {
		return errors.New("role name should have a length between 1 and 64")
	}

	if r.Name != old.Name {
		if has, _ := r.HasUniqueName(); !has {
			return errors.New("duplicated name")
		}
	}

	return nil
}

func (r *Role) IsAdmin() bool {
	return strings.ToLower(r.Name) == adminRole
}

// GetAllRoles returns all roles
func GetAllRoles(count, page int, order, sortBy string) ([]Role, error) {
	var roles []Role

	session := db.Engine.Limit(count, (page-1)*count)
	if order == Desc {
		session = session.Desc(sortBy)
	} else {
		session = session.Asc(sortBy)
	}

	err := session.Find(&roles)
	return roles, err
}

func getAllRoleIDs(count, page int) []int64 {
	roleIDs := make([]int64, 0)

	roles, err := GetAllRoles(count, page, Asc, ID)
	if err != nil {
		return roleIDs
	}

	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	return roleIDs
}

func Count() (int64, error) {
	return db.Engine.Count(new(Role))
}

func Pages(count int64) int64 {
	roleCount, err := Count()
	if err != nil {
		return 1
	}

	pages := roleCount / count
	if roleCount%count != 0 {
		pages += 1
	}

	return pages
}

func IsValidSort(sortBy ...string) bool {
	if len(sortBy) == 0 {
		return false
	}

	for _, by := range sortBy {
		if by != Name && by != ID {
			return false
		}
	}

	return true
}
