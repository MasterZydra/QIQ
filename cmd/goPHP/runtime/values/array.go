package values

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/config"
	"GoPHP/cmd/goPHP/phpError"
	"fmt"
)

type Array struct {
	*abstractValue
	Keys     []RuntimeValue
	Elements map[string]RuntimeValue
	// Keeping track of next key
	nextKey    int64
	nextKeySet bool
}

func NewArray() *Array {
	return &Array{
		abstractValue: newAbstractValue(ArrayValue),
		Keys:          []RuntimeValue{},
		Elements:      map[string]RuntimeValue{},
	}
}

func NewArrayFromMap(elements map[RuntimeValue]RuntimeValue) *Array {
	array := NewArray()
	for key, value := range elements {
		if err := array.SetElement(key, value); err != nil && config.IsDevMode {
			fmt.Println("NewArrayFromMap: " + err.Error())
		}
	}
	return array
}

func (runtimeValue *Array) keyToMapKey(key RuntimeValue) (string, phpError.Error) {
	if key.GetType() == IntValue {
		return fmt.Sprintf("i_%d", key.(*Int).Value), nil
	}
	if key.GetType() == StrValue {
		return fmt.Sprintf("s_%s", key.(*Str).Value), nil
	}
	return "", phpError.NewError("Array.keyToMapKey: Unsupported key type %s", ToPhpType(key))
}

func (runtimeValue *Array) convertKey(key RuntimeValue) (RuntimeValue, phpError.Error) {
	if key == nil {
		return NewVoid(), nil
	}

	if key.GetType() == ArrayValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Arrays and objects can not be used as keys. Doing so will result in a warning: Illegal offset type.
		return NewVoid(), phpError.NewWarning("Illegal offset type %s", ToPhpType(key))
	}

	if key.GetType() == BoolValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Bools are cast to ints, too, i.e. the key true will actually be stored under 1 and the key false under 0.
		if key.(*Bool).Value {
			return NewInt(1), nil
		}
		return NewInt(0), nil
	}

	if key.GetType() == FloatValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Floats are also cast to ints, which means that the fractional part will be truncated. E.g. the key 8.7 will actually be stored under 8.
		return NewInt(int64(key.(*Float).Value)), nil
	}

	if key.GetType() == StrValue && common.IsDecimalLiteral(key.(*Str).Value) {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Strings containing valid decimal ints, unless the number is preceded by a + sign, will be cast to the int type.
		// E.g. the key "8" will actually be stored under 8. On the other hand "08" will not be cast, as it isn't a valid decimal integer.
		return NewInt(common.DecimalLiteralToInt64(key.(*Str).Value)), nil
	}

	if key.GetType() == NullValue {
		// Spec: https://www.php.net/manual/en/language.types.array.php
		// Null will be cast to the empty string, i.e. the key null will actually be stored under "".
		return NewStr(""), nil
	}

	return key, nil
}

func (runtimeValue *Array) getNextKey(key RuntimeValue) (RuntimeValue, phpError.Error) {
	// If a key is passed
	if key != nil {
		var err phpError.Error
		key, err = runtimeValue.convertKey(key)
		if err != nil {
			return key, err
		}

		// If the passed key is an integer after the convertion
		if key.GetType() == IntValue {
			keyValue := key.(*Int).Value
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
			if runtimeValue.Keys[i].GetType() != IntValue {
				continue
			}

			foundInt := runtimeValue.Keys[i].(*Int).Value
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

	key = NewInt(runtimeValue.nextKey)
	runtimeValue.nextKey++
	return key, nil
}

func (runtimeValue *Array) SetElement(key RuntimeValue, value RuntimeValue) phpError.Error {
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

func (runtimeValue *Array) getMapKey(key RuntimeValue, convertKey bool) (string, bool, phpError.Error) {
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

func (runtimeValue *Array) Contains(key RuntimeValue) bool {
	_, found, err := runtimeValue.getMapKey(key, true)
	if err != nil && config.IsDevMode {
		fmt.Println("Array.Contains: " + err.Error())
		return false
	}
	return found
}

func (runtimeValue *Array) GetElement(key RuntimeValue) (RuntimeValue, bool) {
	mapKey, found, err := runtimeValue.getMapKey(key, true)
	if err != nil {
		return NewVoid(), false
	}
	if !found {
		return NewVoid(), false
	}
	return runtimeValue.Elements[mapKey], true
}
