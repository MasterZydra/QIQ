package interpreter

import "testing"

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
	// Simple Query with array
	// TODO fix parseQuery for "123[]=SEGV"
	// runTest(t,
	// 	"123[]=SEGV",
	// 	NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
	// 		NewStringRuntimeValue("123"): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
	// 			NewIntegerRuntimeValue(0): NewStringRuntimeValue("SEGV"),
	// 		}),
	// 	}),
	// )
	// Complex Query with array
	// runTest(t,
	// 	"a[][]=1&a[][]=3&b[a][b][c]=1&b[a][b][d]=1",
	// 	NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
	// 		NewStringRuntimeValue("a"): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
	// 			NewIntegerRuntimeValue(0): NewStringRuntimeValue("1"),
	// 			NewIntegerRuntimeValue(1): NewStringRuntimeValue("3"),
	// 		}),
	// 		NewStringRuntimeValue("b"): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
	// 			NewStringRuntimeValue("a"): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
	// 				NewStringRuntimeValue("b"): NewArrayRuntimeValueFromMap(map[IRuntimeValue]IRuntimeValue{
	// 					NewIntegerRuntimeValue("c"): NewStringRuntimeValue("1"),
	// 					NewIntegerRuntimeValue("d"): NewStringRuntimeValue("1"),
	// 				}),
	// 			}),
	// 		}),
	// 	}),
	// )
}
