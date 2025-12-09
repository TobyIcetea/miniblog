package benchmark

import (
	"math"
	"testing"
)

func BenchmarkExp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = math.Exp(3.5)
	}
}
