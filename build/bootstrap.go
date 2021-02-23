package build

import (
	"context"
	"io/ioutil"
	"strings"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/ipfs-force-community/venus-messager/lib/addrutil"
)

var log = logging.Logger("peermgr")

func BuiltinBootstrap() ([]peer.AddrInfo, error) {
	var out []peer.AddrInfo

	b, err := ioutil.ReadFile("bootstrappers.pi")
	if err != nil {
		log.Info("read bootstrap failed: %s", err)
		return out, nil
	}

	spi := string(b)
	if spi == "" {
		log.Info("no bootstrap")
		return out, nil
	}

	pi, err := addrutil.ParseAddresses(context.TODO(), strings.Split(strings.TrimSpace(spi), "\n"))
	out = append(out, pi...)

	//b := rice.MustFindBox("bootstrap")
	//err := b.Walk("bootstrappers.pi", func(path string, info os.FileInfo, err error) error {
	//	if err != nil {
	//		return xerrors.Errorf("failed to walk box: %w", err)
	//	//	}
	//
	//	if !strings.HasSuffix(path, "bootstrappers.pi") {
	//		return nil
	//	}
	//	spi := b.MustString(path)
	//	if spi == "" {
	//		return nil
	//	}
	//	pi, err := addrutil.ParseAddresses(context.TODO(), strings.Split(strings.TrimSpace(spi), "\n"))
	//	out = append(out, pi...)
	//	return err
	//})

	return out, nil
}

func DrandBootstrap() ([]peer.AddrInfo, error) {
	//addrs := []string{
	//	"/dnsaddr/pl-eu.testnet.drand.sh/",
	//	"/dnsaddr/pl-us.testnet.drand.sh/",
	//	"/dnsaddr/pl-sin.testnet.drand.sh/",
	//}
	//return addrutil.ParseAddresses(context.TODO(), addrs)
	var peers []peer.AddrInfo
	return peers, nil
}
