package config

import (
	"encoding"
	"time"
)

// Common is common config between full node and miner
type Common struct {
	API       API
	Libp2p    Libp2p
	Pubsub    Pubsub
	DbCfg     DbCfg
	RPCServer RPCServer
}

// FullNode is a full node config
type FullNode struct {
	Common
	Metrics Metrics
}

// // Common

// API contains configs for API endpoint
type API struct {
	ListenAddress       string
	RemoteListenAddress string
	Timeout             Duration
}

// db config
type DbCfg struct {
	Conn      string `json:"conn" binding:"required"`
	Type      string `json:"type" binding:"required"`
	DebugMode bool   `json:"debugMode" binding:"required"`
}

// Libp2p contains configs for libp2p
type Libp2p struct {
	ListenAddresses []string
	BootstrapPeers  []string

	ConnMgrLow   uint
	ConnMgrHigh  uint
	ConnMgrGrace Duration
}

type Pubsub struct {
	Bootstrapper bool
	DirectPeers  []string
	RemoteTracer string
}

type RPCServer struct {
	ServerUrls []string
}

// // Full Node

type Metrics struct {
	Nickname   string
	HeadNotifs bool
}

func defCommon() Common {
	return Common{
		API: API{
			ListenAddress: "/ip4/0.0.0.0/tcp/5678/http",
			Timeout:       Duration(30 * time.Second),
		},
		DbCfg: DbCfg{
			Conn:      "",
			Type:      "mysql",
			DebugMode: true,
		},
		Libp2p: Libp2p{
			ListenAddresses: []string{
				"/ip4/0.0.0.0/tcp/0",
				"/ip6/::/tcp/0",
			},
			BootstrapPeers: []string{},

			ConnMgrLow:   150,
			ConnMgrHigh:  180,
			ConnMgrGrace: Duration(20 * time.Second),
		},
		Pubsub: Pubsub{
			Bootstrapper: false,
			DirectPeers:  nil,
			RemoteTracer: "/ip4/147.75.67.199/tcp/4001/p2p/QmTd6UvR47vUidRNZ1ZKXHrAFhqTJAD27rKL9XYghEKgKX",
		},
	}

}

// DefaultFullNode returns the default config
func DefaultFullNode() *FullNode {
	return &FullNode{
		Common: defCommon(),
	}
}

var _ encoding.TextMarshaler = (*Duration)(nil)
var _ encoding.TextUnmarshaler = (*Duration)(nil)

// Duration is a wrapper type for time.Duration
// for decoding and encoding from/to TOML
type Duration time.Duration

// UnmarshalText implements interface for TOML decoding
func (dur *Duration) UnmarshalText(text []byte) error {
	d, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*dur = Duration(d)
	return err
}

func (dur Duration) MarshalText() ([]byte, error) {
	d := time.Duration(dur)
	return []byte(d.String()), nil
}
