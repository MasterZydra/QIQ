package interpreter

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/ini"
	"fmt"
	"net/url"
	"strings"
)

func parseQuery(query string) (*ArrayRuntimeValue, error) {
	result := NewArrayRuntimeValue()

	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if key == "" {
			continue
		}

		// Get parameters without key e.g. ab+cd+ef
		// TODO this is only correct if it is in "phpt mode". "Normal" GET will parse it differently
		// ab+cd+ef => array(1) { ["ab_cd_ef"]=> string(0) "" }
		if !strings.Contains(key, "=") && strings.Contains(key, "+") {
			parts := strings.Split(key, "+")
			for i := 0; i < len(parts); i++ {
				result.SetElement(nil, NewStringRuntimeValue(parts[i]))
			}
			continue
		}

		// Get parameter with key-value-pair
		key, value, _ := strings.Cut(key, "=")
		key, err := url.QueryUnescape(key)
		if err != nil {
			return result, err
		}
		value, err = url.QueryUnescape(value)
		if err != nil {
			return result, err
		}
		if strings.Contains(key, "[") {
			result, err = parseQueryKey(key, value, result)
			if err != nil {
				return result, err
			}
		} else {
			var keyValue IRuntimeValue
			if common.IsIntegerLiteral(key) {
				intValue, _ := common.IntegerLiteralToInt64(key)
				keyValue = NewIntegerRuntimeValue(intValue)
			} else {
				keyValue = NewStringRuntimeValue(key)
			}
			result.SetElement(keyValue, NewStringRuntimeValue(value))
		}
	}

	return result, nil
}

func parseQueryKey(key string, value string, result *ArrayRuntimeValue) (*ArrayRuntimeValue, error) {
	// The parsing of a complex key with arrays is solved by using the interpreter itself:
	// The key and value is transformed into valid PHP code and executed.
	// Example:
	//   Input: 123[][12][de]=abc
	//   Key:   123[][12][de]
	//   Value: abc
	//   PHP:   $array[123][][12]["de"] = "abc";

	firstKey, key, _ := strings.Cut(key, "[")
	key = "[" + key

	phpArrayKeys := []string{firstKey}

	for key != "" {
		if strings.HasPrefix(key, "[]]") {
			phpArrayKeys = append(phpArrayKeys, "")
			key = strings.TrimPrefix(key, "[]]")
			continue
		}

		if strings.HasPrefix(key, "[]") {
			phpArrayKeys = append(phpArrayKeys, "")
			key = strings.TrimPrefix(key, "[]")
			continue
		}

		key = strings.TrimPrefix(key, "[")
		var nextKey string
		nextKey, key, _ = strings.Cut(key, "]")
		phpArrayKeys = append(phpArrayKeys, nextKey)
		for strings.HasPrefix(key, "]") {
			key = strings.TrimPrefix(key, "]")
		}

		if key == "" {
			break
		}
	}

	php := "<?php $array"
	for _, phpArrayKey := range phpArrayKeys {
		if phpArrayKey == "" {
			php += "[]"
		} else if common.IsIntegerLiteral(phpArrayKey) {
			phpArrayKeyInt, _ := common.IntegerLiteralToInt64(phpArrayKey)
			php += fmt.Sprintf("[%d]", phpArrayKeyInt)
		} else {
			php += "['" + phpArrayKey + "']"
		}
	}
	php += " = '" + value + "';"

	interpreter := NewInterpreter(ini.NewDefaultIni(), &Request{}, "")
	interpreter.env.declareVariable("$array", result)
	_, err := interpreter.Process(php)

	return interpreter.env.variables["$array"].(*ArrayRuntimeValue), err
}
