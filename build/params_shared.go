package build

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/libp2p/go-libp2p-core/protocol"

	"github.com/ipfs-force-community/venus-messager/node/modules/dtypes"
)

func BlocksTopic(netName dtypes.NetworkName) string   { return "/fil/blocks/" + string(netName) }
func MessagesTopic(netName dtypes.NetworkName) string { return "/fil/msgs/" + string(netName) }
func DhtProtocolName(netName dtypes.NetworkName) protocol.ID {
	return protocol.ID("/fil/kad/" + string(netName))
}

var DrandConfig = dtypes.DrandConfig{
	Servers: []string{
		"https://pl-eu.testnet.drand.sh",
		"https://pl-us.testnet.drand.sh",
		"https://pl-sin.testnet.drand.sh",
	},
	ChainInfoJSON: `{"public_key":"922a2e93828ff83345bae533f5172669a26c02dc76d6bf59c80892e12ab1455c229211886f35bb56af6d5bea981024df","period":25,"genesis_time":1590445175,"hash":"138a324aa6540f93d0dad002aa89454b1bec2b6e948682cde6bd4db40f4b7c9b"}`,
}

type RepoType int

const (
	_                 = iota // Default is invalid
	FullNode RepoType = iota
)

const FilecoinPrecision = uint64(1_000_000_000_000_000_000)

// ///////
// Limits

// TODO: If this is gonna stay, it should move to specs-actors
const BlockGasLimit = 10_000_000_000

// Actor consts
// TODO: Pull from actors when its made not private
var MinDealDuration = abi.ChainEpoch(180 * builtin.EpochsInDay)

// Blocks (e)
var BlocksPerEpoch = uint64(builtin.ExpectedLeadersPerEpoch)
