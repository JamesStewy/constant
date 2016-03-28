package constant_test

import (
	"github.com/JamesStewy/constant"
	"os"
	"reflect"
	"sort"
	"testing"
)

var pool *constant.Pool
var pool_prefix = "test_"

type constPair struct {
	name  string
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
	{constPair{"string1", "string value"}, "string value", false, "string value", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{"string1", "already exists"}, "string value", true, "string value", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{"string2", ""}, "", false, "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{"byte1", []byte("byte string value")}, "byte string value", false, "byte string value", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{"fmtStringer1", stringer("fmt.Stringer value")}, "fmt.Stringer value", false, "fmt.Stringer value", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{"int1", 12}, "12", false, "12", val_err{12, false}, val_err{12.0, false}, val_err{false, true}, true},
	{constPair{"int2", -112}, "-112", false, "-112", val_err{-112, false}, val_err{-112.0, false}, val_err{false, true}, true},
	{constPair{"float1", 10.043}, "10.043", false, "10.043", val_err{0, true}, val_err{10.043, false}, val_err{false, true}, true},
	{constPair{"float2", -0.02144}, "-0.02144", false, "-0.02144", val_err{0, true}, val_err{-0.02144, false}, val_err{false, true}, true},
	{constPair{"bool1", true}, "true", false, "true", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{"bool2", false}, "false", false, "false", val_err{0, true}, val_err{0.0, true}, val_err{false, false}, true},
	{constPair{"bool3", "1"}, "1", false, "1", val_err{1, false}, val_err{1.0, false}, val_err{true, false}, true},
	{constPair{"bool4", "0"}, "0", false, "0", val_err{0, false}, val_err{0.0, false}, val_err{false, false}, true},
	{constPair{"bool5", "t"}, "t", false, "t", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{"bool6", "f"}, "f", false, "f", val_err{0, true}, val_err{0.0, true}, val_err{false, false}, true},
	{constPair{"bool7", "T"}, "T", false, "T", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{"bool8", "F"}, "F", false, "F", val_err{0, true}, val_err{0.0, true}, val_err{false, false}, true},
	{constPair{"invalid1", invalid("not a valid type")}, "", true, "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{"invalid2", nil}, "", true, "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{"", "empty name"}, "", true, "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{"2name", "name starting with number"}, "", true, "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{"namew1thnum", "name containing number"}, "name containing number", false, "name containing number", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{"invalidChar.1", "invalid char"}, "", true, "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{"invalidChar,2", "invalid char"}, "", true, "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{"invalidChar!&*()3", "invalid char"}, "", true, "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, false},
	{constPair{"_name1", "start with underscore"}, "start with underscore", false, "start with underscore", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{"Name2", "start with capital"}, "start with capital", false, "start with capital", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{"template1", `{{ const "string1" }} - {{ const "int1" }}`}, `{{ const "string1" }} - {{ const "int1" }}`, false, "string value - 12", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{"template2", `{{ isset "doesntexist" }}`}, `{{ isset "doesntexist" }}`, false, "false", val_err{0, true}, val_err{0.0, true}, val_err{false, false}, true},
	{constPair{"template3", `{{ isset "string1" }}`}, `{{ isset "string1" }}`, false, "true", val_err{0, true}, val_err{0.0, true}, val_err{true, false}, true},
	{constPair{"template4", `{{ const "string1"`}, `{{ const "string1"`, false, "", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
	{constPair{"notexist1", `Other stuff: {{ const "doesntexist" }}`}, `Other stuff: {{ const "doesntexist" }}`, false, "Other stuff: ", val_err{0, true}, val_err{0.0, true}, val_err{false, true}, true},
}

func TestMain(m *testing.M) {
	pool = constant.NewPool(pool_prefix)
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	for _, test := range tests {
		err := pool.New(test.pair.name, test.pair.value)
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
	}
}

func TestStr(t *testing.T) {
	for _, test := range tests {
		str := pool.Str(test.pair.name)
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
		val, err := pool.Int(test.pair.name)
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
		val, err := pool.Float(test.pair.name, 64)
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
		val, err := pool.Bool(test.pair.name)
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
		val := pool.IsSet(test.pair.name)
		if val != test.isset {
			t.Error(
				"For", test.pair,
				"expected", test.new_error,
				"got", val,
			)
		}
	}
}

func TestDefault(t *testing.T) {
	for _, test := range tests {
		val := pool.Default(test.pair.name)
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
		if test.new_error == false {
			exp_list = append(exp_list, test.pair.name)
		}
	}

	res_list := pool.List()

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
		if test.new_error == false {
			exp_list = append(exp_list, pool_prefix+test.pair.name)
		}
	}

	res_list := pool.Environment()

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
	for _, test := range tests {
		err := pool.Delete(test.pair.name)
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
	}

	list_len := len(pool.List())
	if list_len != 0 {
		t.Error(
			"expected 0 items left",
			"got", list_len,
		)
	}
}
