package products

import (
	"errors"

	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/domains"
	"github.com/boof/umg/util/validator"
)

func init() {
	db.Sync(new(Product))
}

// Save inserts a new product into the database
func (p *Product) Save() error {
	err := p.validateForInsert()
	if err != nil {
		return err
	}

	// remove id
	p.ID = 0

	_, err = db.Engine.Insert(p)
	return err
}

// GetByID returns a Product with the given id
func (p *Product) GetByID() (*Product, error) {
	prod := &Product{ID: p.ID}
	if has, err := db.Engine.Get(prod); !has || err != nil {
		return prod, errors.New("product not found")
	}

	return prod, nil
}

// GetByName returns a Product with the given Name and DomainID
func (p *Product) GetByName() (*Product, error) {
	prod := &Product{Name: p.Name, DomainID: p.DomainID}
	if has, err := db.Engine.Get(prod); !has || err != nil {
		return prod, errors.New("product not found")
	}

	return prod, nil
}

// RemoveByID removes the product by id
func (p *Product) RemoveByID() error {
	_, err := db.Engine.Id(p.ID).Delete(&Product{})
	return err
}

// RemoveByName removes the product by name
func (p *Product) RemoveByName() error {
	_, err := db.Engine.Delete(&Product{Name: p.Name})
	return err
}

// validateForInsert validates
func (p *Product) validateForInsert() error {
	if p.Name == "" || p.Name == "*" {
		return errors.New("invalid domain name")
	}

	errs := validator.Validate.Var(p.Name, "min=1,max=64")
	if errs != nil {
		return errors.New("product name should have a length between 1 and 64")
	}

	// each product should be subset of a domain
	if has, _ := db.Engine.ID(p.DomainID).Get(&domains.Domain{}); p.DomainID < 1 && !has {
		return errors.New("invalid domain")
	}

	// check for duplicated product name
	has, _ := db.Engine.Get(&Product{Name: p.Name, DomainID: p.DomainID})
	if has {
		return errors.New("duplicated product name")
	}

	return nil
}

// GetAllProducts returns all products
func GetAllProducts() ([]Product, error) {
	var products []Product
	err := db.Engine.Find(&products)

	return products, err
}

// GetProducts returns all products that are subset of current domain
func GetProductsByDomain(domainID int64) ([]Product, error) {
	var products []Product
	err := db.Engine.Where("domain_id = ?", domainID).Find(&products)

	return products, err
}

// GetSoloProducts returns all products that don't belong to any domain
func GetSoloProducts() ([]Product, error) {
	var products []Product
	err := db.Engine.Where("id = 0").Find(&products)

	return products, err
}
