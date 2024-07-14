package md

import (
	"testing"
)

func BenchmarkAllocation(b *testing.B) {

	bytes := []byte("hello, world")

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			s := bytes2Str(bytes)
			s += ""
		}
	})
}
