package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/luminolabs/lumino-go-client/internal/contractsabi"
	"github.com/luminolabs/lumino-go-client/internal/rpcclient"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Client struct {
	rpcClient    *rpcclient.Client
	contractsABI *contractsabi.Manager
	ethClient    *ethclient.Client
	stakeManager *bind.BoundContract
	jobsManager  *bind.BoundContract
	blockManager *bind.BoundContract
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

	client := &Client{
		rpcClient:    rpcClient,
		contractsABI: contractsABI,
		ethClient:    ethClient,
	}

	// Initialize contract bindings
	if err := client.initializeContracts(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) initializeContracts() error {
	var err error

	stakeManagerABI, err := c.contractsABI.GetABI("StakeManager")
	if err != nil {
		return err
	}
	c.stakeManager, err = bind.NewBoundContract(common.HexToAddress("stake_manager_address"), *stakeManagerABI, c.ethClient, c.ethClient, c.ethClient)
	if err != nil {
		return err
	}

	// Initialize other contracts similarly...

	return nil
}

func (c *Client) Stake(ctx context.Context, amount *big.Int) (*types.Transaction, error) {
	opts, err := c.getTransactOpts(ctx)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	err = c.stakeManager.Transact(opts, &tx, "stake", amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (c *Client) Unstake(ctx context.Context, amount *big.Int) (*types.Transaction, error) {
	opts, err := c.getTransactOpts(ctx)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	err = c.stakeManager.Transact(opts, &tx, "unstake", amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (c *Client) getTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	// Implement this method to create transaction options
	// This should include setting the from address, gas price, etc.
	return nil, nil
}

func (c *Client) Close() {
	c.ethClient.Close()
	c.rpcClient.Close()
}

func (c *Client) EthClient() *ethclient.Client {
	return c.ethClient
}
