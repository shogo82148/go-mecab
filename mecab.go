package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

type MeCab struct {
	mecab *C.mecab_t
}

func New(args map[string]string) (*MeCab, error) {
	// build the argument
	opts := make([]*C.char, 0, len(args)+1)
	opt := C.CString("--allocate-sentence")
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

	return &MeCab{
		mecab: mecab,
	}, nil
}

func (m *MeCab) Destroy() {
	C.mecab_destroy(m.mecab)
}

// ParseToNode parses the string and returns the result as string
func (m *MeCab) Parse(s string) (string, error) {
	input := C.CString(s)
	defer C.free(unsafe.Pointer(input))

	result := C.mecab_sparse_tostr(m.mecab, input)
	if result == nil {
		return "", m.Error()
	}
	return C.GoString(result), nil
}

// ParseToString is alias of Parse
func (m *MeCab) ParseToString(s string) (string, error) {
	return m.Parse(s)
}

// ParseToNode parses the string and returns the result as Node
func (m *MeCab) ParseToNode(s string) (*Node, error) {
	input := C.CString(s)
	defer C.free(unsafe.Pointer(input))

	node := C.mecab_sparse_tonode(m.mecab, input)
	if node == nil {
		return nil, m.Error()
	}
	return &Node{node: node}, nil
}

func (m *MeCab) Error() error {
	return errors.New(C.GoString(C.mecab_strerror(m.mecab)))
}
