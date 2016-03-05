package constant

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Pool struct {
	mutex    sync.RWMutex
	prefix   string
	defaults map[string]*string
}

func NewPool(prefix string) *Pool {
	return &Pool{
		prefix:   prefix,
		defaults: make(map[string]*string),
	}
}

func (pool *Pool) New(name string, default_val interface{}) error {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if pool.defaults[name] != nil {
		return errors.New("Constant already exists")
	}

	var str_val string
	switch t := default_val.(type) {
	case string:
		if val, ok := default_val.(string); ok {
			str_val = val
		} else {
			return errors.New("Unabled to assert type string on default value")
		}
	case fmt.Stringer:
		if val, ok := default_val.(fmt.Stringer); ok {
			str_val = val.String()
		} else {
			return errors.New("Unabled to assert type fmt.Stringer on default value")
		}
	case int:
		if val, ok := default_val.(int); ok {
			str_val = strconv.Itoa(val)
		} else {
			return errors.New("Unabled to assert type int on default value")
		}
	case float64:
		if val, ok := default_val.(float64); ok {
			str_val = strconv.FormatFloat(val, 'f', -1, 64)
		} else {
			return errors.New("Unabled to assert type float64 on default value")
		}
	case bool:
		if val, ok := default_val.(bool); ok {
			str_val = strconv.FormatBool(val)
		} else {
			return errors.New("Unabled to assert type bool on default value")
		}
	default:
		return errors.New(fmt.Sprintf("Unexpected type %T", t))
	}

	pool.defaults[name] = new(string)
	*pool.defaults[name] = str_val

	return nil
}

func (pool *Pool) Delete(name string) error {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if pool.defaults[name] == nil {
		return errors.New("Constant doesn't exists")
	}

	delete(pool.defaults, name)
	return nil
}

func (pool *Pool) Prefix() string {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	return pool.prefix
}

func (pool *Pool) List() []string {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	consts := make([]string, 0, len(pool.defaults))
	for c := range pool.defaults {
		consts = append(consts, c)
	}
	return consts
}

func (pool *Pool) Environment() []string {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	consts := pool.List()
	for i := 0; i < len(consts); i++ {
		consts[i] = pool.env_name(consts[i])
	}
	return consts
}
