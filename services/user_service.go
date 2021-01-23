package services

import (
	"fmt"
	"time"

	"github.com/boof/umg/db"
	"github.com/boof/umg/email"
	"github.com/boof/umg/rbac/access"
	"github.com/boof/umg/rbac/domains"
	"github.com/boof/umg/rbac/policies"
	"github.com/boof/umg/rbac/products"
	"github.com/boof/umg/rbac/roles"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/rest_errors"
	"github.com/boof/umg/settings"
	"github.com/boof/umg/util/datetime"
)

func GetUserByID(userID int64) (*users.User, rest_errors.Error) {
	user, err := (&users.User{ID: int64(userID)}).GetByID()
	if err != nil {
		return nil, rest_errors.NewNotFoundError("User not found")
	}

	// don't expose password
	user.Password = ""

	return user, nil
}

func GetSimpleUsers(count, page int64, order, sortBy string) ([]*users.SimpleUser, int64, rest_errors.Error) {
	res := make([]*users.SimpleUser, 0)

	all, err := users.GetAll(count, page, order, sortBy)
	if err != nil {
		return nil, 0, rest_errors.NewNotFoundError("Unable to find any user")
	}

	pages := int64(1)
	if usersCount, err := users.Count(); err == nil {
		pages = usersCount / count

		if usersCount%count != 0 {
			pages += 1
		}
	}

	for _, user := range all {
		user := &users.SimpleUser{
			ID:        user.ID,
			Username:  user.Username,
			Name:      user.Name,
			Email:     user.Email,
			Online:    db.IsOnline(user.ID),
			LastLogin: user.LastLogin.Format(settings.UserDTLayout),
		}

		expire, err := (&access.Expire{UserID: user.ID}).GetByUserID()
		if err == nil {
			user.ExpireAt = expire.ExpireAt.Format(settings.UserDTLayout)
		} else {
			user.ExpireAt = "undefined"
		}

		res = append(res, user)
	}

	return res, pages, nil
}

func GetRoleUsers(count, page int64, order, sortBy string) ([]*users.RoleUser, int64, rest_errors.Error) {
	res := make([]*users.RoleUser, 0)

	all, err := users.GetAll(count, page, order, sortBy)
	if err != nil {
		return res, 0, rest_errors.NewNotFoundError("Unable to find any user")
	}

	pages := int64(1)
	if usersCount, err := users.Count(); err == nil {
		pages = usersCount / count

		if usersCount%count != 0 {
			pages += 1
		}
	}

	for _, user := range all {
		ru := &users.RoleUser{
			ID:        user.ID,
			Username:  user.Username,
			Name:      user.Name,
			Email:     user.Email,
			Online:    db.IsOnline(user.ID),
			CreatedAt: user.CreatedAt.Format("Jan 02, 2006"),
			Roles:     make([]*roles.Role, 0),
		}

		if user.LastLogin.Before(time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)) {
			ru.LastLogin = "-"
		} else {
			ru.LastLogin = datetime.DurationString(user.LastLogin)
		}

		for _, r := range user.RoleIDs {
			role, err := (&roles.Role{ID: r}).GetByID()
			if err == nil {
				ru.Roles = append(ru.Roles, role)
			}
		}

		expire, err := (&access.Expire{UserID: user.ID}).GetByUserID()
		if err == nil {
			ru.ExpireAt = expire.ExpireAt.Format(settings.UserDTLayout)
		} else {
			ru.ExpireAt = "undefined"
		}

		res = append(res, ru)
	}

	return res, pages, nil
}

func GetUserDomains(userID int64) ([]domains.Domain, rest_errors.Error) {
	user, err := (&users.User{ID: userID}).GetByID()
	if err != nil {
		return nil, rest_errors.NewNotFoundError("User not found")
	}

	if user.IsAdmin() {
		all, err := domains.GetAllDomains()
		if err != nil {
			return nil, rest_errors.NewNotFoundError("Unable to find any domain")
		}

		return all, nil
	}

	res := make([]domains.Domain, 0)
	for _, roleID := range user.RoleIDs {
		all, err := policies.GetByRole(roleID)
		if err == nil {
			for _, p := range all {
				if p.IsAllProductPolicy() || p.IsDomainPolicy() {
					dom, err := (&domains.Domain{ID: p.DomainID}).GetByID()
					if err == nil {
						res = uniqueInsert(*dom, res)
					}

				} else if p.IsProductPolicy() {
					prod, err := (&products.Product{ID: p.ProductID}).GetByID()
					dom, err := (&domains.Domain{ID: prod.DomainID}).GetByID()
					if err == nil {
						res = uniqueInsert(*dom, res)
					}
				}
			}
		}
	}

	return res, nil
}

// GetProducts returns the products in the givne domain that current user
// can access to it
func GetUserProducts(userID, domainID int64) ([]products.Product, rest_errors.Error) {
	user, err := (&users.User{ID: userID}).GetByID()
	if err != nil {
		return nil, rest_errors.NewNotFoundError("User not found")
	}

	if user.IsAdmin() {
		all, err := products.GetAllProducts()
		if err != nil {
			return nil, rest_errors.NewNotFoundError("Unable to find any product")
		}

		return all, nil
	}

	allProds := make([]products.Product, 0)
	for _, roleID := range user.RoleIDs {
		allPol, err := policies.GetByRole(roleID)
		if err == nil {
			for _, p := range allPol {
				if p.IsAllProductPolicy() && p.DomainID == domainID {
					_, err := (&domains.Domain{ID: p.DomainID}).GetByID()
					if err == nil {
						prods, _ := products.GetProductsByDomain(domainID)
						allProds = append(allProds, prods...)
					}

				} else if p.IsProductPolicy() && p.DomainID == domainID {
					prod, err := (&products.Product{ID: p.ProductID}).GetByID()
					if err == nil {
						allProds = append(allProds, *prod)
					}
				}
			}
		}
	}

	return allProds, nil
}

func AddUser(user *users.User) rest_errors.Error {
	err := user.Save()
	if err != nil {
		return rest_errors.NewNotAcceptableError(err.Error())
	}

	if user.Email != "" {
		err = email.SendWelcome(user.Name, user.Email)
		if err != nil {
			fmt.Printf("Error while sending welcom email: %v \n", err)
		} else {
			fmt.Println("Welcom email sent")
		}
	}

	return nil
}

func AddUserWithRole(user *users.User, sendEmail bool, expireAt *time.Time) rest_errors.Error {
	err := user.SaveWithRoles()
	if err != nil {
		return rest_errors.NewNotAcceptableError(err.Error())
	}

	if expireAt != nil {
		expire := &access.Expire{UserID: user.ID, ExpireAt: *expireAt}

		if err := expire.Save(); err != nil {
			user.RemoveByID()
			return rest_errors.NewNotAcceptableError(err.Error())
		}
	}

	if sendEmail && user.Email != "" {
		err = email.SendWelcome(user.Name, user.Email)
		if err != nil {
			fmt.Printf("Error while sending welcome email: %v \n", err)
		} else {
			fmt.Println("Welcome email sent")
		}
	}

	return nil
}

func UpdateUser(user *users.User) rest_errors.Error {
	err := user.Update()
	if err != nil {
		return rest_errors.NewNotAcceptableError(err.Error())
	}

	return nil
}

func DelUser(userID int64) rest_errors.Error {
	err := (&users.User{ID: int64(userID)}).RemoveByID()
	if err != nil {
		return rest_errors.NewInternalServerError("Database error", err)
	}

	return nil
}

// GetPolices returns all policies that exists in user's roles
func GetPolices(user *users.User) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)

	if user.IsAdmin() {
		return res, nil
	}

	for _, id := range user.RoleIDs {
		rolePolicies, err := (&policies.Policy{RoleID: id}).GetNamedPolicies()
		if err != nil {
			continue
		}

		res = append(res, rolePolicies...)
	}

	return res, nil
}

// uniqueInsert inserts an int64 number into the slice of int64
// only if it's not duplicated
func uniqueInsert(dom domains.Domain, list []domains.Domain) []domains.Domain {
	for _, d := range list {
		if d.ID == dom.ID {
			return list
		}
	}

	return append(list, dom)
}
