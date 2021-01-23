package domains

import (
	"errors"
	"strings"

	"github.com/boof/umg/db"
	"github.com/boof/umg/util/validator"
)

func init() {
	db.Sync(new(Domain))
}

// Save inserts a new domain into the database
func (d *Domain) Save() error {
	err := d.validateForInsert()
	if err != nil {
		return err
	}

	// remove id
	d.ID = 0

	_, err = db.Engine.Insert(d)
	return err
}

// GetByID returns a Demand with the given id
func (d *Domain) GetByID() (*Domain, error) {
	domain := &Domain{ID: d.ID}
	if has, err := db.Engine.Get(domain); !has || err != nil {
		return domain, errors.New("domain not found")
	}

	return domain, nil
}

// GetByName returns a Demand with the given name
func (d *Domain) GetByName() (*Domain, error) {
	domain := &Domain{Name: d.Name}
	if has, err := db.Engine.Get(domain); !has || err != nil {
		return domain, errors.New("domain not found")
	}

	return domain, nil
}

// RemoveByID removes the domain by id
func (d *Domain) RemoveByID() error {
	_, err := db.Engine.Id(d.ID).Delete(&Domain{})
	return err
}

// RemoveByName removes the domain by name
func (d *Domain) RemoveByName() error {
	_, err := db.Engine.Delete(&Domain{Name: d.Name})
	return err
}

// validateForInsert validates the domain
func (d *Domain) validateForInsert() error {
	if d.Name == "*" {
		return errors.New("* is a reserved name")
	}

	errs := validator.Validate.Var(d.Name, "min=1,max=64")
	if errs != nil {
		return errors.New("domain name should have a length between 1 and 64")
	}

	var domains []Domain
	err := db.Engine.SQL("SELECT * FROM domain WHERE LOWER(name) = ?", strings.ToLower(d.Name)).Find(&domains)
	if err != nil {
		return errors.New("database error")
	} else if domains != nil && len(domains) > 0 {
		return errors.New("duplicated domain name")
	}

	return nil
}

// GetDomains returns all domains
func GetAllDomains() ([]Domain, error) {
	var domains []Domain
	err := db.Engine.Find(&domains)

	return domains, err
}
