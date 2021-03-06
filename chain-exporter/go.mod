module github.com/chain-exporter

go 1.13

require (
	github.com/Robbin-Liu/go-binance-sdk v0.0.0-20200728021042-9ef0842abec7
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200406170659-df5badaf4c2b
	github.com/cosmos/ethermint v0.0.0-20190802135314-3f32f9ba8a1f
	github.com/ethereum/go-ethereum v1.9.18
	//github.com/Robbin-Liu/go-binance-sdk v1.2.3

	github.com/go-pg/pg v8.0.6+incompatible
	github.com/go-resty/resty/v2 v2.2.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	//github.com/tendermint/tendermint v0.32.8
	github.com/tendermint/tendermint v0.33.3
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	mellium.im/sasl v0.2.1 // indirect
)

replace github.com/cosmos/ethermint => github.com/landoyjx/ethermint v0.1.2

//replace github.com/tendermint/go-amino => github.com/Robbin-Liu/bnc-go-amino v0.14.1-binance.1
replace github.com/cosmos/cosmos-sdk => github.com/landoyjx/cosmos-sdk v0.34.4-302

//replace github.com/Robbin-Liu/go-binance-sdk => github.com/Robbin-Liu/go-binance-sdk v1.2.3-bscAlpha.0
