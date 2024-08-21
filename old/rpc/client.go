package rpcclient

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Client struct {
	rpcClient *rpc.Client
}

func NewClient(url string) (*Client, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Client{rpcClient: rpcClient}, nil
}

func (c *Client) Call(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.CallContext(ctx, result, method, args...)
}

func (c *Client) Subscribe(ctx context.Context, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error) {
	return c.rpcClient.Subscribe(ctx, namespace, channel, args...)
}

func (c *Client) Close() {
	c.rpcClient.Close()
}
