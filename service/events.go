package service

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/filecoin-project/venus-messager/log"
	v1 "github.com/filecoin-project/venus/venus-shared/api/chain/v1"
	"github.com/filecoin-project/venus/venus-shared/types"
)

type NodeEvents struct {
	client     v1.FullNode
	log        *log.Logger
	msgService *MessageService
}

func (nd *NodeEvents) listenHeadChangesOnce(ctx context.Context) error {
	notifs, err := nd.client.ChainNotify(ctx)
	if err != nil {
		return err
	}
	select {
	case noti := <-notifs:
		if len(noti) != 1 {
			return xerrors.Errorf("expect hccurrent length 1 but for %d", len(noti))
		}

		if noti[0].Type != types.HCCurrent {
			return xerrors.Errorf("expect hccurrent event but got %s ", noti[0].Type)
		}
		//todo do some check or repaire for the first connect
		if err := nd.msgService.ReconnectCheck(ctx, noti[0].Val); err != nil {
			return xerrors.Errorf("reconnect check error: %v", err)
		}
	case <-ctx.Done():
		return ctx.Err()
	}

	for notif := range notifs {
		var apply []*types.TipSet
		var revert []*types.TipSet

		for _, change := range notif {
			switch change.Type {
			case types.HCApply:
				apply = append(apply, change.Val)
			case types.HCRevert:
				revert = append(revert, change.Val)
			}
		}

		if err := nd.msgService.ProcessNewHead(ctx, apply, revert); err != nil {
			return xerrors.Errorf("process new head error: %v", err)
		}
	}
	return nil
}
