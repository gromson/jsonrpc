package jsonrpc

import (
	"encoding/json"
	"reflect"
	"testing"
)

const wrongParametersTestMessage = "Wrong parameters given"

var sumCallback Callback = func(in []byte) (interface{}, Error) {
	input := make([]int, 0)
	err := json.Unmarshal(in, &input)

	if err != nil {
		return nil, ErrorInvalidParams
	}

	r := 0

	for _, v := range input {
		r += v
	}

	return r, nil
}

var mulCallback Callback = func(in []byte) (interface{}, Error) {
	input := make([]int, 0)
	err := json.Unmarshal(in, &input)

	if err != nil {
		return nil, ErrorInvalidParams
	}

	r := 1

	for _, v := range input {
		r *= v
	}

	return r, nil
}

var wrongParamsCallback Callback = func(in []byte) (interface{}, Error) {
	return nil, NewError(wrongParametersTestMessage, ErrCodeInvalidParams)
}

func createTestRpc() *Rpc {
	rpc := NewRpc()
	ns1 := NewNamespace()
	ns2 := NewNamespace()

	ns1.RegisterNS("ns2", ns2)
	ns1.Register("wrongParams", wrongParamsCallback)

	ns2.Register("sum", sumCallback)

	rpc.RegisterNS("ns1", ns1)
	rpc.Register("mul", mulCallback)

	return rpc
}

func TestRpc_Exec(t *testing.T) {
	rpc := createTestRpc()

	rpcExecSumSuccess(rpc, t)
	rpcExecMulSuccess(rpc, t)
	rpcExecMethodNotFound(rpc, t)
	rpcExecWrongParams(rpc, t)
}

func rpcExecSumSuccess(rpc *Rpc, t *testing.T) {
	req := &Request{
		JsonRpc: version,
		Method:  "ns1.ns2.sum",
		Params:  json.RawMessage("[1, 2, 3, 4]"), // []int{1, 2, 3, 4},
		ID:      1,
	}

	res := rpc.Exec(req)

	expectedRes := Response{
		JsonRpc: version,
		Result:  10,
		Error:   nil,
		ID:      1,
	}

	if !reflect.DeepEqual(*res, expectedRes) {
		t.Errorf("expected sum result %+v, %+v given", *res, expectedRes)
	}
}

func rpcExecMulSuccess(rpc *Rpc, t *testing.T) {
	req := &Request{
		JsonRpc: version,
		Method:  "mul",
		Params:  json.RawMessage("[1, 2, 3, 4]"), // []int{1, 2, 3, 4},
		ID:      1,
	}

	res := rpc.Exec(req)

	expectedRes := Response{
		JsonRpc: version,
		Result:  24,
		Error:   nil,
		ID:      1,
	}

	if !reflect.DeepEqual(*res, expectedRes) {
		t.Errorf("expected mul result %+v, %+v given", *res, expectedRes)
	}
}

func rpcExecMethodNotFound(rpc *Rpc, t *testing.T) {
	req := &Request{
		JsonRpc: version,
		Method:  "ns1.mul",
		Params:  json.RawMessage("[1, 2]"), // []int{1, 2, 3, 4},
		ID:      1,
	}

	res := rpc.Exec(req)

	expectedRes := Response{
		JsonRpc: version,
		Result:  nil,
		Error: &ResponseError{
			Code:    ErrCodeMethodNotFound,
			Message: "method \"mul\" not found",
			Data:    nil,
		},
		ID: 1,
	}

	if !reflect.DeepEqual(*res, expectedRes) {
		t.Errorf("expected mul result %+v, %+v given", *res, expectedRes)
	}
}

func rpcExecWrongParams(rpc *Rpc, t *testing.T) {
	req := &Request{
		JsonRpc: version,
		Method:  "ns1.wrongParams",
		Params:  nil,
		ID:      1,
	}

	res := rpc.Exec(req)

	expectedRes := Response{
		JsonRpc: version,
		Result:  nil,
		Error: &ResponseError{
			Code:    ErrCodeInvalidParams,
			Message: wrongParametersTestMessage,
			Data:    nil,
		},
		ID: 1,
	}

	if !reflect.DeepEqual(*res, expectedRes) {
		t.Errorf("expected wrongParams result %+v, %+v given", *res, expectedRes)
	}
}
