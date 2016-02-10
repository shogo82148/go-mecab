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
	model *C.mecab_model_t
}

func New(args map[string]string) (*MeCab, error) {
	var model *C.mecab_model_t

	// build the argument
	opts := make([]*C.char, 0, len(args)+1)
	opt := C.CString("--allocate-sentence")
	defer C.free(unsafe.Pointer(opt))
	opts = append(opts, opt)
	for k, v := range args {
		opt := C.CString(fmt.Sprintf("--%s=%s", k, v))
		defer C.free(unsafe.Pointer(opt))
		opts = append(opts, opt)
	}

	// create new MeCab model
	model = C.mecab_model_new(C.int(len(opts)), (**C.char)(&opts[0]))
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

func (m *MeCab) NewTagger() (*Tagger, error) {
	tagger := C.mecab_model_new_tagger(m.model)
	if tagger == nil {
		return nil, errors.New("mecab_tagger is not created.")
	}
	return &Tagger{tagger: tagger}, nil
}
