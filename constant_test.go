package constant_test

import (
	"github.com/JamesStewy/constant"
	"os"
	"reflect"
	"sort"
	"testing"
)

var tree *constant.Node
var tree_prefix = "test"
var tree_delimiter = "_"

type constPair struct {
	path  []string
	value interface{}
}

type val_err struct {
	val interface{}
	err bool
}

type test_type struct {
	pair      constPair
	default_s string
	new_error bool
	rel_name  string
	str       string
	int_res   val_err
	float_res val_err
	bool_res  val_err
	isset     bool
}

type stringer string

func (s stringer) String() string {
	return string(s)
}

type invalid string

var tests = []test_type{
	{constPair{[]string{"string1"}, "string value"}, "string value", false, "string1", "string value", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"string2"}, ""}, "", false, "string2", "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{[]string{"byte1"}, []byte("byte string value")}, "byte string value", false, "byte1", "byte string value", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"fmtStringer1"}, stringer("fmt.Stringer value")}, "fmt.Stringer value", false, "fmtStringer1", "fmt.Stringer value", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"int1"}, 12}, "12", false, "int1", "12", val_err{12, false}, val_err{12.0, false}, val_err{false, true}, true},
	{constPair{[]string{"int2"}, -112}, "-112", false, "int2", "-112", val_err{-112, false}, val_err{-112.0, false}, val_err{false, true}, true},
	{constPair{[]string{"float1"}, 10.043}, "10.043", false, "float1", "10.043", val_err{0, true}, val_err{10.043, false}, val_err{false, true}, true},
	{constPair{[]string{"float2"}, -0.02144}, "-0.02144", false, "float2", "-0.02144", val_err{0, true}, val_err{-0.02144, false}, val_err{false, true}, true},
	{constPair{[]string{"bool1"}, true}, "true", false, "bool1", "true", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{[]string{"bool2"}, false}, "false", false, "bool2", "false", val_err{0, true}, val_err{0.0, true}, val_err{false, false}, true},
	{constPair{[]string{"bool3"}, "1"}, "1", false, "bool3", "1", val_err{1, false}, val_err{1.0, false}, val_err{true, false}, true},
	{constPair{[]string{"bool4"}, "0"}, "0", false, "bool4", "0", val_err{0, false}, val_err{0.0, false}, val_err{false, false}, true},
	{constPair{[]string{"bool5"}, "t"}, "t", false, "bool5", "t", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{[]string{"bool6"}, "f"}, "f", false, "bool6", "f", val_err{0, true}, val_err{0.0, true}, val_err{false, false}, true},
	{constPair{[]string{"bool7"}, "T"}, "T", false, "bool7", "T", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{[]string{"bool8"}, "F"}, "F", false, "bool8", "F", val_err{0, true}, val_err{0.0, true}, val_err{false, false}, true},
	{constPair{[]string{"invalid1"}, invalid("not a valid type")}, "", true, "", "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{[]string{"nil_value"}, nil}, "", false, "nil_value", "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{[]string{"2name"}, "name starting with number"}, "", true, "", "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{[]string{"namew1thnum"}, "name containing number"}, "name containing number", false, "namew1thnum", "name containing number", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"invalidChar.1"}, "invalid char"}, "", true, "", "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{[]string{"invalidChar,2"}, "invalid char"}, "", true, "", "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{[]string{"invalidChar!&*()3"}, "invalid char"}, "", true, "", "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{[]string{"name1"}, "start with underscore"}, "start with underscore", false, "name1", "start with underscore", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"Name2"}, "start with capital"}, "start with capital", false, "Name2", "start with capital", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"template1"}, `{{ const "string1" }} - {{ const "int1" }}`}, `{{ const "string1" }} - {{ const "int1" }}`, false, "template1", "string value - 12", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"template2"}, `{{ isset "doesntexist" }}`}, `{{ isset "doesntexist" }}`, false, "template2", "false", val_err{0, true}, val_err{0.0, true}, val_err{false, false}, true},
	{constPair{[]string{"template3"}, `{{ isset "string1" }}`}, `{{ isset "string1" }}`, false, "template3", "true", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{[]string{"template4"}, `{{ const "string1"`}, `{{ const "string1"`, false, "template4", "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"notexist1"}, `Other stuff: {{ const "doesntexist" }}`}, `Other stuff: {{ const "doesntexist" }}`, false, "notexist1", "Other stuff: ", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1"}, "node1"}, "node1", false, "node1", "node1", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1", "str1"}, "node1 string"}, "node1 string", false, "node1" + tree_delimiter + "str1", "node1 string", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1", "int1"}, 10}, "10", false, "node1" + tree_delimiter + "int1", "10", val_err{10, false}, val_err{10.0, false}, val_err{false, true}, true},
	{constPair{[]string{"node1", "int2"}, `{{ const "int1" }}`}, `{{ const "int1" }}`, false, "node1" + tree_delimiter + "int2", "10", val_err{10, false}, val_err{10.0, false}, val_err{false, true}, true},
	{constPair{[]string{"node1", "bool1"}, true}, "true", false, "node1" + tree_delimiter + "bool1", "true", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{[]string{"node1", "subnode1", "noexist"}, "subnode does not exist"}, "", true, "", "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{[]string{"node1", "subnode1"}, nil}, "", false, "node1" + tree_delimiter + "subnode1", "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{[]string{"node1", "subnode1", "val1"}, "dhfsjl"}, "dhfsjl", false, "node1" + tree_delimiter + "subnode1" + tree_delimiter + "val1", "dhfsjl", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1", "subnode1", "val2"}, "afaFASDf"}, "afaFASDf", false, "node1" + tree_delimiter + "subnode1" + tree_delimiter + "val2", "afaFASDf", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1", "subnode1", "root"}, `{{ const "" }}`}, `{{ const "" }}`, false, "node1" + tree_delimiter + "subnode1" + tree_delimiter + "root", "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1", "subnode1", "list"}, `{{ list }}`}, `{{ list }}`, false, "node1" + tree_delimiter + "subnode1" + tree_delimiter + "list", "[root val1 val2]", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1", "subnode2"}, true}, "true", false, "node1" + tree_delimiter + "subnode2", "true", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{[]string{"node1", "subnode2", "val1"}, "aighwksf"}, "aighwksf", false, "node1" + tree_delimiter + "subnode2" + tree_delimiter + "val1", "aighwksf", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1", "subnode2", "val2"}, "uiuhfbaf"}, "uiuhfbaf", false, "node1" + tree_delimiter + "subnode2" + tree_delimiter + "val2", "uiuhfbaf", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1", "subnode2", "root"}, `{{ const "" }}`}, `{{ const "" }}`, false, "node1" + tree_delimiter + "subnode2" + tree_delimiter + "root", "true", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{[]string{"node1", "subnode2", "list"}, `{{ list }}`}, `{{ list }}`, false, "node1" + tree_delimiter + "subnode2" + tree_delimiter + "list", "[ root val1 val2]", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1", "root"}, `{{ const "" }}`}, `{{ const "" }}`, false, "node1" + tree_delimiter + "root", "node1", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1", "subnode2val1"}, `{{ const "subnode2" "val1" }}`}, `{{ const "subnode2" "val1" }}`, false, "node1" + tree_delimiter + "subnode2val1", "aighwksf", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{[]string{"node1", "isset1"}, `{{ isset "subnode1" "val1" }}`}, `{{ isset "subnode1" "val1" }}`, false, "node1" + tree_delimiter + "isset1", "true", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{[]string{"node1", "isset2"}, `{{ isset "subnode1" "noexist" }}`}, `{{ isset "subnode1" "noexist" }}`, false, "node1" + tree_delimiter + "isset2", "false", val_err{0, true}, val_err{0.0, true}, val_err{false, false}, true},
	{constPair{[]string{"node1", "list"}, `{{ list }}`}, `{{ list }}`, false, "node1" + tree_delimiter + "list", "[ bool1 int1 int2 isset1 isset2 root str1 subnode1_list subnode1_root subnode1_val1 subnode1_val2 subnode2 subnode2_list subnode2_root subnode2_val1 subnode2_val2 subnode2val1]", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
}

func TestMain(m *testing.M) {
	tree = constant.NewTree(tree_prefix, tree_delimiter)
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	for _, test := range tests {
		node := tree.Node(test.pair.path[:len(test.pair.path)-1]...)
		if node != nil {
			_, err := node.New(test.pair.path[len(test.pair.path)-1], test.pair.value)
			if (err != nil) != test.new_error {
				no_str := ""
				if !test.new_error {
					no_str = "no "
				}
				t.Error(
					"For", test.pair,
					"expected "+no_str+"error",
					"got", err,
				)
			}
		} else {
			if !test.new_error {
				t.Error(
					"For", test.pair,
					"expected error",
					"got no error",
				)
			}
		}
	}
}

func TestNode(t *testing.T) {
	for _, test := range tests {
		node := tree.Node(test.pair.path...)
		if (node == nil) != test.new_error {
			t.Error(
				"For", test.pair,
				"Expected", test.new_error,
				"got", node == nil,
			)
		}
	}
}

func TestDelimiter(t *testing.T) {
	res := tree.Delimiter()
	if res != tree_delimiter {
		t.Error(
			"For root node",
			"Expected", tree_delimiter,
			"got", res,
		)
	}

	for _, test := range tests {
		if node := tree.Node(test.pair.path...); node != nil {
			res = node.Delimiter()
			if res != tree_delimiter {
				t.Error(
					"For", test.pair,
					"expected", tree_delimiter,
					"got", res,
				)
			}
		}
	}
}

func TestName(t *testing.T) {
	name := tree.Name()
	if name != tree_prefix {
		t.Error(
			"For root node",
			"expected", tree_prefix,
			"got", name,
		)
	}

	for _, test := range tests {
		if node := tree.Node(test.pair.path...); node != nil {
			name = node.Name()
			if name != test.pair.path[len(test.pair.path)-1] {
				t.Error(
					"For", test.pair,
					"expected", test.pair.path[len(test.pair.path)-1],
					"got", name,
				)
			}
		}
	}
}

func TestFullName(t *testing.T) {
	full_name := tree.FullName()
	if full_name != tree_prefix {
		t.Error(
			"For root node",
			"expected", tree_prefix,
			"got", full_name,
		)
	}

	for _, test := range tests {
		if node := tree.Node(test.pair.path...); node != nil {
			full_name = node.FullName()
			if full_name != tree_prefix+tree_delimiter+test.rel_name {
				t.Error(
					"For", test.pair,
					"expected", tree_prefix+tree_delimiter+test.rel_name,
					"got", full_name,
				)
			}
		}
	}
}

func TestStr(t *testing.T) {
	for _, test := range tests {
		str := tree.Str(test.pair.path...)
		if str != test.str {
			t.Error(
				"For", test.pair,
				"expected", test.str,
				"got", str,
			)
		}
	}
}

func TestInt(t *testing.T) {
	for _, test := range tests {
		val, err := tree.Int(test.pair.path...)
		if val != test.int_res.val || (err != nil) != test.int_res.err {
			t.Error(
				"For", test.pair,
				"expected", test.int_res,
				"got (", val, err, ")",
			)
		}
	}
}

func TestFloat(t *testing.T) {
	for _, test := range tests {
		val, err := tree.Float(64, test.pair.path...)
		if val != test.float_res.val || (err != nil) != test.float_res.err {
			t.Error(
				"For", test.pair,
				"expected", test.float_res,
				"got (", val, err, ")",
			)
		}
	}
}

func TestBool(t *testing.T) {
	for _, test := range tests {
		val, err := tree.Bool(test.pair.path...)
		if val != test.bool_res.val || (err != nil) != test.bool_res.err {
			t.Error(
				"For", test.pair,
				"expected", test.bool_res,
				"got (", val, err, ")",
			)
		}
	}
}

func TestIsSet(t *testing.T) {
	for _, test := range tests {
		val := tree.IsSet(test.pair.path...)
		if val != test.isset {
			t.Error(
				"For", test.pair,
				"expected", test.isset,
				"got", val,
			)
		}
	}
}

func TestDefault(t *testing.T) {
	for _, test := range tests {
		val := tree.Default(test.pair.path...)
		if val != test.default_s {
			t.Error(
				"For", test.pair,
				"expected", test.default_s,
				"got", val,
			)
		}
	}
}

func TestList(t *testing.T) {
	exp_list := make([]string, 0)
	for _, test := range tests {
		if test.new_error == false && test.pair.value != nil {
			exp_list = append(exp_list, test.rel_name)
		}
	}

	res_list := tree.List()

	sort.StringSlice(exp_list).Sort()
	sort.StringSlice(res_list).Sort()

	if !reflect.DeepEqual(res_list, exp_list) {
		t.Error(
			"expected", exp_list,
			"got", res_list,
		)
	}
}

func TestEnvironment(t *testing.T) {
	exp_list := make([]string, 0)
	for _, test := range tests {
		if test.new_error == false && test.pair.value != nil {
			exp_list = append(exp_list, tree_prefix+tree_delimiter+test.rel_name)
		}
	}

	res_list := tree.Environment()

	sort.StringSlice(exp_list).Sort()
	sort.StringSlice(res_list).Sort()

	if !reflect.DeepEqual(res_list, exp_list) {
		t.Error(
			"expected", exp_list,
			"got", res_list,
		)
	}
}

func TestDelete(t *testing.T) {
	for i := len(tests) - 1; i >= 0; i-- {
		err := tree.Delete(tests[i].pair.path...)
		if (err != nil) != tests[i].new_error {
			no_str := ""
			if !tests[i].new_error {
				no_str = "no "
			}
			t.Error(
				"For", tests[i].pair,
				"expected "+no_str+"error",
				"got", err,
			)
		}
	}

	list_len := len(tree.List())
	if list_len != 0 {
		t.Error(
			"expected 0 items left",
			"got", list_len,
		)
	}
}
