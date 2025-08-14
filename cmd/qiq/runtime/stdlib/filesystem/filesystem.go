package filesystem

import (
	"QIQ/cmd/qiq/common"
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime"
	"QIQ/cmd/qiq/runtime/funcParamValidator"
	"QIQ/cmd/qiq/runtime/values"
	"os"
	"slices"
)

func Register(environment runtime.Environment) {
	// Category: Filesystem Functions
	environment.AddNativeFunction("file_get_contents", nativeFn_file_get_contents)
	environment.AddNativeFunction("is_dir", nativeFn_is_dir)
	environment.AddNativeFunction("is_file", nativeFn_is_file)
	environment.AddNativeFunction("is_uploaded_file", nativeFn_is_uploaded_file)
	environment.AddNativeFunction("file_exists", nativeFn_file_exists)
	environment.AddNativeFunction("rename", nativeFn_rename)
}

// -------------------------------------- file_get_contents -------------------------------------- MARK: file_get_contents

func nativeFn_file_get_contents(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/function.file-get-contents.php

	args, err := funcParamValidator.NewValidator("file_get_contents").AddParam("$filename", []string{"string"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// TODO file_get_contents: Add support for other params
	//  file_get_contents(
	// string $filename,
	// bool $use_include_path = false,
	// ?resource $context = null,
	// int $offset = 0,
	// ?int $length = null
	// ): string|false

	filename := args[0].(*values.Str).Value

	if !common.PathExists(filename) {
		return values.NewBool(false), phpError.NewWarning("could not find file %s", filename)
	}

	content, goErr := common.ReadFile(filename)
	if goErr != nil {
		return values.NewBool(false), phpError.NewError("%s", goErr)
	}
	return values.NewStr(content), nil
}

// -------------------------------------- is_dir -------------------------------------- MARK: is_dir

func nativeFn_is_dir(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.is-dir.php

	args, err := funcParamValidator.NewValidator("is_dir").AddParam("$filename", []string{"string"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	info, goErr := os.Stat(args[0].(*values.Str).Value)
	if goErr != nil {
		return values.NewVoid(), phpError.NewError("%s", goErr)
	}

	return values.NewBool(info.IsDir()), nil
}

// -------------------------------------- is_file -------------------------------------- MARK: is_file

func nativeFn_is_file(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.is-file.php

	args, err := funcParamValidator.NewValidator("is_file").AddParam("$filename", []string{"string"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	info, goErr := os.Stat(args[0].(*values.Str).Value)
	if goErr != nil {
		return values.NewVoid(), phpError.NewError("%s", goErr)
	}

	return values.NewBool(!info.IsDir()), nil
}

// -------------------------------------- is_uploaded_file -------------------------------------- MARK: is_uploaded_file

func nativeFn_is_uploaded_file(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.is-uploaded-file.php

	args, err := funcParamValidator.NewValidator("is_uploaded_file").AddParam("$filename", []string{"string"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewBool(slices.Contains(context.Interpreter.GetRequest().UploadedFiles, args[0].(*values.Str).Value)), nil
}

// -------------------------------------- file_exists -------------------------------------- MARK: file_exists

func nativeFn_file_exists(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.file-exists.php

	args, err := funcParamValidator.NewValidator("file_exists").AddParam("$filename", []string{"string"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewBool(Exists(args[0].(*values.Str).Value)), nil
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// -------------------------------------- rename -------------------------------------- MARK: rename

func nativeFn_rename(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.rename.php

	args, err := funcParamValidator.NewValidator("rename").
		AddParam("$from", []string{"string"}, nil).
		AddParam("$to", []string{"string"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	goErr := os.Rename(args[0].(*values.Str).Value, args[1].(*values.Str).Value)
	return values.NewBool(goErr == nil), nil
}

// TODO basename
// TODO chgrp
// TODO chmod
// TODO chown
// TODO clearstatcache
// TODO copy
// TODO delete
// TODO dirname
// TODO disk_​free_​space
// TODO disk_​total_​space
// TODO diskfreespace
// TODO fclose
// TODO fdatasync
// TODO feof
// TODO fflush
// TODO fgetc
// TODO fgetcsv
// TODO fgets
// TODO fgetss
// TODO file
// TODO file_​get_​contents
// TODO file_​put_​contents
// TODO fileatime
// TODO filectime
// TODO filegroup
// TODO fileinode
// TODO filemtime
// TODO fileowner
// TODO fileperms
// TODO filesize
// TODO filetype
// TODO flock
// TODO fnmatch
// TODO fopen
// TODO fpassthru
// TODO fputcsv
// TODO fputs
// TODO fread
// TODO fscanf
// TODO fseek
// TODO fstat
// TODO fsync
// TODO ftell
// TODO ftruncate
// TODO fwrite
// TODO glob
// TODO is_​executable
// TODO is_​link
// TODO is_​readable
// TODO is_​writable
// TODO is_​writeable
// TODO lchgrp
// TODO lchown
// TODO link
// TODO linkinfo
// TODO lstat
// TODO mkdir
// TODO move_​uploaded_​file
// TODO parse_​ini_​file
// TODO parse_​ini_​string
// TODO pathinfo
// TODO pclose
// TODO popen
// TODO readfile
// TODO readlink
// TODO realpath
// TODO realpath_​cache_​get
// TODO realpath_​cache_​size
// TODO rewind
// TODO rmdir
// TODO set_​file_​buffer
// TODO stat
// TODO symlink
// TODO tempnam
// TODO tmpfile
// TODO touch
// TODO umask
// TODO unlink
