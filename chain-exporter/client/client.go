package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/Robbin-Liu/go-binance-sdk/client/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethcli "github.com/ethereum/go-ethereum/ethclient"

	//"github.com/chain-exporter/codec"
	//"github.com/chain-exporter/config"
	//"github.com/chain-exporter/types"

	amino "github.com/tendermint/go-amino"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/chain-exporter/codec"
	"github.com/chain-exporter/config"
	"github.com/chain-exporter/types"

	resty "github.com/go-resty/resty/v2"
)

const ERC20RawABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_initialAmount","type":"uint256"},{"name":"_tokenName","type":"string"},{"name":"_decimalUnits","type":"uint8"},{"name":"_tokenSymbol","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"_spender","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Approval","type":"event"}]`

var ERC20ABI *abi.ABI

func init() {
	ABI20, err := abi.JSON(strings.NewReader(ERC20RawABI))
	if err != nil {
		fmt.Println("init ABI20 failed : ", err.Error())
	}
	ERC20ABI = &ABI20
}

// Client wraps around both Tendermint RPC client and
// Cosmos SDK LCD REST client that enables to query necessary data
type Client struct {
	acceleratedClient *resty.Client
	apiClient         *resty.Client
	cdc               *amino.Codec
	explorerClient    *resty.Client
	rpcClient         rpc.Client
	ethClient         *ethcli.Client
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
	ethClient, err := ethcli.Dial(cfg.APIServerEndpoint)
	if err != nil {
		fmt.Println(err.Error())
	}

	return &Client{
		acceleratedClient,
		apiClient,
		codec.Codec,
		explorerClient,
		rpcClient,
		ethClient,
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

func (c Client) IsContract(address *common.Address) (bool, error) {
	//addr := common.BytesToAddress(address)
	codeBytes, err := c.ethClient.CodeAt(context.Background(), *address, nil)
	if err != nil {
		return false, err
	}

	if len(codeBytes) > 0 {
		return true, nil
	}

	return false, nil
}

func (c Client) GetABI() *abi.ABI {
	return ERC20ABI
}

func (c Client) GetContractTotal(address *common.Address) (string, error) {
	msg := ethereum.CallMsg{From: *address, To: address, Data: []byte{0x18, 0x16, 0x0d, 0xdd}} //0x18160ddd
	output, err := c.ethClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	result := new(big.Int)
	err = ERC20ABI.Unpack(&result, "totalSupply", output)
	if err != nil {
		return "", err
	}
	fmt.Println("==== result ", result)

	return result.String(), nil
}

func (c Client) GetContractDecimals(address *common.Address) (uint32, error) {
	msg := ethereum.CallMsg{From: *address, To: address, Data: []byte{0x31, 0x3c, 0xe5, 0x67}} //0x313ce567
	output, err := c.ethClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, err
	}

	result := new(uint8)
	err = ERC20ABI.Unpack(&result, "decimals", output)
	if err != nil {
		return 0, err
	}
	fmt.Println("==== result ", *result)

	return uint32(*result), nil
}

func (c Client) GetContractName(address *common.Address) (string, error) {
	msg := ethereum.CallMsg{From: *address, To: address, Data: []byte{0x06, 0xfd, 0xde, 0x03}} //0x06fdde03
	output, err := c.ethClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	var result string
	err = ERC20ABI.Unpack(&result, "name", output)
	if err != nil {
		return "", err
	}
	fmt.Println("==== result ", result)

	return result, nil
}

func (c Client) GetContractSymbol(address *common.Address) (string, error) {
	msg := ethereum.CallMsg{From: *address, To: address, Data: []byte{0x95, 0xd8, 0x9b, 0x41}} //0x95d89b41
	output, err := c.ethClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	var result string
	err = ERC20ABI.Unpack(&result, "symbol", output)
	if err != nil {
		return "", err
	}
	fmt.Println("==== result ", result)

	return result, nil
}
