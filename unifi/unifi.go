package unifi

import (
	"fmt"
	"strings"
	"time"

	"github.com/jodydadescott/home-server/types"
	"github.com/jodydadescott/home-server/util"
	"github.com/jodydadescott/unifi-go-sdk"
	"go.uber.org/zap"
)

type Config = types.UnifiConfig
type ARecord = types.ARecord
type PTRrecord = types.PTRrecord
type Records = types.DomainRecords

type Client struct {
	unifiClient *unifi.Client
	config      *Config
	ignoreMacs  []string
	domain      string
}

const (
	source = "unifi"
)

func New(config *Config) *Client {

	if config == nil {
		panic("config is required")
	}

	config = config.Clone()

	if config.Hostname == "" {
		panic("Hostname is required")
	}

	if config.Username == "" {
		panic("Username is required")
	}

	if config.Password == "" {
		panic("Password is required")
	}

	domain := types.DefaultDomain
	if config.Domain != "" {
		domain = config.Domain
	}

	return &Client{
		config:      config,
		domain:      domain,
		unifiClient: unifi.New(&config.Config),
	}
}

func (t *Client) GetRefreshDuration() time.Duration {
	return t.config.Refresh
}

func (t *Client) GetDomainName() string {
	return t.domain
}

func (t *Client) GetName() string {
	return "unifi"
}

// GetRecords() (*Records, error)
func (t *Client) GetRecords() (*Records, error) {

	clients, err := t.unifiClient.GetClients()
	if err != nil {
		return nil, err
	}

	enrichedConfigs, err := t.unifiClient.GetEnrichedConfiguration()
	if err != nil {
		return nil, err
	}

	records := &Records{}

	for _, client := range clients {

		name := client.Name

		if name == "" {
			name = client.Hostname
		}

		if name == "" {
			zap.L().Debug(fmt.Sprintf("Client with MAC=%s and IP=%s does not have a name", client.Mac, client.IP))
			continue
		}

		ipV6 := getIpV6String(client.Ipv6Address)

		if client.IP == "" && ipV6 == "" {
			zap.L().Debug(fmt.Sprintf("Client %s does not have an IPv4 or IPv6", name))
			continue
		}

		addRecord := func(ip, iptype string) {

			arpa, err := util.GetARPA(ip)

			if err != nil {
				zap.L().Debug(fmt.Sprintf("Client %s has an invalid IP; error %s", name, err.Error()))
				return
			}

			a := &ARecord{
				Hostname: name,
				Domain:   t.domain,
				IP:       ip,
				SRC:      source + ":unifi-client",
			}

			switch iptype {
			case "A":
				records.AddARecords(a)

			case "AAAA":
				records.AddAAAARecords(a)

			default:
				panic("this should not happen")

			}

			p := &PTRrecord{
				ARPA:     arpa,
				Hostname: name,
				Domain:   t.domain,
				SRC:      source + ":unifi-client",
			}

			records.AddPtrRecords(p)
		}

		if client.IP != "" {
			addRecord(client.IP, "A")
		}

		if ipV6 != "" {
			addRecord(ipV6, "AAAA")
		}

	}

	for _, enrichedConfig := range enrichedConfigs {

		// TODO add IPv6 support here
		// enrichedConfig.Configuration.Ipv6ClientAddressAssignment

		name := strings.ToLower(enrichedConfig.Configuration.Name)
		ip := strings.Split(enrichedConfig.Configuration.IPSubnet, "/")[0]

		if name == "" {
			zap.L().Debug("Interface is missing its name")
			continue
		}

		if name == "default" {
			zap.L().Debug("Skipping default interface")
		}

		if ip == "" {
			zap.L().Debug(fmt.Sprintf("Interface %s is missing its ip", enrichedConfig.Configuration.Name))
			continue
		}

		interfaceName := "inf-" + name + "-"
		interfaceName += strings.Replace(ip, ".", "-", -1)

		arpa, err := util.GetARPA(ip)
		if err != nil {
			zap.L().Debug(fmt.Sprintf("Interface %s has an invalid IP; error %s", name, err.Error()))
			continue
		}

		a := &ARecord{
			Hostname: interfaceName,
			Domain:   t.domain,
			IP:       ip,
			SRC:      source + ":unifi-interface",
		}

		records.AddARecords(a)

		p := &PTRrecord{
			ARPA:     arpa,
			Hostname: interfaceName,
			Domain:   t.domain,
			SRC:      source + ":unifi-interface",
		}

		records.AddPtrRecords(p)
	}

	devices, err := t.unifiClient.GetDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices.NetworkDevices {

		if device.Name == "" {
			zap.L().Debug("Device is missing its name")
			continue
		}

		if device.IP == "" {
			zap.L().Debug(fmt.Sprintf("Interface %s is missing its ip", device.Name))
			continue
		}

		if t.config.IgnoreMac(device.Mac) {
			zap.L().Debug(fmt.Sprintf("Ignoring mac %s with name %s", device.Mac, device.Name))
			continue
		}

		arpa, err := util.GetARPA(device.IP)
		if err != nil {
			zap.L().Debug(fmt.Sprintf("Device %s has an invalid IP; error %s", device.Name, err.Error()))
			continue
		}

		a := &ARecord{
			Hostname: device.Name,
			Domain:   t.domain,
			IP:       device.IP,
			SRC:      source + ":unifi-device",
		}

		records.AddARecords(a)

		p := &PTRrecord{
			ARPA:     arpa,
			Hostname: device.Name,
			Domain:   t.domain,
			SRC:      source + ":unifi-device",
		}

		records.AddPtrRecords(p)
	}

	return records, nil

}

func getIpV6String(input []string) string {

	ipV6 := ""
	ipV6Len := len(input)

	if ipV6Len > 0 {
		for i, v := range input {
			if i == ipV6Len-1 {
				ipV6 = ipV6 + v
			} else {
				ipV6 = ipV6 + v + ":"
			}
		}
	}

	return ipV6
}
