package constant

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"sync"
)

// Pool represents a collection of constants.
// All constants stored in the same pool can be accessed using the template system.
type Pool struct {
	mutex    sync.RWMutex
	prefix   string
	defaults map[string]*string
}

// Creates a new pool.
//
// Prefix sets the environment variable prefix which is prepended to constants names when searching the runtime environment.
// For example if a pool has a prefix 'MYSQL_' and a constant named 'HOST' then constant 'HOST' would be set to the value of the environment variable 'MYSQL_HOST'.
func NewPool(prefix string) *Pool {
	return &Pool{
		prefix:   prefix,
		defaults: make(map[string]*string),
	}
}

/*
Adds a new constant to the pool.

name: Name of the constant.
Must follow variable naming convention.
Lower case letters, uppercase letters, numbers and underscore.
Can't start with a number.

default_val: The default value for the constant if no environment variable is available.

default_val must be one of the following types:
	string
	[]byte
	fmt.Stringer (https://golang.org/pkg/fmt/#Stringer)
	int
	float64
	bool
*/
func (pool *Pool) New(name string, default_val interface{}) error {
	if !valid_name(name) {
		return errors.New("Invalid Name")
	}

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
	case []byte:
		if val, ok := default_val.([]byte); ok {
			str_val = string(val)
		} else {
			return errors.New("Unabled to assert type []byte on default value")
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

func valid_name(name string) bool {
	var validName = regexp.MustCompile(`^[a-zA-Z_]+[a-zA-Z0-9_]*$`)
	return validName.MatchString(name)
}

// Deletes constant with name 'name' from the pool.
func (pool *Pool) Delete(name string) error {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if pool.defaults[name] == nil {
		return errors.New("Constant doesn't exists")
	}

	delete(pool.defaults, name)
	return nil
}

// Returns the prefix for the pool.
func (pool *Pool) Prefix() string {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	return pool.prefix
}

// Returns a slice of all constants in the pool.
func (pool *Pool) List() []string {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	consts := make([]string, 0, len(pool.defaults))
	for c := range pool.defaults {
		consts = append(consts, c)
	}
	return consts
}

// Returns a slice of all constants in the pool with the pool's prefix prepended.
func (pool *Pool) Environment() []string {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	consts := pool.List()
	for i := 0; i < len(consts); i++ {
		consts[i] = pool.env_name(consts[i])
	}
	return consts
}
