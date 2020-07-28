package exporter

import (
	"github.com/chain-exporter/schema"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getBlock parses block information and wrap into Block schema struct
func (ex *Exporter) getBlock(block *tmctypes.ResultBlock) ([]*schema.Block, error) {
	blocks := make([]*schema.Block, 0)

	tempBlock := &schema.Block{
		Height:        block.Block.Height,
		Proposer:      block.Block.ProposerAddress.String(),
		Moniker:       ex.db.QueryMoniker(block.Block.ProposerAddress.String()),
		//BlockHash:     block.BlockMeta.BlockID.Hash.String(),
		BlockHash:     block.BlockID.Hash.String(),
		//ParentHash:    block.BlockMeta.Header.LastBlockID.Hash.String(),
		ParentHash:    block.Block.Header.LastBlockID.Hash.String(),
		//NumPrecommits: int64(len(block.Block.LastCommit.Precommits)),
		//NumTxs:        block.Block.NumTxs,
		//TotalTxs:      block.Block.TotalTxs,
		NumPrecommits: 0,
		NumTxs:        int64(len(block.Block.Txs)) ,
		TotalTxs:      int64(len(block.Block.Txs)),
		Timestamp:     block.Block.Time,
	}

	blocks = append(blocks, tempBlock)

	return blocks, nil
}
