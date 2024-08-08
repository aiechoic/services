package ioc

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewVipers(t *testing.T) {
	dir := t.TempDir()
	env := ConfigEnvTest

	// Create test config files
	jsonContent := `{"key": "value"}`

	err := os.WriteFile(filepath.Join(dir, "config.test.json"), []byte(jsonContent), 0644)
	assert.NoError(t, err)

	config, err := NewVipers(dir, env)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Len(t, config.encoders, 1)
}

func TestVipers_Unmarshal(t *testing.T) {
	dir := t.TempDir()
	env := ConfigEnvTest

	// Create test config file
	jsonContent := `{"key": "value"}`
	err := os.WriteFile(filepath.Join(dir, "config.test.json"), []byte(jsonContent), 0644)
	assert.NoError(t, err)

	config, err := NewVipers(dir, env)
	assert.NoError(t, err)

	var result map[string]string
	err = config.Unmarshal("config", &result, nil)
	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])
}

func TestVipers_UnmarshalWithDefault(t *testing.T) {
	dir := t.TempDir()
	env := ConfigEnvTest

	config, err := NewVipers(dir, env)
	assert.NoError(t, err)

	defaultContent := []byte(`{"key": "default"}`)
	var result map[string]string
	err = config.Unmarshal("config", &result, defaultContent)
	assert.NoError(t, err)
	assert.Equal(t, "default", result["key"])

	// Verify that the default config file is created
	createdFile := filepath.Join(dir, "config.test.json")
	_, err = os.Stat(createdFile)
	assert.NoError(t, err)
}

func TestVipers_getContentEncoder(t *testing.T) {
	config := &Vipers{}

	jsonContent := []byte(`{"key": "value"}`)
	ct := config.getContentType(jsonContent)
	assert.Equal(t, "json", ct)

	yamlContent := []byte(`key: value`)
	ct = config.getContentType(yamlContent)
	assert.Equal(t, "yaml", ct)

	tomlContent := []byte(`key = "value"`)
	ct = config.getContentType(tomlContent)
	assert.Equal(t, "toml", ct)
}

func TestVipers_UnmarshalWithEnv(t *testing.T) {
	dir := t.TempDir()
	env := ConfigEnvTest

	type Nested struct {
		Age int `mapstructure:"age"`
	}

	type T struct {
		Key    string  `mapstructure:"key"`
		Nested *Nested `mapstructure:"nested"`
	}

	// Create test config file
	yamlContent := `
key: "value"
nested:
  age: 21
`
	err := os.WriteFile(filepath.Join(dir, "config.test.yaml"), []byte(yamlContent), 0644)
	assert.NoError(t, err)

	config, err := NewVipers(dir, env)
	assert.NoError(t, err)

	// Set environment variables
	err = os.Setenv("CONFIG_KEY", "env_value")
	assert.NoError(t, err)
	err = os.Setenv("CONFIG_NESTED_AGE", "18")
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("CONFIG_KEY")
		_ = os.Unsetenv("CONFIG_NESTED_AGE")
	}()

	var result T
	err = config.Unmarshal("config", &result, nil)
	assert.NoError(t, err)
	assert.Equal(t, "env_value", result.Key)
	assert.Equal(t, 18, result.Nested.Age)
}

func TestVipers_UnmarshalAndWatch(t *testing.T) {
	dir := t.TempDir()
	env := ConfigEnvTest

	// Create test config file
	jsonContent := `{"key": "value"}`
	err := os.WriteFile(filepath.Join(dir, "config.test.json"), []byte(jsonContent), 0644)
	assert.NoError(t, err)

	config, err := NewVipers(dir, env)
	assert.NoError(t, err)

	var result map[string]string

	callback := func(v *viper.Viper) {
		err := v.Unmarshal(&result)
		assert.NoError(t, err)
	}

	err = config.UnmarshalAndWatch("config", []byte(`{"key": "default"}`), callback)
	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])

	// Modify the config file
	newJsonContent := `{"key": "new_value"}`
	err = os.WriteFile(filepath.Join(dir, "config.test.json"), []byte(newJsonContent), 0644)
	assert.NoError(t, err)

	// Wait for the callback to be called
	time.Sleep(1 * time.Second)
	assert.Equal(t, "new_value", result["key"])
}
