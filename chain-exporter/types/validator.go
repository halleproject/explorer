package types

import "encoding/json"

type (
	// Validator defines the structure for validator information
	Validator struct {
		//AccountAddress     string          `json:"account_address" sql:",notnull, unique"`
		AccountAddress    string          `json:"account_address"`
		OperatorAddress   string          `json:"operator_address" sql:",notnull, unique"`
		ConsensusPubKey   json.RawMessage `json:"consensus_pubkey" sql:",notnull, unique"`
		ConsensusAddress  string          `json:"consensus_address" sql:",notnull, unique"`
		Jailed            bool            `json:"jailed"`
		Status            int64           `json:"status"`
		Tokens            string          `json:"tokens"`
		DelegatorShares   string          `json:"delegator_shares"`
		Description       Description     `json:"description"`
		UnbondingHeight   int64           `json:"unbonding_height"`
		UnbondingTime     string          `json:"unbonding_time"`
		Commission        Commission      `json:"commission"`
		MinSelfDelegation string          `json:"tokens"`
	}

	// Description wraps description of a validator
	Description struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	}

	// Commission wrpas general commission information about a validator
	Commission struct {
		CommissionRates CommissionRates `json:"commission_rates"`
		//MaxRate       string `json:"max_rate"`
		//MaxChangeRate string `json:"max_change_rate"`
		UpdateTime string `json:"update_time"`
	}
	CommissionRates struct {
		MaxRate       string `json:"max_rate"`
		MaxChangeRate string `json:"max_change_rate"`
		Rate          string `json:"rate"`
	}

	HttpBody struct {
		Height    string       `json:"height"`
		Validator []*Validator `json:"result"`
	}
)
