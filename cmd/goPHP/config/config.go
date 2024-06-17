package config

import (
	"fmt"
	"strconv"
)

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
