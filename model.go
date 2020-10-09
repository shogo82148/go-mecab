package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"runtime"
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

	// C.mecab_model_new sets an error in the thread local storage.
	// so C.mecab_model_new and C.mecab_strerror must be call in same thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// create new MeCab model
	model := C.mecab_model_new(C.int(len(opts)), (**C.char)(&opts[0]))
	if model == nil {
		return Model{}, newError(nil)
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
	// C.mecab_model_new_tagger sets an error in the thread local storage.
	// so C.mecab_model_new_tagger and C.mecab_strerror must be call in same thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	mecab := C.mecab_model_new_tagger(m.model)
	if mecab == nil {
		return MeCab{}, newError(nil)
	}
	return MeCab{mecab: mecab}, nil
}

// NewLattice returns a new lattice.
func (m Model) NewLattice() (Lattice, error) {
	// C.mecab_model_new_lattice sets an error in the thread local storage.
	// so C.mecab_model_new_lattice and C.mecab_strerror must be call in same thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	lattice := C.mecab_model_new_lattice(m.model)
	if lattice == nil {
		return Lattice{}, newError(nil)
	}
	return Lattice{lattice: lattice}, nil
}

// Swap replaces the model by the other model.
func (m Model) Swap(m2 Model) error {
	// C.mecab_model_swap sets an error in the thread local storage.
	// so C.mecab_model_swap and C.mecab_strerror must be call in same thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	C.mecab_model_swap(m.model, m2.model)
	return newError(nil)
}
