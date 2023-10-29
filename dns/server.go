package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"go.uber.org/zap"

	"github.com/jodydadescott/home-server/types"
	"github.com/jodydadescott/home-server/types/proto"
)

type Server struct {
	listeners    []*NetPort
	domainNames  []string
	udpDnsClient *dns.Client
	tcpDnsClient *dns.Client
	clients      []*Client
	nameservers  []*NetPort
	trace        bool
}

func New(config *Config) *Server {

	if config == nil {
		panic("config is required")
	}

	if len(config.Listeners) <= 0 {
		config.Listeners = append(config.Listeners, &NetPort{
			Port:  types.DefaultDnsPort,
			Proto: types.DefaultDnsProto,
		})
	} else {
		for _, listener := range config.Listeners {
			switch listener.Proto {

			case proto.UDP, proto.TCP:

			case proto.Empty:
				listener.Proto = proto.UDP

			default:
				panic("Proto Invalid")
			}

			if listener.Port <= 0 {
				listener.Port = types.DefaultDnsPort
			}
		}
	}

	var nameservers []*NetPort
	for _, nameserver := range config.Nameservers {

		switch nameserver.Proto {

		case proto.UDP, proto.TCP:

		case proto.Empty:
			nameserver.Proto = proto.UDP

		default:
			panic("Proto Invalid")
		}

		if nameserver.Port <= 0 {
			nameserver.Port = types.DefaultDnsPort
		}

		nameservers = append(nameservers, nameserver)
	}

	c := &Server{
		listeners:    config.Listeners,
		udpDnsClient: &dns.Client{Net: "udp", SingleInflight: true},
		tcpDnsClient: &dns.Client{Net: "tcp", SingleInflight: true},
		nameservers:  nameservers,
		trace:        config.Trace,
	}

	for _, provider := range config.Providers {
		if provider == nil {
			panic("nil provider")
		}
		c.clients = append(c.clients, newClient(provider, c.trace))
	}

	return c
}

func (t *Server) Run(ctx context.Context) error {

	getARecord := func(name string) *ARecord {
		for _, client := range t.clients {
			r := client.getARecord(name)
			if r != nil {
				return r
			}
		}
		return nil
	}

	getAAAARecord := func(name string) *ARecord {
		for _, client := range t.clients {
			r := client.getAAAARecord(name)
			if r != nil {
				return r
			}
		}
		return nil
	}

	getPTRRecord := func(name string) *PTRrecord {
		for _, client := range t.clients {
			r := client.getPTRRecord(name)
			if r != nil {
				return r
			}
		}
		return nil
	}

	getCNameRecord := func(name string) *CNameRecord {
		for _, client := range t.clients {
			r := client.getCNameRecord(name)
			if r != nil {
				return r
			}
		}
		return nil
	}

	handleLocal := func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false

		switch r.Opcode {
		case dns.OpcodeQuery:

			for _, q := range m.Question {

				switch q.Qtype {

				case dns.TypeA:
					lookup := getARecord(q.Name)
					if lookup != nil {
						record := fmt.Sprintf("%s A %s", q.Name, lookup.GetValue())

						if t.trace {
							zap.L().Debug(fmt.Sprintf("success -> %s, source=%s", record, lookup.SRC))
						}

						rr, err := dns.NewRR(record)

						if err == nil {
							m.Answer = append(m.Answer, rr)
						} else {
							zap.L().Error(err.Error())
						}
					} else {
						if t.trace {
							zap.L().Debug(fmt.Sprintf("fail -> %s has no A record", q.Name))
						}
					}

				case dns.TypeAAAA:
					lookup := getAAAARecord(q.Name)
					if lookup != nil {
						record := fmt.Sprintf("%s AAAA %s", q.Name, lookup.GetValue())

						if t.trace {
							zap.L().Debug(fmt.Sprintf("success -> %s, source=%s", record, lookup.SRC))
						}

						rr, err := dns.NewRR(record)

						if err == nil {
							m.Answer = append(m.Answer, rr)
						} else {
							zap.L().Error(err.Error())
						}
					} else {
						if t.trace {
							zap.L().Debug(fmt.Sprintf("fail -> %s has no AAAA record", q.Name))
						}
					}

				case dns.TypePTR:
					lookup := getPTRRecord(q.Name)
					if lookup != nil {
						record := fmt.Sprintf("%s PTR %s", q.Name, lookup.GetValue())

						if t.trace {
							zap.L().Debug(fmt.Sprintf("success -> %s, source=%s", record, lookup.SRC))
						}

						rr, err := dns.NewRR(record)

						if err == nil {
							m.Answer = append(m.Answer, rr)
						} else {
							zap.L().Error(err.Error())
						}
					} else {
						if t.trace {
							zap.L().Debug((fmt.Sprintf("fail -> %s has no PTR record", q.Name)))
						}
					}

				case dns.TypeCNAME:
					lookup := getCNameRecord(q.Name)
					if lookup != nil {
						record := fmt.Sprintf("%s CNAME %s", q.Name, lookup.GetValue())

						if t.trace {
							zap.L().Debug(fmt.Sprintf("success -> %s, source=%s", record, lookup.SRC))
						}

						rr, err := dns.NewRR(record)

						if err == nil {
							m.Answer = append(m.Answer, rr)
						} else {
							zap.L().Error(err.Error())
						}
					} else {
						if t.trace {
							zap.L().Debug((fmt.Sprintf("fail -> %s has no CNAME record", q.Name)))
						}
					}

				}
			}

		}

		w.WriteMsg(m)
	}

	handleRemote := func(w dns.ResponseWriter, r *dns.Msg) {

		dnsClient := t.tcpDnsClient

		for _, nameserver := range t.nameservers {

			switch nameserver.Proto {

			case proto.TCP:
				dnsClient = t.tcpDnsClient

			case proto.UDP:
				dnsClient = t.udpDnsClient

			}

			r, _, err := dnsClient.Exchange(r, nameserver.GetIPColonPort())

			if err == nil {
				rString, _ := json.Marshal(r)

				if r.Rcode == dns.RcodeSuccess {
					r.Compress = true
					w.WriteMsg(r)

					if t.trace {
						zap.L().Debug(fmt.Sprintf("Remote Nameserver %s responded with %s", nameserver.GetIPColonPort(), rString))
					}

					return
				}
			} else {
				if t.trace {
					zap.L().Debug(fmt.Sprintf("Remote Nameserver %s responded with error %s", nameserver.GetIPColonPort(), err.Error()))
				}
			}
		}

		if t.trace {
			zap.L().Debug("failure to forward request")
		}

		m := new(dns.Msg)
		m.SetReply(r)
		m.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(m)
	}

	addDomainName := func(domainName string) {
		domainName = strings.ToLower(domainName)
		for _, existingDomain := range t.domainNames {
			if domainName == existingDomain {
				return
			}
		}
		t.domainNames = append(t.domainNames, domainName)
	}

	for _, client := range t.clients {
		addDomainName(client.GetDomainName())
		err := client.run()
		if err != nil {
			return err
		}
	}

	for _, v := range t.domainNames {
		zap.L().Debug(fmt.Sprintf("Adding domain %s to be handled locally", v))
		dns.HandleFunc(v+".", handleLocal)
	}

	dns.HandleFunc("10.in-addr.arpa.", handleLocal)
	dns.HandleFunc("168.192.in-addr.arpa.", handleLocal)
	dns.HandleFunc("0.0.16.127.in-addr.arpa.", handleLocal)
	dns.HandleFunc("0.0.168.192.in-addr.arpa.", handleLocal)

	if len(t.nameservers) > 0 {
		for _, v := range t.nameservers {
			if t.trace {
				zap.L().Debug(fmt.Sprintf("Forwarding to nameserver %s : %s", v.IP, string(v.Proto)))
			}
		}

		dns.HandleFunc(".", handleRemote)

	} else {
		zap.L().Debug("Forwarding to nameservers is not enabled")
	}

	errs := make(chan error, len(t.listeners))

	var servers []*dns.Server

	for _, listener := range t.listeners {
		zap.L().Info(fmt.Sprintf("Starting server on %s/%s", listener.IP+":"+strconv.Itoa(listener.Port), string(listener.Proto)))
		server := &dns.Server{Addr: listener.IP + ":" + strconv.Itoa(listener.Port), Net: string(listener.Proto)}
		servers = append(servers, server)

		go func() {
			err := server.ListenAndServe()
			if err != nil {
				errs <- err
			}
		}()

	}

	var err error

	select {

	case err = <-errs:
		zap.L().Info("Shutting down or error")

	case <-ctx.Done():
		zap.L().Info("Shutting down or signal")

	}

	for _, client := range t.clients {
		client.shutdown()
	}

	return err
}
