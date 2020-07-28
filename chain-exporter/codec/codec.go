package codec

import (
	amino "github.com/tendermint/go-amino"

	"github.com/Robbin-Liu/go-binance-sdk/types"

	//"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

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

}
