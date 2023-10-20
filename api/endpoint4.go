package api

import (
	"fmt"
	"strings"
)

func NewHeaderInfo(active, endpath, section string, path ...string) HeaderInfo {
	if len(endpath) > 50 {
		endpath = endpath[:50]
	}
	return HeaderInfo{
		Active:  active,
		Path:    strings.Join(path, " / "),
		EndPath: endpath,
		Section: section,
	}
}

func LimitStringSize(s string, size int) string {
	if len(s) > size {
		return fmt.Sprintf("%v...", s[:size])
	}
	return s
}
