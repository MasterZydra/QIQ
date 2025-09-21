package optionsInfo

import (
	"QIQ/cmd/qiq/common/os"
	"QIQ/cmd/qiq/config"
	"QIQ/cmd/qiq/ini"
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime"
	"QIQ/cmd/qiq/runtime/funcParamValidator"
	"QIQ/cmd/qiq/runtime/stdlib/variableHandling"
	"QIQ/cmd/qiq/runtime/values"
	"fmt"
	goOs "os"
	goRuntime "runtime"
	"strconv"
	"strings"
)

func Register(environment runtime.Environment) {
	// Category: Options/Info Functions
	environment.AddNativeFunction("getenv", nativeFn_getenv)
	environment.AddNativeFunction("getmygid", nativeFn_getmygid)
	environment.AddNativeFunction("getmypid", nativeFn_getmypid)
	environment.AddNativeFunction("getmyuid", nativeFn_getmyuid)
	environment.AddNativeFunction("ini_get", nativeFn_ini_get)
	environment.AddNativeFunction("ini_set", nativeFn_ini_set)
	environment.AddNativeFunction("phpinfo", nativeFn_phpinfo)
	environment.AddNativeFunction("phpversion", nativeFn_phpversion)
	environment.AddNativeFunction("zend_thread_id", nativeFn_zend_thread_id)
	environment.AddNativeFunction("zend_version", nativeFn_zend_version)

	// Const Category: Options/Info Constants
	// Spec: https://www.php.net/manual/en/info.constants.php
	environment.AddPredefinedConstant("INI_USER", values.NewInt(int64(ini.INI_USER)))
	environment.AddPredefinedConstant("INI_PERDIR", values.NewInt(int64(ini.INI_PERDIR)))
	environment.AddPredefinedConstant("INI_SYSTEM", values.NewInt(int64(ini.INI_SYSTEM)))
	environment.AddPredefinedConstant("INI_ALL", values.NewInt(int64(ini.INI_ALL)))
}

// -------------------------------------- getenv -------------------------------------- MARK: getenv

func nativeFn_getenv(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.getenv.php

	//  getenv(?string $name = null, bool $local_only = false): string|array|false

	// Returns the value of the environment variable name, or false if the environment variable name does not exist.
	// If name is null, all environment variables are returned as an associative array.

	// TODO getenv - add support for $local_only
	args, err := funcParamValidator.NewValidator("getenv").AddParam("$name", []string{"string"}, values.NewNull()).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	if args[0].GetType() == values.NullValue {
		slot, err := context.Env.LookupVariable("$_ENV")
		return slot.Value, err
	}

	envVars, err := context.Env.LookupVariable("$_ENV")
	if err != nil {
		return envVars.Value, err
	}
	envArray := envVars.Value.(*values.Array)
	value, found := envArray.GetElement(args[0])
	if !found {
		return values.NewBool(false), nil
	}
	return value.Value, nil
}

// -------------------------------------- getmygid -------------------------------------- MARK: getmygid

func nativeFn_getmygid(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.getmygid

	_, err := funcParamValidator.NewValidator("getmygid").Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	if os.IS_WIN {
		return values.NewInt(0), nil
	}

	return values.NewInt(int64(goOs.Getgid())), nil
}

// -------------------------------------- getmypid -------------------------------------- MARK: getmypid

func nativeFn_getmypid(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.getmypid

	_, err := funcParamValidator.NewValidator("getmypid").Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewInt(int64(goOs.Getpid())), nil
}

// -------------------------------------- getmyuid -------------------------------------- MARK: getmyuid

func nativeFn_getmyuid(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.getmyuid

	_, err := funcParamValidator.NewValidator("getmyuid").Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	if os.IS_WIN {
		return values.NewInt(0), nil
	}

	return values.NewInt(int64(goOs.Getuid())), nil
}

// -------------------------------------- ini_get -------------------------------------- MARK: ini_get

func nativeFn_ini_get(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ini-get

	args, err := funcParamValidator.NewValidator("ini_get").AddParam("$option", []string{"string"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	value, iniErr := context.Interpreter.GetIni().Get(args[0].(*values.Str).Value)
	if iniErr != nil {
		return values.NewBool(false), nil
	}
	return values.NewStr(value), nil
}

// -------------------------------------- ini_set -------------------------------------- MARK: ini_set

func nativeFn_ini_set(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ini-set

	args, err := funcParamValidator.NewValidator("ini_set").
		AddParam("$option", []string{"string"}, nil).
		AddParam("$value", []string{"string", "int", "float", "bool", "null"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	value, err := variableHandling.StrVal(args[1])
	if err != nil {
		return values.NewVoid(), err
	}

	option := args[0].(*values.Str).Value

	oldValue, err := context.Interpreter.GetIni().Get(option)
	if err != nil {
		return values.NewBool(false), nil
	}
	err = context.Interpreter.GetIni().Set(option, value, ini.INI_USER)
	if err != nil {
		return values.NewBool(false), nil
	}

	return values.NewStr(oldValue), nil
}

// -------------------------------------- phpinfo -------------------------------------- MARK: phpinfo

func nativeFn_phpinfo(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.phpinfo.php

	_, err := funcParamValidator.NewValidator("phpinfo").
		Validate(args)
		// TODO phpinfo param $flags
	if err != nil {
		return values.NewVoid(), err
	}

	printPhpInfo(context)

	return values.NewBool(true), nil
}

// -------------------------------------- phpversion -------------------------------------- MARK: phpversion

func nativeFn_phpversion(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.phpversion.php

	_, err := funcParamValidator.NewValidator("phpversion").
		AddParam("$extension", []string{"null", "string"}, values.NewNull()).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// TODO phpversion param $extension

	return values.NewStr(config.Version), nil
}

// -------------------------------------- zend_thread_id -------------------------------------- MARK: zend_thread_id

func nativeFn_zend_thread_id(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.zend-thread-id.php

	_, phpErr := funcParamValidator.NewValidator("zend_thread_id").Validate(args)
	if phpErr != nil {
		return values.NewVoid(), phpErr
	}

	// Source: https://dev.to/leapcell/how-to-get-the-goroutine-id-1h5o
	var (
		buf [64]byte
		n   = goRuntime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}

	return values.NewInt(int64(id)), nil
}

// -------------------------------------- zend_version -------------------------------------- MARK: zend_version

func nativeFn_zend_version(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.zend-version.php

	_, err := funcParamValidator.NewValidator("zend_version").Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewStr(config.QIQVersion), nil
}

// TODO assert
// TODO cli_get_process_title
// TODO cli_set_process_title
// TODO dl
// TODO extension_loaded
// TODO gc_collect_cycles
// TODO gc_disable
// TODO gc_enable
// TODO gc_enabled
// TODO gc_mem_caches
// TODO gc_status
// TODO get_cfg_var
// TODO get_current_user
// TODO get_defined_constants
// TODO get_extension_funcs
// TODO get_include_path
// TODO get_included_files
// TODO get_loaded_extensions
// TODO get_required_files
// TODO get_resources
// TODO getlastmod
// TODO getmyinode
// TODO getopt
// TODO getrusage
// TODO ini_alter
// TODO ini_get_all
// TODO ini_parse_quantity
// TODO ini_restore
// TODO memory_get_peak_usage
// TODO memory_get_usage
// TODO memory_reset_peak_usage
// TODO php_ini_loaded_file
// TODO php_ini_scanned_files
// TODO php_sapi_name
// TODO php_uname
// TODO phpcredits
// TODO putenv
// TODO set_include_path
// TODO set_time_limit
// TODO sys_get_temp_dir
// TODO version_compare
