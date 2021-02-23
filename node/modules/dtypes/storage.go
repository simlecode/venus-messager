package dtypes

import (
	bserv "github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
)

// MetadataDS stores metadata
// dy default it's namespaced under /metadata in main repo datastore
type MetadataDS datastore.Batching

type ChainBlockstore blockstore.Blockstore

type ChainGCLocker blockstore.GCLocker
type ChainGCBlockstore blockstore.GCBlockstore
type ChainExchange exchange.Interface
type ChainBlockService bserv.BlockService

type ClientDatastore datastore.Batching
