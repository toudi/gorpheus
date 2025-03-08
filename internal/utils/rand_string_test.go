package utils_test

import (
	"testing"

	"github.com/toudi/gorpheus/v2/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestRandString(t *testing.T) {
	t.Run("RandStringBytes", func(t *testing.T) {
		// well, you know the famous joke about "3, 3, 3, 3" and the fact
		// that you can never tell whether it is random or not.
		// what we can do, however is just to check if the next generated
		// string is different from the previous one and that they all
		// are of length n.

		var previousString string
		for i := 0; i < 10; i += 1 {
			randomString := utils.RandStringBytes(10)
			require.Len(t, randomString, 10)
			require.NotEqual(t, previousString, randomString)
			previousString = randomString
		}
	})
}
