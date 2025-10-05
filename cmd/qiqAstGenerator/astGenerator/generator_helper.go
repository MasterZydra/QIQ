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
		result += fmt.Sprintf(`{Name: "%s", Type: %s`, param.Name, toStringSlice(param.Type))
		if param.DefaultValue != nil {
			result += ", DefaultValue: " + basicTypesToStr(param.DefaultValue)
		}
		result += "}"
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

func basicTypesToStr(expr ast.IExpression) string {
	switch expr.GetKind() {
	case ast.StringLiteralExpr:
		str := expr.(*ast.StringLiteralExpression)
		strType := ""
		switch str.StringType {
		case ast.SingleQuotedString:
			strType = "ast.SingleQuotedString"
		case ast.DoubleQuotedString:
			strType = "ast.DoubleQuotedString"
		case ast.HeredocString:
			strType = "ast.HeredocString"
		}
		return fmt.Sprintf(`ast.NewStringLiteralExpr(0, nil, "%s", %s)`, str.Value, strType)
	case ast.IntegerLiteralExpr:
		intExpr := expr.(*ast.IntegerLiteralExpression)
		return fmt.Sprintf(`ast.NewIntegerLiteralExpr(0, nil, %d)`, intExpr.Value)
	case ast.ConstantAccessExpr:
		constant := expr.(*ast.ConstantAccessExpression)
		return fmt.Sprintf(`ast.NewConstantAccessExpr(0, nil, "%s")`, constant.ConstantName)
	default:
		panic("basicTypesToStr: Unsupported type " + ast.ToString(expr))
	}
}
