package dns

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	errRefreshDuration = time.Second * 30
)

type Client struct {
	mutex  sync.RWMutex
	ticker *xticker
	Provider
	done         chan bool
	aRecords     map[string]*ARecord
	aaaaRecords  map[string]*ARecord
	ptrRecords   map[string]*PTRrecord
	cnameRecords map[string]*CNameRecord
	trace        bool
}

func newClient(provider Provider, trace bool) *Client {
	if provider == nil {
		panic("provider is nil")
	}

	return &Client{
		Provider: provider,
		trace:    trace,
	}
}

type xticker struct {
	duration time.Duration
	ticker   *time.Ticker
}

func (t *xticker) reset(duration time.Duration) {

	if t.ticker == nil {
		t.ticker = time.NewTicker(duration)
		return
	}

	if t.duration == duration {
		return
	}

	t.duration = duration
	t.ticker.Reset(t.duration)
}

func (t *xticker) stop() {
	if t.ticker == nil {
		return
	}
	t.ticker.Stop()
}

func (t *Client) run() error {

	tick := func() {
		zap.L().Debug(fmt.Sprintf("Running refresh for %s", t.GetName()))
		err := t.refresh()

		if err == nil {
			t.ticker.reset(t.GetRefreshDuration())
		} else {
			t.ticker.reset(errRefreshDuration)
		}
	}

	init := func() {
		t.done = make(chan bool)
		t.ticker = &xticker{}
	}

	if t.GetRefreshDuration() <= 0 {
		zap.L().Info(fmt.Sprintf("Refresh for %s is not enabled", t.GetName()))
		return t.refresh()
	}

	zap.L().Info(fmt.Sprintf("Refresh for %s is %s", t.GetName(), t.GetRefreshDuration().String()))

	init()
	tick()

	go func() {
		for {
			select {
			case <-t.done:
				return

			case <-t.ticker.ticker.C:
				tick()

			}
		}
	}()

	return nil
}

func (t *Client) shutdown() {

	zap.L().Info(fmt.Sprintf("Shutting down %s", t.GetName()))

	if t.ticker != nil {
		t.ticker.stop()
	}

	if t.done != nil {
		t.done <- true
	}

}

func (t *Client) getARecords() []*ARecord {

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	var records []*ARecord

	for _, record := range t.aRecords {
		records = append(records, record)
	}

	return records
}

func (t *Client) getAAAARecords() []*ARecord {

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	var records []*ARecord

	for _, record := range t.aaaaRecords {
		records = append(records, record)
	}

	return records
}

func (t *Client) getARecord(name string) *ARecord {

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.aRecords == nil {
		return nil
	}

	return t.aRecords[name]
}

func (t *Client) getAAAARecord(name string) *ARecord {

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.aaaaRecords == nil {
		return nil
	}

	return t.aaaaRecords[name]
}

func (t *Client) getPTRRecord(name string) *PTRrecord {

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.ptrRecords == nil {
		return nil
	}

	return t.ptrRecords[name]
}

func (t *Client) getCNameRecord(name string) *CNameRecord {

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.cnameRecords == nil {
		return nil
	}

	return t.cnameRecords[name]
}

func (t *Client) refresh() error {

	aRecords := make(map[string]*ARecord)
	aaaRecords := make(map[string]*ARecord)
	ptrRecords := make(map[string]*PTRrecord)
	cnameRecords := make(map[string]*CNameRecord)

	records, err := t.GetRecords()
	if err != nil {
		return err
	}

	if records.ARecords != nil {
		for _, r := range records.ARecords {
			r = r.Clone()
			if r.Domain == "" {
				r.Domain = t.GetDomainName()
			}
			if t.trace {
				zap.L().Debug(fmt.Sprintf("Loading %s record with key %s and value %s", "A", r.GetKey(), r.GetValue()))
			}
			aRecords[r.GetKey()] = r
		}
	}

	if records.AAAARecords != nil {
		for _, r := range records.AAAARecords {
			r = r.Clone()
			if r.Domain == "" {
				r.Domain = t.GetDomainName()
			}
			if t.trace {
				zap.L().Debug(fmt.Sprintf("Loading %s record with key %s and value %s", "AAAA", r.GetKey(), r.GetValue()))
			}
			aaaRecords[r.GetKey()] = r
		}
	}

	if records.PtrRecords != nil {
		for _, r := range records.PtrRecords {
			r = r.Clone()
			if r.Domain == "" {
				r.Domain = t.GetDomainName()
			}
			if t.trace {
				zap.L().Debug(fmt.Sprintf("Loading %s record with key %s and value %s", "PTR", r.GetKey(), r.GetValue()))
			}
			ptrRecords[r.GetKey()] = r
		}
	}

	if records.CnameRecords != nil {
		for _, r := range records.CnameRecords {
			r = r.Clone()
			if r.AliasDomain == "" {
				r.AliasDomain = t.GetDomainName()
			}
			if r.AliasDomain == "" {
				r.AliasDomain = t.GetDomainName()
			}

			if t.trace {
				zap.L().Debug(fmt.Sprintf("Loading %s record with key %s and value %s", "CNAME", r.GetKey(), r.GetValue()))
			}
			cnameRecords[r.GetKey()] = r
		}
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.aRecords = aRecords
	t.aaaaRecords = aaaRecords
	t.ptrRecords = ptrRecords
	t.cnameRecords = cnameRecords

	return nil
}
