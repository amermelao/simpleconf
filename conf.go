package simpleconf

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type (
	ConfValidator interface {
		Valid() bool
		Enabled() bool
		setEnabled()
	}

	EnableConf struct {
		enable bool
	}
)

func (e *EnableConf) setEnabled() {
	e.enable = true
}

func (e EnableConf) Enabled() bool {
	return e.enable
}

var ErrInvalidConfiguration = errors.New("invalid configuration")

func readKey(conf ConfValidator) error {
	keyName := strings.ToLower(getType(conf))

	if viper.IsSet(keyName) {
		conf.setEnabled()
		return readSetConfig(conf, keyName)
	} else {
		return nil
	}
}

func readSetConfig[T ConfValidator](conf T, keyName string) error {
	if err := viper.UnmarshalKey(keyName, conf); err != nil {
		return fmt.Errorf("reading config [%s]: %w", keyName, err)
	} else if !conf.Valid() {
		return fmt.Errorf(
			"config [%s] is invalid: %w",
			keyName,
			ErrInvalidConfiguration,
		)
	} else {
		return nil
	}
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

type (
	readKeyConfiguration func(v ConfValidator) error
	auxReader            func(readKeyConfiguration) error
)

var envPrefix string

func SetEnvPrefix(s string) {
	envPrefix = s
}

func ReadConfiguration(f ...auxReader) error {
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		viper.Set(key, val)
	}

	for _, aF := range f {
		if err := aF(readKey); err != nil {
			return fmt.Errorf("eading key: %w", err)
		}
	}

	return nil
}

func OneSubset[T ConfValidator](conf T) auxReader {
	return func(f readKeyConfiguration) error {
		return f(conf)
	}
}
