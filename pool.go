package constant

import (
	"errors"
	"fmt"
	"strconv"
)

type Pool struct {
	prefix   string
	defaults map[string]string
}

func NewPool(prefix string) *Pool {
	return &Pool{
		prefix:   prefix,
		defaults: make(map[string]string),
	}
}

func (pool *Pool) New(name string, default_val interface{}) error {
	if pool.defaults[name] != "" {
		return errors.New("Constant already exists")
	}

	switch t := default_val.(type) {
	case string:
		if val, ok := default_val.(string); ok {
			pool.defaults[name] = val
		} else {
			return errors.New("Unabled to assert type string on default value")
		}
	case fmt.Stringer:
		if val, ok := default_val.(fmt.Stringer); ok {
			pool.defaults[name] = val.String()
		} else {
			return errors.New("Unabled to assert type fmt.Stringer on default value")
		}
	case int:
		if val, ok := default_val.(int); ok {
			pool.defaults[name] = strconv.Itoa(val)
		} else {
			return errors.New("Unabled to assert type int on default value")
		}
	case float64:
		if val, ok := default_val.(float64); ok {
			pool.defaults[name] = strconv.FormatFloat(val, 'f', -1, 64)
		} else {
			return errors.New("Unabled to assert type float64 on default value")
		}
	case bool:
		if val, ok := default_val.(bool); ok {
			pool.defaults[name] = strconv.FormatBool(val)
		} else {
			return errors.New("Unabled to assert type bool on default value")
		}
	default:
		return errors.New(fmt.Sprintf("Unexpected type %T\n", t))
	}

	return nil
}

func (pool *Pool) Delete(name string) error {
	if pool.defaults[name] == "" {
		return errors.New("Constant doesn't exists")
	}

	delete(pool.defaults, name)
	return nil
}

func (pool *Pool) Prefix() string {
	return pool.prefix
}

func (pool *Pool) List() []string {
	consts := make([]string, 0, len(pool.defaults))
	for c := range pool.defaults {
		consts = append(consts, c)
	}
	return consts
}

func (pool *Pool) Environment() []string {
	consts := pool.List()
	for i := 0; i < len(consts); i++ {
		consts[i] = pool.env_name(consts[i])
	}
	return consts
}
