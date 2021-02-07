package jsonrpc

import (
	"strings"
)

const (
	version      = "2.0"
	defaultNsSep = "."
)

// Rpc RPC service handling a JSON-RPC request
type Rpc struct {
	Namespace
	NsSep string
}

// NewRpc creates new rpc service
func NewRpc() *Rpc {
	ns := NewNamespace()
	return &Rpc{
		Namespace: ns,
		NsSep:     defaultNsSep,
	}
}

// Exec handles JSON-RPC request and calls an implementation of a requested method
func (r *Rpc) Exec(req *Request) *Response {
	callback, err := r.getMethod(req)

	if err != nil {
		return NewErrResponse(req.ID, err)
	}

	result, err := callback(req.Params)

	if err != nil {
		return NewErrResponse(req.ID, err)
	}

	return &Response{
		JsonRpc: version,
		Result:  result,
		Error:   nil,
		ID:      req.ID,
	}
}

func (r *Rpc) getMethod(req *Request) (Callback, Error) {
	stack := strings.Split(req.Method, defaultNsSep)
	return findMethod(r.Namespace, stack)
}

func findMethod(ns Namespace, stack []string) (Callback, Error) {
	if len(stack) == 0 {
		return nil, NewMethodNotFoundError("", nil)
	}

	if len(stack) == 1 {
		m, err := ns.Callback(stack[0])

		if err != nil {
			return nil, err
		}

		return m, nil
	}

	c, err := ns.Namespace(stack[0])

	if err != nil {
		return nil, err
	}

	return findMethod(c, stack[1:])
}
