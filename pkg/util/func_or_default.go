package util

func FuncOrDefault(fn func(), defaultFn func()) {
	if fn != nil {
		fn()
	} else if defaultFn != nil {
		defaultFn()
	}
}
