package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/phpError"
)

// ProcessStmt implements Visitor.
func (interpreter *Interpreter) ProcessStmt(stmt *ast.Statement, _ any) (any, error) {
	panic("ProcessStmt should never be called")
}

// ProcessConstDeclarationStmt implements Visitor.
func (interpreter *Interpreter) ProcessConstDeclarationStmt(stmt *ast.ConstDeclarationStatement, env any) (any, error) {
	value, err := interpreter.processStmt(stmt.Value, env)
	if err != nil {
		return value, err
	}
	return env.(*Environment).declareConstant(stmt.Name, value)
}

// ProcessCompoundStmt implements Visitor.
func (interpreter *Interpreter) ProcessCompoundStmt(stmt *ast.CompoundStatement, env any) (any, error) {
	for _, statement := range stmt.Statements {
		runtimeValue, err := interpreter.processStmt(statement, env)
		if err != nil {
			return runtimeValue, err
		}
	}
	return NewVoidRuntimeValue(), nil
}

// ProcessEchoStmt implements Visitor.
func (interpreter *Interpreter) ProcessEchoStmt(stmt *ast.EchoStatement, env any) (any, error) {
	for _, expr := range stmt.Expressions {
		if runtimeValue, err := interpreter.processStmt(expr, env); err != nil {
			return runtimeValue, err
		} else {
			var str string
			str, err = lib_strval(runtimeValue)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}
			interpreter.print(str)
		}
	}
	return NewVoidRuntimeValue(), nil
}

// ProcessExpressionStmt implements Visitor.
func (interpreter *Interpreter) ProcessExpressionStmt(stmt *ast.ExpressionStatement, env any) (any, error) {
	return interpreter.processStmt(stmt.Expr, env)
}

// ProcessFunctionCallExpr implements Visitor.
func (interpreter *Interpreter) ProcessFunctionDefinitionStmt(stmt *ast.FunctionDefinitionStatement, env any) (any, error) {
	// Check if this function definition was already processed before interpreting the code
	if interpreter.isCached(stmt) {
		return NewVoidRuntimeValue(), nil
	}

	if err := env.(*Environment).defineUserFunction(stmt); err != nil {
		return NewVoidRuntimeValue(), err
	}

	return interpreter.writeCache(stmt, NewVoidRuntimeValue()), nil
}

// ProcessReturnStmt implements Visitor.
func (interpreter *Interpreter) ProcessReturnStmt(stmt *ast.ReturnStatement, env any) (any, error) {
	if stmt.Expr == nil {
		return NewVoidRuntimeValue(), phpError.NewEvent(phpError.ReturnEvent)
	}
	runtimeValue, err := interpreter.processStmt(stmt.Expr, env)
	if err != nil {
		return runtimeValue, err
	}
	return runtimeValue, phpError.NewEvent(phpError.ReturnEvent)
}

// ProcessContinueStmt implements Visitor.
func (interpreter *Interpreter) ProcessContinueStmt(stmt *ast.ContinueStatement, env any) (any, error) {
	if stmt.Expr == nil {
		return NewVoidRuntimeValue(), phpError.NewEvent(phpError.ReturnEvent)
	}
	runtimeValue, err := interpreter.processStmt(stmt.Expr, env)
	if err != nil {
		return runtimeValue, err
	}

	if runtimeValue.GetType() != IntegerValue {
		return runtimeValue, phpError.NewError("Breakout level must be an integer value. Got %s", runtimeValue.GetType())
	}

	return runtimeValue, phpError.NewContinueEvent(runtimeValue.(*IntegerRuntimeValue).Value)
}

// ProcessBreakStmt implements Visitor.
func (interpreter *Interpreter) ProcessBreakStmt(stmt *ast.BreakStatement, env any) (any, error) {
	if stmt.Expr == nil {
		return NewVoidRuntimeValue(), phpError.NewEvent(phpError.ReturnEvent)
	}
	runtimeValue, err := interpreter.processStmt(stmt.Expr, env)
	if err != nil {
		return runtimeValue, err
	}

	if runtimeValue.GetType() != IntegerValue {
		return runtimeValue, phpError.NewError("Breakout level must be an integer value. Got %s", runtimeValue.GetType())
	}

	return runtimeValue, phpError.NewBreakEvent(runtimeValue.(*IntegerRuntimeValue).Value)
}

// ProcessIfStmt implements Visitor.
func (interpreter *Interpreter) ProcessIfStmt(stmt *ast.IfStatement, env any) (any, error) {
	conditionRuntimeValue, err := interpreter.processStmt(stmt.Condition, env)
	if err != nil {
		return conditionRuntimeValue, err
	}

	condition, err := lib_boolval(conditionRuntimeValue)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	if condition {
		runtimeValue, err := interpreter.processStmt(stmt.IfBlock, env)
		if err != nil {
			return runtimeValue, err
		}
		return NewVoidRuntimeValue(), nil
	}

	if len(stmt.ElseIf) > 0 {
		for _, elseIf := range stmt.ElseIf {
			conditionRuntimeValue, err := interpreter.processStmt(elseIf.Condition, env)
			if err != nil {
				return conditionRuntimeValue, err
			}

			condition, err := lib_boolval(conditionRuntimeValue)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}

			if !condition {
				continue
			}

			runtimeValue, err := interpreter.processStmt(elseIf.IfBlock, env)
			if err != nil {
				return runtimeValue, err
			}
			return NewVoidRuntimeValue(), nil
		}
	}

	if stmt.ElseBlock != nil {
		runtimeValue, err := interpreter.processStmt(stmt.ElseBlock, env)
		if err != nil {
			return runtimeValue, err
		}
		return NewVoidRuntimeValue(), nil
	}

	return NewVoidRuntimeValue(), nil
}

// ProcessWhileStmt implements Visitor.
func (interpreter *Interpreter) ProcessWhileStmt(stmt *ast.WhileStatement, env any) (any, error) {
	var condition bool = true
	for condition {
		conditionRuntimeValue, err := interpreter.processStmt(stmt.Condition, env)
		if err != nil {
			return conditionRuntimeValue, err
		}

		condition, err = lib_boolval(conditionRuntimeValue)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}

		if condition {
			runtimeValue, err := interpreter.processStmt(stmt.IfBlock, env)
			if err != nil {
				if err.GetErrorType() == phpError.EventError && err.GetMessage() == "break" {
					breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
					if breakoutLevel == 1 {
						return NewVoidRuntimeValue(), nil
					}
					return NewVoidRuntimeValue(), phpError.NewBreakEvent(breakoutLevel - 1)
				}
				if err.GetErrorType() == phpError.EventError && err.GetMessage() == "continue" {
					breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
					if breakoutLevel == 1 {
						continue
					}
					return NewVoidRuntimeValue(), phpError.NewContinueEvent(breakoutLevel - 1)
				}
				return runtimeValue, err
			}
		}
	}
	return NewVoidRuntimeValue(), nil
}

// ProcessDoStmt implements Visitor.
func (interpreter *Interpreter) ProcessDoStmt(stmt *ast.DoStatement, env any) (any, error) {
	var condition bool = true
	for condition {
		runtimeValue, err := interpreter.processStmt(stmt.IfBlock, env)
		if err != nil {
			if err.GetErrorType() == phpError.EventError && err.GetMessage() == "break" {
				breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
				if breakoutLevel == 1 {
					return NewVoidRuntimeValue(), nil
				}
				return NewVoidRuntimeValue(), phpError.NewBreakEvent(breakoutLevel - 1)
			}
			if err.GetErrorType() == phpError.EventError && err.GetMessage() == "continue" {
				breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
				if breakoutLevel == 1 {
					continue
				}
				return NewVoidRuntimeValue(), phpError.NewContinueEvent(breakoutLevel - 1)
			}
			return runtimeValue, err
		}

		conditionRuntimeValue, err := interpreter.processStmt(stmt.Condition, env)
		if err != nil {
			return conditionRuntimeValue, err
		}

		condition, err = lib_boolval(conditionRuntimeValue)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		if !condition {
			break
		}
	}
	return NewVoidRuntimeValue(), nil
}
