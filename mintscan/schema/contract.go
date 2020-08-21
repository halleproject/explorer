package schema

// Contract defines the structure for Contract information
type Contract struct {
	ID              int32  `json:"id" sql:",pk"`
	Height          int64  `json:"height" sql:",notnull"`
	TxHash          string `json:"tx_hash" sql:",notnull,unique"`
	FromAddress     string `json:"from_address" sql:",notnull"`
	ContractAddress string `json:"contract_address" sql:"default:''"`
	Memo            string `json:"memo"`

	TotalSupply string `json:"total_supply" sql:",notnull,unique"`
	Decimals    uint32 `json:"decimals"  sql:",notnull"`
	Name        string `json:"name" sql:"notnull"`
	Symbol      string `json:"symbol" sql:"notnull"`
}
