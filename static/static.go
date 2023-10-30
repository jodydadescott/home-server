package static

import (
	"fmt"
	"time"

	"github.com/jodydadescott/home-server/types"
	"github.com/jodydadescott/home-server/util"
)

type Config = types.StaticConfig
type ARecord = types.ARecord
type PTRrecord = types.PTRrecord
type Records = types.DomainRecords
type Domain = types.Domain

const (
	source = "config"
)

type Client struct {
	domain *Domain
}

func New(config *Config) []*Client {

	if config == nil {
		panic("config is required")
	}

	config = config.Clone()

	var clients []*Client

	for _, domain := range config.Domains {

		domain = domain.Clone()

		if domain.Domain == "" {
			domain.Domain = types.DefaultDomain
		}

		clients = append(clients, &Client{domain: domain})
	}

	return clients
}

func (t *Client) GetName() string {
	return "static"
}

func (t *Client) GetDomainName() string {
	return t.domain.Domain
}

func (t *Client) GetRefreshDuration() time.Duration {
	return 0
}

func (t *Client) GetRecords() (*Records, error) {

	ptrRecordsMap := make(map[string]*PTRrecord)

	addPTR := func(a *ARecord, iptype string) error {

		if a.Hostname == "" {
			return fmt.Errorf("Record must have a hostname")
		}

		if a.IP == "" {
			return fmt.Errorf("Record must have a IP")
		}

		if a.Domain == "" {
			a.Domain = t.domain.Domain
		}

		a.SRC = source + ":static"

		arpa, err := util.GetARPA(a.IP)
		if err != nil {
			return err
		}

		p := &PTRrecord{
			ARPA:     arpa,
			Hostname: a.Hostname,
			Domain:   a.Domain,
			SRC:      source + ":dynamic",
		}

		ptrRecordsMap[p.GetKey()] = p

		return nil
	}

	for _, record := range t.domain.Records.ARecords {
		err := addPTR(record, "A")
		if err != nil {
			return nil, err
		}
	}

	for _, record := range t.domain.Records.AAAARecords {
		err := addPTR(record, "AAAA")
		if err != nil {
			return nil, err
		}
	}

	for _, r := range t.domain.Records.CnameRecords {

		if r.AliasHostname == "" {
			return nil, fmt.Errorf("CNAME must have AliasHostname")
		}

		if r.TargetHostname == "" {
			return nil, fmt.Errorf("CNAME must have TargetHostname")
		}

		if r.AliasDomain == "" {
			r.AliasDomain = t.domain.Domain
		}

		if r.TargetDomain == "" {
			r.TargetDomain = t.domain.Domain
		}

		r.SRC = source + ":static"

	}

	for _, p := range t.domain.Records.PtrRecords {
		existing := ptrRecordsMap[p.GetKey()]
		if existing == nil {

			arpa, err := util.GetARPA(p.ARPA)
			if err != nil {
				return nil, err
			}

			p.ARPA = arpa
			p.SRC = source + ":static"
			ptrRecordsMap[p.GetKey()] = p

		} else {
			p.SRC = source + ":static-and-dynamic"
		}
	}

	var ptrRecords []*PTRrecord

	for _, v := range ptrRecordsMap {
		ptrRecords = append(ptrRecords, v)
	}

	t.domain.Records.PtrRecords = ptrRecords
	return &t.domain.Records, nil
}
