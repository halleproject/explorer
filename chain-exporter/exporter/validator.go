package exporter

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	"strconv"

	"github.com/chain-exporter/schema"
	"github.com/chain-exporter/types"
)

// getValidators parses validators information and wrap into Precommit schema struct
func (ex *Exporter) getValidators(vals []*types.Validator, valSet *tmctypes.ResultValidators) ([]*schema.Validator, error) {
	validators := make([]*schema.Validator, 0)

	// Looping through validators and insert them if not already exists in database
	for _, val := range vals {
		if len(val.ConsensusAddress) == 0 {
			consensusAddress := sdk.ConsAddress(val.ConsensusPubKey)
			val.ConsensusAddress = consensusAddress.String()
		}
		ok, err := ex.db.ExistValidator(val.ConsensusAddress)
		//ok, err := ex.db.ExistValidator(val.OperatorAddress)
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
