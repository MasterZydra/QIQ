//go:build windows

package common

import (
	"syscall"
)

func FileExistsCaseSensitive(path string) (bool, error) {
	// Convert to UTF-16 for Windows API
	p16, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, err
	}

	// First check existence
	attrs, err := syscall.GetFileAttributes(p16)
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok && errno == syscall.ERROR_FILE_NOT_FOUND {
			return false, nil
		}
		return false, err
	}
	if attrs == syscall.INVALID_FILE_ATTRIBUTES {
		return false, nil
	}

	// Get canonical case-preserved path using GetLongPathNameW
	buf := make([]uint16, syscall.MAX_PATH)
	n, err := syscall.GetLongPathName(p16, &buf[0], uint32(len(buf)))
	if err != nil {
		return false, err
	}
	if n == 0 {
		// failed, but file exists
		return true, nil
	}

	actual := syscall.UTF16ToString(buf[:n])

	// Compare exact case-sensitive
	return actual == path, nil
}
