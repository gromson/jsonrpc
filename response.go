package jsonrpc

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const internalServerErrorReason = "Internal Server Error"

// Response represents JSON-RPC response
type Response struct {
	JsonRpc string          `json:"jsonrpc"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
	ID      interface{}     `json:"id"`
	Log     LogRespondError `json:"-"`
}

// ResponseError represents JSON-RPC response error
type ResponseError struct {
	Code           int      `json:"code"`
	Message        string   `json:"message"`
	Data           []string `json:"data,omitempty"`
	httpStatusCode int
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
		ID:  id,
		Log: logError,
	}
}

func (r *Response) SetHttpStatusCode(code int) {
	if r.Error == nil {
		return
	}

	r.Error.httpStatusCode = code
}

func (r *Response) Respond(w http.ResponseWriter) {
	httpStatusCode := getHttpStatusCode(r)

	if httpStatusCode >= http.StatusInternalServerError {
		r.handleInternalServerError()
	}

	serialized, err := json.Marshal(r)
	if err != nil {
		r.Log("Couldn't serialize JSON-RPC response", err, nil)
		http.Error(w, internalServerErrorReason, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatusCode)

	_, err = w.Write(serialized)
	if err != nil {
		r.Log("Error while trying to write json response body", err, nil)
		http.Error(w, internalServerErrorReason, http.StatusInternalServerError)
		return
	}
}

func (r *Response) handleInternalServerError() {
	if r.Error == nil {
		r.Log(internalServerErrorReason, errors.New("unknown error"), nil)
		return
	}

	r.Log(internalServerErrorReason, errors.New(r.Error.Message), r.Error.Data)

	r.Error.Message = internalServerErrorReason
	r.Error.Data = nil
}

func getHttpStatusCode(r *Response) int {
	if r.Error == nil {
		return http.StatusOK
	}

	code := http.StatusInternalServerError

	if r.Error.httpStatusCode != 0 {
		code = r.Error.httpStatusCode
	}

	return code
}

type LogRespondError func(title string, err error, additional interface{})

func logError(title string, err error, additional interface{}) {
	entry := log.WithError(err)

	if additional != nil {
		entry = entry.WithField("additional", additional)
	}

	entry.Error(title)
}
