package dateTime

import (
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime"
	"QIQ/cmd/qiq/runtime/funcParamValidator"
	"QIQ/cmd/qiq/runtime/values"
	"fmt"
	"math"
	"strings"
	"time"
)

func Register(environment runtime.Environment) {
	// Category: Date/Time Functions
	environment.AddNativeFunction("checkdate", nativeFn_checkdate)
	environment.AddNativeFunction("date", nativeFn_date)
	environment.AddNativeFunction("getdate", nativeFn_getdate)
	environment.AddNativeFunction("localtime", nativeFn_localtime)
	environment.AddNativeFunction("microtime", nativeFn_microtime)
	environment.AddNativeFunction("mktime", nativeFn_mktime)
	environment.AddNativeFunction("time", nativeFn_time)
}

// -------------------------------------- checkdate -------------------------------------- MARK: checkdate

func nativeFn_checkdate(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("checkdate").
		AddParam("$month", []string{"int"}, nil).AddParam("$day", []string{"int"}, nil).AddParam("$year", []string{"int"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.checkdate.php
	year := args[2].(*values.Int).Value
	month := args[0].(*values.Int).Value
	day := args[1].(*values.Int).Value

	// Spec: https://www.php.net/manual/en/function.checkdate.php
	// The year is between 1 and 32767 inclusive.
	if year < 1 || year > 32767 {
		return values.NewBool(false), nil
	}

	// Spec: https://www.php.net/manual/en/function.checkdate.php
	// The month is between 1 and 12 inclusive.
	if month < 1 || month > 12 {
		return values.NewBool(false), nil
	}

	// Spec: https://www.php.net/manual/en/function.checkdate.php
	// The day is within the allowed number of days for the given month. Leap years are taken into consideration.
	return values.NewBool(day >= 1 && day <= int64(DaysIn(time.Month(month), int(year)))), nil
}

// -------------------------------------- date -------------------------------------- MARK: date

func nativeFn_date(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("date").
		AddParam("$format", []string{"string"}, nil).AddParam("$timestamp", []string{"int"}, values.NewNull()).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	timestamp := time.Now()
	if args[1].GetType() != values.NullValue {
		timestamp = time.Unix(args[1].(*values.Int).Value, 0)
	}

	format := args[0].(*values.Str).Value

	// Spec: https://www.php.net/manual/en/datetime.format.php
	// Day
	if strings.Contains(format, "d") {
		// d 	Day of the month, 2 digits with leading zeros 	01 to 31
		format = strings.ReplaceAll(format, "d", fmt.Sprintf("%02d", timestamp.Day()))
	}
	if strings.Contains(format, "j") {
		// j 	Day of the month without leading zeros 	1 to 31
		format = strings.ReplaceAll(format, "j", fmt.Sprintf("%d", timestamp.Day()))
	}
	if strings.Contains(format, "z") {
		// z 	The day of the year (starting from 0) 	0 through 365
		format = strings.ReplaceAll(format, "z", fmt.Sprintf("%d", timestamp.YearDay()-1))
	}
	if strings.Contains(format, "w") {
		// w 	Numeric representation of the day of the week 	0 (for Sunday) through 6 (for Saturday)
		format = strings.ReplaceAll(format, "w", fmt.Sprintf("%d", timestamp.Weekday()))
	}
	if strings.Contains(format, "N") {
		// N 	ISO 8601 numeric representation of the day of the week 	1 (for Monday) through 7 (for Sunday)
		format = strings.ReplaceAll(format, "N", fmt.Sprintf("%d", Iso8601Weekday(timestamp.Weekday())))
	}

	// Spec: https://www.php.net/manual/en/datetime.format.php
	// Week
	if strings.Contains(format, "W") {
		// W 	ISO 8601 week number of year, weeks starting on Monday 	Example: 42 (the 42nd week in the year)
		_, week := timestamp.ISOWeek()
		format = strings.ReplaceAll(format, "W", fmt.Sprintf("%d", week))
	}

	// Spec: https://www.php.net/manual/en/datetime.format.php
	// Month
	if strings.Contains(format, "m") {
		// m 	Numeric representation of a month, with leading zeros 	01 through 12
		format = strings.ReplaceAll(format, "m", fmt.Sprintf("%02d", timestamp.Month()))
	}
	if strings.Contains(format, "n") {
		// n 	Numeric representation of a month, without leading zeros 	1 through 12
		format = strings.ReplaceAll(format, "n", fmt.Sprintf("%d", timestamp.Month()))
	}
	if strings.Contains(format, "t") {
		// t 	Number of days in the given month 	28 through 31
		format = strings.ReplaceAll(format, "t", fmt.Sprintf("%d", DaysIn(timestamp.Month(), timestamp.Year())))
	}

	// Spec: https://www.php.net/manual/en/datetime.format.php
	// Year
	if strings.Contains(format, "L") {
		// L 	Whether it's a leap year 	1 if it is a leap year, 0 otherwise.
		leap := "0"
		if IsLeapYear(timestamp.Year()) {
			leap = "1"
		}
		format = strings.ReplaceAll(format, "L", leap)
	}
	if strings.Contains(format, "Y") {
		// Y 	A full numeric representation of a year, at least 4 digits, with - for years BCE. 	Examples: -0055, 0787, 1999, 2003, 10191
		format = strings.ReplaceAll(format, "Y", fmt.Sprintf("%d", timestamp.Year()))
	}
	if strings.Contains(format, "y") {
		// y 	A two digit representation of a year 	Examples: 99 or 03
		format = strings.ReplaceAll(format, "y", fmt.Sprintf("%02d", timestamp.Year()%100))
	}

	// Spec: https://www.php.net/manual/en/datetime.format.php
	// Time
	if strings.Contains(format, "i") {
		// i 	Minutes with leading zeros 	00 to 59
		format = strings.ReplaceAll(format, "i", fmt.Sprintf("%02d", timestamp.Minute()))
	}
	if strings.Contains(format, "s") {
		// s 	Seconds with leading zeros 	00 through 59
		format = strings.ReplaceAll(format, "s", fmt.Sprintf("%02d", timestamp.Second()))
	}
	if strings.Contains(format, "G") {
		// G 	24-hour format of an hour without leading zeros 	0 through 23
		format = strings.ReplaceAll(format, "G", fmt.Sprintf("%d", timestamp.Hour()))
	}
	if strings.Contains(format, "g") {
		// g 	12-hour format of an hour without leading zeros 	1 through 12
		format = strings.ReplaceAll(format, "g", timestamp.Format("3"))
	}
	if strings.Contains(format, "h") {
		// h 	12-hour format of an hour with leading zeros 	01 through 12
		format = strings.ReplaceAll(format, "h", timestamp.Format("03"))
	}
	if strings.Contains(format, "H") {
		// H 	24-hour format of an hour with leading zeros 	00 through 23
		format = strings.ReplaceAll(format, "H", fmt.Sprintf("%02d", timestamp.Hour()))
	}

	return values.NewStr(format), nil

	// TODO date() missing formats
	/*
		Day 	--- 	---
		D 	A textual representation of a day, three letters 	Mon through Sun
		l (lowercase 'L') 	A full textual representation of the day of the week 	Sunday through Saturday
		S 	English ordinal suffix for the day of the month, 2 characters 	st, nd, rd or th. Works well with j

		Month 	--- 	---
		F 	A full textual representation of a month, such as January or March 	January through December
		M 	A short textual representation of a month, three letters 	Jan through Dec

		Year 	--- 	---
		o 	ISO 8601 week-numbering year. This has the same value as Y, except that if the ISO week number (W) belongs to the previous or next year, that year is used instead. 	Examples: 1999 or 2003
		X 	An expanded full numeric representation of a year, at least 4 digits, with - for years BCE, and + for years CE. 	Examples: -0055, +0787, +1999, +10191
		x 	An expanded full numeric representation if required, or a standard full numeral representation if possible (like Y). At least four digits. Years BCE are prefixed with a -. Years beyond (and including) 10000 are prefixed by a +. 	Examples: -0055, 0787, 1999, +10191

		Time 	--- 	---
		a 	Lowercase Ante meridiem and Post meridiem 	am or pm
		A 	Uppercase Ante meridiem and Post meridiem 	AM or PM
		B 	Swatch Internet time 	000 through 999
		u 	Microseconds. Note that date() will always generate 000000 since it takes an int parameter, whereas DateTimeInterface::format() does support microseconds if an object of type DateTimeInterface was created with microseconds. 	Example: 654321
		v 	Milliseconds. Same note applies as for u. 	Example: 654
		Timezone 	--- 	---
		e 	Timezone identifier 	Examples: UTC, GMT, Atlantic/Azores
		I (capital i) 	Whether or not the date is in daylight saving time 	1 if Daylight Saving Time, 0 otherwise.
		O 	Difference to Greenwich time (GMT) without colon between hours and minutes 	Example: +0200
		P 	Difference to Greenwich time (GMT) with colon between hours and minutes 	Example: +02:00
		p 	The same as P, but returns Z instead of +00:00 (available as of PHP 8.0.0) 	Examples: Z or +02:00
		T 	Timezone abbreviation, if known; otherwise the GMT offset. 	Examples: EST, MDT, +05
		Z 	Timezone offset in seconds. The offset for timezones west of UTC is always negative, and for those east of UTC is always positive. 	-43200 through 50400
		Full Date/Time 	--- 	---
		c 	ISO 8601 date 	2004-02-12T15:19:21+00:00
		r 	» RFC 2822/» RFC 5322 formatted date 	Example: Thu, 21 Dec 2000 16:01:07 +0200
		U 	Seconds since the Unix Epoch (January 1 1970 00:00:00 GMT) 	See also time()
	*/
}

// -------------------------------------- getdate -------------------------------------- MARK: getdate

func nativeFn_getdate(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("getdate").AddParam("$timestamp", []string{"int"}, values.NewNull()).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.getdate.php

	// If timestamp is omitted, use the default value `time()`
	if args[0].GetType() == values.NullValue {
		args[0] = lib_time()
	}

	timestamp := time.Unix(args[0].(*values.Int).Value, 0)
	array := values.NewArray()
	array.SetElement(values.NewStr("seconds"), values.NewInt(int64(timestamp.UTC().Second())))
	array.SetElement(values.NewStr("minutes"), values.NewInt(int64(timestamp.UTC().Minute())))
	array.SetElement(values.NewStr("hours"), values.NewInt(int64(timestamp.UTC().Hour())))
	array.SetElement(values.NewStr("mday"), values.NewInt(int64(timestamp.UTC().Day())))
	array.SetElement(values.NewStr("wday"), values.NewInt(int64(timestamp.UTC().Weekday())))
	array.SetElement(values.NewStr("mon"), values.NewInt(int64(timestamp.UTC().Month())))
	array.SetElement(values.NewStr("year"), values.NewInt(int64(timestamp.UTC().Year())))
	array.SetElement(values.NewStr("yday"), values.NewInt(int64(timestamp.UTC().YearDay()-1)))
	array.SetElement(values.NewStr("weekday"), values.NewStr(timestamp.UTC().Weekday().String()))
	array.SetElement(values.NewStr("month"), values.NewStr(timestamp.UTC().Month().String()))
	array.SetElement(nil, values.NewInt(timestamp.UTC().Unix()))

	return array, nil
}

// -------------------------------------- localtime -------------------------------------- MARK: localtime

func nativeFn_localtime(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("localtime").
		AddParam("$timestamp", []string{"int"}, values.NewNull()).
		AddParam("associative", []string{"bool"}, values.NewBool(false)).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.localtime.php

	// If timestamp is omitted, use the default value `time()`
	if args[0].GetType() == values.NullValue {
		args[0] = lib_time()
	}

	timestamp := time.Unix(args[0].(*values.Int).Value, 0)
	array := values.NewArray()
	var isDst int64
	if timestamp.Local().IsDST() {
		isDst = 1
	}
	year := int64(timestamp.Local().Year()) - 1900

	if args[1].(*values.Bool).Value {
		// Associative array
		array.SetElement(values.NewStr("tm_sec"), values.NewInt(int64(timestamp.Local().Second())))
		array.SetElement(values.NewStr("tm_min"), values.NewInt(int64(timestamp.Local().Minute())))
		array.SetElement(values.NewStr("tm_hour"), values.NewInt(int64(timestamp.Local().Hour())))
		array.SetElement(values.NewStr("tm_mday"), values.NewInt(int64(timestamp.Local().Day())))
		array.SetElement(values.NewStr("tm_mon"), values.NewInt(int64(timestamp.Local().Month())))
		array.SetElement(values.NewStr("tm_year"), values.NewInt(year))
		array.SetElement(values.NewStr("tm_wday"), values.NewInt(int64(timestamp.Local().Weekday())))
		array.SetElement(values.NewStr("tm_yday"), values.NewInt(int64(timestamp.Local().YearDay()-1)))
		array.SetElement(values.NewStr("tm_isdst"), values.NewInt(isDst))
	} else {
		//Numerically index array
		array.SetElement(nil, values.NewInt(int64(timestamp.Local().Second())))
		array.SetElement(nil, values.NewInt(int64(timestamp.Local().Minute())))
		array.SetElement(nil, values.NewInt(int64(timestamp.Local().Hour())))
		array.SetElement(nil, values.NewInt(int64(timestamp.Local().Day())))
		array.SetElement(nil, values.NewInt(int64(timestamp.Local().Month())))
		array.SetElement(nil, values.NewInt(year))
		array.SetElement(nil, values.NewInt(int64(timestamp.Local().Weekday())))
		array.SetElement(nil, values.NewInt(int64(timestamp.Local().YearDay()-1)))
		array.SetElement(nil, values.NewInt(isDst))
	}

	return array, nil
}

// -------------------------------------- microtime -------------------------------------- MARK: microtime

func nativeFn_microtime(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("microtime").AddParam("$as_float", []string{"bool"}, values.NewBool(false)).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.microtime.php

	now := time.Now()
	micro := float64(now.UnixMicro()) / math.Pow(10, 6)

	// As float
	if args[0].(*values.Bool).Value {
		return values.NewFloat(micro), nil
	}
	// As string
	return values.NewStr(fmt.Sprintf("%f %d", micro-float64(now.Unix()), now.Unix())), nil
}

// -------------------------------------- mktime -------------------------------------- MARK: mktime

func nativeFn_mktime(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("mktime").
		AddParam("$hour", []string{"int"}, nil).
		AddParam("$minute", []string{"int"}, values.NewNull()).
		AddParam("$second", []string{"int"}, values.NewNull()).
		AddParam("$month", []string{"int"}, values.NewNull()).
		AddParam("$day", []string{"int"}, values.NewNull()).
		AddParam("$year", []string{"int"}, values.NewNull()).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.mktime.php

	now := time.Now().Local()

	hour := int(args[0].(*values.Int).Value)

	minute := now.Minute()
	if args[1].GetType() != values.NullValue {
		minute = int(args[1].(*values.Int).Value)
	}

	second := now.Second()
	if args[2].GetType() != values.NullValue {
		second = int(args[2].(*values.Int).Value)
	}

	month := now.Month()
	if args[3].GetType() != values.NullValue {
		month = time.Month(args[3].(*values.Int).Value)
	}

	day := now.Day()
	if args[4].GetType() != values.NullValue {
		day = int(args[4].(*values.Int).Value)
	}

	year := now.Year()
	if args[5].GetType() != values.NullValue {
		year = int(args[5].(*values.Int).Value)
	}
	if year >= 0 && year <= 69 {
		year = 2000 + year
	}
	if year >= 70 && year <= 100 {
		year = 1900 + year
	}

	timestamp := time.Date(year, month, day, hour, minute, second, 0, time.Local)

	return values.NewInt(timestamp.Unix()), nil
}

// -------------------------------------- time -------------------------------------- MARK: time

func nativeFn_time(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	_, err := funcParamValidator.NewValidator("time").Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return lib_time(), nil
}

func lib_time() *values.Int {
	// Spec: https://www.php.net/manual/en/function.time.php
	return values.NewInt(time.Now().UTC().Unix())
}

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
