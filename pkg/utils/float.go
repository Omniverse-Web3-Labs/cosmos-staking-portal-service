package utils

import (
	"fmt"
	"strconv"
)

var FloatUtil = newFloatUtil()

type floatUtil struct {
}

func newFloatUtil() *floatUtil {
	return &floatUtil{}
}

func (*floatUtil) Decimal(value float64, radix int) float64 {
	s := "%." + strconv.Itoa(radix) + "f"
	value, _ = strconv.ParseFloat(fmt.Sprintf(s, value), 64)
	return value
}
