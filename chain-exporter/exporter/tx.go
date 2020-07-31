package exporter

import (
	"fmt"

	"github.com/chain-exporter/schema"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authx "github.com/cosmos/cosmos-sdk/x/auth"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	escrow "github.com/cosmos/ethermint/x/escrow/types"
	ethtypes "github.com/cosmos/ethermint/x/evm/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getTxs parses transactions and wrap into Transaction schema struct
func (ex *Exporter) getTxs(block *tmctypes.ResultBlock) ([]*schema.Transaction, error) {
	transactions := make([]*schema.Transaction, 0)

	txs, err := ex.client.Txs(block)
	if err != nil {
		return nil, err
	}

	if len(txs) > 0 {
		for _, tx := range txs {
			var msgsBz, memo, fromAddress, toAddress string
			// 先解析 halle ，如果解析不出来，再按 eth交易类型交易格式解析
			halleTx, halleTxErr := ex.parseHalleTx(tx.Tx)
			if halleTxErr != nil {
				ethTx, ethTxErr := ex.parseEthTx(tx.Tx)
				if ethTxErr != nil {
					return nil, err
				}
				from, _ := ethTx.VerifySig(ethTx.ChainID())
				fromAddress = from.String()

				fromAccAddr := sdk.AccAddress(from.Bytes())
				toAccAddr := sdk.AccAddress(ethTx.To().Bytes())
				msgEther := ethtypes.NewMsgEthermint(ethTx.Data.AccountNonce, &toAccAddr,
					sdk.NewIntFromBigInt(ethTx.Data.Amount), ethTx.Data.GasLimit,
					sdk.NewIntFromBigInt(ethTx.Data.Price), ethTx.Data.Payload, fromAccAddr)
				msgsBz0, err := ex.cdc.MarshalJSON(msgEther) //ethTx.GetMsgs())
				if err != nil {
					return nil, err
				}
				msgsBz = string(msgsBz0)
				memo = string(ethTx.Data.Payload)

				toAddress = ethTx.To().String()
			} else {
				msgsBz0, err := ex.cdc.MarshalJSON(halleTx.GetMsgs())
				if err != nil {
					return nil, err
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
				Height:      tx.Height,
				TxHash:      tx.Hash.String(),
				FromAddress: fromAddress,
				ToAddress:   toAddress,
				Code:        tx.TxResult.Code, // 0 is success
				Messages:    msgsBz,
				Signatures:  string("{\"sigsBz\": \"unsolved\"}"),
				Memo:        memo,
				GasWanted:   tx.TxResult.GasWanted,
				GasUsed:     tx.TxResult.GasUsed,
				Timestamp:   block.Block.Time,
			}

			transactions = append(transactions, tempTransaction)
		}
	}

	return transactions, nil
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
