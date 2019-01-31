package envutils

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-errors/errors"
)

var ErrInvalidType = errors.New("Invalid type")

func IsInvalidType(err error) bool {
	return err == ErrInvalidType
}

type Env map[string]string

func New() Env {
	return make(Env)
}

func Pair(s string) (key, value string) {
	parts := strings.SplitN(s, "=", 2)
	key, value = parts[0], parts[1]
	if value[0] == '"' {
		value = value[1 : len(value)-1]
	}
	return
}

func (e Env) SetToSys(key ...string) {
	for _, key := range key {
		os.Setenv(key, e[key])
	}
}

func (e *Env) ParseValue(s string, filter ...func(key, value string) bool) {
	e.ParseValues([]string{s}, filter...)
}

func (e *Env) ParseValues(values []string, filter ...func(key, value string) bool) {
	f := func(key, value string) bool {
		return true
	}
	if len(filter) > 0 && filter[0] != nil {
		f = filter[0]
	}

	var key, value string
	for _, s := range values {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		key, value = Pair(s)
		if f(key, value) {
			(*e)[key] = value
		}
	}
}

func (e *Env) Parse(data interface{}, filter ...func(key, value string) bool) (err error) {
	switch dt := data.(type) {
	case string:
		e.ParseString(dt, filter...)
	case []byte:
		e.ParseString(string(dt), filter...)
	case []string:
		e.ParseValues(dt, filter...)
	case io.Reader:
		var data []byte
		if data, err = ioutil.ReadAll(dt); err != nil {
			return err
		}
		e.ParseString(string(data), filter...)
	default:
		err = ErrInvalidType
	}
	return
}

func (e *Env) ParseString(data string, filter ...func(key, value string) bool) {
	e.ParseValues(strings.Split(data, "\n"), filter...)
}

func Get(key string, defaul ...string) (value string) {
	if value = os.Getenv(key); value != "" {
		return
	}
	for _, value = range defaul {
		if value != "" {
			return
		}
	}
	return
}

func FistEnv(key ...string) (value string) {
	for _, key := range key {
		if value = os.Getenv(key); value != "" {
			return
		}
	}
	return
}
