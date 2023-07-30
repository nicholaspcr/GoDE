package fastrand

import (
	"hash/maphash"
	"math/rand"
	"time"
)

// Rand64 returns a pseudo-random uint64. It can be used concurrently and is
// lock-free. Effectively, it calls runtime.fastrand.
func Rand64() uint64 {
	return new(maphash.Hash).Sum64()
}

// NewRand returns a properly seeded *rand.Rand. It has *slightly* higher
// overhead than Rand64 (as it has to allocate), but the resulting PRNG can be
// re-used to offset that cost. Use this if you can't just mask off bits from a
// uint64 (e.g. if you need to use Intn() with something that is not a power of
// 2).
func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
	// return rand.New(rand.NewSource(int64(Rand64())))
}
