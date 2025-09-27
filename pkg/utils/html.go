package utils

import (
	"bytes"
	"html/template"
)

var HtmlUtil = newHtmlUtil()

type htmlUtil struct {
}

func newHtmlUtil() *htmlUtil {
	return &htmlUtil{}
}

// Generate html
func (util *htmlUtil) HTML(html string, data interface{}) string {
	t := template.New("html").Delims("{{", "}}")
	t, _ = t.Parse(html)
	bf := bytes.NewBuffer(nil)
	_ = t.Execute(bf, data)
	return bf.String()
}
