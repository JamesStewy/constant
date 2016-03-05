package constant

import (
	"os"
	"strconv"
)

func (pool *Pool) Str(name string) (val string) {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	if pool.defaults[name] == nil {
		return ""
	}

	if env := os.Getenv(pool.env_name(name)); env == "" {
		val = *pool.defaults[name]
	} else {
		val = env
	}

	return
}

func (pool *Pool) Int(name string) (val int, err error) {
	val, err = strconv.Atoi(pool.Str(name))
	return
}

func (pool *Pool) Float(name string, bitSize int) (val float64, err error) {
	val, err = strconv.ParseFloat(pool.Str(name), bitSize)
	return
}

func (pool *Pool) Bool(name string) (val bool, err error) {
	val, err = strconv.ParseBool(pool.Str(name))
	return
}

func (pool *Pool) IsSet(name string) bool {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	return pool.defaults[name] != nil
}

func (pool *Pool) env_name(name string) string {
	return pool.prefix + name
}