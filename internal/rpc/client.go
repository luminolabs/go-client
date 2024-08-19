package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Client struct {
	rpc *rpc.Client
}

func NewClient(url string) (*Client, error) {
	rpc, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Client{rpc: rpc}, nil
}

func (c *Client) Call(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return c.rpc.CallContext(ctx, result, method, args...)
}

func (c *Client) Close() {
	c.rpc.Close()
}
