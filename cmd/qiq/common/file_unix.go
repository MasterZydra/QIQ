//go:build !windows

package common

func FileExistsCaseSensitive(path string) (bool, error) { return PathExists(path), nil }
