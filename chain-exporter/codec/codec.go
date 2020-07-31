package codec

import (
	"github.com/Robbin-Liu/go-binance-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	escrow "github.com/cosmos/ethermint/x/escrow/types"
	ethtypes "github.com/cosmos/ethermint/x/evm/types"
	amino "github.com/tendermint/go-amino"
)

// Codec is amino codec to serialize Binance Chain interfaces and data
var Codec *amino.Codec

// initializes upon package loading
func init() {
	Codec = types.NewCodec()

	sdk.RegisterCodec(Codec)

	// register test types
	//cdc.RegisterConcrete(&txTest{}, "cosmos-sdk/baseapp/txTest", nil)
	Codec.RegisterConcrete(bank.MsgSend{}, "cosmos-sdk/MsgSend", nil)
	Codec.RegisterConcrete(bank.MsgMultiSend{}, "cosmos-sdk/MsgMultiSend", nil)
	Codec.RegisterConcrete(auth.StdTx{}, "cosmos-sdk/StdTx", nil)

	Codec.RegisterConcrete(escrow.MsgSendWithUnlock{}, "escrow/MsgSendWithUnlock", nil)
	Codec.RegisterConcrete(escrow.MsgPayout{}, "escrow/Payout", nil)

	Codec.RegisterConcrete(ethtypes.MsgEthereumTx{}, "ethermint/MsgEthereumTx", nil)
	Codec.RegisterConcrete(ethtypes.MsgEthermint{}, "ethermint/MsgEthermint", nil)
	Codec.RegisterConcrete(ethtypes.TxData{}, "ethermint/TxData", nil)
}
