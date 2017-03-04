package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

type MeCab struct {
	mecab *C.mecab_t
}

func New(args map[string]string) (*MeCab, error) {
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

	// create new MeCab
	mecab := C.mecab_new(C.int(len(opts)), (**C.char)(&opts[0]))
	if mecab == nil {
		return nil, errors.New("mecab is not created.")
	}
	m := &MeCab{
		mecab: mecab,
	}
	runtime.SetFinalizer(m, (*MeCab).Destroy)

	return m, nil
}

func (m *MeCab) Destroy() {
	if m == nil || m.mecab == nil {
		return
	}
	C.mecab_destroy(m.mecab)
	m.mecab = nil
}

// ParseToNode parses the string and returns the result as string
func (m *MeCab) Parse(s string) (string, error) {
	length := C.size_t(len(s))
	if s == "" {
		s = "dummy"
	}
	input := *(**C.char)(unsafe.Pointer(&s))
	result := C.mecab_sparse_tostr2(m.mecab, input, length)
	if result == nil {
		return "", m.Error()
	}
	return C.GoString(result), nil
}

// ParseToString is alias of Parse
func (m *MeCab) ParseToString(s string) (string, error) {
	return m.Parse(s)
}

// ParseToNode parses the lattice and returns the result as string.
func (m *MeCab) ParseLattice(lattice *Lattice) error {
	if C.mecab_parse_lattice(m.mecab, lattice.lattice) == 0 {
		return errors.New("parse error")
	}
	return nil
}

// ParseToNode parses the string and returns the result as Node
func (m *MeCab) ParseToNode(s string) (Node, error) {
	length := C.size_t(len(s))
	if s == "" {
		s = "dummy"
	}
	input := *(**C.char)(unsafe.Pointer(&s))

	node := C.mecab_sparse_tonode2(m.mecab, input, length)
	if node == nil {
		return Node{}, m.Error()
	}
	return Node{node: node}, nil
}

func (m *MeCab) Error() error {
	return errors.New(C.GoString(C.mecab_strerror(m.mecab)))
}
