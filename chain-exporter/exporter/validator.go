package exporter

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/chain-exporter/schema"
	"github.com/chain-exporter/types"
	"github.com/chain-exporter/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// getValidators parses validators information and wrap into Precommit schema struct
func (ex *Exporter) getValidators(vals []*types.Validator) ([]*schema.Validator, error) {
	validators := make([]*schema.Validator, 0)

	// Looping through validators and insert them if not already exists in database
	for _, val := range vals {
		if len(val.ConsensusAddress) == 0 {
			var bech32Pubkey string
			err := json.Unmarshal(val.ConsensusPubKey, &bech32Pubkey)
			if err != nil {
				fmt.Errorf("failed to get ConsensusPubKey: %s", err)
				continue
			}
			val.ConsensusAddress = utils.GenHexAddrFromPubKey(bech32Pubkey)
		}
		if len(val.AccountAddress) == 0 {
			val.AccountAddress = utils.Convert(sdk.GetConfig().GetBech32AccountAddrPrefix(), val.OperatorAddress)
			//fmt.Println(sdk.GetConfig().GetBech32AccountAddrPrefix(), val.AccountAddress)
		}
		ok, err := ex.db.ExistValidator(val.ConsensusAddress)
		if !ok {

			tempVal := &schema.Validator{
				Moniker:          val.Description.Moniker,
				AccountAddress:   val.AccountAddress,
				OperatorAddress:  val.OperatorAddress,
				ConsensusAddress: val.ConsensusAddress,
				Jailed:           val.Jailed,
				Status:           strconv.FormatInt(val.Status, 10),
				Tokens:           val.Tokens,
				//VotingPower:             val.Power,
				DelegatorShares: val.DelegatorShares,
				//BondHeight:              val.BondHeight,
				//BondIntraTxCounter:      val.BondIntraTxCounter,
				UnbondingHeight: val.UnbondingHeight,
				UnbondingTime:   val.UnbondingTime,
				//CommissionRate:          val.Commission.Rate,
				CommissionMaxRate:       val.Commission.CommissionRates.MaxRate,
				CommissionMaxChangeRate: val.Commission.CommissionRates.MaxChangeRate,
				CommissionUpdateTime:    val.Commission.UpdateTime,
			}

			validators = append(validators, tempVal)
		}

		if err != nil {
			return nil, fmt.Errorf("unexpected error when checking validator existence: %s", err)
		}
	}

	return validators, nil
}
