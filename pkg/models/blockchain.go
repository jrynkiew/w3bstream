package models

import (
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"
)

// Blockchain database model blockchain
// @def primary                   ID
// @def unique_index UI_chain_id ChainID
//
//go:generate toolkit gen model Blockchain --database DB
type Blockchain struct {
	datatypes.PrimaryID
	RelBlockchain
	BlockchainInfo
	datatypes.OperationTimes
}

type RelBlockchain struct {
	ChainID uint64 `db:"f_chain_id"             json:"chainID"`
}

type BlockchainInfo struct {
	Address string `db:"f_chain_address"         json:"chainAddress"`
}