package constant

import (
	"bytes"
	"os"
	"strconv"
	"text/template"
)

// Returns the value of the constant in the pool named 'name' as a string.
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

// Returns the value of the constant in the pool named 'name' as an integer.
//
// Follows convention of strconv.Atoi (https://golang.org/pkg/strconv/#Atoi).
func (pool *Pool) Int(name string) (val int, err error) {
	val, err = strconv.Atoi(pool.Str(name))
	return
}

// Run Int but ignore errors
func (pool *Pool) IntI(name string) (val int) {
	val, _ = pool.Int(name)
	return
}

// Returns the value of the constant in the pool named 'name' as a float64.
//
// Follows convention of strconv.ParseFloat (https://golang.org/pkg/strconv/#ParseFloat).
func (pool *Pool) Float(name string, bitSize int) (val float64, err error) {
	val, err = strconv.ParseFloat(pool.Str(name), bitSize)
	return
}

// Run Float but ignore errors
func (pool *Pool) FloatI(name string, bitSize int) (val float64) {
	val, _ = pool.Float(name, bitSize)
	return
}

// Returns the value of the constant in the pool named 'name' as a boolean.
//
// Follows convention of strconv.ParseBool (https://golang.org/pkg/strconv/#ParseBool).
func (pool *Pool) Bool(name string) (val bool, err error) {
	val, err = strconv.ParseBool(pool.Str(name))
	return
}

// Run Bool but ignore errors
func (pool *Pool) BoolI(name string) (val bool) {
	val, _ = pool.Bool(name)
	return
}

// Returns if the constant named 'name' is set in the pool.
func (pool *Pool) IsSet(name string) bool {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	return pool.defaults[name] != nil
}

// Returns the default value of the constant in the pool named 'name'.
// If the default value contains templates the templates will not be parsed.
// If the default value is not a string it will be converted to a string as per the strconv package (https://golang.org/pkg/strconv/).
func (pool *Pool) Default(name string) string {
	if pool.defaults[name] == nil {
		return ""
	}

	return *pool.defaults[name]
}

func (pool *Pool) env_name(name string) string {
	return pool.prefix + name
}
