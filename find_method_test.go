package jsonrpc

import (
	"errors"
	"testing"
)

var callbackMock Callback = func(in []byte) (out interface{}, err Error) {
	return nil, nil
}

var findMethodTestCases = []struct {
	stack []string
	err   error
}{
	{
		stack: []string{"ns1", "ns2", "sum"},
		err:   nil,
	},
	{
		stack: []string{"mul"},
		err:   nil,
	},
	{
		stack: []string{"mock"},
		err:   nil,
	},
	{
		stack: []string{"root_non_existing"},
		err:   errMethodNotFound{},
	},
	{
		stack: []string{"ns1", "non_existing_ns", "non_existing_method"},
		err:   errNamespaceNotFound{},
	},
	{
		stack: []string{},
		err:   errMethodNotFound{},
	},
}

func Test_findMethod(t *testing.T) {
	ns := createTestRpc()
	ns.Register("mock", callbackMock)

	for _, c := range findMethodTestCases {
		callback, err := findMethod(ns.Namespace, c.stack)

		var errMethodNotFound errMethodNotFound
		if errors.As(c.err, &errMethodNotFound) && !errors.As(err, &errMethodNotFound) {
			t.Errorf("errMethodNotFound error expected, %T given", err)
		}

		var errNsNotFound errNamespaceNotFound
		if errors.As(c.err, &errNsNotFound) && !errors.As(err, &errNsNotFound) {
			t.Errorf("errNamespaceNotFound error expected, %T given", err)
		}

		if err != nil && c.err == nil {
			t.Errorf("no error expected, given %T: %v", err, err)
		}

		if err == nil && callback == nil {
			t.Error("callback expected to be a *Callback type, nil returned")
		}
	}
}
