package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"
)

type Lattice struct {
	lattice *C.mecab_lattice_t
}

func NewLattice() (Lattice, error) {
	lattice := C.mecab_lattice_new()

	if lattice == nil {
		return Lattice{}, errors.New("mecab_lattice is not created")
	}
	return Lattice{lattice: lattice}, nil
}

func (l Lattice) Destroy() {
	C.mecab_lattice_destroy(l.lattice)
}

func (l Lattice) Clear() {
	C.mecab_lattice_clear(l.lattice)
}

func (l Lattice) IsAvailable() bool {
	return C.mecab_lattice_is_available(l.lattice) != 0
}

func (l Lattice) BOSNode() Node {
	return Node{node: C.mecab_lattice_get_bos_node(l.lattice)}
}

func (l Lattice) EOSNode() Node {
	return Node{node: C.mecab_lattice_get_eos_node(l.lattice)}
}

func (l Lattice) Sentence() string {
	return C.GoString(C.mecab_lattice_get_sentence(l.lattice))
}

func (l Lattice) SetSentence(s string) {
	length := C.size_t(len(s))
	if s == "" {
		s = "dummy"
	}
	input := *(**C.char)(unsafe.Pointer(&s))
	C.mecab_lattice_add_request_type(l.lattice, 64) // MECAB_ALLOCATE_SENTENCE = 64
	C.mecab_lattice_set_sentence2(l.lattice, input, length)
}

func (l Lattice) String() string {
	return C.GoString(C.mecab_lattice_tostr(l.lattice))
}
