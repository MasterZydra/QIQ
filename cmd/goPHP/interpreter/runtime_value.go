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
	Elements map[string]IRuntimeValue
	// Keeping track of next key
	nextKey    int64
	nextKeySet bool
}

func NewArrayRuntimeValue() *ArrayRuntimeValue {
	return &ArrayRuntimeValue{
		RuntimeValue: NewRuntimeValue(ArrayValue),
		Keys:         []IRuntimeValue{},
		Elements:     map[string]IRuntimeValue{},
	}
}

func NewArrayRuntimeValueFromMap(elements map[IRuntimeValue]IRuntimeValue) *ArrayRuntimeValue {
	array := NewArrayRuntimeValue()
	for key, value := range elements {
		if err := array.SetElement(key, value); err != nil {
			printDev("NewArrayRuntimeValueFromMap: " + err.Error())
		}
	}
	return array
}

func (runtimeValue *ArrayRuntimeValue) keyToMapKey(key IRuntimeValue) (string, phpError.Error) {
	if key.GetType() == IntegerValue {
		return fmt.Sprintf("i_%d", key.(*IntegerRuntimeValue).Value), nil
	}
	if key.GetType() == StringValue {
		return fmt.Sprintf("s_%s", key.(*StringRuntimeValue).Value), nil
	}
	return "", phpError.NewError("ArrayRuntimeValue.keyToMapKey: Unsupported key type %s", paramTypeRuntimeValue[key.GetType()])
}

func (runtimeValue *ArrayRuntimeValue) convertKey(key IRuntimeValue) (IRuntimeValue, phpError.Error) {
	if key == nil {
		return NewVoidRuntimeValue(), nil
	}

	if key.GetType() == ArrayValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Arrays and objects can not be used as keys. Doing so will result in a warning: Illegal offset type.
		return NewVoidRuntimeValue(), phpError.NewWarning("Illegal offset type %s", paramTypeRuntimeValue[key.GetType()])
	}

	if key.GetType() == BooleanValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Bools are cast to ints, too, i.e. the key true will actually be stored under 1 and the key false under 0.
		if key.(*BooleanRuntimeValue).Value {
			return NewIntegerRuntimeValue(1), nil
		}
		return NewIntegerRuntimeValue(0), nil
	}

	if key.GetType() == FloatingValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Floats are also cast to ints, which means that the fractional part will be truncated. E.g. the key 8.7 will actually be stored under 8.
		keyInt, err := lib_intval(key)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return NewIntegerRuntimeValue(keyInt), nil
	}

	if key.GetType() == StringValue && common.IsDecimalLiteral(key.(*StringRuntimeValue).Value) {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Strings containing valid decimal ints, unless the number is preceded by a + sign, will be cast to the int type.
		// E.g. the key "8" will actually be stored under 8. On the other hand "08" will not be cast, as it isn't a valid decimal integer.
		return NewIntegerRuntimeValue(common.DecimalLiteralToInt64(key.(*StringRuntimeValue).Value)), nil
	}

	if key.GetType() == NullValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Null will be cast to the empty string, i.e. the key null will actually be stored under "".
		return NewStringRuntimeValue(""), nil
	}

	return key, nil
}

func (runtimeValue *ArrayRuntimeValue) getNextKey(key IRuntimeValue) (IRuntimeValue, phpError.Error) {
	// If a key is passed
	if key != nil {
		var err phpError.Error
		key, err = runtimeValue.convertKey(key)
		if err != nil {
			return key, err
		}

		// If the passed key is an integer after the convertion
		if key.GetType() == IntegerValue {
			keyValue := key.(*IntegerRuntimeValue).Value
			// If no key is stored yet or the passed key is greater than nextKey
			if !runtimeValue.nextKeySet ||
				(runtimeValue.nextKeySet && keyValue > runtimeValue.nextKey) {
				// Store value + 1 as next key
				runtimeValue.nextKey = keyValue + 1
				runtimeValue.nextKeySet = true
			}
		}

		return key, nil
	}

	// If no key is passed and no nextKey is set yet
	if !runtimeValue.nextKeySet {
		// lastFoundInt is -1 because at the end of the loop it will be increased by one
		var lastFoundInt int64 = -1
		found := false
		// Iterate all keys and search for integer values
		for i := len(runtimeValue.Keys) - 1; i >= 0; i-- {
			if runtimeValue.Keys[i].GetType() != IntegerValue {
				continue
			}

			foundInt := runtimeValue.Keys[i].(*IntegerRuntimeValue).Value
			if !found {
				lastFoundInt = foundInt
				found = true
				continue
			}

			if foundInt > lastFoundInt {
				lastFoundInt = foundInt
			}
		}
		runtimeValue.nextKey = lastFoundInt + 1
		runtimeValue.nextKeySet = true
	}

	key = NewIntegerRuntimeValue(runtimeValue.nextKey)
	runtimeValue.nextKey++
	return key, nil
}

func (runtimeValue *ArrayRuntimeValue) SetElement(key IRuntimeValue, value IRuntimeValue) phpError.Error {
	key, err := runtimeValue.getNextKey(key)
	if err != nil {
		return err
	}

	mapKey, found, err := runtimeValue.getMapKey(key, false)
	if err != nil {
		return err
	}

	if !found {
		runtimeValue.Keys = append(runtimeValue.Keys, key)
	}
	runtimeValue.Elements[mapKey] = value

	return nil
}

func (runtimeValue *ArrayRuntimeValue) getMapKey(key IRuntimeValue, convertKey bool) (string, bool, phpError.Error) {
	if convertKey {
		var err phpError.Error
		key, err = runtimeValue.convertKey(key)
		if err != nil {
			return "", false, err
		}
	}

	mapKey, err := runtimeValue.keyToMapKey(key)
	if err != nil {
		return "", false, err
	}

	_, found := runtimeValue.Elements[mapKey]

	return mapKey, found, nil
}

func (runtimeValue *ArrayRuntimeValue) Contains(key IRuntimeValue) bool {
	_, found, err := runtimeValue.getMapKey(key, true)
	if err != nil {
		printDev("ArrayRuntimeValue.Contains: " + err.Error())
		return false
	}
	return found
}

func (runtimeValue *ArrayRuntimeValue) GetElement(key IRuntimeValue) (IRuntimeValue, bool) {
	mapKey, found, err := runtimeValue.getMapKey(key, true)
	if err != nil {
		return NewVoidRuntimeValue(), false
	}
	if !found {
		return NewVoidRuntimeValue(), false
	}
	return runtimeValue.Elements[mapKey], true
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
			value, _ := arrayValue.GetElement(key)
			DumpRuntimeValue(value)

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
