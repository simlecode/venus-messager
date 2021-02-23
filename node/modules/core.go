package modules

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"

	"github.com/filecoin-project/go-jsonrpc/auth"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/gbrlsnchs/jwt/v3"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/peerstore"
	record "github.com/libp2p/go-libp2p-record"
	"golang.org/x/xerrors"

	"github.com/ipfs-force-community/venus-messager/api/apistruct"
	"github.com/ipfs-force-community/venus-messager/build"
	"github.com/ipfs-force-community/venus-messager/chain/types"
	"github.com/ipfs-force-community/venus-messager/lib/addrutil"
	"github.com/ipfs-force-community/venus-messager/node/modules/dtypes"
	"github.com/ipfs-force-community/venus-messager/node/repo"
)

var log = logging.Logger("modules")

// RecordValidator provides namesys compatible routing record validator
func RecordValidator(ps peerstore.Peerstore) record.Validator {
	return record.NamespacedValidator{
		"pk": record.PublicKeyValidator{},
	}
}

const JWTSecretName = "auth-jwt-private"

type jwtPayload struct {
	Allow []auth.Permission
}

func APISecret(keystore types.KeyStore, lr repo.LockedRepo) (*dtypes.APIAlg, error) {
	key, err := keystore.Get(JWTSecretName)

	if errors.Is(err, types.ErrKeyInfoNotFound) {
		log.Warn("Generating new API secret")

		sk, err := ioutil.ReadAll(io.LimitReader(rand.Reader, 32))
		if err != nil {
			return nil, err
		}

		key = types.KeyInfo{
			Type:       "jwt-hmac-secret",
			PrivateKey: sk,
		}

		if err := keystore.Put(JWTSecretName, key); err != nil {
			return nil, xerrors.Errorf("writing API secret: %w", err)
		}

		// TODO: make this configurable
		p := jwtPayload{
			Allow: apistruct.AllPermissions,
		}

		cliToken, err := jwt.Sign(&p, jwt.NewHS256(key.PrivateKey))
		if err != nil {
			return nil, err
		}

		if err := lr.SetAPIToken(cliToken); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, xerrors.Errorf("could not get JWT Token: %w", err)
	}

	return (*dtypes.APIAlg)(jwt.NewHS256(key.PrivateKey)), nil
}

func ConfigBootstrap(peers []string, local bool) func() (dtypes.BootstrapPeers, error) {
	return func() (dtypes.BootstrapPeers, error) {
		var bp []peer.AddrInfo
		if len(peers) > 0 {
			tbp, err := addrutil.ParseAddresses(context.TODO(), peers)
			if err != nil {
				return nil, err
			}
			bp = tbp
		}

		// 公网的+官方默认节点
		if local == false {
			tbp, err := BuiltinBootstrap()
			if err != nil {
				return nil, err
			}
			bp = append(bp, tbp...)
		}
		return bp, nil
	}
}

func ConfigBlacklist(peers []string) func() (dtypes.BlacklistPeers, error) {
	return func() (dtypes.BlacklistPeers, error) {
		return peers, nil
	}
}

func BuiltinBootstrap() (dtypes.BootstrapPeers, error) {
	return build.BuiltinBootstrap()
}

func DrandBootstrap() (dtypes.DrandBootstrap, error) {
	return build.DrandBootstrap()
}
