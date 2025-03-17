package interpreter

import (
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/request"
	"GoPHP/cmd/goPHP/runtime/values"
	"testing"
)

func TestParsePost(t *testing.T) {
	array, err := parsePost(
		`Content-Type: multipart/form-data; boundary=---------------------------20896060251896012921717172737
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name=name1

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name=name\4

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name=name\\5

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name=name\'6

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name=name\"7

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name='name\8'

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name='name\\9'

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name='name\'10'

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name='name\"11'

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name="name\12"

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name="name\\13"

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name="name\'14"

testname
-----------------------------20896060251896012921717172737
Content-Disposition: form-data; name="name\"15"

testname
-----------------------------20896060251896012921717172737--`,
		NewInterpreter(ini.NewDevIni(), request.NewRequest(), TEST_FILE_NAME),
	)
	if err != nil {
		t.Errorf("Parsing post data failed: %s", err)
	}
	expected := `{ArrayValue: 
Key: {StrValue: name1}
Value: {StrValue: testname}
Key: {StrValue: name\4}
Value: {StrValue: testname}
Key: {StrValue: name\5}
Value: {StrValue: testname}
Key: {StrValue: name\'6}
Value: {StrValue: testname}
Key: {StrValue: name\"7}
Value: {StrValue: testname}
Key: {StrValue: name\8}
Value: {StrValue: testname}
Key: {StrValue: name\9}
Value: {StrValue: testname}
Key: {StrValue: name'10}
Value: {StrValue: testname}
Key: {StrValue: name\"11}
Value: {StrValue: testname}
Key: {StrValue: name\12}
Value: {StrValue: testname}
Key: {StrValue: name\13}
Value: {StrValue: testname}
Key: {StrValue: name\'14}
Value: {StrValue: testname}
Key: {StrValue: name"15}
Value: {StrValue: testname}
}
`
	if expected != values.ToString(array) {
		t.Errorf("Parsing post data failed.\nExpected:\n%s\nGot:\n%s", expected, values.ToString(array))
	}
}

func TestParseQuery(t *testing.T) {
	runTest := func(t *testing.T, input string, expected *values.Array) {
		actual, err := parseQuery(input, NewInterpreter(ini.NewDefaultIni(), request.NewRequest(), TEST_FILE_NAME))
		if err != nil {
			t.Errorf("Unexpected error: \"%s\"", err)
			return
		}
		equal, err := compare(actual, "===", expected)
		if err != nil {
			t.Errorf("Unexpected error while comparing: \"%s\"", err)
			return
		}
		if !equal.Value {
			t.Errorf("Wrong result for query \"%s\"", input)
		}
	}

	// Query with only values without keys
	runTest(t,
		"ab+cd+ef",
		values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
			values.NewInt(0): values.NewStr("ab"),
			values.NewInt(1): values.NewStr("cd"),
			values.NewInt(2): values.NewStr("ef"),
		}),
	)

	// Query with only key-value-pairs
	runTest(t,
		"b=Hello+Again+World&c=Hi+Mom",
		values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
			values.NewStr("b"): values.NewStr("Hello Again World"),
			values.NewStr("c"): values.NewStr("Hi Mom"),
		}),
	)

	// Query with special characters
	runTest(t,
		"a+-_!.%22%C2%A7$/()%=a",
		values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
			values.NewStr(`a_-_!_"§$/()%`): values.NewStr("a"),
		}),
	)
	runTest(t,
		"a[+-_!%22%C2%A7$/()%=a",
		values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
			values.NewStr(`a__-_!"§$/()%`): values.NewStr("a"),
		}),
	)
	runTest(t,
		"a]+-_!%22%C2%A7$/()%=a",
		values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
			values.NewStr(`a]_-_!"§$/()%`): values.NewStr("a"),
		}),
	)
	runTest(t,
		"a[+-_!]%22%C2%A7$/()%=a",
		values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
			values.NewStr("a"): values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
				values.NewStr(" -_!"): values.NewStr("a"),
			}),
		}),
	)

	// Simple Query with array
	runTest(t,
		"123[]=SEGV",
		values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
			values.NewInt(123): values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
				values.NewInt(0): values.NewStr("SEGV"),
			}),
		}),
	)
	// Query with too many closing "]"
	runTest(t,
		"123[]]]]]]]]]=SEGV",
		values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
			values.NewInt(123): values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
				values.NewInt(0): values.NewStr("SEGV"),
			}),
		}),
	)
	// Simple Query with array that overwrites old value
	runTest(t,
		"a[]=1&a[0]=5",
		values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
			values.NewStr("a"): values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
				values.NewInt(0): values.NewStr("5"),
			}),
		}),
	)
	// Complex Query with array
	runTest(t,
		"a[][]=1&a[][]=3&b[a][b][c]=1&b[a][b][d]=1",
		values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
			values.NewStr("a"): values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
				values.NewInt(0): values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
					values.NewInt(0): values.NewStr("1"),
				}),
				values.NewInt(1): values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
					values.NewInt(0): values.NewStr("3"),
				}),
			}),
			values.NewStr("b"): values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
				values.NewStr("a"): values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
					values.NewStr("b"): values.NewArrayFromMap(map[values.RuntimeValue]values.RuntimeValue{
						values.NewStr("c"): values.NewStr("1"),
						values.NewStr("d"): values.NewStr("1"),
					}),
				}),
			}),
		}),
	)
}
