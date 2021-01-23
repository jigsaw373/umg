package users

import (
	"errors"
	"fmt"
	"strings"

	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/roles"
	"github.com/boof/umg/util/password"
	"github.com/boof/umg/util/validator"
)

func init() {
	db.Sync(new(User))
}

// Save inserts a new user into the database
func (u *User) Save() error {
	err := u.validateForInsert()
	if err != nil {
		return err
	}

	// remove id
	u.ID = 0

	// remove roles
	u.RoleIDs = nil

	// Hash password
	hash, err := password.HashPassword(u.Password)
	if err != nil {
		return errors.New("unable to hash password")
	}

	u.Password = hash

	_, err = db.Engine.Insert(u)
	return err
}

func (u *User) Update() error {
	if err := u.validateFroUpdate(); err != nil {
		return err
	}

	_, err := db.Engine.Id(u.ID).AllCols().Update(u)
	return err
}

func (u *User) ChangePassword(newPass string) error {
	errs := validator.Validate.Var(newPass, "min=4,max=64")
	if errs != nil {
		return errors.New("password should have a length between 1 and 64")
	}

	// Hash password
	hash, err := password.HashPassword(newPass)
	if err != nil {
		return errors.New("unable to hash password")
	}

	u.Password = hash

	_, err = db.Engine.Id(u.ID).Cols("password").Update(u)
	return err
}

func (u *User) UpdateLastLogin() error {
	_, err := db.Engine.Id(u.ID).Cols("last_login").Update(u)
	return err
}

func (u *User) SaveWithRoles() error {
	err := u.validateForInsertWithRoles()
	if err != nil {
		return err
	}

	// remove id
	u.ID = 0

	// Hash password
	hash, err := password.HashPassword(u.Password)
	if err != nil {
		return errors.New("unable to hash password")
	}

	u.Password = hash

	_, err = db.Engine.Insert(u)
	return err
}

// GetByID returns a User with the given id
func (u *User) GetByID() (*User, error) {
	user := &User{ID: u.ID}
	if has, err := db.Engine.Get(user); !has || err != nil {
		return user, errors.New("user not found")
	}

	return user, nil
}

// GetByEmail returns a User with the given email
func (u *User) GetByEmail() (*User, error) {
	user := &User{Email: u.Email}
	if has, err := db.Engine.Get(user); !has || err != nil {
		return user, errors.New("user not found")
	}

	return user, nil
}

func (u *User) GetByUsername() (*User, error) {
	user := &User{Username: u.Username}
	if has, err := db.Engine.Get(user); !has || err != nil {
		return user, errors.New("user not found")
	}

	return user, nil
}

// RemoveByID removes the user by id
func (u *User) RemoveByID() error {
	_, err := db.Engine.Id(u.ID).Delete(&User{})
	return err
}

// RemoveByUserName removes the user by username
func (u *User) RemoveByUserName() error {
	_, err := db.Engine.Delete(&User{Username: u.Username})
	return err
}

// validateForInsert validates
func (u *User) validateForInsert() error {
	if err := password.ValidateUsername(u.Username); err != nil {
		return err
	}

	errs := validator.Validate.Var(u.Password, "min=4,max=64")
	if errs != nil {
		return errors.New("password should have a length between 1 and 64")
	}

	errs = validator.Validate.Var(u.Email, "email,required")
	if errs != nil {
		return errors.New("invalid email")
	}

	if ok, err := u.HasUniqueUsername(); !ok || err != nil {
		// Todo: log error message
		return errors.New("username should be unique")
	}

	if ok, err := u.HasUniqueEmail(); !ok || err != nil {
		// Todo: log error message
		return errors.New("email should be unique")
	}

	// Todo: validate phone numbers
	// Todo: validate website addresses

	return nil
}

func (u *User) validateFroUpdate() error {
	old, err := u.GetByID()
	if err != nil {
		return errors.New("invalid user")
	}

	if u.Username != old.Username {
		if err := password.ValidateUsername(u.Username); err != nil {
			return err
		}

		if has, _ := u.HasUniqueUsername(); !has {
			return errors.New("username is taken")
		}
	}

	if u.Email != old.Email {
		errs := validator.Validate.Var(u.Email, "email,required")
		if errs != nil {
			return errors.New("invalid email")
		}

		if has, _ := u.HasUniqueEmail(); !has {
			return errors.New("email is taken")
		}
	}

	u.RoleIDs = old.RoleIDs
	u.Password = old.Password

	return nil
}

// validateForInsertWithRoles validates
func (u *User) validateForInsertWithRoles() error {
	if err := password.ValidateUsername(u.Username); err != nil {
		return err
	}

	errs := validator.Validate.Var(u.Password, "min=4,max=64")
	if errs != nil {
		return errors.New("password should have a length between 1 and 64")
	}

	errs = validator.Validate.Var(u.Email, "email,required")
	if errs != nil {
		return errors.New("invalid email")
	}

	if ok, err := u.HasUniqueUsername(); !ok || err != nil {
		// Todo: log error message
		return errors.New("username should be unique")
	}

	if ok, err := u.HasUniqueEmail(); !ok || err != nil {
		// Todo: log error message
		return errors.New("email should be unique")
	}

	// Todo: validate phone numbers
	// Todo: validate website addresses

	// validate roles
	for _, id := range u.RoleIDs {
		_, err := (&roles.Role{ID: id}).GetByID()
		if err != nil {
			return err
		}
	}

	return nil
}

// HasUniqueUsername indicates that user's username is unique or not
func (u *User) HasUniqueUsername() (bool, error) {
	var user User
	ok, err := db.Engine.SQL("SELECT * FROM \"user\" WHERE LOWER(username) = ?", strings.ToLower(u.Username)).Get(&user)
	if err != nil {
		return false, errors.New("database error")
	} else if ok {
		return false, errors.New("duplicated domain name")
	}

	return true, nil
}

// HasUniqueEmail indicates that user's email is unique or not
func (u *User) HasUniqueEmail() (bool, error) {
	var user User
	ok, err := db.Engine.SQL("SELECT * FROM \"user\" WHERE LOWER(email) = ?", strings.ToLower(u.Email)).Get(&user)
	if err != nil {
		return false, errors.New("database error")
	} else if ok {
		return false, errors.New("duplicated email address")
	}

	return true, nil
}

// AssignRole assigns a role to the current user
func (u *User) AssignRole(roleID int64) error {
	// validate role
	if has, _ := db.Engine.Get(&roles.Role{ID: roleID}); !has {
		return errors.New("invalid role")
	}

	if u.RoleIDs == nil {
		u.RoleIDs = make([]int64, 0)
	}

	for _, id := range u.RoleIDs {
		if id == roleID {
			return nil
		}
	}
	u.RoleIDs = append(u.RoleIDs, roleID)

	// Todo: validate other roles inside a lock
	// Todo: fetch users role from database

	// update roles
	_, err := db.Engine.ID(u.ID).Update(&User{RoleIDs: u.RoleIDs})
	return err
}

func (u *User) DisallowRole(roleID int64) error {
	// validate role
	if has, _ := db.Engine.Get(&roles.Role{ID: roleID}); !has {
		return errors.New("invalid role")
	}

	roleIDs := make([]int64, 0)
	for _, id := range u.RoleIDs {
		if id != roleID {
			roleIDs = append(roleIDs, id)
		}
	}

	// Todo: validate other roles inside a lock
	// Todo: fetch users role from database

	// update roles
	_, err := db.Engine.ID(u.ID).Cols("role_ids").Update(&User{RoleIDs: roleIDs})
	return err
}

// IsAdmin indicates that current user has admin permission or not
func (u *User) IsAdmin() bool {
	for _, roleID := range u.RoleIDs {
		r, err := (&roles.Role{ID: roleID}).GetByID()
		if err == nil && r.Name == "admin" {
			return true
		}
	}

	return false
}

// GetAll returns all users sorted by given field, limited to count
// and offset by page
func GetAll(count, page int64, order string, sortBy string) ([]User, error) {
	var users []User
	var query string

	if sortBy == LastLogin {
		order = reverseOrder(order)
		query = "SELECT * FROM \"user\" WHERE last_login IS NOT NULL ORDER BY %s %s LIMIT %d OFFSET %d;"
	} else if sortBy == Created {
		order = reverseOrder(order)
		query = "SELECT * FROM \"user\" WHERE created_at IS NOT NULL ORDER BY %s %s LIMIT %d OFFSET %d;"
	} else if sortBy == ID {
		query = "SELECT * FROM \"user\" ORDER BY %s %s LIMIT %d OFFSET %d;"
	} else if sortBy == Name {
		query = "SELECT * FROM \"user\" WHERE NOT name = '' AND name IS NOT NULL ORDER BY %s %s LIMIT %d OFFSET %d;"
	} else {
		query = "SELECT * FROM \"user\" ORDER BY LOWER(%s) %s LIMIT %d OFFSET %d;"
	}

	// populate query
	query = fmt.Sprintf(query, sortBy, order, count, (page-1)*count)

	err := db.Engine.SQL(query).Find(&users)
	return users, err
}

func Count() (int64, error) {
	return db.Engine.Count(new(User))
}

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

func reverseOrder(order string) string {
	if order == Desc {
		return Asc
	}

	return Desc
}
