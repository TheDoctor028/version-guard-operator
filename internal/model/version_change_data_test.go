package model

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSelectorParse(t *testing.T) {
	t.Run("should parse selector map to string", func(t *testing.T) {
		selector := map[string]string{
			"app": "nginx",
			"env": "dev",
		}

		result := ParseSelector(selector)

		assert.True(t, strings.Contains(result, "app=nginx"))
		assert.True(t, strings.Contains(result, "env=dev"))
	})
}
