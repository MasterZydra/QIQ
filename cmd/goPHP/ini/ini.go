package ini

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/phpError"
	"slices"
	"strings"
)

// Spec: https://www.php.net/manual/en/info.constants.php#constant.ini-system
const (
	// Entry can be set in user scripts (like with ini_set()) or in the Windows registry. Entry can be set in .user.ini
	INI_USER int = 1
	// Entry can be set in php.ini, .htaccess, httpd.conf or .user.ini
	INI_PERDIR int = 2
	// Entry can be set in php.ini or httpd.conf
	INI_SYSTEM int = 4
	// Entry can be set anywhere
	INI_ALL int = 7
)

var allowedDirectives = map[string]int{
	"arg_separator.input":  INI_SYSTEM,
	"arg_separator.output": INI_ALL,
	"default_charset":      INI_ALL,
	"error_reporting":      INI_ALL,
	"input_encoding":       INI_ALL,
	"internal_encoding":    INI_ALL,
	"output_encoding":      INI_ALL,
	"register_argc_argv":   INI_PERDIR,
	"short_open_tag":       INI_PERDIR,
	"variables_order":      INI_PERDIR,
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
			"arg_separator.input":  "&",
			"arg_separator.output": "&",
			"default_charset":      "UTF-8",
			"error_reporting":      "0",
			"input_encoding":       "",
			"internal_encoding":    "",
			"output_encoding":      "",
			"register_argc_argv":   "",
			"short_open_tag":       "",
			"variables_order":      "EGPCS",
		},
	}
}

func NewDevIni() *Ini {
	defaultIni := NewDefaultIni()
	defaultIni.Set("error_reporting", "32767", INI_ALL)
	return defaultIni
}

func NewIniFromArray(ini []string) *Ini {
	defaultIni := NewDefaultIni()

	for _, setting := range ini {
		parts := strings.SplitN(setting, "=", 2)
		defaultIni.Set(parts[0], parts[1], INI_ALL)
	}

	return defaultIni
}

func (ini *Ini) Set(directive string, value string, source int) phpError.Error {
	changeable, found := allowedDirectives[directive]
	if !found {
		return phpError.NewError("Directive not found")
	}

	if changeable&source == 0 {
		return phpError.NewError("Not allowed to change %s", directive)
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

	ini.directives[directive] = value
	return nil
}

func (ini *Ini) Get(directive string) (string, phpError.Error) {
	if _, found := allowedDirectives[directive]; !found {
		return "", phpError.NewError("Directive not found")
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

func (ini *Ini) GetStr(directive string) string {
	value, err := ini.Get(directive)
	if err != nil {
		return ""
	}
	return value
}
