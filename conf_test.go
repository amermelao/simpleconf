package simpleconf

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readTOMLViper(text string) {
	viper.Reset()
	viper.SetConfigType("TOML")

	viper.ReadConfig(bytes.NewBuffer([]byte(text)))
}

type ValidJWT struct {
	Port int    `mapstructure:"port"`
	Url  string `mapstructure:"url"`

	EnableConf
}

func (v ValidJWT) Valid() bool {

	return true
}

type Other struct {
	Port int    `mapstructure:"port"`
	Url  string `mapstructure:"url"`

	EnableConf
}

func (v Other) Valid() bool {

	return true
}

func TestValidJWT(t *testing.T) {
	readTOMLViper(`
[validjwt]
port = 8081
url = "localhost:3000"

[other]
port = 80800
url = "localhost:13000"
`)

	var validConf ValidJWT
	var otherConf Other

	err := ReadConfiguration(
		OneSubset(&validConf),
		OneSubset(&otherConf),
	)

	require.Nil(t, err)

	assert.Equal(t, 8081, validConf.Port)
	assert.Equal(t, "localhost:3000", validConf.Url)
	assert.True(t, validConf.Enabled(), "should be enabled")

	assert.Equal(t, 80800, otherConf.Port)
	assert.Equal(t, "localhost:13000", otherConf.Url)
	assert.True(t, otherConf.Enabled(), "should be enabled")
}

func TestJWTDisable(t *testing.T) {
	readTOMLViper("")

	var conf ValidJWT
	err := readKey(&conf)

	require.Nil(t, err)

	assert.False(t, conf.Enabled())
}

type invalid struct {
	EnableConf
}

func (invalid) Valid() bool {
	return false
}

func TestInvalidConf(t *testing.T) {
	readTOMLViper(`
[invalid]
`)

	var conf invalid
	err := readKey(&conf)

	assert.True(t, errors.Is(err, ErrInvalidConfiguration))
}
