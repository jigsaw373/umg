package services

import (
	"strings"

	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/access"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/rest_errors"
	"github.com/boof/umg/util/password"
)

// Login if there was a user with the current username and password
// this functions fetch and returned it otherwise returns error
func Login(username, pass string) (*users.User, rest_errors.Error) {
	if pass == "" || username == "" {
		return nil, rest_errors.NewUnauthorizedError("Incorrect username or password.")
	}

	var user users.User
	ok, err := db.Engine.SQL("SELECT * FROM \"user\" WHERE LOWER(username) = ?", strings.ToLower(username)).Get(&user)
	if !ok || err != nil {
		return nil, rest_errors.NewUnauthorizedError("Incorrect username or password.")
	}

	if password.IsValidPass(pass, user.Password) {
		if ok, _ := access.Expired(user.ID); ok {
			return nil, rest_errors.NewUnauthorizedError("Your access time is expired!")
		}
	} else {
		return nil, rest_errors.NewUnauthorizedError("Incorrect username or password.")
	}

	return &user, nil
}
