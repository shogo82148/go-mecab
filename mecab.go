package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"

import (
	"errors"
	"unsafe"
)

type MeCab struct {
	mecab *C.mecab_t
}

func (m *MeCab) Destroy() {
	C.mecab_destroy(m.mecab)
}

func (m *MeCab) Parse(s string) (string, error) {
	input := C.CString(s)
	defer C.free(unsafe.Pointer(input))

	result := C.mecab_sparse_tostr(m.mecab, input)
	if result == nil {
		return "", m.Error()
	}
	return C.GoString(result), nil
}

func (m *MeCab) Error() error {
	return errors.New(C.GoString(C.mecab_strerror(m.mecab)))
}
