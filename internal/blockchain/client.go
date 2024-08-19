package blockchain

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/luminolabs/go-client/internal/contractsabi"
	"github.com/luminolabs/go-client/internal/rpcclient"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Client struct {
	rpcClient    *rpcclient.Client
	contractsABI *contractsabi.Manager
	ethClient    *ethclient.Client
}

func NewClient(rpcURL string) (*Client, error) {
	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	rpcClient, err := rpcclient.NewClient(rpcURL)
	if err != nil {
		return nil, err
	}

	contractsABI := contractsabi.NewManager()

	return &Client{
		rpcClient:    rpcClient,
		contractsABI: contractsABI,
		ethClient:    ethClient,
	}, nil
}

func (c *Client) Close() {
	c.ethClient.Close()
	c.rpcClient.Close()
}
