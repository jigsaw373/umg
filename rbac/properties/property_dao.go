package properties

import (
	"errors"

	"github.com/boof/umg/db"
	"github.com/boof/umg/rest_errors"
	"github.com/boof/umg/util/validator"
)

func init() {
	db.Sync(new(Property))
}

func (p *Property) Save() error {
	if _, err := p.ValidateForInsert(); err != nil {
		return err
	}

	// remove id
	p.ID = 0

	_, err := db.Engine.Insert(p)
	return err
}

func (p *Property) GetByID() (*Property, error) {
	property := &Property{ID: p.ID}
	if has, err := db.Engine.Get(property); !has || err != nil {
		return nil, errors.New("property not found")
	}

	return property, nil
}

func (p *Property) GetByMeteringID() (*Property, error) {
	property := &Property{MeteringID: p.MeteringID, Type: p.Type}
	if has, err := db.Engine.Get(property); !has || err != nil {
		return nil, errors.New("property not found")
	}

	return property, nil
}

func (p *Property) ValidateForInsert() (*Property, error) {
	if !p.HasValidType() {
		return nil, errors.New("invalid property type")
	}

	errs := validator.Validate.Var(p.Name, "min=1,max=64")
	if errs != nil {
		return nil, errors.New("password must have a length between 1 and 64")
	}

	if _, err := p.GetByName(); err == nil {
		return nil, errors.New("duplicated property name")
	}

	return p, nil
}

func (p *Property) GetByName() (*Property, error) {
	property := &Property{Type: p.Type, Name: p.Name}
	if has, err := db.Engine.Get(property); !has || err != nil {
		return property, errors.New("property not found")
	}

	return property, nil
}

func (p *Property) HasValidType() bool {
	if p.Type == "CARMA" || p.Type == "METER" {
		return true
	}

	return false
}

// RemoveByID removes the property by id
func (p *Property) RemoveByID() error {
	_, err := db.Engine.Id(p.ID).Delete(&Property{})
	return err
}

// GetProperties returns all domains
func GetAllProperties() ([]Property, rest_errors.Error) {
	var properties []Property
	err := db.Engine.Find(&properties)
	if err != nil {
		return nil, rest_errors.NewNotFoundError("Unable to find any properties")

	}

	return properties, nil
}
