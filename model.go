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

type Model struct {
	model *C.mecab_model_t
}

func NewModel(args map[string]string) (*Model, error) {
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
		return nil, errors.New("mecab_model is not created.")
	}
	m := &Model{model: model}
	runtime.SetFinalizer(m, (*Model).Destroy)

	return m, nil
}

func (m *Model) Destroy() {
	if m == nil || m.model == nil {
		return
	}
	C.mecab_model_destroy(m.model)
	m.model = nil
}

func (m *Model) NewMeCab() (*MeCab, error) {
	mecab := C.mecab_model_new_tagger(m.model)
	if mecab == nil {
		return nil, errors.New("mecab is not created.")
	}
	m2 := &MeCab{mecab: mecab}
	runtime.SetFinalizer(m2, (*MeCab).Destroy)

	return m2, nil
}

func (m *Model) NewLattice() (*Lattice, error) {
	lattice := C.mecab_model_new_lattice(m.model)
	if lattice == nil {
		return nil, errors.New("lattice is not created")
	}
	l := &Lattice{lattice: lattice}
	runtime.SetFinalizer(l, (*Lattice).Destroy)

	return l, nil
}

// Swap replaces the model by the othor model.
func (m *Model) Swap(m2 *Model) {
	C.mecab_model_swap(m.model, m2.model)
}
