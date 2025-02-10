package interpreter

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/ini"
	"fmt"
	"net/url"
	"regexp"
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

		key, err := url.QueryUnescape(fixPercentEscaping(key))
		if err != nil {
			return result, err
		}
		// fmt.Println(key)

		value, err = url.QueryUnescape(value)
		if err != nil {
			return result, err
		}
		if strings.Contains(key, "[") && strings.Contains(key, "]") {
			result, err = parseQueryKey(key, value, result)
			if err != nil {
				return result, err
			}
		} else {
			key = replaceSpecialCharacters(key)

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
		key = strings.TrimPrefix(key, "[")
		var nextKey string
		nextKey, key, _ = strings.Cut(key, "]")
		phpArrayKeys = append(phpArrayKeys, nextKey)
		for key != "" && !strings.HasPrefix(key, "[") {
			key = key[1:]
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

// This fix is required because "url.QueryUnescape()" cannot handle an unescaped percent
func fixPercentEscaping(key string) string {
	re, _ := regexp.Compile("%([^0-9A-Fa-f]|$)")
	// Replace only the '%' character with '%25' without affecting the following character
	return re.ReplaceAllStringFunc(key, func(match string) string {
		return "%25" + match[1:]
	})
}

func replaceSpecialCharacters(key string) string {
	return strings.NewReplacer(
		" ", "_",
		"+", "_",
		"[", "_",
		".", "_",
	).Replace(key)
}
