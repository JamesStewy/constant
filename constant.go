package constant

import (
	"os"
	"strconv"
)

type cnst struct {
	val string
	set bool
}

func (pool *Pool) Str(name string) (val string) {
	pool.RLock()
	defer pool.RUnlock()

	if !pool.defaults[name].set {
		return ""
	}

	if env := os.Getenv(pool.env_name(name)); env == "" {
		val = pool.defaults[name].val
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
	pool.RLock()
	defer pool.RUnlock()

	return pool.defaults[name].set
}

func (pool *Pool) env_name(name string) string {
	return pool.prefix + name
}
