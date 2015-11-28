package jsonclient

import (
	"bytes"
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

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	buf, err := json2.EncodeClientRequest(serviceMethod, args)
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
