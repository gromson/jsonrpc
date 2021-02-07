package jsonrpc

// Response represents JSON-RPC response
type Response struct {
	JsonRpc string         `json:"jsonrpc"`
	Result  interface{}    `json:"result,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
	ID      interface{}    `json:"id"`
}

// ResponseError represents JSON-RPC response error
type ResponseError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    []string `json:"data,omitempty"`
}

// NewErrResponse creates JSON-RPC response
func NewErrResponse(id interface{}, rpcErr Error) *Response {
	return &Response{
		JsonRpc: version,
		Result:  nil,
		Error: &ResponseError{
			Code:    rpcErr.Code(),
			Message: rpcErr.Error(),
			Data:    nil,
		},
		ID: id,
	}
}
