package jsonrpc

import "errors"

const (
	ErrCodeParseError     = -32700
	ErrCodeInvalidRequest = -32600
	ErrCodeMethodNotFound = -32601
	ErrCodeInvalidParams  = -32602
	ErrCodeInternalError  = -32603
)

// ErrorInvalidParams states that parameters in the JSON-RPC request are invalid
var ErrorInvalidParams = NewError("Wrong input parameters", ErrCodeInvalidParams)

// RpcError interface representing an rpc error
type Error interface {
	error
	// returns an RPC error code
	Code() int
}

type rpcError struct {
	code int
	err  error
}

// NewError creates new RPC error
func NewError(message string, code int) Error {
	err := errors.New(message)
	return NewFromError(err, code)
}

// NewFromError creates new RPC error from a regular error
func NewFromError(err error, code int) Error {
	return rpcError{
		code: code,
		err:  err,
	}
}

// Error returns string representation of an error
func (e rpcError) Error() string {
	return e.Unwrap().Error()
}

// Unwrap returns a wrapped error
func (e rpcError) Unwrap() error {
	return e.err
}

// Code returns an RPC error code
func (e rpcError) Code() int {
	return e.code
}

type errParamMissing struct {
	// name of the missing parameter
	name string
}

// NewParamMissingError creates an RPC error describing missing parameter in the RPC response
func NewParamMissingError(paramName string) Error {
	return errParamMissing{
		name: paramName,
	}
}

// Error returns error description
func (e errParamMissing) Error() string {
	return "parameter " + e.name + " expected but not found"
}

// Code returns RPC error code
func (e errParamMissing) Code() int {
	return ErrCodeInvalidParams
}

type errMethodNotFound struct {
	// method name
	name string
	// wrapped error
	err error
}

// NewMethodNotFoundError creates an RPC error for non existing method, first argument is a name of a method,
// the second argument is a wrapped error
func NewMethodNotFoundError(name string, err error) Error {
	return errMethodNotFound{
		name: name,
		err:  err,
	}
}

// Error returns string representation of an error
func (e errMethodNotFound) Error() string {
	msg := "method not found"

	if e.name != "" {
		msg = "method \"" + e.name + "\" not found"
	}

	if err := e.Unwrap(); err != nil {
		msg += ": " + err.Error()
	}

	return msg
}

// Unwrap returns a wrapped error
func (e errMethodNotFound) Unwrap() error {
	return e.err
}

// Code returns an RPC error code
func (e errMethodNotFound) Code() int {
	return ErrCodeMethodNotFound
}

type errNamespaceNotFound struct {
	// namespace name
	name string
	// wrapped error
	err error
}

// NewNsNotFoundError creates an RPC error for non existing namespace, first argument is a name of a namespace,
// the second argument is a wrapped error
func NewNsNotFoundError(name string, err error) Error {
	return errNamespaceNotFound{
		name: name,
		err:  err,
	}
}

// Error returns string representation of an error
func (e errNamespaceNotFound) Error() string {
	msg := "namespace " + e.name + " not found"

	if err := e.Unwrap(); err != nil {
		msg += ": " + e.err.Error()
	}

	return msg
}

// Unwrap returns a wrapped error
func (e errNamespaceNotFound) Unwrap() error {
	return e.err
}

func (e errNamespaceNotFound) Code() int {
	return ErrCodeMethodNotFound
}
