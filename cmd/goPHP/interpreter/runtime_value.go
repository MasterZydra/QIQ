package interpreter

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/phpError"
	"fmt"
	"strconv"
)

type ValueType string

const (
	VoidValue     ValueType = "Void"
	NullValue     ValueType = "Null"
	ArrayValue    ValueType = "Array"
	BooleanValue  ValueType = "Boolean"
	IntegerValue  ValueType = "Integer"
	FloatingValue ValueType = "Floating"
	StringValue   ValueType = "String"
)

// MARK: RuntimeValue

type IRuntimeValue interface {
	GetType() ValueType
}

type RuntimeValue struct {
	valueType ValueType
}

func NewRuntimeValue(valueType ValueType) *RuntimeValue {
	return &RuntimeValue{valueType: valueType}
}

func (runtimeValue *RuntimeValue) GetType() ValueType {
	return runtimeValue.valueType
}

// MARK: VoidValue

var voidRuntimeValue = &RuntimeValue{valueType: VoidValue}

func NewVoidRuntimeValue() *RuntimeValue {
	return voidRuntimeValue
}

// MARK: NullValue

var nullRuntimeValue = &RuntimeValue{valueType: NullValue}

func NewNullRuntimeValue() *RuntimeValue {
	return nullRuntimeValue
}

// MARK: ArrayRuntimeValue

type ArrayRuntimeValue struct {
	*RuntimeValue
	Keys     []IRuntimeValue
	Elements map[IRuntimeValue]IRuntimeValue
}

func NewArrayRuntimeValue() *ArrayRuntimeValue {
	return &ArrayRuntimeValue{
		RuntimeValue: NewRuntimeValue(ArrayValue),
		Keys:         []IRuntimeValue{},
		Elements:     map[IRuntimeValue]IRuntimeValue{},
	}
}

func NewArrayRuntimeValueFromMap(elements map[IRuntimeValue]IRuntimeValue) *ArrayRuntimeValue {
	keys := []IRuntimeValue{}
	for key := range elements {
		keys = append(keys, key)
	}
	return &ArrayRuntimeValue{
		RuntimeValue: NewRuntimeValue(ArrayValue),
		Keys:         keys,
		Elements:     elements,
	}
}

func (runtimeValue *ArrayRuntimeValue) SetElement(key IRuntimeValue, value IRuntimeValue) phpError.Error {
	if key == nil {
		found := false
		var lastInt int64 = -1
		for i := len(runtimeValue.Keys) - 1; i >= 0; i-- {
			if runtimeValue.Keys[i].GetType() != IntegerValue {
				continue
			}

			if !found {
				lastInt = runtimeValue.Keys[i].(*IntegerRuntimeValue).Value
				found = true
			} else {
				foundInt := runtimeValue.Keys[i].(*IntegerRuntimeValue).Value
				if foundInt > lastInt {
					lastInt = foundInt
				}
			}
		}
		key = NewIntegerRuntimeValue(lastInt + 1)
	}

	if key.GetType() == ArrayValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Arrays and objects can not be used as keys. Doing so will result in a warning: Illegal offset type.
		return phpError.NewWarning("Illegal offset type %s", paramTypeRuntimeValue[key.GetType()])
	}

	if key.GetType() == BooleanValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Bools are cast to ints, too, i.e. the key true will actually be stored under 1 and the key false under 0.
		if key.(*BooleanRuntimeValue).Value {
			key = NewIntegerRuntimeValue(1)
		} else {
			key = NewIntegerRuntimeValue(0)
		}
	}

	if key.GetType() == FloatingValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Floats are also cast to ints, which means that the fractional part will be truncated. E.g. the key 8.7 will actually be stored under 8.
		keyInt, err := lib_intval(key)
		if err != nil {
			return err
		}
		key = NewIntegerRuntimeValue(keyInt)
	}

	if key.GetType() == StringValue && common.IsDecimalLiteral(key.(*StringRuntimeValue).Value) {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Strings containing valid decimal ints, unless the number is preceded by a + sign, will be cast to the int type.
		// E.g. the key "8" will actually be stored under 8. On the other hand "08" will not be cast, as it isn't a valid decimal integer.
		key = NewIntegerRuntimeValue(common.DecimalLiteralToInt64(key.(*StringRuntimeValue).Value))
	}

	if key.GetType() == NullValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Null will be cast to the empty string, i.e. the key null will actually be stored under "".
		key = NewStringRuntimeValue("")
	}

	existingKey, exists := runtimeValue.findKey(key)
	if !exists {
		runtimeValue.Keys = append(runtimeValue.Keys, key)
		runtimeValue.Elements[key] = value
	} else {
		runtimeValue.Elements[existingKey] = value
	}

	return nil
}

func (runtimeValue *ArrayRuntimeValue) findKey(key IRuntimeValue) (IRuntimeValue, bool) {
	if key == nil {
		return NewVoidRuntimeValue(), false
	}

	for k := range runtimeValue.Elements {
		if k.GetType() != key.GetType() {
			continue
		}
		boolean, err := compare(key, "===", k)
		if err != nil {
			return NewVoidRuntimeValue(), false
		}
		if boolean.Value {
			return k, true
		}
	}
	return NewVoidRuntimeValue(), false
}

func (runtimeValue *ArrayRuntimeValue) GetElement(key IRuntimeValue) (IRuntimeValue, bool) {
	key, found := runtimeValue.findKey(key)
	if !found {
		return NewVoidRuntimeValue(), false
	}
	return runtimeValue.Elements[key], true
}

// MARK: BooleanRuntimeValue

type BooleanRuntimeValue struct {
	*RuntimeValue
	Value bool
}

var trueRuntimeValue = &BooleanRuntimeValue{RuntimeValue: NewRuntimeValue(BooleanValue), Value: true}
var falseRuntimeValue = &BooleanRuntimeValue{RuntimeValue: NewRuntimeValue(BooleanValue), Value: false}

func NewBooleanRuntimeValue(value bool) *BooleanRuntimeValue {
	if value {
		return trueRuntimeValue
	}
	return falseRuntimeValue
}

// MARK: IntegerRuntimeValue

type IntegerRuntimeValue struct {
	*RuntimeValue
	Value int64
}

func NewIntegerRuntimeValue(value int64) *IntegerRuntimeValue {
	return &IntegerRuntimeValue{RuntimeValue: NewRuntimeValue(IntegerValue), Value: value}
}

// MARK: FloatingRuntimeValue

type FloatingRuntimeValue struct {
	*RuntimeValue
	Value float64
}

func NewFloatingRuntimeValue(value float64) *FloatingRuntimeValue {
	return &FloatingRuntimeValue{RuntimeValue: NewRuntimeValue(FloatingValue), Value: value}
}

func (value *FloatingRuntimeValue) ToPhpString() string {
	return strconv.FormatFloat(value.Value, 'f', -1, 64)
}

// MARK: StringRuntimeValue

type StringRuntimeValue struct {
	*RuntimeValue
	Value string
}

func NewStringRuntimeValue(value string) *StringRuntimeValue {
	return &StringRuntimeValue{RuntimeValue: NewRuntimeValue(StringValue), Value: value}
}

// MARK: Dump

func DumpRuntimeValue(value IRuntimeValue) {
	switch value.GetType() {
	case ArrayValue:
		fmt.Printf("{ArrayValue: ")
		arrayValue := value.(*ArrayRuntimeValue)
		for _, key := range arrayValue.Keys {
			fmt.Print("Key: ")
			DumpRuntimeValue(key)
			fmt.Print("Value: ")
			DumpRuntimeValue(arrayValue.Elements[key])

		}
		fmt.Printf("}\n")
	case BooleanValue:
		valueStr := "true"
		if value.(*BooleanRuntimeValue).Value {
			valueStr = "false"
		}
		fmt.Printf("{BooleanValue: %s}\n", valueStr)
	case IntegerValue:
		fmt.Printf("{IntegerValue: %d}\n", value.(*IntegerRuntimeValue).Value)
	case FloatingValue:
		fmt.Printf("{FloatingValue: %f}\n", value.(*FloatingRuntimeValue).Value)
	case StringValue:
		fmt.Printf("{StringValue: %s}\n", value.(*StringRuntimeValue).Value)
	default:
		fmt.Printf("Unsupported RuntimeValue type %s\n", value.GetType())
	}
}
