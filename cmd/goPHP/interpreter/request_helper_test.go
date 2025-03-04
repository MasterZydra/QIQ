package interpreter

import (
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/runtime/values"
	"testing"
)

func TestParseQuery(t *testing.T) {
	runTest := func(t *testing.T, input string, expected *values.Array) {
		actual, err := parseQuery(input, ini.NewDefaultIni())
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
