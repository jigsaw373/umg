package initialize

import (
	"log"

	"github.com/boof/umg/rbac/roles"
	"github.com/boof/umg/rbac/users"
)

func createAdmin() {
	role, err := createRole("admin")
	if err != nil {
		log.Fatalf("unable to create admin role: %v", err)
	}

	admin, err := createUser("admin", "pass", "admin@edgecomenergy.ca")
	if err != nil {
		log.Fatalf("unable to create admin user: %v", err)
	}

	err = admin.AssignRole(role.ID)
	if err != nil {
		log.Fatalf("unable to assign a admin role to admin user: %v", err)
	}
}

// createUser creates a user if not exists
func createUser(username, pass, email string) (*users.User, error) {
	user, err := (&users.User{Username: username}).GetByUsername()
	if err != nil {
		// try to create
		user = &users.User{
			Username: username,
			Password: pass,
			Email:    email,
		}
		if err = user.Save(); err != nil {
			return user, err
		}
	} else {
		// Todo: check that admin user has admin role or not
	}

	return user, nil
}

func createRole(name string) (*roles.Role, error) {
	role, err := (&roles.Role{Name: name}).GetByName()
	if err != nil {
		// try to create role
		role = &roles.Role{Name: name}
		if err = role.Save(); err != nil {
			return role, err
		}
	}

	return role, err
}
