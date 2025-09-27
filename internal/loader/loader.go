package loader

import (
	"time"

	"github.com/json-iterator/go/extra"
)

func Initialize() {
	extra.RegisterFuzzyDecoders()
	extra.RegisterTimeAsInt64Codec(time.Millisecond)
}
