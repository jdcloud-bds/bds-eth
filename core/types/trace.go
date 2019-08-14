package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
)

type TraceResult struct {
	BlockHash        common.Hash     `json:"blockHash"`
	TransactionHash  common.Hash     `json:"transactionHash" gencodec:"required"`
	TransactionIndex uint            `json:"transactionIndex" gencodec:"required"`
	Result           json.RawMessage `json:"traceResult"`
}
