module github.com/cosmostation/mintscan-binance-dex-backend/mintscan

go 1.13

require (
	github.com/binance-chain/go-sdk v1.2.2
	github.com/ethereum/go-ethereum v1.9.24
	github.com/go-pg/pg v8.0.6+incompatible
	github.com/go-resty/resty/v2 v2.2.0
	github.com/gorilla/mux v1.7.3
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.6.0
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.32.3
	mellium.im/sasl v0.2.1 // indirect
)

replace github.com/binance-chain/go-sdk => github.com/wade-liwei/go-sdk v1.2.20
