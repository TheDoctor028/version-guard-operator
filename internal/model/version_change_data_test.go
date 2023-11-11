package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectorParse(t *testing.T) {
	t.Run("should parse selector map to string", func(t *testing.T) {
		selector := map[string]string{
			"app": "nginx",
			"env": "dev",
		}

		result := ParseSelector(selector)

		assert.Equal(t, "app=nginx,env=dev", result)
	})
}
