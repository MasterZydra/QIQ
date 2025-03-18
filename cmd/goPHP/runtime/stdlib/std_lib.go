package stdlib

import (
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/stdlib/array"
	"GoPHP/cmd/goPHP/runtime/stdlib/dateTime"
	"GoPHP/cmd/goPHP/runtime/stdlib/errorHandling"
	"GoPHP/cmd/goPHP/runtime/stdlib/filesystem"
	"GoPHP/cmd/goPHP/runtime/stdlib/math"
	"GoPHP/cmd/goPHP/runtime/stdlib/misc"
	"GoPHP/cmd/goPHP/runtime/stdlib/optionsInfo"
	"GoPHP/cmd/goPHP/runtime/stdlib/outputControl"
	"GoPHP/cmd/goPHP/runtime/stdlib/strings"
	"GoPHP/cmd/goPHP/runtime/stdlib/variableHandling"
)

func Register(environment runtime.Environment) {
	array.Register(environment)
	dateTime.Register(environment)
	errorHandling.Register(environment)
	filesystem.Register(environment)
	math.Register(environment)
	misc.Register(environment)
	optionsInfo.Register(environment)
	outputControl.Register(environment)
	strings.Register(environment)
	variableHandling.Register(environment)
}
