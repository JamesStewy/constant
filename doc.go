/*
Package constant provides an interface for creating and storing constants in a key based tree structure.
If available, constants are read from the environment, to provide dynamic configuration.

Package constant also provides a templating system to create constants based of the values of other constants.

Template

Package constant provides a template system using text/template (https://golang.org/pkg/text/template/) to allow nodes to be expressed as a combintation of other nodes in the same context.

The following methods are available to use in a node.

	{{ const "path1" ["path2" ...] }}
		Returns the value of of another node in the same context as defined by
		path1[, path2 ...]. Self reference is not allowed and returns an empty string.
		Nodes that don't exist return an empty string.

		Example:
			`{{ const "host" }}:{{ const "port" }}`
			If host=`localhost` and port=`3306` then the above template would return
			`localhost:3306`.

		Caution:
			Although there is a check in place to test if a node references	itself,
			there is no check for cyclic dependancy. If a cyclic dependancy is
			created then the program will enter an infinite loop.

	{{ list }}
		Returns a slice of all nodes in the context except itself.

		Example:
			`{{ range list }}{{ . }}={{ const . }}; {{ end }}`
			If host=`localhost` and port=`3306` then the above template would return
			`host=localhost; port=3306; `.

	{{ isset "path1" ["path2" ...] }}
		Returns whether or not a node as defined by path1[, path2 ...] is in the
		current context.

		Example:
			`{{ const "protocol" }}://{{ const "domain" }}{{if isset "port"}}:{{ const "port" }}{{end}}/{{ const "page" }}`

			If protocol=`http`, domain=`localhost` and page=`index.html` then the
			above template would return `http://localhost/index.html`.

			Or if the same constants are set as well as port=`8080` then the above
			template would return `http://localhost:8080/index.html`.

Template Context

The context for a node includes the context's root node and all of its children recursively.
The root node for a context is the parent of the node.
Therefore the context for a node includes itself, its parent node, its parent's child nodes (siblings of the original node) and all the recursive children.
A node is excluded from the context if its default value is nil.

Other nodes in the same context are referenced by their path relative to the root of the context.
This starts with the root of the context which is referenced as an empty string.
Siblings are referenced by their name ("sibling name").
Recursive children are referenced by the names of their parents followed by their name ("sibling name", "recursive child").

For example if the folling tree structure is created

	                                  ___   name: MYAPP  ___
	                                /       value: nil       \
	                               /            |             \
	                    name: LOG         name: RUNTIME        name: DATABASE
	                    value: nil        value: `dev`         value: `true`
	               /        |                               /        |        \
	 name: LEVEL       name: FILE            name: HOST         name: PORT       name: ADDRESS
	  value: 5       value: `stdout`     value: `localhost`     value: 3306        value: ?
	                                             |
	                                      name: PROVIDER
	                                     value: `internal`

then 'ADDRESS' could be set to and would return

	`{{ const "" }}`                 ->  `true`                        (value of parent)
	`{{ const "HOST" }}`             ->  `localhost`
	`{{ const "HOST" "PROVIDER" }}`  ->  `internal`
	`{{ const "PORT" }}`             ->  `true`
	`{{ const "ADDRESS" }}`          ->  ``                            (self reference not allowed)

	`{{ list }}`                     ->  `[ HOST HOST_PROVIDER PORT]`  (includes an empty string at the start)
	`{{ isset "HOST" }}`             ->  `true`
	`{{ isset "SOMETHING" }}`        ->  `false`


Example

In the following example the tree from the above section (Template Context) is created.
In this example 'ADDRESS' is set by default to equal `HOST + ":" + PORT`.

After creation, HOST is updated via an envionment variable which, in turn, also updates the value of ADDRESS.

	package main

	import (
		"fmt"
		"github.com/JamesStewy/constant"
		"os"
	)

	var tree *constant.Node

	func main() {
		// Create new tree for my app
		tree = constant.NewTree("MYAPP", "_")

		// Create LOG node
		tree.New("LOG", nil)
		tree.Node("LOG").New("LEVEL", 5)
		tree.Node("LOG").New("FILE", "stdout")

		// Create RUNTIME node
		tree.New("RUNTIME", "dev")

		// Create DATABASE node (using alternate method to LOG node)
		mysql_tree, _ := tree.New("DATABASE", true)
		mysql_tree.New("HOST", "localhost")
		mysql_tree.Node("HOST").New("PROVIDER", "internal")
		mysql_tree.New("PORT", 3306)

		// Set ADDRESS to be equal to HOST + ":" + PORT
		mysql_tree.New("ADDRESS", `{{ const "HOST" }}:{{ const "PORT" }}`)

		display_pool()

		// Update the MySQL host
		os.Setenv("MYAPP_DATABASE_HOST", "mydomain.com")
		fmt.Println("\nChanged MYAPP_DATABASE_HOST to mydomain.com\n")

		display_pool()
	}

	func display_pool() {
		// Loop through each constant in the pool and display its value
		for _, node := range tree.Nodes() {
			// Call node.Str(name) to retrieve the node's value
			fmt.Printf("%s=%s\n", node.FullName(), node.Str())
		}
	}

The above example returns:

	MYAPP_LOG_LEVEL=5
	MYAPP_LOG_FILE=stdout
	MYAPP_RUNTIME=dev
	MYAPP_DATABASE=true
	MYAPP_DATABASE_HOST=localhost
	MYAPP_DATABASE_HOST_PROVIDER=internal
	MYAPP_DATABASE_PORT=3306
	MYAPP_DATABASE_ADDRESS=localhost:3306

	Changed MYAPP_DATABASE_HOST to mydomain.com

	MYAPP_LOG_LEVEL=5
	MYAPP_LOG_FILE=stdout
	MYAPP_RUNTIME=dev
	MYAPP_DATABASE=true
	MYAPP_DATABASE_HOST=mydomain.com
	MYAPP_DATABASE_HOST_PROVIDER=internal
	MYAPP_DATABASE_PORT=3306
	MYAPP_DATABASE_ADDRESS=mydomain.com:3306
*/
package constant
