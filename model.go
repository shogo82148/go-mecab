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

var errModelNotAvailable = errors.New("mecab: model is not available")

// to introduce garbage-collection while maintaining backwards compatibility.
type model struct {
	model *C.mecab_model_t
}

func newModel(m *C.mecab_model_t) *model {
	ret := &model{
		model: m,
	}
	runtime.SetFinalizer(ret, finalizeModel)
	return ret
}

// It is a marker that a model must not be copied after the first use.
// See https://github.com/golang/go/issues/8005#issuecomment-190753527
// for details.
func (*model) Lock() {}

func finalizeModel(m *model) {
	if m.model != nil {
		C.mecab_model_destroy(m.model)
	}
	m.model = nil
}

// Model is a dictionary model of MeCab.
type Model struct {
	m *model
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
	m := C.mecab_model_new(C.int(len(opts)), (**C.char)(&opts[0]))
	if m == nil {
		return Model{}, newError(nil)
	}

	return Model{
		m: newModel(m),
	}, nil
}

// Destroy frees the model.
func (m Model) Destroy() {
	runtime.SetFinalizer(m.m, nil) // clear the finalizer
	if m.m.model != nil {
		C.mecab_model_destroy(m.m.model)
	}
	m.m.model = nil
}

// NewMeCab returns a new mecab.
func (m Model) NewMeCab() (MeCab, error) {
	if m.m.model == nil {
		panic(errModelNotAvailable)
	}

	// C.mecab_model_new_tagger sets an error in the thread local storage.
	// so C.mecab_model_new_tagger and C.mecab_strerror must be call in same thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	mm := C.mecab_model_new_tagger(m.m.model)
	if mm == nil {
		return MeCab{}, newError(nil)
	}
	runtime.KeepAlive(m.m)
	return MeCab{m: newMeCab(mm)}, nil
}

// NewLattice returns a new lattice.
func (m Model) NewLattice() (Lattice, error) {
	if m.m.model == nil {
		panic(errModelNotAvailable)
	}

	// C.mecab_model_new_lattice sets an error in the thread local storage.
	// so C.mecab_model_new_lattice and C.mecab_strerror must be call in same thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	lattice := C.mecab_model_new_lattice(m.m.model)
	if lattice == nil {
		return Lattice{}, newError(nil)
	}
	return Lattice{l: newLattice(lattice)}, nil
}

// Swap replaces the model by the other model.
func (m Model) Swap(m2 Model) error {
	if m.m.model == nil || m2.m.model == nil {
		panic(errModelNotAvailable)
	}

	// C.mecab_model_swap sets an error in the thread local storage.
	// so C.mecab_model_swap and C.mecab_strerror must be call in same thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	C.mecab_model_swap(m.m.model, m2.m.model)
	return newError(nil)
}
