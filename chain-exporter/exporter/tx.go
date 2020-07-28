package exporter

import (
	"encoding/base64"
	"fmt"
	"github.com/chain-exporter/schema"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	//sdktypes "github.com/cosmos/cosmos-sdk/types"
	authx "github.com/cosmos/cosmos-sdk/x/auth"
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


			fmt.Println("----tx.TxResult-----------")
			fmt.Println(tx.TxResult)
			fmt.Println("----tx.Tx-----------")
			fmt.Println(tx.Tx)



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

			var authstdTx authx.StdTx
			//ex.cdc.UnmarshalBinaryLengthPrefixed([]byte(tx.Tx), &stdTx)
			ex.cdc.UnmarshalBinaryBare([]byte(tx.Tx), &authstdTx)

			fmt.Printf("----test -UnmarshalBinaryBare----------")
			fmt.Println(authstdTx)
			fmt.Printf("----authstdTx -GetMsgs----------")
			fmt.Println(authstdTx.GetMsgs())

			fmt.Printf("----authstdTx -GetSignatures----------")
			fmt.Println(authstdTx.GetSignatures())


			msgsBz, err := ex.cdc.MarshalJSON(authstdTx.GetMsgs())
			if err != nil {
				return nil, err
			}

			//sigs := make([]types.Signature, len(stdTx.Signatures), len(stdTx.Signatures))

			//for i, sig := range stdTx.Signatures {
			//	consPubKey, err := ctypes.Bech32ifyConsPub(sig.PubKey)
			//	if err != nil {
			//		return nil, err
			//	}
			//
			//	sigs[i] = types.Signature{
			//		Address:       sig.Address().String(), // hex string
			//		AccountNumber: sig.AccountNumber,
			//		Pubkey:        consPubKey,
			//		Sequence:      sig.Sequence,
			//		Signature:     base64.StdEncoding.EncodeToString(sig.Signature), // encode base64
			//	}
			//}
			//

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
				Code:       tx.TxResult.Code, // 0 is success
				Messages:   string(msgsBz),
				//Signatures: string(sigsBz),
				Signatures: string("{\"sigsBz\": \"unsolved\"}"),
				Memo:       authstdTx.Memo,
				GasWanted:  tx.TxResult.GasWanted,
				GasUsed:    tx.TxResult.GasUsed,
				Timestamp:  block.Block.Time,
			}

			transactions = append(transactions, tempTransaction)

		}
	}

	return transactions, nil
}
