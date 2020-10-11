package mecab

import (
	"runtime"
	"testing"
)

func TestLatticeFinalizer(t *testing.T) {
	for i := 0; i < 10000; i++ {
		NewLattice()
	}
	runtime.GC()
	runtime.GC()
	runtime.GC()
}
