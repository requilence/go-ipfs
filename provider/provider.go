package provider

import (
	"context"
	"time"

	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
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
