package jsonrpc

import "encoding/json"

// Request represents JSON-RPC request
type Request struct {
	JsonRpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}
