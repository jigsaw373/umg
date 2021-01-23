package services

import (
	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/access"
	"github.com/boof/umg/rbac/domains"
	"github.com/boof/umg/rbac/policies"
	"github.com/boof/umg/rbac/products"
	"github.com/boof/umg/rbac/properties"
	"github.com/boof/umg/rbac/roles"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/rest_errors"
)

func HasProdPerm(user *users.User, domName, prodName, action string) bool {
	if ok, _ := Expired(user.ID); ok {
		return false
	}

	dom, err := (&domains.Domain{Name: domName}).GetByName()
	if err != nil {
		return false
	}

	prod := &products.Product{DomainID: dom.ID, Name: prodName}
	if has, err := db.Engine.Get(prod); !has || err != nil {
		return false
	}

	has, _ := hasProdPerm(user, dom.ID, prod.ID, action)
	return has
}

func HasDomPerm(user *users.User, domName string, action string) bool {
	if ok, _ := access.Expired(user.ID); ok {
		return false
	}

	dom, err := (&domains.Domain{Name: domName}).GetByName()
	if err != nil {
		return false
	}

	has, _ := hasDomPerm(user, dom.ID, action)
	return has
}

func HasPropertyPerm(user *users.User, meteringID int64, pType string) bool {
	property, err := (&properties.Property{MeteringID: meteringID, Type: pType}).GetByMeteringID()
	if err != nil {
		return false
	}

	for _, roleID := range user.RoleIDs {
		r, err := (&roles.Role{ID: roleID}).GetByID()
		if err == nil && r.IsAdmin() {
			return true
		}

		//policies, err := (&roles.Role{ID: roleID}).GetPolicies()
		pols, err := (&policies.Policy{RoleID: roleID}).GetRolePolicies()
		if err == nil {
			for _, policy := range pols {
				for _, pID := range policy.Properties {
					if pID == property.ID {
						return true
					}
				}
			}
		}
	}

	return false
}

func hasProdPerm(user *users.User, domID, prodID int64, action string) (bool, rest_errors.Error) {
	// validate product
	if has, _ := db.Engine.ID(prodID).Get(&products.Product{}); !has {
		return false, rest_errors.NewNotFoundError("Invalid product")
	}

	// check for permission
	// allow-override algorithm
	for _, roleID := range user.RoleIDs {
		r, err := (&roles.Role{ID: roleID}).GetByID()
		if err == nil && r.IsAdmin() {
			return true, nil
		}

		allPolicies, err := policies.GetByRole(roleID)
		if err == nil {
			for _, p := range allPolicies {
				if p.LetProdAct(domID, prodID, action) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func hasDomPerm(user *users.User, domID int64, action string) (bool, error) {
	// validate domain
	if has, _ := db.Engine.ID(domID).Get(&domains.Domain{}); !has {
		return false, rest_errors.NewNotFoundError("Invalid domain")
	}

	// check for permission
	// allow-override algorithm
	for _, roleID := range user.RoleIDs {
		r, err := (&roles.Role{ID: roleID}).GetByID()
		if err == nil && r.IsAdmin() {
			return true, nil
		}

		rolePolicies, err := policies.GetByRole(roleID)
		if err == nil {
			for _, p := range rolePolicies {
				if p.LetDomAct(domID, action) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
