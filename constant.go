package constant

import (
	"bytes"
	"os"
	"strconv"
	"text/template"
)

func (pool *Pool) Str(name string) string {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	if pool.defaults[name] == nil {
		return ""
	}

	var tmpl string
	if env := os.Getenv(pool.env_name(name)); env == "" {
		tmpl = *pool.defaults[name]
	} else {
		tmpl = env
	}

	t, err := template.New("constant").Funcs(template.FuncMap{
		"const": func(in_name string) string {
			if in_name == name {
				return ""
			}
			return pool.Str(in_name)
		},
		"list": func() []string {
			consts := pool.List()
			for i, cnst := range consts {
				if cnst == name {
					consts = append(consts[:i], consts[i+1:]...)
				}
			}
			return consts
		},
		"isset": func(in_name string) bool {
			return pool.IsSet(in_name)
		},
	}).Parse(tmpl)

	if err != nil {
		return ""
	}

	var byte_string bytes.Buffer
	if err = t.Execute(&byte_string, nil); err != nil {
		return ""
	}

	return byte_string.String()
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
