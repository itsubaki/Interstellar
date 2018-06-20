package cache

import "github.com/itsubaki/interstellar/broker"

type CacheBroker struct {
}

func NewCacheBroker() *CacheBroker {
	return &CacheBroker{}
}

func (b *CacheBroker) Config() *broker.Config {
	return &broker.Config{}
}

func (b *CacheBroker) Catalog() *broker.Catalog {
	return &broker.Catalog{}
}

func (b *CacheBroker) Binding(in *broker.BindingInput) *broker.BindingOutput {
	return &broker.BindingOutput{}
}

func (b *CacheBroker) Unbinding(in *broker.UnbindingInput) *broker.UnbindingOutput {
	return &broker.UnbindingOutput{}
}

func (b *CacheBroker) Create(in *broker.CreateInput) *broker.CreateOutput {
	return &broker.CreateOutput{}
}

func (b *CacheBroker) Delete(in *broker.DeleteInput) *broker.DeleteOutput {
	return &broker.DeleteOutput{}
}

func (b *CacheBroker) Update(in *broker.UpdateInput) *broker.UpdateOutput {
	return &broker.UpdateOutput{}
}
