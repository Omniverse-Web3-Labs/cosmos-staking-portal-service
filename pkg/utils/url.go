package utils

import (
	"net/url"
	"strings"
)

var UrlUtil = newUrlUtil()

type urlUtil struct {
}

func newUrlUtil() *urlUtil {
	return &urlUtil{}
}

func (*urlUtil) ParseQuery(query string) (map[string]string, error) {
	data, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for k, v := range data {
		m[k] = ""
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m, nil
}

func (*urlUtil) IsHttpURL(path string) bool {
	return StringUtil.StartsWith(path, "http://") || StringUtil.StartsWith(path, "https://")
}

func (*urlUtil) IsDataURL(path string) bool {
	return StringUtil.StartsWith(path, "blob:") || StringUtil.StartsWith(path, "data:")
}

func (util *urlUtil) RealURL(urlstr string, path string) string {
	path = StringUtil.Trim(path)
	if path == "" {
		return path
	}
	if util.IsHttpURL(path) || util.IsDataURL(path) {
		return path
	}
	u, _ := url.Parse(urlstr)
	if u != nil {
		if StringUtil.StartsWith(path, "//") {
			return u.Scheme + ":" + path
		} else if StringUtil.StartsWith(path, "/") {
			return u.Scheme + ":" + u.Host + path
		}
	}
	return PathUtil.Join(urlstr, path)
}

func (util *urlUtil) Suffix(urlstr string) string {
	suffix := ""
	tmp := strings.Split(strings.Split(urlstr, "?")[0], ".")
	l := len(tmp)
	if l > 1 {
		suffix = tmp[l-1]
		var i int
		for i = range suffix {
			if 'a' <= suffix[i] && suffix[i] <= 'z' || 'A' <= suffix[i] && suffix[i] <= 'Z' {
				continue
			} else {
				i--
				break
			}
		}
		suffix = suffix[0 : i+1]
	}
	return strings.ToLower(suffix)
}
