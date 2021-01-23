package services

import (
	"github.com/boof/umg/rbac/domains"
	"github.com/boof/umg/rbac/policies"
	"github.com/boof/umg/rbac/properties"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/rest_errors"
)

func GetProperties(user *users.User) ([]properties.Property, rest_errors.Error) {
	if user.IsAdmin() {
		return properties.GetAllProperties()
	}

	props := make([]properties.Property, 0)
	for _, roleID := range user.RoleIDs {
		policies, err := (&policies.Policy{RoleID: roleID}).GetRolePolicies()
		if err == nil {
			for _, policy := range policies {
				for _, pID := range policy.Properties {
					if property, err := (&properties.Property{ID: pID}).GetByID(); err == nil {
						has := false
						for _, p := range props {
							if p.ID == property.ID {
								has = true
							}
						}
						if !has {
							props = append(props, *property)
						}
					}
				}
			}
		}
	}

	return props, nil
}

func GetPropertiesForDomain(user *users.User, domainName string) ([]properties.Property, rest_errors.Error) {
	if user.IsAdmin() {
		return properties.GetAllProperties()
	}

	domain, err := (&domains.Domain{Name: domainName}).GetByName()
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	props := make([]properties.Property, 0)
	for _, roleID := range user.RoleIDs {
		pols, err := (&policies.Policy{RoleID: roleID}).GetRolePolicies()
		if err == nil {
			for _, policy := range pols {
				if policy.DomainID == domain.ID {
					for _, pID := range policy.Properties {
						if property, err := (&properties.Property{ID: pID}).GetByID(); err == nil {
							has := false
							for _, p := range props {
								if p.ID == property.ID {
									has = true
								}
							}
							if !has {
								props = append(props, *property)
							}
						}
					}
				}
			}
		}
	}

	return props, nil
}
