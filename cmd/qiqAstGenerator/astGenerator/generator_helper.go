package astGenerator

import (
	"QIQ/cmd/qiq/ast"
	"fmt"
)

func (generator *AstGenerator) println(format string, a ...any) { generator.print(format+"\n", a...) }

func (generator *AstGenerator) print(format string, a ...any) {
	generator.output += fmt.Sprintf(format, a...)
}

func toGoVarName(name string) string {
	result := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			result += string(r)
		}
	}
	return result
}

func toStringSlice(s []string) string {
	result := "[]string{"
	for index, str := range s {
		if index > 0 {
			result += ", "
		}
		result += `"` + str + `"`
	}
	result += "}"
	return result
}

func funcParamArrayToStr(params []ast.FunctionParameter) string {
	result := "[]ast.FunctionParameter{"
	for index, param := range params {
		if index > 0 {
			result += ", "
		}
		result += fmt.Sprintf(`{Name: "%s", Type: %s}`, param.Name, toStringSlice(param.Type))
	}
	result += "}"
	return result
}

func toBoolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
