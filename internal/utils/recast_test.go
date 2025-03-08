package utils_test

import (
	"testing"

	"github.com/toudi/gorpheus/v2/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestRecast(t *testing.T) {
	type TestStruct struct {
		SomeField string `yaml:"some_field"`
	}
	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		var inputMap = map[string]interface{}{
			"some_field": "some value",
		}
		var testStruct TestStruct
		err := utils.Recast(inputMap, &testStruct)
		require.NoError(t, err)
		require.Equal(t, "some value", testStruct.SomeField)
	})
}
