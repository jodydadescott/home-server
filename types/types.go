package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	logger "github.com/jodydadescott/jody-go-logger"
	"github.com/jodydadescott/unifi-go-sdk"

	"github.com/jodydadescott/home-server/types/proto"
)

type Logger = logger.Config

// Netport is the IP, Port and Protocol type
type NetPort struct {
	IP          string      `json:"ip,omitempty" yaml:"ip,omitempty"`
	Port        int         `json:"port,omitempty" yaml:"port,omitempty"`
	Proto       proto.Proto `json:"proto,omitempty" yaml:"proto,omitempty"`
	ipColonPort string      `json:"-"`
}

// Clone return copy
func (t *NetPort) Clone() *NetPort {
	c := &NetPort{}
	copier.Copy(&c, &t)
	return c
}

// SetProtoTCP sets proto type to TCP
func (t *NetPort) SetProtoTCP() {
	t.Proto = proto.TCP
}

// SetProtoUDP sets proto type to UDP
func (t *NetPort) SetProtoUDP() {
	t.Proto = proto.UDP
}

// GetIPColonPort returns the IP + colong + port as a string
func (t *NetPort) GetIPColonPort() string {
	if t.ipColonPort == "" {
		t.ipColonPort = t.IP + ":" + fmt.Sprint(t.Port)
	}
	return t.ipColonPort
}

// ARecord is a DNS A Record
type ARecord struct {
	Domain   string `json:"domain,omitempty" yaml:"domain,omitempty"`
	Hostname string `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	IP       string `json:"ip,omitempty" yaml:"ip,omitempty"`
	SRC      string `json:"src,omitempty" yaml:"src,omitempty"`
	fqdn     string `json:"-"`
}

// Clone return copy
func (t *ARecord) Clone() *ARecord {
	c := &ARecord{}
	copier.Copy(&c, &t)
	return c
}

// GetKey returns the key for the record type
func (t *ARecord) GetKey() string {
	if t.fqdn == "" {
		t.fqdn = cleanHostname(t.Hostname) + "." + t.Domain + "."
	}
	return t.fqdn
}

// GetValue returns the value for the record type
func (t *ARecord) GetValue() string {
	return t.IP
}

// CNameRecord is a DNS CNAME Record
type CNameRecord struct {
	AliasHostname  string `json:"aliasHostname,omitempty" yaml:"aliasHostname,omitempty"`
	AliasDomain    string `json:"aliasDomain,omitempty" yaml:"aliasDomain,omitempty"`
	TargetHostname string `json:"targetHostname,omitempty" yaml:"targetHostname,omitempty"`
	TargetDomain   string `json:"targetDomain,omitempty" yaml:"targetDomain,omitempty"`
	SRC            string `json:"src,omitempty" yaml:"src,omitempty"`
	fqdnAlias      string `json:"-"`
	fqdnTarget     string `json:"-"`
}

// Clone return copy
func (t *CNameRecord) Clone() *CNameRecord {
	c := &CNameRecord{}
	copier.Copy(&c, &t)
	return c
}

// GetKey returns the key for the record type
func (t *CNameRecord) GetKey() string {
	if t.fqdnAlias == "" {
		t.fqdnAlias = cleanHostname(t.AliasHostname) + "." + t.AliasDomain + "."
	}
	return t.fqdnAlias
}

// GetValue returns the value for the record type
func (t *CNameRecord) GetValue() string {
	if t.fqdnTarget == "" {
		t.fqdnTarget = cleanHostname(t.TargetHostname) + "." + t.TargetDomain + "."
	}
	return t.fqdnTarget
}

// PTRrecord is a DNS PTR Record
type PTRrecord struct {
	ARPA     string `json:"arpa,omitempty" yaml:"arpa,omitempty"`
	Hostname string `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Domain   string `json:"domain,omitempty" yaml:"domain,omitempty"`
	SRC      string `json:"src,omitempty" yaml:"src,omitempty"`
	fqdn     string `json:"-"`
}

// Clone return copy
func (t *PTRrecord) Clone() *PTRrecord {
	c := &PTRrecord{}
	copier.Copy(&c, &t)
	return c
}

// GetKey returns the key for the record type
func (t *PTRrecord) GetKey() string {
	return t.ARPA
}

// GetValue returns the value for the record type
func (t *PTRrecord) GetValue() string {
	if t.fqdn == "" {
		t.fqdn = cleanHostname(t.Hostname) + "." + t.Domain + "."
	}
	return t.fqdn
}

// Config is the main user level config
type Config struct {
	Notes       string        `json:"notes,omitempty" yaml:"notes,omitempty"`
	Unifi       *UnifiConfig  `json:"unifiConfig,omitempty" yaml:"unifiConfig,omitempty"`
	Listeners   []*NetPort    `json:"listeners,omitempty" yaml:"listeners,omitempty"`
	Static      *StaticConfig `json:"static,omitempty" yaml:"static,omitempty"`
	Nameservers []*NetPort    `json:"nameservers,omitempty" yaml:"nameservers,omitempty"`
	Logging     *Logger       `json:"logging,omitempty" yaml:"logging,omitempty"`
	HttpConfig  *HttpConfig   `json:"httpConfig,omitempty" yaml:"httpConfig,omitempty"`
}

// HttpConfig is the config for HTTP servers
type HttpConfig struct {
	Listener *NetPort `json:"listener,omitempty" yaml:"listener,omitempty"`
	Enabled  bool     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

// Clone return copy
func (t *HttpConfig) Clone() *HttpConfig {
	c := &HttpConfig{}
	copier.Copy(&c, &t)
	return c
}

// Clone return copy
func (t *Config) Clone() *Config {
	c := &Config{}
	copier.Copy(&c, &t)
	return c
}

// AddNameserver adds the specified nameserver to the config
func (t *Config) AddNameservers(nameservers ...*NetPort) *Config {
	for _, v := range nameservers {
		t.Nameservers = append(t.Nameservers, v)
	}
	return t
}

// AddNameserver adds the specified nameserver to the config
func (t *Config) AddListeners(listeners ...*NetPort) *Config {
	for _, v := range listeners {
		t.Listeners = append(t.Listeners, v)
	}
	return t
}

// UnifiConfig is the config for Unifi servers
type UnifiConfig struct {
	unifi.Config
	Refresh    time.Duration `json:"refresh,omitempty" yaml:"refresh,omitempty"`
	Enabled    bool          `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Domain     string        `json:"domain,omitempty" yaml:"domain,omitempty"`
	IgnoreMacs []string      `json:"ignoreMacs,omitempty" yaml:"ignoreMacs,omitempty"`
}

// Clone return copy
func (t *UnifiConfig) Clone() *UnifiConfig {
	c := &UnifiConfig{}
	copier.Copy(&c, &t)
	return c
}

// AddIgnoreMac add a MAC that will be ignored when process the Unifi config
func (t *UnifiConfig) AddIgnoreMacs(macs ...string) *UnifiConfig {
	for _, v := range macs {
		t.IgnoreMacs = append(t.IgnoreMacs, v)
	}
	return t
}

// IgnoreMac is a convenience function that returns true if the MAC exist in the ignore slice
func (t *UnifiConfig) IgnoreMac(mac string) bool {
	mac = strings.ToLower(mac)
	for _, m := range t.IgnoreMacs {
		if mac == strings.ToLower(m) {
			return true
		}
	}
	return false
}

// StaticConfig are records from config that are statically defined
type StaticConfig struct {
	Enabled bool      `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Domains []*Domain `json:"domains,omitempty" yaml:"domains,omitempty"`
}

// Clone return copy
func (t *StaticConfig) Clone() *StaticConfig {
	c := &StaticConfig{}
	copier.Copy(&c, &t)
	return c
}

// Domain is a collection of A & CNAME records with a common domain. If the domain is
// not set then a default domain will be used. The same default domain will be used
// of CNAME target domains if not configured. It is not normally required to add PTR
// records as they will be automatically generated when the A record is created.
type Domain struct {
	Domain  string        `json:"domain,omitempty" yaml:"dnsDomain,omitempty"`
	Records DomainRecords `json:"records,omitempty" yaml:"records,omitempty"`
}

// Clone return copy
func (t *Domain) Clone() *Domain {
	c := &Domain{}
	copier.Copy(&c, &t)
	return c
}

type DomainRecords struct {
	ARecords     []*ARecord     `json:"aRecords,omitempty" yaml:"aRecords,omitempty"`
	AAAARecords  []*ARecord     `json:"aaaRecords,omitempty" yaml:"aaaRecords,omitempty"`
	CnameRecords []*CNameRecord `json:"cnameRecords,omitempty" yaml:"cnameRecords,omitempty"`
	PtrRecords   []*PTRrecord   `json:"ptrRecords,omitempty" yaml:"ptrRecords,omitempty"`
}

// AddDomain is a convenience function that adds the specified Domaain to the StaticConfig
func (t *StaticConfig) AddDomains(domains ...*Domain) *StaticConfig {
	for _, v := range domains {
		t.Domains = append(t.Domains, v)
	}
	return t
}

// AddARecord is a convenience that adds the specified ARecord to the Domain
func (t *DomainRecords) AddARecords(records ...*ARecord) *DomainRecords {
	for _, v := range records {
		t.ARecords = append(t.ARecords, v)
	}
	return t
}

// AddARecord is a convenience that adds the specified ARecord to the Domain
func (t *DomainRecords) AddAAAARecords(records ...*ARecord) *DomainRecords {
	for _, v := range records {
		t.AAAARecords = append(t.AAAARecords, v)
	}
	return t
}

func (t *DomainRecords) AddPtrRecords(records ...*PTRrecord) *DomainRecords {
	for _, v := range records {
		t.PtrRecords = append(t.PtrRecords, v)
	}
	return t
}

// CNameRecord is a convenience that adds the specified CNameRecord to the Domain
func (t *DomainRecords) AddCNameRecords(records ...*CNameRecord) *DomainRecords {
	for _, v := range records {
		t.CnameRecords = append(t.CnameRecords, v)
	}
	return t
}
