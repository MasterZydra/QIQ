package interpreter

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/phpError"
	"fmt"
	"math"
	"time"
)

func registerNativeDateTimeFunctions(environment *Environment) {
	environment.nativeFunctions["checkdate"] = nativeFn_checkdate
	environment.nativeFunctions["getdate"] = nativeFn_getdate
	environment.nativeFunctions["localtime"] = nativeFn_localtime
	environment.nativeFunctions["microtime"] = nativeFn_microtime
	environment.nativeFunctions["mktime"] = nativeFn_mktime
	environment.nativeFunctions["time"] = nativeFn_time
}

// ------------------- MARK: checkdate -------------------

func nativeFn_checkdate(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("checkdate").
		addParam("$month", []string{"int"}, nil).addParam("$day", []string{"int"}, nil).addParam("$year", []string{"int"}, nil).
		validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://www.php.net/manual/en/function.checkdate.php
	year := args[2].(*IntegerRuntimeValue).Value
	month := args[0].(*IntegerRuntimeValue).Value
	day := args[1].(*IntegerRuntimeValue).Value

	// Spec: https://www.php.net/manual/en/function.checkdate.php
	// The year is between 1 and 32767 inclusive.
	if year < 1 || year > 32767 {
		return NewBooleanRuntimeValue(false), nil
	}

	// Spec: https://www.php.net/manual/en/function.checkdate.php
	// The month is between 1 and 12 inclusive.
	if month < 1 || month > 12 {
		return NewBooleanRuntimeValue(false), nil
	}

	// Spec: https://www.php.net/manual/en/function.checkdate.php
	// The day is within the allowed number of days for the given month. Leap years are taken into consideration.
	return NewBooleanRuntimeValue(day >= 1 && day <= int64(common.DaysIn(time.Month(month), int(year)))), nil
}

// ------------------- MARK: getdate -------------------

func nativeFn_getdate(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("getdate").addParam("$timestamp", []string{"int"}, NewNullRuntimeValue()).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://www.php.net/manual/en/function.getdate.php

	// If timestamp is omitted, use the default value `time()`
	if args[0].GetType() == NullValue {
		args[0] = lib_time()
	}

	timestamp := time.Unix(args[0].(*IntegerRuntimeValue).Value, 0)
	array := NewArrayRuntimeValue()
	array.SetElement(NewStringRuntimeValue("seconds"), NewIntegerRuntimeValue(int64(timestamp.UTC().Second())))
	array.SetElement(NewStringRuntimeValue("minutes"), NewIntegerRuntimeValue(int64(timestamp.UTC().Minute())))
	array.SetElement(NewStringRuntimeValue("hours"), NewIntegerRuntimeValue(int64(timestamp.UTC().Hour())))
	array.SetElement(NewStringRuntimeValue("mday"), NewIntegerRuntimeValue(int64(timestamp.UTC().Day())))
	array.SetElement(NewStringRuntimeValue("wday"), NewIntegerRuntimeValue(int64(timestamp.UTC().Weekday())))
	array.SetElement(NewStringRuntimeValue("mon"), NewIntegerRuntimeValue(int64(timestamp.UTC().Month())))
	array.SetElement(NewStringRuntimeValue("year"), NewIntegerRuntimeValue(int64(timestamp.UTC().Year())))
	array.SetElement(NewStringRuntimeValue("yday"), NewIntegerRuntimeValue(int64(timestamp.UTC().YearDay()-1)))
	array.SetElement(NewStringRuntimeValue("weekday"), NewStringRuntimeValue(timestamp.UTC().Weekday().String()))
	array.SetElement(NewStringRuntimeValue("month"), NewStringRuntimeValue(timestamp.UTC().Month().String()))
	array.SetElement(NewIntegerRuntimeValue(0), NewIntegerRuntimeValue(timestamp.UTC().Unix()))

	return array, nil
}

// ------------------- MARK: localtime -------------------

func nativeFn_localtime(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("localtime").
		addParam("$timestamp", []string{"int"}, NewNullRuntimeValue()).
		addParam("associative", []string{"bool"}, NewBooleanRuntimeValue(false)).
		validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://www.php.net/manual/en/function.localtime.php

	// If timestamp is omitted, use the default value `time()`
	if args[0].GetType() == NullValue {
		args[0] = lib_time()
	}

	timestamp := time.Unix(args[0].(*IntegerRuntimeValue).Value, 0)
	array := NewArrayRuntimeValue()
	var isDst int64
	if timestamp.Local().IsDST() {
		isDst = 1
	}
	year := int64(timestamp.Local().Year()) - 1900

	if args[1].(*BooleanRuntimeValue).Value {
		// Associative array
		array.SetElement(NewStringRuntimeValue("tm_sec"), NewIntegerRuntimeValue(int64(timestamp.Local().Second())))
		array.SetElement(NewStringRuntimeValue("tm_min"), NewIntegerRuntimeValue(int64(timestamp.Local().Minute())))
		array.SetElement(NewStringRuntimeValue("tm_hour"), NewIntegerRuntimeValue(int64(timestamp.Local().Hour())))
		array.SetElement(NewStringRuntimeValue("tm_mday"), NewIntegerRuntimeValue(int64(timestamp.Local().Day())))
		array.SetElement(NewStringRuntimeValue("tm_mon"), NewIntegerRuntimeValue(int64(timestamp.Local().Month())))
		array.SetElement(NewStringRuntimeValue("tm_year"), NewIntegerRuntimeValue(year))
		array.SetElement(NewStringRuntimeValue("tm_wday"), NewIntegerRuntimeValue(int64(timestamp.Local().Weekday())))
		array.SetElement(NewStringRuntimeValue("tm_yday"), NewIntegerRuntimeValue(int64(timestamp.Local().YearDay()-1)))
		array.SetElement(NewStringRuntimeValue("tm_isdst"), NewIntegerRuntimeValue(isDst))
	} else {
		//Numerically index array
		array.SetElement(NewIntegerRuntimeValue(0), NewIntegerRuntimeValue(int64(timestamp.Local().Second())))
		array.SetElement(NewIntegerRuntimeValue(1), NewIntegerRuntimeValue(int64(timestamp.Local().Minute())))
		array.SetElement(NewIntegerRuntimeValue(2), NewIntegerRuntimeValue(int64(timestamp.Local().Hour())))
		array.SetElement(NewIntegerRuntimeValue(3), NewIntegerRuntimeValue(int64(timestamp.Local().Day())))
		array.SetElement(NewIntegerRuntimeValue(4), NewIntegerRuntimeValue(int64(timestamp.Local().Month())))
		array.SetElement(NewIntegerRuntimeValue(5), NewIntegerRuntimeValue(year))
		array.SetElement(NewIntegerRuntimeValue(6), NewIntegerRuntimeValue(int64(timestamp.Local().Weekday())))
		array.SetElement(NewIntegerRuntimeValue(7), NewIntegerRuntimeValue(int64(timestamp.Local().YearDay()-1)))
		array.SetElement(NewIntegerRuntimeValue(8), NewIntegerRuntimeValue(isDst))
	}

	return array, nil
}

// ------------------- MARK: microtime -------------------

func nativeFn_microtime(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("microtime").addParam("$as_float", []string{"bool"}, NewBooleanRuntimeValue(false)).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://www.php.net/manual/en/function.microtime.php

	now := time.Now()
	micro := float64(now.UnixMicro()) / math.Pow(10, 6)

	// As float
	if args[0].(*BooleanRuntimeValue).Value {
		return NewFloatingRuntimeValue(micro), nil
	}
	// As string
	return NewStringRuntimeValue(fmt.Sprintf("%f %d", micro-float64(now.Unix()), now.Unix())), nil
}

// ------------------- MARK: mktime -------------------

func nativeFn_mktime(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("mktime").
		addParam("$hour", []string{"int"}, nil).
		addParam("$minute", []string{"int"}, NewNullRuntimeValue()).
		addParam("$second", []string{"int"}, NewNullRuntimeValue()).
		addParam("$month", []string{"int"}, NewNullRuntimeValue()).
		addParam("$day", []string{"int"}, NewNullRuntimeValue()).
		addParam("$year", []string{"int"}, NewNullRuntimeValue()).
		validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://www.php.net/manual/en/function.mktime.php

	now := time.Now().Local()

	hour := int(args[0].(*IntegerRuntimeValue).Value)

	minute := now.Minute()
	if args[1].GetType() != NullValue {
		minute = int(args[1].(*IntegerRuntimeValue).Value)
	}

	second := now.Second()
	if args[2].GetType() != NullValue {
		second = int(args[2].(*IntegerRuntimeValue).Value)
	}

	month := now.Month()
	if args[3].GetType() != NullValue {
		month = time.Month(args[3].(*IntegerRuntimeValue).Value)
	}

	day := now.Day()
	if args[4].GetType() != NullValue {
		day = int(args[4].(*IntegerRuntimeValue).Value)
	}

	year := now.Year()
	if args[5].GetType() != NullValue {
		year = int(args[5].(*IntegerRuntimeValue).Value)
	}
	if year >= 0 && year <= 69 {
		year = 2000 + year
	}
	if year >= 70 && year <= 100 {
		year = 1900 + year
	}

	timestamp := time.Date(year, month, day, hour, minute, second, 0, time.Local)

	return NewIntegerRuntimeValue(timestamp.Unix()), nil
}

// ------------------- MARK: time -------------------

func nativeFn_time(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	_, err := NewFuncParamValidator("time").validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return lib_time(), nil
}

func lib_time() *IntegerRuntimeValue {
	// Spec: https://www.php.net/manual/en/function.time.php
	return NewIntegerRuntimeValue(time.Now().UTC().Unix())
}

// TODO date
// TODO date_add
// TODO date_create
// TODO date_create_from_format
// TODO date_create_immutable
// TODO date_create_immutable_from_format
// TODO date_date_set
// TODO date_default_timezone_get
// TODO date_default_timezone_set
// TODO date_diff
// TODO date_format
// TODO date_get_last_errors
// TODO date_interval_create_from_date_string
// TODO date_interval_format
// TODO date_isodate_set
// TODO date_modify
// TODO date_offset_get
// TODO date_parse
// TODO date_parse_from_format
// TODO date_sub
// TODO date_sun_info
// TODO date_sunrise
// TODO date_sunset
// TODO date_time_set
// TODO date_timestamp_get
// TODO date_timestamp_set
// TODO date_timezone_get
// TODO date_timezone_set
// TODO gettimeofday
// TODO gmdate
// TODO gmmktime
// TODO gmstrftime
// TODO idate
// TODO strftime - DEPRECATED
// TODO strptime - DEPRECATED
// TODO strtotime
// TODO timezone_abbreviations_list
// TODO timezone_identifiers_list
// TODO timezone_location_get
// TODO timezone_name_from_abbr
// TODO timezone_name_get
// TODO timezone_offset_get
// TODO timezone_open
// TODO timezone_transitions_get
// TODO timezone_version_get
