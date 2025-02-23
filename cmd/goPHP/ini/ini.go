package ini

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/phpError"
	"slices"
	"strings"
)

// Spec: https://www.php.net/manual/en/info.constants.php#constant.ini-system
const (
	INI_USER   int = 1 // Entry can be set in user scripts (like with ini_set()) or in the Windows registry. Entry can be set in .user.ini
	INI_PERDIR int = 2 // Entry can be set in php.ini, .htaccess, httpd.conf or .user.ini
	INI_SYSTEM int = 4 // Entry can be set in php.ini or httpd.conf
	INI_ALL    int = 7 // Entry can be set anywhere
)

var allowedDirectives = map[string]int{
	"error_reporting":    INI_ALL,
	"register_argc_argv": INI_PERDIR,
	"short_open_tag":     INI_PERDIR,
	"variables_order":    INI_PERDIR,
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
			"variables_order":    "EGPCS",
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
		parts := strings.Split(setting, "=")
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
