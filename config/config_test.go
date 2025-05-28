package config_test

import (
	"fmt"
	"testing"

	"github.com/jlrosende/go-agents/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("load config", func(t *testing.T) {

		config, err := config.LoadConfig()

		fmt.Printf("%+v\n", config)

		assert.NoError(t, err)
	})
}
