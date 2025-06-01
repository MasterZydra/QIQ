package array

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/stdlib/variableHandling"
	"GoPHP/cmd/goPHP/runtime/values"
	"slices"
)

func RemoveByKey(array *values.Array, key values.RuntimeValue) phpError.Error {
	if key == nil {
		return phpError.NewError("Array.RemoveByKey: Key can not be nil")
	}
	mapKey, found, err := array.GetMapKey(key, true)
	if err != nil {
		return err
	}
	if !found {
		return phpError.NewError("Array.RemoveByKey: Key %s not found", values.ToPhpType(key))
	}

	delete(array.Elements, mapKey)
	array.Keys = slices.DeleteFunc(array.Keys, func(k values.RuntimeValue) bool {
		match, _ := variableHandling.Compare(k, "===", key)
		return match.Value
	})
	return nil
}
