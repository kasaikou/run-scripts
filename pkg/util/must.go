package util

// Must returns v1 if err is nil, otherwise it calls panic.
func Must[T1 any](v1 T1, err error) T1 {
	if err != nil {
		panic(err)
	}

	return v1
}
