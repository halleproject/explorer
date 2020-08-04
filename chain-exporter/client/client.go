package client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Robbin-Liu/go-binance-sdk/client/rpc"

	//"github.com/chain-exporter/codec"
	//"github.com/chain-exporter/config"
	//"github.com/chain-exporter/types"

	"github.com/chain-exporter/codec"
	"github.com/chain-exporter/config"
	"github.com/chain-exporter/types"

	amino "github.com/tendermint/go-amino"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	resty "github.com/go-resty/resty/v2"
)

// Client wraps around both Tendermint RPC client and
// Cosmos SDK LCD REST client that enables to query necessary data
type Client struct {
	acceleratedClient *resty.Client
	apiClient         *resty.Client
	cdc               *amino.Codec
	explorerClient    *resty.Client
	rpcClient         rpc.Client
}

// NewClient creates a new client with the given config
func NewClient(cfg config.NodeConfig) *Client {

	acceleratedClient := resty.New().
		SetHostURL(cfg.AcceleratedNode).
		SetTimeout(time.Duration(5 * time.Second))

	apiClient := resty.New().
		SetHostURL(cfg.APIServerEndpoint).
		SetTimeout(time.Duration(5 * time.Second))

	explorerClient := resty.New().
		SetHostURL(cfg.ExplorerServerEndpoint).
		SetTimeout(time.Duration(30 * time.Second))

	rpcClient := rpc.NewRPCClient(cfg.RPCNode, cfg.NetworkType)

	return &Client{
		acceleratedClient,
		apiClient,
		codec.Codec,
		explorerClient,
		rpcClient,
	}
}

// Block queries for a block by height. An error is returned if the query fails.
func (c Client) Block(height int64) (*tmctypes.ResultBlock, error) {
	return c.rpcClient.Block(&height)
}

// LatestBlockHeight returns the latest block height on the active chain
func (c Client) LatestBlockHeight() (int64, error) {
	status, err := c.rpcClient.Status()
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight

	return height, nil
}

// Txs queries for all the transactions in a block height.
// It uses `Tx` RPC method to query for the transaction
func (c Client) Txs(block *tmctypes.ResultBlock) ([]*rpc.ResultTx, error) {
	txs := make([]*rpc.ResultTx, len(block.Block.Txs), len(block.Block.Txs))

	for i, tmTx := range block.Block.Txs {
		tx, err := c.rpcClient.Tx(tmTx.Hash(), true)
		if err != nil {
			return nil, err
		}
		fmt.Println("-----rpcClient.Tx()--------")
		fmt.Println(tx)
		txs[i] = tx
	}

	return txs, nil
}

// ValidatorSet returns all the known Tendermint validators for a given block
// height. An error is returned if the query fails.
func (c Client) ValidatorSet(height int64) (*tmctypes.ResultValidators, error) {
	return c.rpcClient.Validators(&height)
}

// Validators returns validators detail information in Tendemrint validators in active chain
// An error returns if the query fails.
func (c Client) Validators() ([]*types.Validator, error) {
	resp, err := c.apiClient.R().Get("/staking/validators")
	//resp, err := c.apiClient.R().Get("/validators")
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nResponse Body: %v", resp)

	///var vals []*types.Validator
	var vals *types.HttpBody

	err = json.Unmarshal(resp.Body(), &vals)
	if err != nil {
		return nil, err
	}

	return vals.Validator, nil
	//return nil,nil
}

// Tokens returns information about existing tokens in active chain
func (c Client) Tokens(limit int, offset int) ([]*types.Token, error) {
	resp, err := c.apiClient.R().Get("/tokens?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}

	var tokens []*types.Token
	err = json.Unmarshal(resp.Body(), &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
