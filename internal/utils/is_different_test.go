package utils_test

import (
	"fmt"
	"testing"

	"github.com/toudi/gorpheus/v2/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestIsDifferent(t *testing.T) {
	type testCase struct {
		a map[string]interface{}
		b map[string]interface{}
	}
	t.Run("test same maps", func(t *testing.T) {
		for i, test := range []testCase{
			{
				a: map[string]interface{}{},
				b: map[string]interface{}{},
			},
			{
				a: map[string]interface{}{"a": "a"},
				b: map[string]interface{}{"a": "a"},
			},
			{
				a: map[string]interface{}{"a": 1},
				b: map[string]interface{}{"a": 1},
			},
		} {
			t.Run(fmt.Sprintf("equality case %d", i), func(t *testing.T) {
				t.Parallel()

				require.False(t, utils.IsDifferent(test.a, test.b))
			})
		}
	})

	t.Run("test different maps", func(t *testing.T) {
		for i, test := range []testCase{
			{
				a: map[string]interface{}{},
				b: map[string]interface{}{"a": "a"},
			},
			{
				a: map[string]interface{}{"a": "a"},
				b: map[string]interface{}{},
			},
			{
				a: map[string]interface{}{"a": "a"},
				b: map[string]interface{}{"a": "b"},
			},
			{
				a: map[string]interface{}{"a": 1},
				b: map[string]interface{}{"a": 2},
			},
			{
				a: map[string]interface{}{"a": "a"},
				b: map[string]interface{}{"a": 1},
			},
			{
				a: map[string]interface{}{"a": "a"},
				b: map[string]interface{}{},
			},
		} {
			t.Run(fmt.Sprintf("inequality case %d", i), func(t *testing.T) {
				t.Parallel()

				require.True(t, utils.IsDifferent(test.a, test.b))
			})
		}
	})
}
