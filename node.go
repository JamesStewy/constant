package constant

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"sync"
)

// A Node represents one node in a tree of constants.
// A Node can have a value and/or child nodes associated with it.
type Node struct {
	mutex     sync.RWMutex
	name      string
	delimiter string
	def_val   *string
	parent    *Node
	nodes     map[string]*Node
}

// Creates the root node for a new tree.
//
// Prefix sets the environment variable prefix which is prepended to node names when searching the runtime environment.
// For example if a tree has a prefix 'MYSQL', a delimiter of '_' and a child node named 'HOST' then constant 'HOST' would be set to the value of the environment variable 'MYSQL_HOST'.
func NewTree(prefix, delimiter string) *Node {
	return &Node{
		name:      prefix,
		delimiter: delimiter,
		nodes:     make(map[string]*Node),
	}
}

/*
Adds a new child node to the node 'n'.
Returns the newly created child node if successful.

name: Name of the constant.
Must follow variable naming convention.
Lower case letters, uppercase letters, numbers and underscores.
Can't start with a number.

def_val: The default value for the constant if no environment variable is available.

def_val must be one of the following types:
	string
	[]byte
	fmt.Stringer (https://golang.org/pkg/fmt/#Stringer)
	int
	float64
	bool
	nil (no default value: the new child node will act purely as a node)
*/
func (n *Node) New(name string, def_val interface{}) (*Node, error) {
	if !valid_name(name) {
		return nil, errors.New("Invalid Name")
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.nodes[name] != nil {
		return nil, errors.New("Already exists")
	}

	new_node := &Node{
		name:      name,
		delimiter: n.delimiter,
		parent:    n,
		nodes:     make(map[string]*Node),
	}

	if def_val != nil {
		new_node.def_val = new(string)

		var str_val string
		switch t := def_val.(type) {
		case string:
			if val, ok := def_val.(string); ok {
				str_val = val
			} else {
				return nil, errors.New("Unabled to assert type string on default value")
			}
		case []byte:
			if val, ok := def_val.([]byte); ok {
				str_val = string(val)
			} else {
				return nil, errors.New("Unabled to assert type []byte on default value")
			}
		case fmt.Stringer:
			if val, ok := def_val.(fmt.Stringer); ok {
				str_val = val.String()
			} else {
				return nil, errors.New("Unabled to assert type fmt.Stringer on default value")
			}
		case int:
			if val, ok := def_val.(int); ok {
				str_val = strconv.Itoa(val)
			} else {
				return nil, errors.New("Unabled to assert type int on default value")
			}
		case float64:
			if val, ok := def_val.(float64); ok {
				str_val = strconv.FormatFloat(val, 'f', -1, 64)
			} else {
				return nil, errors.New("Unabled to assert type float64 on default value")
			}
		case bool:
			if val, ok := def_val.(bool); ok {
				str_val = strconv.FormatBool(val)
			} else {
				return nil, errors.New("Unabled to assert type bool on default value")
			}
		default:
			return nil, errors.New(fmt.Sprintf("Unexpected type %T", t))
		}

		*new_node.def_val = str_val
	}

	n.nodes[name] = new_node
	return n.nodes[name], nil
}

func valid_name(name string) bool {
	var validName = regexp.MustCompile(`^[a-zA-Z_]+[a-zA-Z0-9_]*$`)
	return validName.MatchString(name)
}

// Starting at n, Node traverses down the tree of constant as defined by path.
// If no node is found at path then nil is returned.
//
// If no path is specified Node returns itself (n).
// If an element of path is an empty string it is ignored.
// For example n.Node("LOGGING", "LEVEL") is the same as n.Node("", "LOGGING", "LEVEL") or n.Node("LOGGING", "", "LEVEL").
// As a result if n.Node("") is called, Node returns itself (n).
func (n *Node) Node(path ...string) *Node {
	if len(path) == 0 {
		return n
	}

	if path[0] == "" {
		return n.Node(path[1:]...)
	}

	n.mutex.RLock()
	defer n.mutex.RUnlock()

	if n.nodes[path[0]] == nil {
		return nil
	}
	return n.nodes[path[0]].Node(path[1:]...)
}

// Returns the name for the node.
func (n *Node) Name() string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	return n.name
}

// Returns the full name for the node.
func (n *Node) FullName() string {
	return n.nameOffset(0)
}

// Returns the delimiter for the node.
func (n *Node) Delimiter() string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	return n.delimiter
}

// Returns a slice of itself and all child nodes in the node that have a non nil default value.
func (n *Node) Nodes() []*Node {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	var nodes []*Node
	if n.def_val == nil {
		nodes = make([]*Node, 0)
	} else {
		nodes = []*Node{n}
	}

	for _, node := range n.nodes {
		nodes = append(nodes, node.Nodes()...)
	}

	return nodes
}

// Orphans the node as defined by path.
// Sets the default value for the node to nil.
func (n *Node) Delete(path ...string) error {
	node := n.Node(path...)
	if node == nil {
		return errors.New("Does not exist")
	}

	node_full_name := node.FullName()

	node.mutex.Lock()
	defer node.mutex.Unlock()

	parent := node.parent
	if parent == nil {
		return errors.New("Can't delete root node")
	}

	parent.mutex.Lock()
	defer parent.mutex.Unlock()

	delete(parent.nodes, node.name)
	node.name = node_full_name
	node.parent = nil
	node.def_val = nil

	return nil
}

// Returns a sorted slice of names relative to the node 'n' for itself and all child nodes in the node that have non nil default values.
func (n *Node) List() []string {
	return n.environmentOffset(len(n.path()))
}

// Returns a sorted slice of full names for itself and all child nodes in the node that have non nil default values.
func (n *Node) Environment() []string {
	return n.environmentOffset(0)
}

func (n *Node) path() []string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	if n.parent == nil {
		return []string{n.name}
	}

	return append(n.parent.path(), n.name)
}

func (n *Node) pathJoin(path ...string) string {
	if len(path) == 0 {
		return ""
	}
	if len(path) == 1 {
		return path[0]
	}

	if path[1] == "" {
		path[1] = path[0]
	} else if path[0] != "" {
		path[1] = path[0] + n.Delimiter() + path[1]
	}

	return n.pathJoin(path[1:]...)
}

func (n *Node) nameOffset(offset int) string {
	return n.pathJoin(n.path()[offset:]...)
}

func (n *Node) environmentOffset(offset int) []string {
	nodes := n.Nodes()

	env := make([]string, len(nodes))
	for i, node := range nodes {
		env[i] = n.pathJoin(node.path()[offset:]...)
	}

	sort.StringSlice(env).Sort()
	return env
}
