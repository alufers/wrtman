package main

import (
	"fmt"
	"sync"
	"time"
)

type DHCPLeasesService struct {
	conn          *OpenWrtConnection
	CacheDuration time.Duration
	lastRequest   time.Time
	cachedLeases  []*DHCPLease
	mutex         sync.Mutex
	Hooks         []func([]*DHCPLease)
}

func NewDHCPLeasesService(conn *OpenWrtConnection) *DHCPLeasesService {
	return &DHCPLeasesService{
		conn:          conn,
		CacheDuration: 10 * time.Second,
		Hooks:         []func([]*DHCPLease){},
	}
}

func (d *DHCPLeasesService) GetDHCPLeases() ([]*DHCPLease, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.cachedLeases == nil || time.Since(d.lastRequest) > d.CacheDuration {
		var err error
		d.cachedLeases, err = d.getDHCPLeasesUncached()
		if err != nil {
			d.cachedLeases = nil
			return nil, fmt.Errorf("failed to fetch DHCP leasess for %v: %w", d.conn.Hostname, err)
		}
		d.lastRequest = time.Now()
		for _, h := range d.Hooks {
			h(d.cachedLeases)
		}
	}

	return d.cachedLeases, nil
}

func (d *DHCPLeasesService) AddHook(h func([]*DHCPLease)) {
	d.Hooks = append(d.Hooks, h)
}

func (d *DHCPLeasesService) getDHCPLeasesUncached() ([]*DHCPLease, error) {

	fileContents, err := d.conn.RunCommandAndGetString("cat /tmp/dhcp.leases")
	if err != nil {
		return nil, fmt.Errorf("failed to read /tmp/dhcp.leases: %w", err)
	}

	leases, err := ParseDHCPLeases(fileContents)
	if err != nil {

		return nil, fmt.Errorf("failed to parse /tmp/dhcp.leases: %w", err)
	}

	return leases, nil

}
