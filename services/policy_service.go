package services

import (
	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/domains"
	"github.com/boof/umg/rbac/policies"
	"github.com/boof/umg/rbac/products"
	"github.com/boof/umg/rbac/properties"
	"github.com/boof/umg/rest_errors"
)

func GetNamedPolicies(roleID int64) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)

	all, err := policies.GetByRole(roleID)
	if err != nil {
		return res, err
	}

	for _, policy := range all {
		p := make(map[string]interface{})

		if policy.IsProductPolicy() {
			prod, err := (&products.Product{ID: policy.ProductID}).GetByID()
			if err != nil {
				continue
			}
			dom, err := (&domains.Domain{ID: prod.DomainID}).GetByID()
			if err != nil {
				continue
			}

			p["product"] = prod.Name
			p["domain"] = dom.Name
			p["actions"] = policy.Actions

		} else if policy.IsDomainPolicy() {
			dom, err := (&domains.Domain{ID: policy.DomainID}).GetByID()
			if err != nil {
				continue
			}

			p["domain"] = dom.Name
			p["actions"] = policy.Actions

		} else if policy.IsAllProductPolicy() {
			dom, err := (&domains.Domain{ID: policy.DomainID}).GetByID()
			if err != nil {
				continue
			}

			p["domain"] = dom.Name
			p["actions"] = policy.Actions
		} else {
			continue
		}

		p["id"] = policy.ID
		p["type"] = policy.Type

		props := make([]*properties.Property, 0)
		for _, pID := range policy.Properties {
			if property, err := (&properties.Property{ID: pID}).GetByID(); err == nil {
				props = append(props, property)
			}
		}

		p["properties"] = props

		res = append(res, p)
	}

	return res, nil
}

func removePoliciesByRole(roleID int64) rest_errors.Error {
	_, err := db.Engine.Delete(&policies.Policy{RoleID: roleID})
	if err != nil {
		return rest_errors.NewInternalServerError(err.Error(), err)
	}

	return nil
}
