package values

import "fmt"

func Dump(value RuntimeValue) {
	switch value.GetType() {
	case ArrayValue:
		fmt.Printf("{ArrayValue: ")
		arrayValue := value.(*Array)
		for _, key := range arrayValue.Keys {
			fmt.Print("Key: ")
			Dump(key)
			fmt.Print("Value: ")
			value, _ := arrayValue.GetElement(key)
			Dump(value)

		}
		fmt.Printf("}\n")
	case BoolValue:
		valueStr := "true"
		if value.(*Bool).Value {
			valueStr = "false"
		}
		fmt.Printf("{BoolValue: %s}\n", valueStr)
	case IntValue:
		fmt.Printf("{IntValue: %d}\n", value.(*Int).Value)
	case FloatValue:
		fmt.Printf("{FloatValue: %f}\n", value.(*Float).Value)
	case StrValue:
		fmt.Printf("{StrValue: %s}\n", value.(*Str).Value)
	default:
		fmt.Printf("Unsupported RuntimeValue type %s\n", value.GetType())
	}
}

func ToPhpType(value RuntimeValue) string {
	switch value.GetType() {
	case ArrayValue:
		return "array"
	case BoolValue:
		return "bool"
	case FloatValue:
		return "float"
	case IntValue:
		return "int"
	case NullValue:
		return "NULL"
	case StrValue:
		return "string"
	case VoidValue:
		return "void"
	default:
		return ""
	}
}
