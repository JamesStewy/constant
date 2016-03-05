/*
Package constant provides a simple interface for creating an storing constants by a key in an application.
If available constants are read from the environment, to provide dynamic configuration.

Package constant also provides a simple templating system to create constants based of the values of other constants.

Template

Package constant provides a simple template system using text/template (https://golang.org/pkg/text/template/) to allow constants to be expressed a combintation of other constants in the same pool.

The following methods are available to use in a constant.

	{{ const "name" }}
		Returns the value of of another constant in the same pool named 'name'.
		Self reference is not allowed and returns an empty string.

		Example:
			`{{ const "host" }}:{{ const "port" }}`
			If host=`localhost` and port=`3306` then the above template would return `localhost:3306`.

		Caution:
			Although there is a check in place to test if a constant references itself,
			there is no check for cyclic dependancy. If a cyclic dependancy is created
			then the program will enter an infinite loop.

	{{ list }}
		Returns a slice of all constants in the pool except itself.

		Example:
			`{{ range list }}{{ . }}={{ const . }}; {{ end }}`
			If host=`localhost` and port=`3306` then the above template would return `host=localhost; port=3306; `.

	{{ isset "name" }}
		Returns whether or not a constant named 'name' is in the current pool. Same as pool.IsSet(name).

		Example:
			`{{ const "protocol" }}://{{ const "domain" }}{{if isset "port"}}:{{ const "port" }}{{end}}/{{ const "page" }}`

			If protocol=`http`, domain=`localhost` and page=`index.html` then the above template would return `http://localhost/index.html`.

			Or if the same constants are set as well as port=`8080` then the above template would return `http://localhost:8080/index.html`.

Example

In the following example a pool is created to store constants related to MySQL.
HOST, PORT, USER and PASSWORD have standard default values while ADDRESS is set by default to equal `HOST + ":" + PORT`.

Near the end HOST is updated via an envionment variable which, in turn, also updates ADDRESS.

	package main

	import (
		"github.com/JamesStewy/constant"
		"fmt"
		"os"
	)

	var mysql_const *constant.Pool

	func main() {
		// Create new pool to store constants related to MySQL
		mysql_const = constant.NewPool("MYSQL_")

		// Set default values for HOST, PORT, USER and PASSWORD
		mysql_const.New("HOST", "localhost")
		mysql_const.New("PORT", 3306)
		mysql_const.New("USER", "root")
		mysql_const.New("PASSWORD", "")

		// Set ADDRESS to be equal to HOST + ":" + PORT
		mysql_const.New("ADDRESS", `{{ const "HOST" }}:{{ const "PORT" }}`)

		display_pool()

		// Update the MySQL host
		os.Setenv("MYSQL_HOST", "mydomain.com")
		fmt.Println("\nChanged MYSQL_HOST\n")

		display_pool()
	}

	func display_pool() {
		// Loop through each constant in the pool and display its value
		for _, name := range mysql_const.List() {
			// Call mysql_const.Str(name) to retrieve a constants value
			fmt.Printf("%s=%s\n", name, mysql_const.Str(name))
		}
	}

The above example returns the following:

	ADDRESS=localhost:3306
	HOST=localhost
	PORT=3306
	USER=root
	PASSWORD=

	Changed MYSQL_HOST

	USER=root
	PASSWORD=
	ADDRESS=mydomain.com:3306
	HOST=mydomain.com
	PORT=3306
*/
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

// Returns the value of the constant in the pool named 'name' as a float64.
//
// Follows convention of strconv.ParseFloat (https://golang.org/pkg/strconv/#ParseFloat).
func (pool *Pool) Float(name string, bitSize int) (val float64, err error) {
	val, err = strconv.ParseFloat(pool.Str(name), bitSize)
	return
}

// Returns the value of the constant in the pool named 'name' as a boolean.
//
// Follows convention of strconv.ParseBool (https://golang.org/pkg/strconv/#ParseBool).
func (pool *Pool) Bool(name string) (val bool, err error) {
	val, err = strconv.ParseBool(pool.Str(name))
	return
}

// Returns if the constant named 'name' is set in the pool.
func (pool *Pool) IsSet(name string) bool {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	return pool.defaults[name] != nil
}

func (pool *Pool) env_name(name string) string {
	return pool.prefix + name
}
