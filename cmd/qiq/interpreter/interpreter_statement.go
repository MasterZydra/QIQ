package interpreter

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime/stdlib/variableHandling"
	"QIQ/cmd/qiq/runtime/values"
)

// ProcessStmt implements Visitor.
func (interpreter *Interpreter) ProcessStmt(stmt *ast.Statement, _ any) (any, error) {
	panic("ProcessStmt should never be called")
}

// ProcessClassDeclarationStmt implements Visitor.
func (visitor *Interpreter) ProcessClassDeclarationStmt(stmt *ast.ClassDeclarationStatement, _ any) (any, error) {
	visitor.AddClass(stmt.GetPosition().File.GetNamespaceStr()+stmt.Name, stmt)
	return values.NewVoid(), nil
}

// ProcessConstDeclarationStmt implements Visitor.
func (interpreter *Interpreter) ProcessConstDeclarationStmt(stmt *ast.ConstDeclarationStatement, env any) (any, error) {
	value := must(interpreter.processStmt(stmt.Value, env))
	return env.(*Environment).declareConstant(stmt.Name, value)
}

// ProcessCompoundStmt implements Visitor.
func (interpreter *Interpreter) ProcessCompoundStmt(stmt *ast.CompoundStatement, env any) (any, error) {
	for _, statement := range stmt.Statements {
		must(interpreter.processStmt(statement, env))
	}
	return values.NewVoid(), nil
}

// ProcessEchoStmt implements Visitor.
func (interpreter *Interpreter) ProcessEchoStmt(stmt *ast.EchoStatement, env any) (any, error) {
	for _, expr := range stmt.Expressions {
		runtimeValue := must(interpreter.processStmt(expr, env))

		str := mustOrVoid(variableHandling.StrVal(runtimeValue))
		interpreter.Print(str)
	}
	return values.NewVoid(), nil
}

// ProcessExpressionStmt implements Visitor.
func (interpreter *Interpreter) ProcessExpressionStmt(stmt *ast.ExpressionStatement, env any) (any, error) {
	return interpreter.processStmt(stmt.Expr, env)
}

// ProcessFunctionCallExpr implements Visitor.
func (interpreter *Interpreter) ProcessFunctionDefinitionStmt(stmt *ast.FunctionDefinitionStatement, env any) (any, error) {
	// Check if this function definition was already processed before interpreting the code
	if interpreter.isCached(stmt) {
		return values.NewVoid(), nil
	}

	mustOrVoid(0, env.(*Environment).defineUserFunction(stmt))

	return interpreter.writeCache(stmt, values.NewVoid()), nil
}

// ProcessReturnStmt implements Visitor.
func (interpreter *Interpreter) ProcessReturnStmt(stmt *ast.ReturnStatement, env any) (any, error) {
	if stmt.Expr == nil {
		return values.NewVoid(), phpError.NewEvent(phpError.ReturnEvent)
	}
	runtimeValue := must(interpreter.processStmt(stmt.Expr, env))
	return runtimeValue, phpError.NewEvent(phpError.ReturnEvent)
}

// ProcessContinueStmt implements Visitor.
func (interpreter *Interpreter) ProcessContinueStmt(stmt *ast.ContinueStatement, env any) (any, error) {
	if stmt.Expr == nil {
		return values.NewVoid(), phpError.NewEvent(phpError.ReturnEvent)
	}
	runtimeValue := must(interpreter.processStmt(stmt.Expr, env))

	if runtimeValue.GetType() != values.IntValue {
		return runtimeValue, phpError.NewError("Breakout level must be an integer value. Got %s", runtimeValue.GetType())
	}

	return runtimeValue, phpError.NewContinueEvent(runtimeValue.(*values.Int).Value)
}

// ProcessBreakStmt implements Visitor.
func (interpreter *Interpreter) ProcessBreakStmt(stmt *ast.BreakStatement, env any) (any, error) {
	if stmt.Expr == nil {
		return values.NewVoid(), phpError.NewEvent(phpError.ReturnEvent)
	}
	runtimeValue := must(interpreter.processStmt(stmt.Expr, env))

	if runtimeValue.GetType() != values.IntValue {
		return runtimeValue, phpError.NewError("Breakout level must be an integer value. Got %s", runtimeValue.GetType())
	}

	return runtimeValue, phpError.NewBreakEvent(runtimeValue.(*values.Int).Value)
}

// ProcessForStmt implements Visitor.
func (interpreter *Interpreter) ProcessForStmt(stmt *ast.ForStatement, env any) (any, error) {
	// Spec: https://phplang.org/spec/11-statements.html#grammar-for-statement
	// If for-initializer is omitted, no action is taken at the start of the loop processing.
	if stmt.Initializer != nil {
		// Spec: https://phplang.org/spec/11-statements.html#grammar-for-statement
		// The group of expressions in for-initializer is evaluated once, left-to-right, for their side effects.
		for _, statement := range stmt.Initializer.Statements {
			mustOrVoid(interpreter.processStmt(statement, env))
		}
	}

	for {
		// Spec: https://phplang.org/spec/11-statements.html#grammar-for-statement
		// If for-control is omitted, this is treated as if for-control was an expression with the value TRUE.
		condition := true

		// Spec: https://phplang.org/spec/11-statements.html#grammar-for-statement
		// Then the group of expressions in for-control is evaluated left-to-right (with all but the right-most one for their side
		// effects only), with the right-most expression’s value being converted to type bool.
		if stmt.Control != nil {
			var conditionRuntimeValue values.RuntimeValue
			for _, statement := range stmt.Control.Statements {
				conditionRuntimeValue = mustOrVoid(interpreter.processStmt(statement, env))
			}
			condition = mustOrVoid(variableHandling.BoolVal(conditionRuntimeValue))
		}

		executeEndOfLoop := func() phpError.Error {
			// Spec: https://phplang.org/spec/11-statements.html#grammar-for-statement
			// If for-end-of-loop is omitted, no action is taken at the end of each iteration.
			if stmt.EndOfLoop != nil {
				for _, statement := range stmt.EndOfLoop.Statements {
					_, err := interpreter.processStmt(statement, env)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}

		// Spec: https://phplang.org/spec/11-statements.html#grammar-for-statement
		// If the result is TRUE, statement is executed, ...
		if condition {
			_, err := interpreter.processStmt(stmt.Block, env)
			if err != nil {
				if err.GetErrorType() == phpError.EventError && err.GetMessage() == "break" {
					breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
					if breakoutLevel == 1 {
						break
					}
					return values.NewVoid(), phpError.NewBreakEvent(breakoutLevel - 1)
				}
				if err.GetErrorType() == phpError.EventError && err.GetMessage() == "continue" {
					breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
					if breakoutLevel == 1 {
						// Execute end-of-loop logic
						mustOrVoid(0, executeEndOfLoop())
						continue
					}
					return values.NewVoid(), phpError.NewContinueEvent(breakoutLevel - 1)
				}
				return values.NewVoid(), err
			}
		}

		// Spec: https://phplang.org/spec/11-statements.html#grammar-for-statement
		// ... and the group of expressions in for-end-of-loop is evaluated left-to-right, for their side effects only.
		mustOrVoid(0, executeEndOfLoop())

		// Spec: https://phplang.org/spec/11-statements.html#grammar-for-statement
		// Once the right-most expression in for-control is FALSE, control transfers to the point immediately following the end of the for statement.
		if !condition {
			break
		}
	}

	return values.NewVoid(), nil
}

// ProcessIfStmt implements Visitor.
func (interpreter *Interpreter) ProcessIfStmt(stmt *ast.IfStatement, env any) (any, error) {
	conditionRuntimeValue := must(interpreter.processStmt(stmt.Condition, env))
	condition := mustOrVoid(variableHandling.BoolVal(conditionRuntimeValue))
	if condition {
		must(interpreter.processStmt(stmt.IfBlock, env))
		return values.NewVoid(), nil
	}

	if len(stmt.ElseIf) > 0 {
		for _, elseIf := range stmt.ElseIf {
			conditionRuntimeValue := must(interpreter.processStmt(elseIf.Condition, env))
			condition := mustOrVoid(variableHandling.BoolVal(conditionRuntimeValue))
			if !condition {
				continue
			}

			must(interpreter.processStmt(elseIf.IfBlock, env))
			return values.NewVoid(), nil
		}
	}

	if stmt.ElseBlock != nil {
		must(interpreter.processStmt(stmt.ElseBlock, env))
		return values.NewVoid(), nil
	}

	return values.NewVoid(), nil
}

// ProcessWhileStmt implements Visitor.
func (interpreter *Interpreter) ProcessWhileStmt(stmt *ast.WhileStatement, env any) (any, error) {
	for {
		conditionRuntimeValue := must(interpreter.processStmt(stmt.Condition, env))
		condition := mustOrVoid(variableHandling.BoolVal(conditionRuntimeValue))
		if !condition {
			break
		}

		runtimeValue, err := interpreter.processStmt(stmt.Block, env)
		if err != nil {
			if err.GetErrorType() == phpError.EventError && err.GetMessage() == "break" {
				breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
				if breakoutLevel == 1 {
					return values.NewVoid(), nil
				}
				return values.NewVoid(), phpError.NewBreakEvent(breakoutLevel - 1)
			}
			if err.GetErrorType() == phpError.EventError && err.GetMessage() == "continue" {
				breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
				if breakoutLevel == 1 {
					continue
				}
				return values.NewVoid(), phpError.NewContinueEvent(breakoutLevel - 1)
			}
			return runtimeValue, err
		}
	}
	return values.NewVoid(), nil
}

// ProcessDoStmt implements Visitor.
func (interpreter *Interpreter) ProcessDoStmt(stmt *ast.DoStatement, env any) (any, error) {
	var condition bool = true
	for condition {
		runtimeValue, err := interpreter.processStmt(stmt.Block, env)
		if err != nil {
			if err.GetErrorType() == phpError.EventError && err.GetMessage() == "break" {
				breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
				if breakoutLevel == 1 {
					return values.NewVoid(), nil
				}
				return values.NewVoid(), phpError.NewBreakEvent(breakoutLevel - 1)
			}
			if err.GetErrorType() == phpError.EventError && err.GetMessage() == "continue" {
				breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
				if breakoutLevel == 1 {
					continue
				}
				return values.NewVoid(), phpError.NewContinueEvent(breakoutLevel - 1)
			}
			return runtimeValue, err
		}

		conditionRuntimeValue := must(interpreter.processStmt(stmt.Condition, env))
		condition = mustOrVoid(variableHandling.BoolVal(conditionRuntimeValue))
		if !condition {
			break
		}
	}
	return values.NewVoid(), nil
}

// ProcessGlobalDeclarationStmt implements Visitor.
func (interpreter *Interpreter) ProcessGlobalDeclarationStmt(stmt *ast.GlobalDeclarationStatement, env any) (any, error) {
	for _, variable := range stmt.Variables {
		variableName, err := interpreter.varExprToVarName(variable, env.(*Environment))
		if err != nil {
			return values.NewVoid(), err
		}
		env.(*Environment).addGlobalVariable(variableName)
	}
	return values.NewVoid(), nil
}

// ProcessThrowStmt implements Visitor.
func (interpreter *Interpreter) ProcessThrowStmt(stmt *ast.ThrowStatement, env any) (any, error) {
	return values.NewVoid(), phpError.NewError("ProcessThrowStmt is not implemented")
}

// ProcessDeclareStmt implements Visitor.
func (interpreter *Interpreter) ProcessDeclareStmt(stmt *ast.DeclareStatement, env any) (any, error) {
	if stmt.Directive == "strict_types" {
		if stmt.Literal.(*ast.IntegerLiteralExpression).Value == 1 {
			stmt.GetPosition().File.IsStrictType = true
		} else {
			stmt.GetPosition().File.IsStrictType = false
		}
		return values.NewVoid(), nil
	}

	return values.NewVoid(), phpError.NewError("ProcessDeclareStmt: Directive '%s' is not implemented", stmt.Directive)
}

// ProcessForeachStmt implements Visitor.
func (interpreter *Interpreter) ProcessForeachStmt(stmt *ast.ForeachStatement, env any) (any, error) {
	runtimeValue, err := interpreter.processStmt(stmt.Collection, env)
	if err != nil {
		return values.NewVoid(), err
	}

	environment := env.(*Environment)

	// Array
	if runtimeValue.GetType() == values.ArrayValue {
		runtimeArray := runtimeValue.(*values.Array)
		for _, keyValue := range runtimeArray.Keys {
			// Set key and value variable
			if stmt.Key != nil {
				keyName := mustOrVoid(interpreter.varExprToVarName(stmt.Key, environment))
				environment.declareVariable(keyName, keyValue)
			}
			value, _ := runtimeArray.GetElement(keyValue)
			valueName := mustOrVoid(interpreter.varExprToVarName(stmt.Value, environment))
			environment.declareVariable(valueName, value)

			// Execute body
			runtimeValue, err := interpreter.processStmt(stmt.Block, env)
			if err != nil {
				if err.GetErrorType() == phpError.EventError && err.GetMessage() == "break" {
					breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
					if breakoutLevel == 1 {
						return values.NewVoid(), nil
					}
					return values.NewVoid(), phpError.NewBreakEvent(breakoutLevel - 1)
				}
				if err.GetErrorType() == phpError.EventError && err.GetMessage() == "continue" {
					breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
					if breakoutLevel == 1 {
						continue
					}
					return values.NewVoid(), phpError.NewContinueEvent(breakoutLevel - 1)
				}
				return runtimeValue, err
			}
		}
		return values.NewVoid(), nil
	}

	// Object
	if runtimeValue.GetType() == values.ObjectValue {
		runtimeObject := runtimeValue.(*values.Object)

		for _, propertyName := range runtimeObject.PropertyNames {
			if runtimeObject.Class.Properties[propertyName].Visibility != "public" {
				continue
			}

			// Set key and value variable
			if stmt.Key != nil {
				keyName := mustOrVoid(interpreter.varExprToVarName(stmt.Key, environment))
				environment.declareVariable(keyName, values.NewStr(propertyName[1:]))
			}
			valueName := mustOrVoid(interpreter.varExprToVarName(stmt.Value, environment))
			environment.declareVariable(valueName, runtimeObject.Properties[propertyName])

			// Execute body
			runtimeValue, err := interpreter.processStmt(stmt.Block, env)
			if err != nil {
				if err.GetErrorType() == phpError.EventError && err.GetMessage() == "break" {
					breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
					if breakoutLevel == 1 {
						return values.NewVoid(), nil
					}
					return values.NewVoid(), phpError.NewBreakEvent(breakoutLevel - 1)
				}
				if err.GetErrorType() == phpError.EventError && err.GetMessage() == "continue" {
					breakoutLevel := err.(*phpError.ContinueEventError).GetBreakoutLevel()
					if breakoutLevel == 1 {
						continue
					}
					return values.NewVoid(), phpError.NewContinueEvent(breakoutLevel - 1)
				}
				return runtimeValue, err
			}
		}
		return values.NewVoid(), nil
	}

	givenType := values.ToPhpType(runtimeValue)
	if runtimeValue.GetType() == values.BoolValue {
		if runtimeValue.(*values.Bool).Value {
			givenType = "true"
		} else {
			givenType = "false"
		}
	}
	if runtimeValue.GetType() == values.NullValue {
		givenType = "null"
	}
	interpreter.PrintError(phpError.NewWarning("foreach() argument must be of type array|object, %s given in %s", givenType, stmt.Collection.GetPosString()))

	return values.NewVoid(), nil
}
