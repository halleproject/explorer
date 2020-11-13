package handlers

import (
	"testing"

	//sdk "github.com/binance-chain/go-sdk/types"
	sdk "github.com/binance-chain/go-sdk/common/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
)

func TestExchange(t *testing.T) {

	addr, err := sdk.AccAddressFromBech32("halle1nea868ahrhvlj66zlug4e5ayf2hr4eff62qh49")
	if err != nil {
		t.Error(err)
	}

	t.Logf("addr:  %s\n", addr)
	address := ethcmn.BytesToAddress(addr.Bytes())
	t.Logf("address: %x \n", address)

	acc := sdk.AccAddress(address.Bytes())
	t.Logf("address: %s \n", acc)

	//acc = sdk.AccAddress([]byte("0x9e7a7d1fb71dd9f96b42ff115cd3a44aae3ae529"))
	acc1, err := sdk.AccAddressFromHex("9e7a7d1fb71dd9f96b42ff115cd3a44aae3ae529")
	if err != nil {
		t.Error(err)
	}

	t.Logf("address: %s \n", acc1)

}
