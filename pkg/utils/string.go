package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
	"unicode"
)

var randObject *rand.Rand

func init() {
	randObject = rand.New(rand.NewSource(time.Now().UnixNano()))
}

var StringUtil = newStringUtil()

type stringUtil struct {
	Alpha         string
	Number        string
	Alnum         string
	AlnumUserRead string
}

func newStringUtil() *stringUtil {
	bf := bytes.NewBuffer(nil)
	var i byte
	for i = 'a'; i <= 'z'; i++ {
		bf.WriteByte(i)
	}
	for i = 'A'; i <= 'Z'; i++ {
		bf.WriteByte(i)
	}
	number := "1234567890"
	alpha := bf.String()
	alnum := alpha + number
	alnumUserRead := alnum
	alnumUserRead = strings.ReplaceAll(alnumUserRead, "0", "")
	alnumUserRead = strings.ReplaceAll(alnumUserRead, "O", "")
	alnumUserRead = strings.ReplaceAll(alnumUserRead, "l", "")
	alnumUserRead = strings.ReplaceAll(alnumUserRead, "1", "")
	return &stringUtil{
		Alpha:         alpha,
		Number:        number,
		Alnum:         alnum,
		AlnumUserRead: alnumUserRead,
	}
}

func (util *stringUtil) CamelCase(s string) string {
	b := &strings.Builder{}
	length := len(s)
	for i := 0; i < length; i++ {
		if s[i] == '_' {
			if i == length-1 {

			} else {
				i += 1
				c := s[i]
				if 'a' <= s[i] && s[i] <= 'z' {
					c = s[i] - 32
				}
				b.WriteByte(c)
			}
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

func (util *stringUtil) SnakeCase(s string) string {
	b := &strings.Builder{}
	length := len(s)
	for i := 0; i < length; i++ {
		if 'A' <= s[i] && s[i] <= 'Z' {
			b.WriteByte('_')
			b.WriteByte(s[i] + 32)
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

func (util *stringUtil) LcFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	if unicode.IsUpper([]rune(s)[0]) {
		sr := []rune(s)
		sr[0] = unicode.ToLower(sr[0])
		return string(sr)
	}
	return s
}
func (util *stringUtil) UcFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	if unicode.IsLower([]rune(s)[0]) {
		sr := []rune(s)
		sr[0] = unicode.ToUpper(sr[0])
		return string(sr)
	}
	return s
}

func (util *stringUtil) Random(s string, n int) string {
	bf := bytes.NewBuffer(nil)
	for i := 0; i < n; i++ {
		idx := randObject.Intn(len(s))
		bf.WriteByte(s[idx])
	}
	return bf.String()
}

func (util *stringUtil) StartsWith(s, prefix string) bool {
	if len(prefix) < len(s) {
		return prefix == s[:len(prefix)]
	}
	return false
}

func (util *stringUtil) Sub(s string, start int, length int) string {
	var a = []rune(s)
	n := len(a)
	if n == 0 {
		return ""
	}
	if start >= n {
		return ""
	} else if start < 0 {
		start = n + start
		if start < 0 {
			start = 0
		}
	}
	if start+length > n {
		length = n - start
	}
	return string(a[start : start+length])
}

func (util *stringUtil) Concat(sep string, a ...string) string {
	bf := bytes.NewBuffer(nil)
	for i, s := range a {
		if i == 0 {
			bf.WriteString(strings.TrimRight(s, sep))
		} else if i == len(a)-1 {
			bf.WriteString(strings.TrimLeft(s, sep))
		} else {
			bf.WriteString(strings.Trim(s, sep))
		}
		if i != len(a)-1 {
			bf.WriteString(sep)
		}
	}
	return bf.String()
}

func (util *stringUtil) ToTime(s string) int64 {
	timeLayout := "2006-01-02 15:04:05"
	timestamp, err := time.ParseInLocation(timeLayout, s, TimeUtil.Location)
	if err != nil {
		return 0
	}
	return timestamp.Unix()
}

func (util *stringUtil) Trim(s string) string {
	return strings.Trim(s, " ")
}

func (util *stringUtil) UniqueID() string {
	t := time.Now().UnixNano()
	return strconv.FormatInt(t, 36) + util.Random(util.Alnum, 4)
}

func (util *stringUtil) ToString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case *string:
		return *v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case *float32, *float64:
		return fmt.Sprintf("%f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (util *stringUtil) TrimAndTolower(s string) string {
	return strings.ToLower(strings.Trim(s, " "))
}
