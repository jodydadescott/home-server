package http

import (
	"github.com/jinzhu/copier"

	"github.com/jodydadescott/home-server/types"
)

type DomainRecords = types.DomainRecords

type Config struct {
	Listener       *NetPort
	RecordProvider RecordProvider
}

type RecordProvider interface {
	GetRecords() *DomainRecords
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}
