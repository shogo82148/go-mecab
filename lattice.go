package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

// RequestType is a request type.
type RequestType int

const (
	// RequestTypeOneBest is a request type for one best result.
	RequestTypeOneBest RequestType = 1

	// RequestTypeNBest is a request type for N-best results.
	RequestTypeNBest RequestType = 2

	// RequestTypePartial enables a partial parsing mode.
	// When this flag is set, the input |sentence| needs to be written
	// in partial parsing format.
	RequestTypePartial RequestType = 4

	// RequestTypeMarginalProb is a request type for marginal probability.
	// Set this flag if you want to obtain marginal probabilities.
	// Marginal probability is set in [Node.Prob].
	// The parsing speed will get 3-5 times slower than the default mode.
	RequestTypeMarginalProb RequestType = 8

	// RequestTypeMorphsToNBest is a request type for alternative results.
	// Set this flag if you want to obtain alternative results.
	// Not implemented.
	RequestTypeAlternative RequestType = 16

	// RequestTypeAllMorphs is a request type for all morphs.
	RequestTypeAllMorphs RequestType = 32

	// RequestTypeAllocateSentence is a request type for allocating sentence.
	// When this flag is set, tagger internally copies the body of passed
	// sentence into internal buffer.
	RequestTypeAllocateSentence RequestType = 64
)

var errLatticeNotAvailable = errors.New("mecab: lattice is not available")

type lattice struct {
	lattice *C.mecab_lattice_t
}

func newLattice(l *C.mecab_lattice_t) *lattice {
	ret := &lattice{
		lattice: l,
	}
	runtime.SetFinalizer(ret, finalizeLattice)
	return ret
}

// It is a marker that a lattice must not be copied after the first use.
// See https://github.com/golang/go/issues/8005#issuecomment-190753527
// for details.
func (*lattice) Lock() {}

func finalizeLattice(l *lattice) {
	if l.lattice != nil {
		C.mecab_lattice_destroy(l.lattice)
	}
	l.lattice = nil
}

// Lattice is a lattice.
type Lattice struct {
	l *lattice
}

// NewLattice creates new lattice.
func NewLattice() (Lattice, error) {
	// C.mecab_lattice_new sets an error in the thread local storage.
	// so C.mecab_lattice_new and C.mecab_strerror must be call in same thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	l := C.mecab_lattice_new()
	if l == nil {
		return Lattice{}, newError(nil)
	}
	return Lattice{l: newLattice(l)}, nil
}

// Destroy frees the lattice.
func (l Lattice) Destroy() {
	runtime.SetFinalizer(l.l, nil) // clear the finalizer
	if l.l.lattice != nil {
		C.mecab_lattice_destroy(l.l.lattice)
	}
	l.l.lattice = nil
}

// Clear set empty string to the lattice.
func (l Lattice) Clear() {
	if l.l.lattice == nil {
		panic(errLatticeNotAvailable)
	}
	C.mecab_lattice_clear(l.l.lattice)
	runtime.KeepAlive(l.l)
}

// IsAvailable returns the lattice is available.
func (l Lattice) IsAvailable() bool {
	if l.l.lattice == nil {
		return false
	}
	available := C.mecab_lattice_is_available(l.l.lattice) != 0
	runtime.KeepAlive(l.l)
	return available
}

// BOSNode returns the Begin Of Sentence node.
func (l Lattice) BOSNode() Node {
	if l.l.lattice == nil {
		panic(errLatticeNotAvailable)
	}
	return Node{
		node:    C.mecab_lattice_get_bos_node(l.l.lattice),
		lattice: l.l,
	}
}

// EOSNode returns the End Of Sentence node.
func (l Lattice) EOSNode() Node {
	if l.l.lattice == nil {
		panic(errLatticeNotAvailable)
	}
	return Node{
		node:    C.mecab_lattice_get_eos_node(l.l.lattice),
		lattice: l.l,
	}
}

// Sentence returns the sentence in the lattice.
func (l Lattice) Sentence() string {
	if l.l.lattice == nil {
		panic(errLatticeNotAvailable)
	}
	s := C.GoString(C.mecab_lattice_get_sentence(l.l.lattice))
	runtime.KeepAlive(l.l)
	return s
}

// SetSentence set the sentence in the lattice.
func (l Lattice) SetSentence(s string) {
	if l.l.lattice == nil {
		panic(errLatticeNotAvailable)
	}
	length := C.size_t(len(s))
	input := C.CString(s)
	defer C.free(unsafe.Pointer(input))

	C.mecab_lattice_add_request_type(l.l.lattice, C.int(RequestTypeAllocateSentence)) // MECAB_ALLOCATE_SENTENCE = 64
	C.mecab_lattice_set_sentence2(l.l.lattice, input, length)
	runtime.KeepAlive(l.l)
}

func (l Lattice) String() string {
	if l.l.lattice == nil {
		panic(errLatticeNotAvailable)
	}
	s := C.GoString(C.mecab_lattice_tostr(l.l.lattice))
	runtime.KeepAlive(l.l)
	return s
}

// Next obtains next-best result. The internal linked list structure is updated.
// You should set [RequestTypeNBest] in advance.
// Return false if no more results are available or [RequestType] is invalid.
func (l Lattice) Next() bool {
	if l.l.lattice == nil {
		panic(errLatticeNotAvailable)
	}
	next := C.mecab_lattice_next(l.l.lattice) != 0
	runtime.KeepAlive(l)
	return next
}

// RequestType returns the request type.
func (l Lattice) RequestType() RequestType {
	if l.l.lattice == nil {
		panic(errLatticeNotAvailable)
	}
	return RequestType(C.mecab_lattice_get_request_type(l.l.lattice))
}

// SetRequestType sets the request type.
func (l Lattice) SetRequestType(t RequestType) {
	if l.l.lattice == nil {
		panic(errLatticeNotAvailable)
	}
	C.mecab_lattice_add_request_type(l.l.lattice, C.int(t))
	runtime.KeepAlive(l)
}

// AddRequestType adds the request type.
func (l Lattice) AddRequestType(t RequestType) {
	if l.l.lattice == nil {
		panic(errLatticeNotAvailable)
	}
	C.mecab_lattice_add_request_type(l.l.lattice, C.int(t))
	runtime.KeepAlive(l)
}
