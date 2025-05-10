package util

import (
	"sync"

	"github.com/ExtraWhy/internal-libs/mt19937"
)

var (
	do sync.Once
)

// specialize for those only
type PrimitiveSpecialize interface {
	int | int32 | int64 | uint16 | uint64 | uint32
}

func RandMT[T PrimitiveSpecialize](x T) T {

	do.Do(func() {
		mt19937.InitEX()
	})

	r := mt19937.Genrand64_int64() % uint64(x)
	return T(r)
}
