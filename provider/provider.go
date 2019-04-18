package provider

import (
	cid "github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("provider")

// Provider announces blocks to the network
type Provider interface {
	// Run is used to begin processing the provider work
	Run()
	// Provide takes a cid and makes an attempt to announce it to the network
	Provide(cid.Cid) error
	// Close stops the provider
	Close() error
}
