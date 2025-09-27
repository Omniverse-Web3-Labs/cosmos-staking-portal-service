package utils

import (
	"bytes"
	"path"
	"strings"
)

var PathUtil = newPathUtil()

type pathUtil struct {
}

func newPathUtil() *pathUtil {
	return &pathUtil{}
}

func (util *pathUtil) Join(a ...string) string {
	bs := bytes.NewBuffer(nil)
	for i, s := range a {
		if i == 0 {
			s = strings.TrimRight(s, "/")
			bs.WriteString(s)
		} else if i == len(a)-1 {
			s = strings.TrimLeft(s, "/")
			bs.WriteString("/")
			bs.WriteString(s)
		} else {
			s = strings.Trim(s, "/")
			bs.WriteString("/")
			bs.WriteString(s)
		}
	}
	return bs.String()
}

func (util *pathUtil) Base(s string) string {
	return s[:len(s)-len(path.Ext(s))]
}
