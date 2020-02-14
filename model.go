package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// Model is a dictionary model of MeCab.
type Model struct {
	model *C.mecab_model_t
}

// NewModel returns a new model.
func NewModel(args map[string]string) (Model, error) {
	// build the argument
	opts := make([]*C.char, 0, len(args)+1)
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

	// create new MeCab model
	model := C.mecab_model_new(C.int(len(opts)), (**C.char)(&opts[0]))
	if model == nil {
		return Model{}, errors.New("mecab: mecab_model is not created")
	}

	return Model{
		model: model,
	}, nil
}

// Destroy frees the model.
func (m Model) Destroy() {
	C.mecab_model_destroy(m.model)
}

// NewMeCab returns a new mecab.
func (m Model) NewMeCab() (MeCab, error) {
	mecab := C.mecab_model_new_tagger(m.model)
	if mecab == nil {
		return MeCab{}, errors.New("mecab: mecab is not created")
	}
	return MeCab{mecab: mecab}, nil
}

// NewLattice returns a new lattice.
func (m Model) NewLattice() (Lattice, error) {
	lattice := C.mecab_model_new_lattice(m.model)
	if lattice == nil {
		return Lattice{}, errors.New("mecab: lattice is not created")
	}
	return Lattice{lattice: lattice}, nil
}

// Swap replaces the model by the other model.
func (m Model) Swap(m2 Model) {
	C.mecab_model_swap(m.model, m2.model)
}
