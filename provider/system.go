package provider

import (
	"context"
	"github.com/ipfs/go-cid"
	"time"
)

type System interface {
	Run(time.Duration)
	Close() error
	Provide(cid.Cid) error
	Reprovide(context.Context) error
}

type system struct {
    provider   Provider
    reprovider Reprovider
}

func NewSystem(provider Provider, reprovider Reprovider) System {
	return &system{provider, reprovider}
}

func (s *system) Run(reprovideInteval time.Duration) {
	go s.provider.Run()
	go s.reprovider.Run(reprovideInteval)
}

func (s *system) Close() error {
	return s.provider.Close()
}

func (s *system) Provide(cid cid.Cid) error {
	return s.provider.Provide(cid)
}

func (s *system) Reprovide(ctx context.Context) error {
	return s.reprovider.Trigger(ctx)
}
