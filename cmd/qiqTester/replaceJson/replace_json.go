package replacejson

import (
	"QIQ/cmd/qiq/common/os"
	"strings"
)

type ReplaceEntry struct {
	File    string `json:"file"`
	Section string `json:"section"`
	Search  string `json:"search"`
	Replace string `json:"replace"`
}

type ReplaceJson struct {
	Replace []ReplaceEntry `json:"replace"`
}

func (replaceJson *ReplaceJson) GetEntry(filename string) (ReplaceEntry, bool) {
	for _, entry := range replaceJson.Replace {
		entryFile := entry.File
		if os.IS_WIN {
			entryFile = strings.ReplaceAll(entryFile, "/", "\\")
		}
		if entryFile == filename {
			return entry, true
		}
	}
	return ReplaceEntry{}, false
}
