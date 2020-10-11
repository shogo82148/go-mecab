package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

var errMeCabNotAvailable = errors.New("mecab: mecab is not available")

// to introduce garbage-collection while maintaining backwards compatibility.
type mecab struct {
	mecab *C.mecab_t
}

func newMeCab(m *C.mecab_t) *mecab {
	ret := &mecab{
		mecab: m,
	}
	runtime.SetFinalizer(ret, finalizeMeCab)
	return ret
}

// It is a marker that a mecab must not be copied after the first use.
// See https://github.com/golang/go/issues/8005#issuecomment-190753527
// for details.
func (*mecab) Lock() {}

func finalizeMeCab(m *mecab) {
	if m.mecab != nil {
		C.mecab_destroy(m.mecab)
	}
	m.mecab = nil
}

// MeCab is a morphological parser.
type MeCab struct {
	m *mecab
}

// New returns new MeCab parser.
func New(args map[string]string) (MeCab, error) {
	// build the argument
	opts := make([]*C.char, 0, len(args)+2)
	opt := C.CString("mecab")
	defer C.free(unsafe.Pointer(opt))
	opts = append(opts, opt)
	opt = C.CString("--allocate-sentence")
	defer C.free(unsafe.Pointer(opt))
	opts = append(opts, opt)
	for k, v := range args {
		var goopt string
		if v != "" {
			goopt = fmt.Sprintf("--%s=%s", k, v)
		} else {
			goopt = "--" + k
		}
		opt := C.CString(goopt)
		defer C.free(unsafe.Pointer(opt))
		opts = append(opts, opt)
	}

	// C.mecab_new sets an error in the thread local storage.
	// so C.mecab_new and C.mecab_strerror must be call in same thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// create new MeCab
	m := C.mecab_new(C.int(len(opts)), (**C.char)(&opts[0]))
	if m == nil {
		return MeCab{}, newError(nil)
	}

	return MeCab{
		m: newMeCab(m),
	}, nil
}

// Destroy frees the MeCab parser.
func (m MeCab) Destroy() {
	runtime.SetFinalizer(m.m, nil) // clear the finalizer
	if m.m.mecab != nil {
		C.mecab_destroy(m.m.mecab)
	}
	m.m.mecab = nil
}

// Parse parses the string and returns the result as string
func (m MeCab) Parse(s string) (string, error) {
	if m.m.mecab == nil {
		panic(errMeCabNotAvailable)
	}
	length := C.size_t(len(s))
	if s == "" {
		s = "dummy"
	}
	header := (*reflect.StringHeader)(unsafe.Pointer(&s))
	input := (*C.char)(unsafe.Pointer(header.Data))

	result := C.mecab_sparse_tostr2(m.m.mecab, input, length)
	if result == nil {
		return "", newError(m.m.mecab)
	}
	runtime.KeepAlive(s)
	runtime.KeepAlive(m.m)
	return C.GoString(result), nil
}

// ParseToString is alias of Parse
func (m MeCab) ParseToString(s string) (string, error) {
	if m.m.mecab == nil {
		panic(errMeCabNotAvailable)
	}
	return m.Parse(s)
}

// ParseLattice parses the lattice and returns the result as string.
func (m MeCab) ParseLattice(lattice Lattice) error {
	if m.m.mecab == nil {
		panic(errMeCabNotAvailable)
	}
	if C.mecab_parse_lattice(m.m.mecab, lattice.l.lattice) == 0 {
		return newError(m.m.mecab)
	}
	runtime.KeepAlive(m.m)
	return nil
}

// ParseToNode parses the string and returns the result as Node
func (m MeCab) ParseToNode(s string) (Node, error) {
	if m.m.mecab == nil {
		panic(errMeCabNotAvailable)
	}
	length := C.size_t(len(s))
	if s == "" {
		s = "dummy"
	}
	header := (*reflect.StringHeader)(unsafe.Pointer(&s))
	input := (*C.char)(unsafe.Pointer(header.Data))

	node := C.mecab_sparse_tonode2(m.m.mecab, input, length)
	if node == nil {
		return Node{}, newError(m.m.mecab)
	}
	runtime.KeepAlive(s)
	return Node{
		node:  node,
		mecab: m.m,
	}, nil
}

func (m MeCab) Error() error {
	if m.m.mecab == nil {
		panic(errMeCabNotAvailable)
	}
	return newError(m.m.mecab)
}
