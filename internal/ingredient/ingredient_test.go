package ingredient_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"executive-chef/internal/ingredient"
)

func TestLoadFromFile(t *testing.T) {
	data := `- name: Chicken
  role: Protein
- name: Rice
  role: Carb
`
	tmpFile, err := os.CreateTemp(t.TempDir(), "ingredients-*.yaml")
	require.NoError(t, err)
	_, err = tmpFile.Write([]byte(data))
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	ingredients, err := ingredient.LoadFromFile(tmpFile.Name())
	require.NoError(t, err)
	expected := []ingredient.Ingredient{
		{Name: "Chicken", Role: ingredient.Protein},
		{Name: "Rice", Role: ingredient.Carb},
	}
	assert.Equal(t, expected, ingredients)
}
