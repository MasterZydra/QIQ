package os

import "runtime"

var OS = Os()

func Os() string {
	switch runtime.GOOS {
	case "android":
		return "Android"
	case "darwin":
		return "Darwin"
	case "dragonfly":
		return "DragonFly"
	case "freebsd":
		return "FreeBSD"
	case "illumos":
		return "IllumOS"
	case "linux":
		return "Linux"
	case "netbsd":
		return "NetBSD"
	case "openbsd":
		return "OpenBSD"
	case "solaris":
		return "Solaris"
	case "windows":
		return "Windows"
	default:
		return "Unkown"
	}
}

var OS_FAMILY = OsFamily()

func OsFamily() string {
	switch runtime.GOOS {
	case "android", "linux":
		return "Linux"
	case "darwin":
		return "Darwin"
	case "dragonfly", "freebsd", "netbsd", "openbsd":
		return "BSD"
	case "solaris":
		return "Solaris"
	case "windows":
		return "Windows"
	default:
		return "Unkown"
	}
}

// Is OS Windows?
var IS_WIN = IsWindows()

func IsWindows() bool { return Os() == "Windows" }

var EOL string = Eol()

func Eol() string {
	if Os() == "Windows" {
		return "\r\n"
	} else {
		return "\n"
	}
}

var DIR_SEP = DirSep()

func DirSep() string {
	if Os() == "Windows" {
		return `\`
	} else {
		return "/"
	}
}
