package service

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/ipfs-force-community/venus-messager/models/repo"
	"github.com/ipfs-force-community/venus-messager/types"
)

type WalletService struct {
	repo repo.Repo
	log  *logrus.Logger

	walletClients map[string]*WalletClient

	l sync.Mutex
}

func NewWalletService(repo repo.Repo, logger *logrus.Logger) (*WalletService, error) {
	ws := &WalletService{
		repo:          repo,
		log:           logger,
		walletClients: make(map[string]*WalletClient),
	}

	walletList, err := ws.ListWallet(context.TODO())
	if err != nil {
		return nil, err
	}

	for _, w := range walletList {
		cli, _, err := newWalletClient(context.Background(), w.Url, w.Token)
		if err != nil {
			return nil, err
		}

		if _, ok := ws.walletClients[w.Name]; !ok {
			ws.walletClients[w.Name] = &cli
		}
	}

	return ws, err
}

func (walletService WalletService) SaveWallet(ctx context.Context, wallet *types.Wallet) (string, error) {
	return walletService.repo.WalletRepo().SaveWallet(wallet)
}

func (walletService WalletService) GetWallet(ctx context.Context, uuid string) (*types.Wallet, error) {
	return walletService.repo.WalletRepo().GetWallet(uuid)
}

func (walletService WalletService) ListWallet(ctx context.Context) ([]*types.Wallet, error) {
	return walletService.repo.WalletRepo().ListWallet()
}

func (walletService WalletService) updateWalletClient(ctx context.Context, wallet *types.Wallet) error {
	cli, _, err := newWalletClient(context.Background(), wallet.Url, wallet.Token)
	if err != nil {
		return err
	}
	walletService.l.Lock()
	defer walletService.l.Unlock()
	if _, ok := walletService.walletClients[wallet.Name]; !ok {
		walletService.walletClients[wallet.Name] = &cli
	}

	return nil
}