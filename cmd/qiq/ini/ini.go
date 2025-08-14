package ini

import (
	"QIQ/cmd/qiq/common"
	"QIQ/cmd/qiq/phpError"
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
	"always_populate_raw_post_data": INI_ALL,
	"arg_separator.input":           INI_SYSTEM,
	"arg_separator.output":          INI_ALL,
	"default_charset":               INI_ALL,
	"error_reporting":               INI_ALL,
	"expose_php":                    INI_SYSTEM,
	"file_uploads":                  INI_SYSTEM,
	"filter.default":                INI_PERDIR,
	"input_encoding":                INI_ALL,
	"internal_encoding":             INI_ALL,
	"max_input_nesting_level":       INI_PERDIR,
	"max_input_vars":                INI_PERDIR,
	"mbstring.encoding_translation": INI_PERDIR,
	"open_basedir":                  INI_ALL,
	"output_encoding":               INI_ALL,
	"post_max_size":                 INI_PERDIR,
	"register_argc_argv":            INI_PERDIR,
	"short_open_tag":                INI_PERDIR,
	"session.name":                  INI_ALL,
	"session.save_path":             INI_ALL,
	"upload_max_filesize":           INI_PERDIR,
	"upload_tmp_dir":                INI_SYSTEM,
	"variables_order":               INI_PERDIR,
}

var boolDirectives = []string{
	"expose_php", "file_uploads", "mbstring.encoding_translation",
	"register_argc_argv", "short_open_tag",
}

var intDirectives = []string{
	"error_reporting", "max_input_nesting_level", "max_input_vars",
}

type Ini struct {
	directives map[string]string
}

func NewDefaultIni() *Ini {
	return &Ini{
		directives: map[string]string{
			// Ini Directives:
			"always_populate_raw_post_data": "",
			"arg_separator.input":           "&",
			"arg_separator.output":          "&",
			"default_charset":               "UTF-8",
			"error_reporting":               "0",
			"expose_php":                    "",
			"file_uploads":                  "1",
			"filter.default":                "unsafe_raw",
			"input_encoding":                "",
			"internal_encoding":             "",
			"max_input_nesting_level":       "64",
			"max_input_vars":                "1000",
			"mbstring.encoding_translation": "",
			"open_basedir":                  "",
			"output_encoding":               "",
			"post_max_size":                 "8M",
			"register_argc_argv":            "",
			"short_open_tag":                "",
			"session.name":                  "PHPSESSID",
			"session.save_path":             "",
			"upload_max_filesize":           "2M",
			"upload_tmp_dir":                "",
			"variables_order":               "EGPCS",
		},
	}
}

func NewDevIni() *Ini {
	defaultIni := NewDefaultIni()
	defaultIni.Set("error_reporting", "32767", INI_ALL)
	defaultIni.Set("expose_php", "1", INI_ALL)
	return defaultIni
}

func NewDevIniFromArray(ini []string) (*Ini, phpError.Error) {
	defaultIni := NewDevIni()
	var resultErr phpError.Error = nil

	for _, setting := range ini {
		parts := strings.SplitN(setting, "=", 2)
		err := defaultIni.Set(parts[0], parts[1], INI_ALL)
		if err != nil {
			if resultErr == nil {
				resultErr = phpError.NewError("%s", err)
			} else {
				resultErr = phpError.NewError("%s\n%s", resultErr, err)
			}
		}
	}

	return defaultIni, resultErr
}

func (ini *Ini) Set(directive string, value string, source int) phpError.Error {
	changeable, found := allowedDirectives[directive]
	if !found {
		return phpError.NewError("Directive %s not found", directive)
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
		if !common.IsIntegerLiteralWithSign(value, false) {
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
	if common.IsIntegerLiteralWithSign(value, false) {
		intVal, _ := common.IntegerLiteralToInt64WithSign(value, false)
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
