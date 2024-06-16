package ini

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/phpError"
	"strings"
)

type Ini struct {
	ErrorReporting int64
	ShortOpenTag   bool
}

func NewDevIni() *Ini {
	defaultIni := NewDefaultIni()
	defaultIni.ErrorReporting = phpError.E_ALL
	return defaultIni
}

func NewIniFromArray(ini []string) *Ini {
	settings := map[string]string{}
	for _, setting := range ini {
		parts := strings.Split(setting, "=")
		settings[parts[0]] = parts[1]
	}

	defaultIni := NewDefaultIni()

	if value, found := settings["error_reporting"]; found {
		defaultIni.ErrorReporting = iniIntStrToInt(value)
	}
	if value, found := settings["short_open_tag"]; found {
		defaultIni.ShortOpenTag = iniBoolStrToBool(value)
	}

	return defaultIni
}

func NewDefaultIni() *Ini {
	return &Ini{
		ErrorReporting: 0,
		ShortOpenTag:   false,
	}
}

func iniBoolStrToBool(str string) bool {
	if str == "1" || strings.ToLower(str) == "on" {
		return true
	}
	return false
}

func iniIntStrToInt(str string) int64 {
	if common.IsIntegerLiteralWithSign(str) {
		intVal, _ := common.IntegerLiteralToInt64WithSign(str)
		return intVal
	}
	return -1
}
