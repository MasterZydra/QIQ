package interpreter

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/values"
)

func registerNativeOutputControlFunctions(environment *Environment) {
	environment.nativeFunctions["ob_clean"] = nativeFn_ob_clean
	environment.nativeFunctions["ob_flush"] = nativeFn_ob_flush
	environment.nativeFunctions["ob_end_clean"] = nativeFn_ob_end_clean
	environment.nativeFunctions["ob_end_flush"] = nativeFn_ob_end_flush
	environment.nativeFunctions["ob_get_clean"] = nativeFn_ob_get_clean
	environment.nativeFunctions["ob_get_flush"] = nativeFn_ob_get_flush
	environment.nativeFunctions["ob_get_contents"] = nativeFn_ob_get_contents
	environment.nativeFunctions["ob_get_level"] = nativeFn_ob_get_level
	environment.nativeFunctions["ob_start"] = nativeFn_ob_start
}

// ------------------- MARK: ob_clean -------------------

func nativeFn_ob_clean(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ob-clean.php

	_, err := NewFuncParamValidator("ob_clean").validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// TODO Call output handler
	// Spec: https://www.php.net/manual/en/function.ob-clean.php
	// This function calls the output handler (with the PHP_OUTPUT_HANDLER_CLEAN flag),
	// discards it's return value and cleans (erases) the contents of the active output buffer.

	// TODO Throw notice if no buffer: e.g. Notice: ob_clean(): Failed to delete buffer. No buffer to delete in /home/user/scripts/code.php on line 5

	if len(interpreter.outputBuffers) == 0 {
		return values.NewBool(false), nil
	}

	interpreter.outputBuffers[len(interpreter.outputBuffers)-1].Content = ""
	return values.NewBool(true), nil
}

// ------------------- MARK: ob_flush -------------------

func nativeFn_ob_flush(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ob-flush.php

	_, err := NewFuncParamValidator("ob_flush").validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// TODO Call output handler
	// Spec: https://www.php.net/manual/en/function.ob-flush.php
	// This function calls the output handler (with the PHP_OUTPUT_HANDLER_FLUSH flag),
	// discards it's return value and flushs (erases) the contents of the active output buffer.

	// TODO Throw notice if no buffer: e.g. Notice: ob_flush(): Failed to delete buffer. No buffer to delete in /home/user/scripts/code.php on line 5

	if len(interpreter.outputBuffers) == 0 {
		return values.NewBool(false), nil
	}

	if len(interpreter.outputBuffers) == 1 {
		interpreter.result += interpreter.outputBuffers[len(interpreter.outputBuffers)-1].Content
	} else {
		interpreter.outputBuffers[len(interpreter.outputBuffers)-2].Content += interpreter.outputBuffers[len(interpreter.outputBuffers)-1].Content
	}

	nativeFn_ob_clean(args, interpreter)
	return values.NewBool(true), nil
}

// ------------------- MARK: ob_end_clean -------------------

func nativeFn_ob_end_clean(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ob-end-clean.php

	_, err := NewFuncParamValidator("ob_end_clean").validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// TODO Call output handler
	// Spec: https://www.php.net/manual/en/function.ob-end-clean.php
	// This function calls the output handler (with the PHP_OUTPUT_HANDLER_CLEAN and PHP_OUTPUT_HANDLER_FINAL flags),
	// discards it's return value, discards the contents of the active output buffer and turns off the active output buffer.

	// TODO Throw notice if no buffer: e.g. Notice: ob_end_clean(): Failed to delete buffer. No buffer to delete in /home/user/scripts/code.php on line 5

	if len(interpreter.outputBuffers) == 0 {
		return values.NewBool(false), nil
	}

	interpreter.outputBuffers = interpreter.outputBuffers[:len(interpreter.outputBuffers)-1]
	return values.NewBool(true), nil
}

// ------------------- MARK: ob_end_flush -------------------

func nativeFn_ob_end_flush(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ob-end-flush.php

	_, err := NewFuncParamValidator("ob_end_flush").validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// TODO Call output handler
	// Spec: https://www.php.net/manual/en/function.ob-end-flush.php
	// This function calls the output handler (with the PHP_OUTPUT_HANDLER_FINAL flag),
	// flushes (sends) it's return value, discards the contents of the active output buffer and turns off the active output buffer.

	// TODO Throw notice if no buffer: e.g. Notice: ob_end_flush(): Failed to delete buffer. No buffer to delete in /home/user/scripts/code.php on line 5

	if len(interpreter.outputBuffers) == 0 {
		return values.NewBool(false), nil
	}

	nativeFn_ob_flush(args, interpreter)
	nativeFn_ob_end_clean(args, interpreter)
	return values.NewBool(true), nil
}

// ------------------- MARK: ob_get_clean -------------------

func nativeFn_ob_get_clean(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ob-get-clean.php

	_, err := NewFuncParamValidator("ob_get_clean").validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// TODO Call output handler
	// Spec: https://www.php.net/manual/en/function.ob-get-clean.php
	// This function calls the output handler (with the PHP_OUTPUT_HANDLER_CLEAN and PHP_OUTPUT_HANDLER_FINAL flags),
	// discards it's return value, returns the contents of the active output buffer and turns off the active output buffer.

	// TODO Throw notice if no buffer: e.g. Notice: ob_get_clean(): Failed to delete buffer. No buffer to delete in /home/user/scripts/code.php on line 5

	if len(interpreter.outputBuffers) == 0 {
		return values.NewBool(false), nil
	}

	content := interpreter.outputBuffers[len(interpreter.outputBuffers)-1].Content
	nativeFn_ob_end_clean(args, interpreter)
	return values.NewStr(content), nil
}

// ------------------- MARK: ob_get_flush -------------------

func nativeFn_ob_get_flush(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ob-get-flush.php

	_, err := NewFuncParamValidator("ob_get_flush").validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// TODO Call output handler
	// Spec: https://www.php.net/manual/en/function.ob-get-flush.php
	// This function calls the output handler (with the PHP_OUTPUT_HANDLER_FINAL flag),
	// flushes (sends) it's return value, returns the contents of the active output buffer and turns off the active output buffer.

	// TODO Throw notice if no buffer: e.g. Notice: ob_get_flush(): Failed to delete buffer. No buffer to delete in /home/user/scripts/code.php on line 5

	if len(interpreter.outputBuffers) == 0 {
		return values.NewBool(false), nil
	}

	content := interpreter.outputBuffers[len(interpreter.outputBuffers)-1].Content
	nativeFn_ob_end_flush(args, interpreter)
	return values.NewStr(content), nil
}

// ------------------- MARK: ob_get_contents -------------------

func nativeFn_ob_get_contents(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ob-get-contents.php

	_, err := NewFuncParamValidator("nativeFn_ob_get_contents").validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	if len(interpreter.outputBuffers) == 0 {
		return values.NewBool(false), nil
	}

	return values.NewStr(interpreter.outputBuffers[len(interpreter.outputBuffers)-1].Content), nil
}

// ------------------- MARK: ob_get_level -------------------

func nativeFn_ob_get_level(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ob-get-level.php

	_, err := NewFuncParamValidator("ob_get_level").validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewInt(int64(len(interpreter.outputBuffers))), nil
}

// ------------------- MARK: ob_start -------------------

func nativeFn_ob_start(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ob-start

	_, err := NewFuncParamValidator("ob_start").validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// TODO ob_start parameters
	//  ob_start(?callable $callback = null, int $chunk_size = 0, int $flags = PHP_OUTPUT_HANDLER_STDFLAGS): bool

	interpreter.outputBuffers = append(interpreter.outputBuffers, NewOutputBuffer())

	return values.NewBool(true), nil
}

// TODO flush
// TODO ob_​get_​length
// TODO ob_​get_​status
// TODO ob_​implicit_​flush
// TODO ob_​list_​handlers
// TODO output_​add_​rewrite_​var
// TODO output_​reset_​rewrite_​vars
