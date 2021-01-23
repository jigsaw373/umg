package initialize

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/boof/umg/rbac/domains"
	"github.com/boof/umg/rbac/products"
)

type Config struct {
	Domains []struct {
		Name     string `yaml:"name"`
		Products []struct {
			Name string `yaml:"name"`
		} `yaml:"products"`
	} `yaml:"domains"`
}

func loadConfig() (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open("./initialize/domains.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func createDomains() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("unable to load config file: %v", err)
	}

	for _, domain := range config.Domains {
		dom, err := createDomain(domain.Name)
		if err != nil {
			log.Fatalf("unable to create %s Domain: %v", domain.Name, err)
		}

		for _, product := range domain.Products {
			_, err = createProduct(dom.ID, product.Name)
			if err != nil {
				log.Fatalf("unable to create %s product: %v", product.Name, err)
			}
		}
	}
}

// createDomain creates a domain if not exists
func createDomain(name string) (*domains.Domain, error) {
	domain, err := (&domains.Domain{Name: name}).GetByName()
	if err != nil {
		// try to create domain
		domain = &domains.Domain{Name: name}
		if err := domain.Save(); err != nil {
			return domain, err
		}
	}

	return domain, nil
}

// createProduct creates a product if not exists
func createProduct(domainID int64, name string) (*products.Product, error) {
	product, err := (&products.Product{Name: name, DomainID: domainID}).GetByName()
	if err != nil {
		// try to create domain
		product = &products.Product{Name: name, DomainID: domainID}
		if err := product.Save(); err != nil {
			return product, err
		}
	}

	return product, nil
}
