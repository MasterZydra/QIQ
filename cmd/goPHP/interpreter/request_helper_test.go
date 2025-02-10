package interpreter

import (
	"testing"
)

func TestParseQuery(t *testing.T) {
	runTest := func(t *testing.T, input string, expected *ArrayRuntimeValue) {
		actual, err := parseQuery(input)
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
		NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
			NewIntegerRuntimeValue(0): NewStringRuntimeValue("ab"),
			NewIntegerRuntimeValue(1): NewStringRuntimeValue("cd"),
			NewIntegerRuntimeValue(2): NewStringRuntimeValue("ef"),
		}),
	)

	// Query with only key-value-pairs
	runTest(t,
		"b=Hello+Again+World&c=Hi+Mom",
		NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
			NewStringRuntimeValue("b"): NewStringRuntimeValue("Hello Again World"),
			NewStringRuntimeValue("c"): NewStringRuntimeValue("Hi Mom"),
		}),
	)

	// Query with special characters
	runTest(t,
		"a+-_!.%22%C2%A7$/()%=a",
		NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
			NewStringRuntimeValue(`a_-_!_"§$/()%`): NewStringRuntimeValue("a"),
		}),
	)
	runTest(t,
		"a[+-_!%22%C2%A7$/()%=a",
		NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
			NewStringRuntimeValue(`a__-_!"§$/()%`): NewStringRuntimeValue("a"),
		}),
	)
	runTest(t,
		"a]+-_!%22%C2%A7$/()%=a",
		NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
			NewStringRuntimeValue(`a]_-_!"§$/()%`): NewStringRuntimeValue("a"),
		}),
	)
	runTest(t,
		"a[+-_!]%22%C2%A7$/()%=a",
		NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
			NewStringRuntimeValue("a"): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
				NewStringRuntimeValue(" -_!"): NewStringRuntimeValue("a"),
			}),
		}),
	)

	// Simple Query with array
	runTest(t,
		"123[]=SEGV",
		NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
			NewIntegerRuntimeValue(123): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
				NewIntegerRuntimeValue(0): NewStringRuntimeValue("SEGV"),
			}),
		}),
	)
	// Query with too many closing "]"
	runTest(t,
		"123[]]]]]]]]]=SEGV",
		NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
			NewIntegerRuntimeValue(123): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
				NewIntegerRuntimeValue(0): NewStringRuntimeValue("SEGV"),
			}),
		}),
	)
	// Simple Query with array that overwrites old value
	runTest(t,
		"a[]=1&a[0]=5",
		NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
			NewStringRuntimeValue("a"): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
				NewIntegerRuntimeValue(0): NewStringRuntimeValue("5"),
			}),
		}),
	)
	// Complex Query with array
	runTest(t,
		"a[][]=1&a[][]=3&b[a][b][c]=1&b[a][b][d]=1",
		NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
			NewStringRuntimeValue("a"): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
				NewIntegerRuntimeValue(0): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
					NewIntegerRuntimeValue(0): NewStringRuntimeValue("1"),
				}),
				NewIntegerRuntimeValue(1): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
					NewIntegerRuntimeValue(0): NewStringRuntimeValue("3"),
				}),
			}),
			NewStringRuntimeValue("b"): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
				NewStringRuntimeValue("a"): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
					NewStringRuntimeValue("b"): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
						NewStringRuntimeValue("c"): NewStringRuntimeValue("1"),
						NewStringRuntimeValue("d"): NewStringRuntimeValue("1"),
					}),
				}),
			}),
		}),
	)
}
