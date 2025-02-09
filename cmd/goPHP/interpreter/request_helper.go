package interpreter

import (
	"GoPHP/cmd/goPHP/common"
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
		// if strings.Contains(key, "[") {
		// 	if parseQueryKey(key, value, result) != nil {
		// 		return result, err
		// 	}
		// } else {
		var keyValue IRuntimeValue
		if common.IsIntegerLiteral(key) {
			intValue, _ := common.IntegerLiteralToInt64(key)
			keyValue = NewIntegerRuntimeValue(intValue)
		} else {
			keyValue = NewStringRuntimeValue(key)
		}
		result.SetElement(keyValue, NewStringRuntimeValue(value))
		// }
	}

	return result, nil
}

// func parseQueryKey(key string, value string, result *ArrayRuntimeValue) error {
// 	// Idea:
// 	// Convert key and value into valid PHP and exec with interpreter

// 	// Input: 123[]=abc
// 	// Key: 123[]
// 	// Value: abc
// 	// PHP: $array[123][] = "abc";

// 	// Input: a23[]=abc
// 	// Key: a23[]
// 	// Value: abc
// 	// PHP: $array["a23"][] = "abc";

// 	// firstKey, key, _ := strings.Cut(key, "[")
// 	// key = "[" + key

// 	// phpArrayKeys := []string{}

// 	// if common.IsDecimalLiteral(firstKey) {
// 	// 	firstKeyInt, _ := common.IntegerLiteralToInt64(firstKey)
// 	// 	phpArrayKeys = append(phpArrayKeys, fmt.Sprintf("[%d]", firstKeyInt))
// 	// } else {
// 	// 	phpArrayKeys = append(phpArrayKeys, "["+firstKey+"]")
// 	// }

// 	// for key != "" {
// 	// 	nextKey, key, _ := strings.Cut(key, "[")
// 	// 	key = "[" + key
// 	// }

// 	return nil
// }
