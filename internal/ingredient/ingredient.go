package ingredient

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Role represents the role of an ingredient in a dish.
type Role string

const (
	Protein   Role = "Protein"
	Carb      Role = "Carb"
	Vegetable Role = "Vegetable"
)

// Ingredient represents a single ingredient with a name and role.
type Ingredient struct {
	Name string `yaml:"name"`
	Role Role   `yaml:"role"`
}

// LoadFromFile reads ingredients from a YAML file at the given path.
func LoadFromFile(path string) ([]Ingredient, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ingredients []Ingredient
	if err := yaml.Unmarshal(data, &ingredients); err != nil {
		return nil, err
	}
	return ingredients, nil
}
