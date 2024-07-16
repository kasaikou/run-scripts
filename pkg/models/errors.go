package models

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var regValidateIdentifierExp = `[a-zA-Z]([a-zA-Z0-9_-]*[a-zA-Z0-9]|[a-zA-Z0-9]*)`
var regValidateStringKeyExp = `'([^'\\]|\\'|\\\\)*'`
var regValidateIntegerKeyExp = `((-|\+)?)([0-9]|[1-9][0-9]+)`
var regValidateKeyExp = strings.Join([]string{regValidateStringKeyExp, regValidateIntegerKeyExp}, "|")
var regValidateMemberExp = fmt.Sprintf(`(%s)(\[(%s)\])?`, regValidateIdentifierExp, regValidateKeyExp)
var regValidateAtExp = fmt.Sprintf(`(%s)(\.(%s))*`, regValidateMemberExp, regValidateMemberExp)
var regValidateErrAt = regexp.MustCompile(fmt.Sprintf(`^%s$`, regValidateAtExp))

// ModelValidateError indicates validation error.
type ModelValidateError struct {
	at          []string
	internalErr error
}

// NewModelValidateError creates ModelValidateError instance.
func NewModelValidateError(err error) error {
	return ModelValidateError{internalErr: err}
}

// WrapModelValidateError wraps internal
func WrapModelValidateError(at string, err error) error {

	at = strings.TrimPrefix(strings.TrimSuffix(at, "."), ".")

	switch err := err.(type) {
	case ModelValidateError:
		return ModelValidateError{at: append(err.at, at), internalErr: err.internalErr}

	default:
		return WrapModelValidateError(at, NewModelValidateError(err))
	}
}

// Error indicates validation error message.
func (e ModelValidateError) Error() string {
	if len(e.at) == 0 {
		return fmt.Sprintf("validation error: %s", e.internalErr.Error())
	} else {
		slices.Reverse(slices.Clone(e.at))
		return fmt.Sprintf("validation error in '%s': %s", strings.Join(e.at, "."), e.internalErr.Error())
	}
}

// Is indicates whether error equals with target or not.
func (e ModelValidateError) Is(target error) bool {
	switch err := target.(type) {
	case ModelValidateError:
		return slices.Equal(e.at, err.at) && errors.Is(e.internalErr, err)
	default:
		return errors.Is(e.internalErr, target)
	}
}

// CommandNotFoundError expressions command path is not found.
type CommandNotFoundError struct{ path string }

// NewCommandNotFoundError creates new CommandNotFoundError.
func NewCommandNotFoundError(path string) error { return CommandNotFoundError{path: path} }

// Error indicates error message.
func (e CommandNotFoundError) Error() string { return fmt.Sprintf("'%s' command not found", e.path) }

// EnvironmentNotFoundError expressions environment variable is not found.
type EnvironmentNotFoundError struct{ key string }

// NewEnvironmentNotFoundError creates new EnvironmentNotFoundError.
func NewEnvironmentNotFoundError(key string) error { return EnvironmentNotFoundError{key: key} }

// Error indicates error message.
func (e EnvironmentNotFoundError) Error() string {
	return fmt.Sprintf("'%s' environment not found", e.key)
}

// NotImplementedForOSError expressions no implementation has been created for specific OS.
type NotImplementedForOSError struct{ os string }

// NotImplementedForWindows expressions no implementation has been created for Windows.
var NotImplementedForWindows = NotImplementedForOSError{os: "windows"}

// Error indicates error message.
func (e NotImplementedForOSError) Error() string { return fmt.Sprintf("not implemented for %s", e.os) }

// ErrExitedProcess is expressions executed task has been already exited.
var ErrExitedProcess = errors.New("exited process")
