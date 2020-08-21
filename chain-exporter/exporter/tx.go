package exporter

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/chain-exporter/schema"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authx "github.com/cosmos/cosmos-sdk/x/auth"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	escrow "github.com/cosmos/ethermint/x/escrow/types"
	ethtypes "github.com/cosmos/ethermint/x/evm/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// getTxs parses transactions and wrap into Transaction schema struct
func (ex *Exporter) getTxs(block *tmctypes.ResultBlock) ([]*schema.Transaction, []*schema.Contract, error) {
	transactions := make([]*schema.Transaction, 0)

	txs, err := ex.client.Txs(block)
	if err != nil {
		return nil, nil, err
	}

	contracts := make([]*schema.Contract, 0)
	if len(txs) > 0 {
		for _, tx := range txs {
			var msgsBz, memo, fromAddress, toAddress, contractAddress string
			// 先解析 halle ，如果解析不出来，再按 eth交易类型交易格式解析
			halleTx, halleTxErr := ex.parseHalleTx(tx.Tx)
			if halleTxErr != nil {
				ethTx, ethTxErr := ex.parseEthTx(tx.Tx)
				if ethTxErr != nil {
					return nil, nil, err
				}
				from, _ := ethTx.VerifySig(ethTx.ChainID())
				fromAccAddr := sdk.AccAddress(from.Bytes())
				fromAddress = fromAccAddr.String()

				//to is empty when create contract
				var toAccAddr *sdk.AccAddress
				var amount big.Int
				amount = *ethTx.Data.Amount
				memo = hex.EncodeToString(ethTx.Data.Payload)
				to := ethTx.To()
				if to != nil {
					newAccAddr := sdk.AccAddress(ethTx.To().Bytes())
					toAccAddr = &newAccAddr
					toAddress = newAccAddr.String()

					if len(ethTx.Data.Payload) > 0 {
						isContract, _ := ex.client.IsContract(to)
						contractAddress = toAddress
						if isContract {
							var fromAccAddress *sdk.AccAddress
							var toAccAddress *sdk.AccAddress
							fromAddress, toAddress, fromAccAddress, toAccAddress = ex.processERC20(
								ethTx.Data.Payload, fromAddress, toAddress, &amount)
							if fromAccAddress != nil {
								fromAccAddr = *fromAccAddress
							}
							if toAccAddress != nil {
								toAccAddr = toAccAddress
							}
						}
					}
				} else {
					//合约部署 toAddress 为空
					contractAddr := crypto.CreateAddress(from, ethTx.Data.AccountNonce)
					newAccAddr := sdk.AccAddress(contractAddr.Bytes())
					contractAddress = newAccAddr.String()

					contract := &schema.Contract{
						Height:          tx.Height,
						TxHash:          tx.Hash.String(),
						FromAddress:     fromAddress,
						ContractAddress: contractAddress,
						Memo:            memo,
					}
					contract.TotalSupply, _ = ex.client.GetContractTotal(&contractAddr)
					contract.Decimals, _ = ex.client.GetContractDecimals(&contractAddr)
					contract.Name, _ = ex.client.GetContractName(&contractAddr)
					contract.Symbol, _ = ex.client.GetContractSymbol(&contractAddr)

					contracts = append(contracts, contract)
				}

				msgEther := ethtypes.NewMsgEthermint(ethTx.Data.AccountNonce, toAccAddr,
					sdk.NewIntFromBigInt(&amount), ethTx.Data.GasLimit,
					sdk.NewIntFromBigInt(ethTx.Data.Price), ethTx.Data.Payload, fromAccAddr)
				msgs := []sdk.Msg{msgEther}
				msgsBz0, err := ex.cdc.MarshalJSON(msgs) //ethTx.GetMsgs())
				if err != nil {
					return nil, nil, err
				}
				msgsBz = string(msgsBz0)

			} else {
				msgsBz0, err := ex.cdc.MarshalJSON(halleTx.GetMsgs())
				if err != nil {
					return nil, nil, err
				}
				msgsBz = string(msgsBz0)
				memo = halleTx.Memo

				fromAddress = halleTx.FeePayer().String()
				msgs := halleTx.GetMsgs()
				switch msgs[0].(type) {
				case bank.MsgSend:
					msg, ok := msgs[0].(bank.MsgSend)
					if ok {
						toAddress = msg.ToAddress.String()
					}
				case escrow.MsgSendWithUnlock:
					msg, ok := msgs[0].(escrow.MsgSendWithUnlock)
					if ok {
						toAddress = msg.ToAddress.String()
					}
				case escrow.MsgPayout:
					msg, ok := msgs[0].(escrow.MsgPayout)
					if ok {
						toAddress = msg.Receiver.String()
					}
				default:
					fmt.Println("not support Msg Type =================")
				}
			}
			fmt.Println("fromAddress", fromAddress)
			fmt.Println("toAddress", toAddress)

			tempTransaction := &schema.Transaction{
				Height:          tx.Height,
				TxHash:          tx.Hash.String(),
				FromAddress:     fromAddress,
				ToAddress:       toAddress,
				ContractAddress: contractAddress,
				Code:            tx.TxResult.Code, // 0 is success
				Messages:        msgsBz,
				Signatures:      string("{\"sigsBz\": \"unsolved\"}"),
				Memo:            memo,
				GasWanted:       tx.TxResult.GasWanted,
				GasUsed:         tx.TxResult.GasUsed,
				Timestamp:       block.Block.Time,
			}

			transactions = append(transactions, tempTransaction)
		}
	}

	return transactions, contracts, nil
}

func (ex *Exporter) processERC20(data []byte, from, to string, amount *big.Int) (fromAddress, toAddress string, fromAccAddr, toAccAddr *sdk.AccAddress) {
	abi := ex.client.GetABI()

	fromAddress = from
	toAddress = "not_erc20_method"

	methodHash := hex.EncodeToString(data[0:4])
	fmt.Println("==== data", methodHash, hex.EncodeToString(data))

	isApprove := false
	switch methodHash {
	case "095ea7b3": //095ea7b3 + approve(_spender address, _value uint256)
		//todo 额度授权，暂不处理，设置转账金额为0；将来，需要区分：不是转账代币，是额度
		isApprove = true
		fallthrough
	case "a9059cbb": //a9059cbb + transfer(_to address, _value uint256)
		values, err := abi.Methods["transfer"].Inputs.UnpackValues(data[4:])
		if err != nil {
			fmt.Println("ERC20ABI.UnpackValues() transfer failed : ", err.Error())
		} else {
			resultsJSON, _ := json.Marshal(values)
			//fmt.Println(string(resultsJSON))
			var to common.Address
			result := []interface{}{&to, amount}
			json.Unmarshal(resultsJSON, &result)
			//fmt.Println(to.String(), amount.String())
			if isApprove {
				amount.SetInt64(int64(0))
			}

			newAccAddr := sdk.AccAddress(to.Bytes())
			toAccAddr = &newAccAddr
			toAddress = toAccAddr.String()
		}
	case "23b872dd": //23b872dd + transferFrom(_from address, _to address, _value uint256)
		values, err := abi.Methods["transferFrom"].Inputs.UnpackValues(data[4:])
		if err != nil {
			fmt.Println("ERC20ABI.UnpackValues() transferFrom failed : ", err.Error())
		} else {
			resultsJSON, _ := json.Marshal(values)
			//fmt.Println(string(resultsJSON))
			var fromAddr common.Address
			var to common.Address
			result := []interface{}{&fromAddr, &to, amount}
			json.Unmarshal(resultsJSON, &result)
			fmt.Println(to.String(), amount.String())

			newAccAddrFrom := sdk.AccAddress(to.Bytes())
			fromAccAddr = &newAccAddrFrom
			fromAddress = fromAccAddr.String() //todo 当前from地址是保存的代币转出方；将来，建议存两条交易，交易发起者也存一条

			newAccAddr := sdk.AccAddress(to.Bytes())
			toAccAddr = &newAccAddr
			toAddress = toAccAddr.String()
		}
	} // switch
	return fromAddress, toAddress, fromAccAddr, toAccAddr
}

func (ex *Exporter) parseHalleTx(tx []byte) (*authx.StdTx, error) {
	var authstdTx authx.StdTx
	err := ex.cdc.UnmarshalBinaryBare(tx, &authstdTx)
	if err != nil {
		return nil, err
	}
	return &authstdTx, nil
}

func (ex *Exporter) parseEthTx(decodeBytes []byte) (*ethtypes.MsgEthereumTx, error) {
	//codec1 := codec.New()
	//sdk.RegisterCodec(codec1)
	//codec1.RegisterConcrete(ethtypes.MsgEthereumTx{}, "ethermint/MsgEthereumTx", nil)
	//codec1.RegisterConcrete(ethtypes.MsgEthermint{}, "ethermint/MsgEthermint", nil)
	//codec1.RegisterConcrete(ethtypes.TxData{}, "ethermint/TxData", nil)
	//codec.RegisterCrypto(codec1)
	//codec1.Seal()

	var ethTxMsg ethtypes.MsgEthereumTx
	err := ex.cdc.UnmarshalBinaryBare(decodeBytes, &ethTxMsg)
	if err != nil {
		fmt.Println("---------parse ethtx err-----")
		fmt.Println(err)
		return nil, err
	} else {
		fmt.Println("---------ethTxMsg-----")
		fmt.Println(ethTxMsg)
		return &ethTxMsg, nil
	}
}
