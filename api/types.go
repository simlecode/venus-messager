package api

import (
	"github.com/filecoin-project/go-state-types/abi"
)

type MessageSendSpec struct {
	MaxFee abi.TokenAmount
}
