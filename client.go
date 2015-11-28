package jsonclient

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/gorilla/rpc/v2/json2"
)

// NewClient returns a new Client usable
func NewClient(address string) Client {
	return Client{addr: address}
}

// Client is a JSON-RPC client.
type Client struct {
	addr string
}

// clientRequest represents a JSON-RPC request sent by a client.
type clientRequest struct {
	// JSON-RPC protocol.
	Version string `json:"jsonrpc"`

	// A String containing the name of the method to be invoked.
	Method string `json:"method"`

	// Object to pass as request parameter to the method.
	Params interface{} `json:"params"`

	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	ID uint32 `json:"id"`
}

// EncodeClientRequest encodes parameters for a JSON-RPC client request with
// an upper bound on the ID of (^uint32(0) - 1) / 2.
func EncodeClientRequest(method string, args interface{}) ([]byte, error) {
	c := &clientRequest{
		Version: "2.0",
		Method:  method,
		Params:  args,
		ID:      rand.Uint32() % ((^uint32(0) - 1) / 2),
	}
	return json.Marshal(c)
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	buf, err := EncodeClientRequest(serviceMethod, args)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.addr, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return json2.DecodeClientResponse(resp.Body, reply)
}
