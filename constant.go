package constant

import (
	"bytes"
	"os"
	"strconv"
	"text/template"
)

// Returns the value of the node as defined by path.
// If the node's default value is nil an empty string is returned.
// If the envionment variable associated with the node is not equal to an empty string that value is used instead of the node's default value.
// Templates in the node's value are parsed (see sections Template, Template Context and Example for details).
func (n *Node) Str(path ...string) string {
	node := n.Node(path...)
	if node == nil {
		return ""
	}

	node_fullname := node.FullName()
	tmpl := node.Default()

	node.mutex.RLock()
	defer node.mutex.RUnlock()

	parent := node.parent

	if env := os.Getenv(node_fullname); env != "" {
		tmpl = env
	}

	t, err := template.New("constant").Funcs(template.FuncMap{
		"const": func(path ...string) string {
			if len(path) == 1 && path[0] == node.name {
				return ""
			}
			return parent.Str(path...)
		},
		"list": func() []string {
			consts := parent.List()
			for i, cnst := range consts {
				if cnst == node.name {
					consts = append(consts[:i], consts[i+1:]...)
				}
			}
			return consts
		},
		"isset": func(path ...string) bool {
			return parent.IsSet(path...)
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

// Alias of n.Str()
func (n *Node) String() string {
	return n.Str()
}

// Returns the value of n.Str(path...) as an integer.
//
// Follows convention of strconv.Atoi (https://golang.org/pkg/strconv/#Atoi).
func (n *Node) Int(path ...string) (val int, err error) {
	val, err = strconv.Atoi(n.Str(path...))
	return
}

// Run n.Int(path...) but ignore errors
func (n *Node) IntI(path ...string) (val int) {
	val, _ = n.Int(path...)
	return
}

// Returns the value of n.Str(path...) as a float64.
//
// Follows convention of strconv.ParseFloat (https://golang.org/pkg/strconv/#ParseFloat).
func (n *Node) Float(bitSize int, path ...string) (val float64, err error) {
	val, err = strconv.ParseFloat(n.Str(path...), bitSize)
	return
}

// Run n.Float(bitSize, path...) but ignore errors
func (n *Node) FloatI(bitSize int, path ...string) (val float64) {
	val, _ = n.Float(bitSize, path...)
	return
}

// Returns the value of n.Str(path...) as a boolean.
//
// Follows convention of strconv.ParseBool (https://golang.org/pkg/strconv/#ParseBool).
func (n *Node) Bool(path ...string) (val bool, err error) {
	val, err = strconv.ParseBool(n.Str(path...))
	return
}

// Run n.Bool(path...) but ignore errors
func (n *Node) BoolI(path ...string) (val bool) {
	val, _ = n.Bool(path...)
	return
}

// Returns false if: the node as defined by path doesn't exist; the node's default value is nil; the node's default value it an empty string.
// Otherwise returns true.
func (n *Node) IsSet(path ...string) bool {
	node := n.Node(path...)
	if node == nil {
		return false
	}

	node.mutex.RLock()
	defer node.mutex.RUnlock()

	return node.def_val != nil && *node.def_val != ""
}

// Returns the default value of node as defined by path.
// If the default value contains templates the templates will not be parsed.
// If the default value is not a string it will be converted to a string as per the strconv package (https://golang.org/pkg/strconv/).
// If the default value is nil an empty string is returned.
func (n *Node) Default(path ...string) string {
	node := n.Node(path...)
	if node == nil {
		return ""
	}

	if !node.IsSet() {
		return ""
	}

	node.mutex.RLock()
	defer node.mutex.RUnlock()

	return *node.def_val
}
