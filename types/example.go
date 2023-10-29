package types

import (
	"github.com/jodydadescott/home-server/types/proto"
	logger "github.com/jodydadescott/jody-go-logger"
)

func NewExampleConfig() *Config {

	d := &Domain{
		Domain: "home",
	}

	d.Records.AddARecords(&ARecord{
		Hostname: "a_record_1",
		IP:       "192.168.1.1",
	})

	d.Records.AddARecords(&ARecord{
		Hostname: "a_record_2",
		IP:       "192.168.1.2",
	})

	d.Records.AddAAAARecords(&ARecord{
		Hostname: "a_record_1",
		IP:       "2001:db8:3333:4444:5555:6666:7777:8888",
	})

	d.Records.AddCNameRecords(&CNameRecord{
		AliasHostname:  "cname_record_1",
		TargetHostname: "a_record_1",
		TargetDomain:   DefaultDomain,
	})

	d.Records.AddCNameRecords(&CNameRecord{
		AliasHostname:  "cname_record_2",
		TargetHostname: "a_record_2",
		TargetDomain:   DefaultDomain,
	})

	static := &StaticConfig{Enabled: true}
	static.AddDomains(d)

	unifiConfig := &UnifiConfig{}
	unifiConfig.Hostname = "https://10.0.1.1"
	unifiConfig.Username = "homeauto"
	unifiConfig.Password = "******"

	unifiConfig.Enabled = true
	unifiConfig.Refresh = DefaultRefresh

	unifiConfig.AddIgnoreMacs("60:22:32:9f:0f:fd")

	listener1 := &NetPort{
		Port:  53,
		Proto: proto.UDP,
	}

	listener2 := &NetPort{
		Port:  53,
		Proto: proto.TCP,
	}

	c := &Config{
		Notes:  "PTR records will automatically be created",
		Unifi:  unifiConfig,
		Static: static,
		Logging: &Logger{
			LogLevel: logger.DebugLevel,
		},
	}

	c.AddListeners(listener1, listener2)

	c.AddNameservers(&NetPort{
		IP:    "8.8.8.8",
		Port:  53,
		Proto: proto.UDP,
	})

	c.AddNameservers(&NetPort{
		IP:    "4.4.4.4",
		Port:  53,
		Proto: proto.TCP,
	})

	c.AddNameservers(&NetPort{
		IP:    "1.1.1.1",
		Port:  4053,
		Proto: proto.TCP,
	})

	return c
}
