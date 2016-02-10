package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"

import (
	"errors"
	"unsafe"
)

type Tagger struct {
	tagger *C.mecab_t
}

func (t *Tagger) Destroy() {
	C.mecab_destroy(t.tagger)
}

func (t *Tagger) Parse(s string) (string, error) {
	input := C.CString(s)
	defer C.free(unsafe.Pointer(input))

	result := C.mecab_sparse_tostr(t.tagger, input)
	if result == nil {
		return "", t.Error()
	}
	return C.GoString(result), nil
}

func (t *Tagger) Error() error {
	return errors.New(C.GoString(C.mecab_strerror(t.tagger)))
}
