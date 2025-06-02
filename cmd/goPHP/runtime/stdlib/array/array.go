package array

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/funcParamValidator"
	"GoPHP/cmd/goPHP/runtime/values"
	"slices"
)

func Register(environment runtime.Environment) {
	// Category: Array Functions
	environment.AddNativeFunction("array_first", nativeFn_array_first)
	environment.AddNativeFunction("array_key_exists", nativeFn_array_key_exists)
	environment.AddNativeFunction("array_key_first", nativeFn_array_key_first)
	environment.AddNativeFunction("array_key_last", nativeFn_array_key_last)
	environment.AddNativeFunction("array_last", nativeFn_array_last)
	environment.AddNativeFunction("array_pop", nativeFn_array_pop)
	environment.AddNativeFunction("array_push", nativeFn_array_push)
	environment.AddNativeFunction("key_exists", nativeFn_array_key_exists)
}

// -------------------------------------- array_first -------------------------------------- MARK: array_first

func nativeFn_array_first(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://php.watch/versions/8.5/array_first-array_last
	args, err := funcParamValidator.NewValidator("array_first").
		AddParam("$array", []string{"array"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	array := args[0].(*values.Array)
	if array.IsEmpty() {
		return values.NewNull(), nil
	}

	value, _ := array.GetElement(FirstKey(array))
	return value, nil
}

// -------------------------------------- array_key_exists -------------------------------------- MARK: array_key_exists

func nativeFn_array_key_exists(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("array_key_exists").
		AddParam("$key", []string{"string", "int", "float", "bool", "resource", "null"}, nil).
		AddParam("$array", []string{"array"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	boolean, err := lib_array_key_exists(args[0], args[1].(*values.Array))
	return values.NewBool(boolean), err
}

func lib_array_key_exists(key values.RuntimeValue, array *values.Array) (bool, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.array-key-exists.php

	// TODO lib_array_key_exists - allowedKeyTypes - resource
	allowedKeyTypes := []values.ValueType{values.StrValue, values.IntValue, values.FloatValue, values.BoolValue, values.NullValue}

	if !slices.Contains(allowedKeyTypes, key.GetType()) {
		return false, phpError.NewError("Values of type %s are not allowed as array key", key.GetType())
	}

	return array.Contains(key), nil
}

// -------------------------------------- array_key_first -------------------------------------- MARK: array_key_first

func nativeFn_array_key_first(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.array-key-first.php
	args, err := funcParamValidator.NewValidator("array_key_first").
		AddParam("$array", []string{"array"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return FirstKey(args[0].(*values.Array)), nil
}

func FirstKey(array *values.Array) values.RuntimeValue {
	// Spec: https://www.php.net/manual/en/function.array-key-first.php
	if array.IsEmpty() {
		return values.NewNull()
	}
	return array.Keys[0]
}

// -------------------------------------- array_key_last -------------------------------------- MARK: array_key_last

func nativeFn_array_key_last(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.array-key-last.php
	args, err := funcParamValidator.NewValidator("array_key_last").
		AddParam("$array", []string{"array"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return LastKey(args[0].(*values.Array)), nil
}

func LastKey(array *values.Array) values.RuntimeValue {
	// Spec: https://www.php.net/manual/en/function.array-key-last.php
	if array.IsEmpty() {
		return values.NewNull()
	}
	return array.Keys[len(array.Keys)-1]
}

// -------------------------------------- array_last -------------------------------------- MARK: array_last

func nativeFn_array_last(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://php.watch/versions/8.5/array_first-array_last
	args, err := funcParamValidator.NewValidator("array_last").
		AddParam("$array", []string{"array"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	array := args[0].(*values.Array)
	if array.IsEmpty() {
		return values.NewNull(), nil
	}

	value, _ := array.GetElement(LastKey(array))
	return value, nil
}

// -------------------------------------- array_pop -------------------------------------- MARK: array_pop

func nativeFn_array_pop(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.array-pop.php
	args, err := funcParamValidator.NewValidator("array_pop").
		AddParam("$array", []string{"array"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	array := args[0].(*values.Array)
	if array.IsEmpty() {
		return values.NewNull(), nil
	}

	lastKey := LastKey(array)
	value, _ := array.GetElement(lastKey)
	if err := RemoveByKey(array, lastKey); err != nil {
		return values.NewVoid(), err
	}

	return value, nil
}

// -------------------------------------- array_push -------------------------------------- MARK: array_push

func nativeFn_array_push(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.array-push.php
	args, err := funcParamValidator.NewValidator("array_push").
		AddParam("$array", []string{"array"}, nil).
		AddVariableLenParam("$values", []string{"mixed"}).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	array := args[0].(*values.Array)
	arrayValues := args[1].(*values.Array)
	for _, key := range arrayValues.Keys {
		argValue, _ := arrayValues.GetElement(key)
		if err := array.SetElement(nil, argValue); err != nil {
			return values.NewVoid(), err
		}
	}

	return values.NewInt(int64(len(array.Keys))), nil
}

// TODO array
// TODO array_all
// TODO array_any
// TODO array_change_key_case
// TODO array_chunk
// TODO array_column
// TODO array_combine
// TODO array_count_values
// TODO array_diff
// TODO array_diff_assoc
// TODO array_diff_key
// TODO array_diff_uassoc
// TODO array_diff_ukey
// TODO array_fill
// TODO array_fill_keys
// TODO array_filter
// TODO array_find
// TODO array_find_key
// TODO array_flip
// TODO array_intersect
// TODO array_intersect_assoc
// TODO array_intersect_key
// TODO array_intersect_uassoc
// TODO array_intersect_ukey
// TODO array_is_list
// TODO array_keys
// TODO array_map
// TODO array_merge
// TODO array_merge_recursive
// TODO array_multisort
// TODO array_pad
// TODO array_product
// TODO array_rand
// TODO array_reduce
// TODO array_replace
// TODO array_replace_recursive
// TODO array_reverse
// TODO array_search
// TODO array_shift
// TODO array_slice
// TODO array_splice
// TODO array_sum
// TODO array_udiff
// TODO array_udiff_assoc
// TODO array_udiff_uassoc
// TODO array_uintersect
// TODO array_uintersect_assoc
// TODO array_uintersect_uassoc
// TODO array_unique
// TODO array_unshift
// TODO array_values
// TODO array_walk
// TODO array_walk_recursive
// TODO arsort
// TODO asort
// TODO compact
// TODO count
// TODO current
// TODO end
// TODO extract
// TODO in_array
// TODO key
// TODO krsort
// TODO ksort
// TODO list
// TODO natcasesort
// TODO natsort
// TODO next
// TODO pos
// TODO prev
// TODO range
// TODO reset
// TODO rsort
// TODO shuffle
// TODO sizeof
// TODO sort
// TODO uasort
// TODO uksort
// TODO usort
