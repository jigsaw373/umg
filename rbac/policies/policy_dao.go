package policies

import (
	"errors"
	"strings"

	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/domains"
	"github.com/boof/umg/rbac/products"
	"github.com/boof/umg/rbac/properties"
	"github.com/boof/umg/rbac/roles"
	"github.com/boof/umg/util/validator"
)

const (
	// ProdPolicy indicates that the policy is for a product
	ProdPolicy = "P"

	// AllProductPolicy indicates that the policy is for the all products in a domain
	AllProdPolicy = "A"

	// DomPolicy indicates that the policy is for a domain
	DomPolicy = "D"
)

func init() {
	db.Sync(new(Policy))
}

// Save inserts a new role to the database
func (p *Policy) Save() error {
	if err := p.ValidateForInsert(); err != nil {
		return err
	}

	// remove id
	p.ID = 0

	_, err := db.Engine.Insert(p)
	return err
}

// RemoveByID removes the policy by id
func (p *Policy) RemoveByID() error {
	_, err := db.Engine.Id(p.ID).Delete(&Policy{})
	return err
}

// RemoveByRoleID deletes all policies of the current role
func (p *Policy) RemoveByRoleID() error {
	_, err := db.Engine.Delete(&Policy{RoleID: p.RoleID})
	return err
}

// IsProductPolicy indicates that the current policy is a product policy or not
func (p *Policy) IsProductPolicy() bool {
	return p.Type == ProdPolicy
}

// IsDomainPolicy indicates that the current policy is a Domain policy or not
func (p *Policy) IsDomainPolicy() bool {
	return p.Type == DomPolicy
}

// IsAllProductPolicy indicates that the current policy is an `All Product Policy` or not
func (p *Policy) IsAllProductPolicy() bool {
	return p.Type == AllProdPolicy
}

// LetProdAct indicates that the current policy can let the given action
// to be carry out on the given product or not
func (p *Policy) LetProdAct(domId int64, prodID int64, action string) bool {
	if p.IsProductPolicy() {
		return p.ProductID == prodID && p.HasAction(action)
	} else if p.IsAllProductPolicy() {
		return p.DomainID == domId && p.HasAction(action)
	}

	return false
}

// LetDomAct indicates that the current policy can let the given action
// to be carry out on the given domain or not
func (p *Policy) LetDomAct(domID int64, action string) bool {
	return p.IsDomainPolicy() && p.DomainID == domID && p.HasAction(action)
}

// ValidateForInsert validates the policy
func (p *Policy) ValidateForInsert() error {
	if p.Type != DomPolicy && p.Type != ProdPolicy && p.Type != AllProdPolicy {
		return errors.New("invalid policy type")
	}

	if len(p.Actions) == 0 {
		return errors.New("actions can't be empty")
	}

	for _, act := range p.Actions {
		if act == "" {
			return errors.New("invalid action")
		} else if act == "*" && len(p.Actions) > 1 {
			return errors.New("you can't combine * action with other ones")
		}
	}

	errs := validator.Validate.Var(p.Actions, "dive,min=1,max=32")
	if errs != nil {
		return errors.New("actions text length should be between 1 and 32")
	}

	role := &roles.Role{ID: p.RoleID}
	if has, _ := db.Engine.Get(role); p.RoleID < 1 || !has {
		return errors.New("invalid role")
	}

	if role.IsAdmin() {
		return errors.New("you can't add policy to admin role")
	}

	if has, _ := db.Engine.Get(&domains.Domain{ID: p.DomainID}); p.DomainID < 1 || !has {
		return errors.New("invalid domain")
	}

	if p.IsProductPolicy() {
		if has, _ := db.Engine.Get(&products.Product{ID: p.ProductID}); p.ProductID < 1 || !has {
			return errors.New("invalid product")
		}
	}

	for _, pID := range p.Properties {
		if _, err := (&properties.Property{ID: pID}).GetByID(); err != nil {
			return errors.New("invalid property")
		}
	}

	return nil
}

// HasAction indicates that policy has given action or not
func (p *Policy) HasAction(action string) bool {
	for _, act := range p.Actions {
		if act == "*" || strings.ToLower(act) == strings.ToLower(action) {
			return true
		}
	}

	return false
}

// GetPolicies returns all policies for the current role
func (p *Policy) GetRolePolicies() ([]Policy, error) {
	var policies []Policy
	err := db.Engine.Where("role_id = ?", p.RoleID).Find(&policies)

	return policies, err
}

// GetNamedPolicies returns all named policies for current role
func (p *Policy) GetNamedPolicies() ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)

	policies, err := p.GetRolePolicies()
	if err != nil {
		return res, err
	}

	for _, policy := range policies {
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

		allProperties := make([]*properties.Property, 0)
		for _, pID := range policy.Properties {
			if property, err := (&properties.Property{ID: pID}).GetByID(); err == nil {
				allProperties = append(allProperties, property)
			}
		}

		p["properties"] = allProperties

		res = append(res, p)
	}

	return res, nil
}

func GetByRole(roleID int64) ([]Policy, error) {
	var policies []Policy
	err := db.Engine.Where("role_id = ?", roleID).Find(&policies)

	return policies, err
}
