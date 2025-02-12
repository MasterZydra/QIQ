package config

import (
	"fmt"
	"runtime"
	"strconv"
)

// GoPHP version
const GoPhpVersion string = "0.1.0"

var SoftwareVersion string = softwareVersion()

func softwareVersion() string {
	return fmt.Sprintf("GoPHP/%s (%s)", GoPhpVersion, runtime.GOOS)
}

// PHP version
const MajorVersion int64 = 8
const MinorVersion int64 = 3
const ReleaseVersion int64 = 0
const ExtraVersion string = ""

var Version = fmt.Sprintf("%d.%d.%d%s", MajorVersion, MinorVersion, ReleaseVersion, ExtraVersion)

var VersionId int64 = versionId()

func versionId() int64 {
	versionId, _ := strconv.ParseInt(fmt.Sprintf("%d%02d%02d", MajorVersion, MinorVersion, ReleaseVersion), 10, 64)
	return versionId
}

// Runtime mode
var ShowStats bool = false
var ShowInterpreterCallStack bool = false
