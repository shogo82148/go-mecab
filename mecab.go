package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

type MeCab struct {
	model *C.mecab_model_t
}

func New(args map[string]string) (*MeCab, error) {
	var model *C.mecab_model_t
	if len(args) > 0 {
		opts := make([]*C.char, 0, len(args))
		for k, v := range args {
			opt := C.CString(fmt.Sprintf("--%s=%s", strings.Replace(k, "_", "-", -1), v))
			defer C.free(unsafe.Pointer(opt))
			opts = append(opts, opt)
		}
		model = C.mecab_model_new(C.int(len(opts)), (**C.char)(&opts[0]))
	} else {
		opt := C.CString("")
		defer C.free(unsafe.Pointer(opt))
		model = C.mecab_model_new2(opt)
	}
	if model == nil {
		return nil, errors.New("mecab_model is not created.")
	}

	return &MeCab{
		model: model,
	}, nil
}

func (m *MeCab) Destroy() {
	C.mecab_model_destroy(m.model)
}
