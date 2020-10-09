package mecab

// Error is an error of MeCab.
type Error struct {
	err string
}

func (e *Error) Error() string {
	return e.err
}

func newError(m *C.mecab_t) error {
	err := C.GoString(C.mecab_strerror(m))
	if err == "" {
		return nil
	}
	return &Error{
		err: err,
	}
}
