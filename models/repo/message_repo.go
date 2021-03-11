package repo

import (
	"time"

	"github.com/filecoin-project/go-state-types/abi"
	venustypes "github.com/filecoin-project/venus/pkg/types"
	"github.com/ipfs-force-community/venus-messager/types"
)

type MessageRepo interface {
	SaveMessage(msg *types.Message) (string, error)
	GetMessage(uuid string) (*types.Message, error)
	UpdateMessageReceipt(cid string, receipt *venustypes.MessageReceipt, height abi.ChainEpoch, state types.MessageState) (string, error)
	ListMessage() ([]*types.Message, error)
	ListUnchainedMsgs() ([]*types.Message, error)
	GetMessageByCid(cid string) (*types.Message, error)
	GetMessageByTime(start time.Time) ([]*types.Message, error)
	UpdateMessageStateByCid(cid string, state types.MessageState) error
	GetActorNonce(fromAddr string) (uint64, error)
	ListUnAssignNonceMessages(fromAddr string) ([]*types.Message, error)
	UpdateMessageNonce(uuid string, nonce uint64) error

	LockFromAddr(fromAddr string)
	UnlockFromAddr(fromAddr string)
}
