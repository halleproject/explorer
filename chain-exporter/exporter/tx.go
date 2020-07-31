package exporter

import (
	"encoding/base64"
	"fmt"
	//"github.com/Robbin-Liu/go-binance-sdk/types/tx"
	"github.com/chain-exporter/schema"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	//tmtypes "github.com/tendermint/tendermint/types"
	//"math/big"

	//sdktypes "github.com/cosmos/cosmos-sdk/types"
	authx "github.com/cosmos/cosmos-sdk/x/auth"
	ethtypes "github.com/cosmos/ethermint/x/evm/types"
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

			decodeBytes, err01 := base64.StdEncoding.DecodeString(tx.Tx.String())
			if err01 != nil {
				fmt.Printf("base64.StdEncoding.DecodeString(encodeString) err:  %v \n", decodeBytes)
			}
			fmt.Printf("----base64-----------")
			fmt.Println(decodeBytes)


			/* 测试了几种方式
			var ptr codec.ProtoMarshaler
			ex.cdc.UnmarshalBinaryLengthPrefixed([]byte(tx.Tx), &ptr)
			fmt.Printf("----test -ProtoMarshaler----------")
			fmt.Println(ptr)

			var stdTx txtypes.StdTx
			ex.cdc.UnmarshalBinaryLengthPrefixed([]byte(tx.Tx), &stdTx)

			var authstdTx auth.StdTx
			//ex.cdc.UnmarshalBinaryLengthPrefixed([]byte(tx.Tx), &stdTx)
			ex.cdc.UnmarshalBinaryBare(decodeBytes, &authstdTx)

			var authstdTx2 auth.StdTx
			//ex.cdc.UnmarshalBinaryLengthPrefixed([]byte(tx.Tx), &stdTx)
			ex.cdc.UnmarshalBinaryLengthPrefixed([]byte(tx.Tx), &authstdTx2)


			var authstdTx3 auth.StdTx
			//ex.cdc.UnmarshalBinaryLengthPrefixed([]byte(tx.Tx), &stdTx)
			ex.cdc.UnmarshalBinaryBare([]byte(tx.Tx), &authstdTx3)

			fmt.Printf("----test -UnmarshalBinaryBare----------")
			fmt.Println(authstdTx3)
			*/


			var msgsBz , memo   ,fromAddress ,toAddress  string
			// 先解析 halle ，如果解析不出来，再按 eth交易类型交易格式解析
			halleTx, halleTxErr :=  parseHalleTx(ex,decodeBytes )
			if halleTxErr != nil {
				ethTx, ethTxErr :=  parseEthTx(decodeBytes )
				if ethTxErr != nil {
					return nil, err
				}
				msgsBz0, err := ex.cdc.MarshalJSON(ethTx.GetMsgs())
				if err != nil {
					return nil, err
				}
				msgsBz = string(msgsBz0)
				//	memo, _ = ethTx.Data.Payload
			} else {
				var authstdTx  = halleTx
				fmt.Printf("----test -UnmarshalBinaryBare----------")
				fmt.Println(authstdTx)
				fmt.Printf("----authstdTx -GetMsgs----------")
				fmt.Println(authstdTx.GetMsgs())

				fmt.Printf("----authstdTx -GetSignatures----------")
				fmt.Println(authstdTx.GetSignatures())

				//fromAddress = authstdTx.GetMsgs()


				msgsBz0, err := ex.cdc.MarshalJSON(authstdTx.GetMsgs())
				if err != nil {
					return nil, err
				}
				msgsBz = string(msgsBz0)
				memo = authstdTx.Memo

			}

				/*
				if ethTxErr != nil {
					var authstdTx authx.StdTx
					//ex.cdc.UnmarshalBinaryLengthPrefixed([]byte(tx.Tx), &stdTx)
					ex.cdc.UnmarshalBinaryBare([]byte(tx.Tx), &authstdTx)

					fmt.Printf("----test -UnmarshalBinaryBare----------")
					fmt.Println(authstdTx)
					fmt.Printf("----authstdTx -GetMsgs----------")
					fmt.Println(authstdTx.GetMsgs())

					fmt.Printf("----authstdTx -GetSignatures----------")
					fmt.Println(authstdTx.GetSignatures())

					//fromAddress = authstdTx.GetMsgs()


					msgsBz0, err := ex.cdc.MarshalJSON(authstdTx.GetMsgs())
					if err != nil {
						return nil, err
					}
					msgsBz = string(msgsBz0)
					memo = authstdTx.Memo
				}else {
					msgsBz0, err := ex.cdc.MarshalJSON(ethTx.GetMsgs())
					if err != nil {
						return nil, err
					}
					msgsBz = string(msgsBz0)
				//	memo, _ = ethTx.Data.Payload
				}
				*/


			//pubkeys := authstdTx.GetPubKeys()
			//for _, pubkey := range pubkeys {
			//	fmt.Println(pubkey)
			//}

			/* @TODO 待处理

			sigs := make([]types.Signature, len(authstdTx.Signatures), len(authstdTx.Signatures))
			for i, sig := range authstdTx.Signatures {
				//stdSignature  :=sig.Signature;
				//fmt.Println(stdSignature)
				//
				//consPubKey01, err01 := bech32.ConvertAndEncode( sdktypes.GetConfig().GetBech32ConsensusPubPrefix(), stdSignature)
				//if err01 != nil {
				//	return nil, err01
				//}
				//fmt.Println("---------consPubKey01-----------------")
				//fmt.Println(consPubKey01)
				//
				//
				//var pk crypto.PubKey
				////ex.cdc.MustUnmarshalBinaryBare(stdSignature, &pk)
				//ex.cdc.MustUnmarshalBinaryBare(decodeBytes, &pk)
				//fmt.Println(pk)


				consPubKey, err := sdktypes.Bech32ifyPubKey(sdktypes.Bech32PubKeyTypeConsPub,sig.GetPubKey())
				if err != nil {
					return nil, err
				}

				sigs[i] = types.Signature{
					Address:       sig.GetPubKey().Address().String(), // hex string
					AccountNumber: int64(sig.Size()),
					Pubkey:        consPubKey,
					Sequence:      int64(tx.TxResult.Code),
					Signature:     base64.StdEncoding.EncodeToString(sig.Signature), // encode base64
				}

			}

			/* @TODO 待处理
			sigsBz, err := ex.cdc.MarshalJSON(sigs)
			if err != nil {
				return nil, err
			}
			*/
			tempTransaction := &schema.Transaction{
				Height:     tx.Height,
				TxHash:     tx.Hash.String(),
				FromAddress:     fromAddress,
				ToAddress:     toAddress,
				Code:       tx.TxResult.Code, // 0 is success
				//Messages:   string(msgsBz),
				Messages:   msgsBz,
				//Signatures: string(sigsBz),
				Signatures: string("{\"sigsBz\": \"unsolved\"}"),
				//Memo:       authstdTx.Memo,
				Memo:       memo,
				GasWanted:  tx.TxResult.GasWanted,
				GasUsed:    tx.TxResult.GasUsed,
				Timestamp:  block.Block.Time,
			}

			transactions = append(transactions, tempTransaction)

		}
	}

	return transactions, nil
}

func parseHalleTx(ex *Exporter ,tx []byte ) ( *authx.StdTx, error) {
	var authstdTx authx.StdTx
	//ex.cdc.UnmarshalBinaryLengthPrefixed([]byte(tx.Tx), &stdTx)
	err :=  ex.cdc.UnmarshalBinaryBare(tx, &authstdTx)
	if err!= nil {
		return nil, err
	}
	return &authstdTx,nil
}



func parseEthTx( decodeBytes []byte ) ( *ethtypes.MsgEthereumTx, error) {

	//var ethTx = "Jaa+VAreAQgWEgwxMDEwMDAwMDAwMDAYjfEDIhT2hWaXd8LnDZE3KBXkPhlUfl23BSoQMTAwMDAwMDAwMDAwMDAwMDoCNTJCTTg1OTk4NzQ1MDQ3MjY0MzI5MjAyODQ2NDYzMzk4NDEyNjQ1MTc4ODk5ODM5NTgwMjU4Mjk2NzMzMjQxMzc1MDA1NTc4OTgxMDU2NzExSk0yODIwNTE5NDYzMzk3MzYxMzk2MjkzNDM3NTYyNzYyMTM0NTU2MTgzMjM4NDIzMzA0MDUwOTQ1MjYxOTE3ODY5NjYzMTU0NTMzMTMyMA=="


	//decodeBytes, err01 := base64.StdEncoding.DecodeString(ethTx)
	//if err01 != nil {
	//	fmt.Printf("base64.StdEncoding.DecodeString(encodeString) err:  %v \n", decodeBytes)
	//}
	//fmt.Println(decodeBytes)


	codec1 := codec.New()
	// register Tx, Msg
	sdk.RegisterCodec(codec1)

	// register test types
	//cdc.RegisterConcrete(&txTest{}, "cosmos-sdk/baseapp/txTest", nil)
	//codec.RegisterConcrete(ethtypes.MsgEthermint{}, " ", nil)

	//codec.RegisterConcrete(bank.MsgSend{}, "cosmos-sdk/MsgSend", nil)
	//codec.RegisterConcrete(bank.MsgMultiSend{}, "cosmos-sdk/MsgMultiSend", nil)
	//codec.RegisterConcrete(auth.StdTx{}, "cosmos-sdk/StdTx", nil)
	codec1.RegisterConcrete(ethtypes.MsgEthereumTx{}, "ethermint/MsgEthereumTx", nil)
	codec1.RegisterConcrete(ethtypes.MsgEthermint{}, "ethermint/MsgEthermint", nil)
	codec1.RegisterConcrete(ethtypes.TxData{}, "ethermint/TxData", nil)

	codec.RegisterCrypto(codec1)
	codec1.Seal()

	var ethTxMsg ethtypes.MsgEthereumTx
	err := codec1.UnmarshalBinaryBare(decodeBytes, &ethTxMsg)
	if err != nil {
		fmt.Println("---------parse ethtx err-----");
		fmt.Println(err)
		return nil,err
	}else {
		fmt.Println("---------ethTxMsg-----");
		fmt.Println(ethTxMsg)
		return &ethTxMsg,nil;
	}
	//address := ethTxMsg.To()
	//fmt.Println("---------to address-----");
	//fmt.Println(address.String())
	//
	//intChainID, _ := new(big.Int).SetString("8", 10)
	//addressFrom, _ :=ethTxMsg.VerifySig(intChainID)
	//fmt.Println("---------from address-----");
	//fmt.Println(addressFrom.String())

}