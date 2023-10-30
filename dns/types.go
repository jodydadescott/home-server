package dns

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/jodydadescott/home-server/types"
	"github.com/jodydadescott/home-server/types/proto"
)

type NetPort = types.NetPort
type Proto = proto.Proto
type Domain = types.Domain
type ARecord = types.ARecord
type PTRrecord = types.PTRrecord
type CNameRecord = types.CNameRecord
type DomainRecords = types.DomainRecords

type Config struct {
	Providers   []Provider
	Trace       bool
	Listeners   []*NetPort
	Nameservers []*NetPort
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

func (t *Config) AddProvider(provider Provider) {
	t.Providers = append(t.Providers, provider)
}

type Provider interface {
	GetName() string
	GetDomainName() string
	GetRecords() (*DomainRecords, error)
	GetRefreshDuration() time.Duration
}
