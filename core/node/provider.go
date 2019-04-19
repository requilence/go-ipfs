package node

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipld-format"
	"github.com/libp2p/go-libp2p-routing"
	"go.uber.org/fx"

	"github.com/ipfs/go-ipfs/pin"
	"github.com/ipfs/go-ipfs/provider"
	"github.com/ipfs/go-ipfs/provider/deprecated"
	q "github.com/ipfs/go-ipfs/provider/queue"
	"github.com/ipfs/go-ipfs/repo"
)

const kReprovideFrequency = time.Hour * 12

func ProviderSysCtor(mctx MetricsCtx, lc fx.Lifecycle, cfg *config.Config, repo repo.Repo, bs BaseBlocks, ds format.DAGService, pinning pin.Pinner, rt routing.IpfsRouting) (provider.System, error) {
	if cfg.Experimental.ProviderSystemEnabled {
        return provider.NewOfflineProvider(), nil
	}

	queue, err := q.NewQueue(lifecycleCtx(mctx, lc), "provider-v1", repo.Datastore())
	if err != nil {
		return nil, err
	}

	p := deprecated.NewProvider(lifecycleCtx(mctx, lc), queue, rt)

	var keyProvider deprecated.KeyChanFunc
	switch cfg.Reprovider.Strategy {
	case "all":
		fallthrough
	case "":
		keyProvider = deprecated.NewBlockstoreProvider(bs)
	case "roots":
		keyProvider = deprecated.NewPinnedProvider(pinning, ds, true)
	case "pinned":
		keyProvider = deprecated.NewPinnedProvider(pinning, ds, false)
	default:
		return nil, fmt.Errorf("unknown reprovider strategy '%s'", cfg.Reprovider.Strategy)
	}
	r := deprecated.NewReprovider(lifecycleCtx(mctx, lc), rt, keyProvider)

    sys := provider.NewSystem(p, r)

	return sys, nil
}

func OnlineProviderSysCtor(mctx MetricsCtx, lc fx.Lifecycle, cfg *config.Config, repo repo.Repo, bs BaseBlocks, ds format.DAGService, pinning pin.Pinner, rt routing.IpfsRouting) (provider.System, error) {
	sys, err := ProviderSysCtor(mctx, lc, cfg, repo, bs, ds, pinning, rt)
	if err != nil {
		return nil, err
	}

	reproviderInterval := kReprovideFrequency
	if cfg.Reprovider.Interval != "" {
		dur, err := time.ParseDuration(cfg.Reprovider.Interval)
		if err != nil {
			return nil, err
		}

		reproviderInterval = dur
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go sys.Run(reproviderInterval)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return sys.Close()
		},
	})

	return sys, nil
}

func OfflineProviderSysCtor(mctx MetricsCtx, lc fx.Lifecycle, cfg *config.Config, repo repo.Repo, bs BaseBlocks, ds format.DAGService, pinning pin.Pinner, rt routing.IpfsRouting) (provider.System, error) {
	return ProviderSysCtor(mctx, lc, cfg, repo, bs, ds, pinning, rt)
}

