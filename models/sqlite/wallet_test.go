package sqlite

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ipfs-force-community/venus-messager/config"
	"github.com/ipfs-force-community/venus-messager/types"
)

func TestWallet(t *testing.T) {
	path := "sqlite_wallet.db"
	repo, err := OpenSqlite(&config.SqliteConfig{Path: path})
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.Remove(path))
	}()
	err = repo.AutoMigrate()
	assert.NoError(t, err)

	walletRepo := repo.WalletRepo()

	w := &types.Wallet{
		ID:   types.NewUUID(),
		Name: "wallet1",
		Url:  "http://127.0.0.1:8080",

		IsDeleted: -1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	w2 := &types.Wallet{
		ID:        types.NewUUID(),
		Name:      "wallet2",
		Url:       "http://127.0.0.1:8082",
		IsDeleted: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := walletRepo.SaveWallet(w)
	assert.NoError(t, err)
	uuid, err := types.ParseUUID(id)
	assert.NoError(t, err)

	id2, err := walletRepo.SaveWallet(w2)
	assert.NoError(t, err)
	uuid2, err := types.ParseUUID(id2)
	assert.NoError(t, err)

	r, err := walletRepo.GetWallet(uuid)
	assert.NoError(t, err)
	assert.Equal(t, w.Name, r.Name)
	assert.Equal(t, w.Url, r.Url)
	t.Logf("%+v", r)

	_, err = walletRepo.GetWallet(uuid2)
	assert.Error(t, err)

	rs, err := walletRepo.ListWallet()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(rs))

	err = walletRepo.DelWallet(uuid)
	assert.NoError(t, err)

	_, err = walletRepo.GetWallet(uuid)
	assert.Error(t, err)
}
