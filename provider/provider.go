package provider

import (
	"context"
	"fmt"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipfs-blockstore"
	"time"

	cid "github.com/ipfs/go-cid"
	pin "github.com/ipfs/go-ipfs/pin"
	deprecated "github.com/ipfs/go-ipfs/provider/deprecated"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log"
	routing "github.com/libp2p/go-libp2p-routing"
	q "github.com/ipfs/go-ipfs/provider/queue"
)

var (
	UseExperimentalSystem = false

	log = logging.Logger("provider")
)

// Provider announces blocks to the network
type Provider interface {
	// Run is used to begin processing the provider work
	Run()
	// Provide takes a cid and makes an attempt to announce it to the network
	Provide(cid.Cid) error
	// Close stops the provider
	Close() error
}

// Reprovider reannounces blocks to the network
type Reprovider interface {
	// Run is used to begin processing the reprovider work and waiting for reprovide triggers
	Run(time.Duration)
	// Trigger a reprovide
	Trigger(context.Context) error
}

// NewConfiguredProvider returns either the deprecated or the experimental provider
func NewConfiguredProvider(ctx context.Context, name string, datastore datastore.Datastore, routing routing.ContentRouting) (Provider, error) {
	if UseExperimentalSystem {
        return nil, nil
	} else {
		queue, err := q.NewQueue(ctx, name, datastore)
		if err != nil {
			return nil, err
		}
		return deprecated.NewProvider(ctx, queue, routing), nil
	}
}

// NewConfiguredReprovider returns either the deprecated or the experimental reprovider
func NewConfiguredReprovider(ctx context.Context, routing routing.ContentRouting, blockstore blockstore.Blockstore, pinning pin.Pinner, dag ipld.DAGService, strategy string) (Reprovider, error) {
	if UseExperimentalSystem {
		return nil, nil
	} else {
		var keyProvider deprecated.KeyChanFunc
		switch strategy {
		case "all":
			fallthrough
		case "":
			keyProvider = deprecated.NewBlockstoreProvider(blockstore)
		case "roots":
			keyProvider = deprecated.NewPinnedProvider(pinning, dag, true)
		case "pinned":
			keyProvider = deprecated.NewPinnedProvider(pinning, dag, false)
		default:
			return nil, fmt.Errorf("unknown reprovider strategy '%s'", strategy)
		}
		return deprecated.NewReprovider(ctx, routing, keyProvider), nil
	}
}