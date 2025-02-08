package ini

import (
	"GoPHP/cmd/goPHP/common"
	"fmt"
	"slices"
	"strings"
)

var allowedDirectives = []string{
	"error_reporting", "register_argc_argv", "short_open_tag",
}

var boolDirectives = []string{
	"register_argc_argv", "short_open_tag",
}

var intDirectives = []string{
	"error_reporting",
}

type Ini struct {
	directives map[string]string
}

func NewDefaultIni() *Ini {
	return &Ini{
		directives: map[string]string{
			"error_reporting":    "0",
			"register_argc_argv": "",
			"short_open_tag":     "",
		},
	}
}

func NewDevIni() *Ini {
	defaultIni := NewDefaultIni()
	defaultIni.Set("error_reporting", "32767")
	return defaultIni
}

func NewIniFromArray(ini []string) *Ini {
	defaultIni := NewDefaultIni()

	for _, setting := range ini {
		parts := strings.Split(setting, "=")
		defaultIni.Set(parts[0], parts[1])
	}

	return defaultIni
}

func (ini *Ini) Set(directive string, value string) error {
	if !slices.Contains(allowedDirectives, directive) {
		return fmt.Errorf("Directive not found")
	}

	if slices.Contains(boolDirectives, directive) {
		if value == "1" || strings.ToLower(value) == "on" {
			ini.directives[directive] = "1"
			return nil
		}
		ini.directives[directive] = ""
		return nil
	}

	if slices.Contains(intDirectives, directive) {
		if !common.IsIntegerLiteralWithSign(value) {
			return nil
		}
		ini.directives[directive] = value
		return nil
	}

	return fmt.Errorf("Ini.Set: Unsupported directive type")
}

func (ini *Ini) Get(directive string) (string, error) {
	if !slices.Contains(allowedDirectives, directive) {
		return "", fmt.Errorf("Directive not found")
	}

	return ini.directives[directive], nil
}

func (ini *Ini) GetBool(directive string) bool {
	value, err := ini.Get(directive)
	if err != nil {
		return false
	}
	return value == "1"
}

func (ini *Ini) GetInt(directive string) int64 {
	value, err := ini.Get(directive)
	if err != nil {
		return -1
	}
	if common.IsIntegerLiteralWithSign(value) {
		intVal, _ := common.IntegerLiteralToInt64WithSign(value)
		return intVal
	}
	return -1
}
