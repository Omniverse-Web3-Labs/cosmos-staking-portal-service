package utils

import (
	"bytes"
	"encoding/json"
	"regexp"

	jsoniter "github.com/json-iterator/go"
)

var JsonUtil = newJsonUtil()

type jsonUtil struct {
	keyMatchRegex    *regexp.Regexp
	wordBarrierRegex *regexp.Regexp
}

func newJsonUtil() *jsonUtil {
	return &jsonUtil{
		keyMatchRegex:    regexp.MustCompile(`\"(\w+)\":`),
		wordBarrierRegex: regexp.MustCompile(`(\w)([A-Z])`),
	}
}

type JsonCamelCase struct {
	Value interface{}
}

func (c JsonCamelCase) MarshalJSON() ([]byte, error) {
	marshalled, err := jsoniter.Marshal(c.Value)
	converted := JsonUtil.keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			matchStr := string(match)
			key := matchStr[1 : len(matchStr)-2]
			resKey := StringUtil.LcFirst(StringUtil.CamelCase(key))
			return []byte(`"` + resKey + `":`)
		},
	)
	return converted, err
}

type JsonSnakeCase struct {
	Value interface{}
}

func (c JsonSnakeCase) MarshalJSON() ([]byte, error) {
	// Regexp definitions
	marshalled, err := json.Marshal(c.Value)
	converted := JsonUtil.keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			return bytes.ToLower(JsonUtil.wordBarrierRegex.ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)
	return converted, err
}
