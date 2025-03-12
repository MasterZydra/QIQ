package filesystem

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/values"
)

func Register(environment runtime.Environment) {
	// Category: Filesystem Functions
	environment.AddNativeFunction("file_exists", nativeFn_file_exists)
}

// ------------------- MARK: file_exists -------------------

func nativeFn_file_exists(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.file-exists.php

	// args, err := funcParamValidator.NewValidator("file_exists").AddParam("$value", []string{"mixed"}, nil).Validate(args)
	// if err != nil {
	// 	return values.NewVoid(), err
	// }

	return values.NewBool(true), nil
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
// TODO file_​exists
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
// TODO is_​dir
// TODO is_​executable
// TODO is_​file
// TODO is_​link
// TODO is_​readable
// TODO is_​uploaded_​file
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
// TODO rename
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
