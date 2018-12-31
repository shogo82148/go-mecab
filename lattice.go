package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"reflect"
	"unsafe"
)

// Lattice is a lattice.
type Lattice struct {
	lattice *C.mecab_lattice_t
}

// NewLattice creates new lattice.
func NewLattice() (Lattice, error) {
	lattice := C.mecab_lattice_new()

	if lattice == nil {
		return Lattice{}, errors.New("mecab: mecab_lattice is not created")
	}
	return Lattice{lattice: lattice}, nil
}

// Destroy frees the lattice.
func (l Lattice) Destroy() {
	C.mecab_lattice_destroy(l.lattice)
}

// Clear set empty string to the lattice.
func (l Lattice) Clear() {
	C.mecab_lattice_clear(l.lattice)
}

// IsAvailable returns the lattice is available.
func (l Lattice) IsAvailable() bool {
	return C.mecab_lattice_is_available(l.lattice) != 0
}

// BOSNode returns the Begin Of Sentence node.
func (l Lattice) BOSNode() Node {
	return Node{node: C.mecab_lattice_get_bos_node(l.lattice)}
}

// EOSNode returns the End Of Sentence node.
func (l Lattice) EOSNode() Node {
	return Node{node: C.mecab_lattice_get_eos_node(l.lattice)}
}

// Sentence returns the sentence in the lattice.
func (l Lattice) Sentence() string {
	return C.GoString(C.mecab_lattice_get_sentence(l.lattice))
}

// SetSentence set the sentence in the lattice.
func (l Lattice) SetSentence(s string) {
	length := C.size_t(len(s))
	if s == "" {
		s = "dummy"
	}
	header := (*reflect.StringHeader)(unsafe.Pointer(&s))
	input := (*C.char)(unsafe.Pointer(header.Data))

	C.mecab_lattice_add_request_type(l.lattice, 64) // MECAB_ALLOCATE_SENTENCE = 64
	C.mecab_lattice_set_sentence2(l.lattice, input, length)
}

func (l Lattice) String() string {
	return C.GoString(C.mecab_lattice_tostr(l.lattice))
}
