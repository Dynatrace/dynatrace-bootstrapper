package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToMap(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		type test struct {
			Beep string `json:"beep"`
			Boop string `json:"boop"`
		}

		expectedMap := map[string]string{
			"beep": "boop",
			"boop": "beep",
		}

		input := test{
			Beep: "boop",
			Boop: "beep",
		}

		output, err := ToMap(input)
		require.NoError(t, err)
		assert.Equal(t, expectedMap, output)
	})

	t.Run("fails, if not only string", func(t *testing.T) {
		type test struct {
			Beep int
			Boop string
		}

		input := test{
			Beep: 1,
			Boop: "beep",
		}

		output, err := ToMap(input)
		require.Error(t, err)
		assert.Nil(t, output)
	})
}

func TestFromMap(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		type test struct {
			Beep string `json:"beep"`
			Boop string `json:"boop"`
		}

		input := map[string]string{
			"beep": "boop",
			"boop": "beep",
		}

		expected := test{
			Beep: "boop",
			Boop: "beep",
		}

		output, err := FromMap[test](input)
		require.NoError(t, err)
		assert.Equal(t, expected, *output)
	})

	t.Run("fails, if not only string", func(t *testing.T) {
		type test struct {
			Beep int    `json:"beep"`
			Boop string `json:"boop"`
		}

		input := map[string]string{
			"beep": "1",
			"boop": "beep",
		}

		output, err := FromMap[test](input)
		require.Error(t, err)
		assert.Nil(t, output)
	})
}
