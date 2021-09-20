package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type App struct {
	Connections           []*OpenWrtConnection
	cachedDHCPLeases      []*DHCPLease
	lastDHCPLeaseFetch    time.Time
	cachedDHCPLeasesMutex sync.Mutex
	OuiHelper             *OuiHelper
}

func NewApp(Connections []*OpenWrtConnection) *App {
	return &App{
		Connections: Connections,
		OuiHelper:   NewOuiHelper(),
	}
}

func (a *App) MountEndpoints(fiberApp *fiber.App) {
	fiberApp.Get("/api/devices", a.getDevices)
	fiberApp.Get("/api/dhcp-leases", a.getDHCPLeases)
}

func (a *App) AutodiscoverDHCPDevices() error {
	leases, err := a.fetchDHCPLeases()
	if err != nil {
		return err
	}
	devicesConfig := viper.Get("devices").(map[string]interface{})

	for k, _ := range devicesConfig {
		for _, l := range leases {
			if l.Hostname == k {
				alreadyExists := false
				for _, conn := range a.Connections {
					if conn.SSHAddress == l.IPAddress+":22" {
						alreadyExists = true
					}
				}
				if alreadyExists {
					continue
				}
				log.Printf("Discovered device %v (%v) at %v", l.Hostname, l.MACAddress, l.IPAddress)
				conn := NewOpenWrtConnection(l.IPAddress+":22", a.deviceWithDHCP().SSHClientConfig)
				conn.Hostname = l.Hostname
				a.Connections = append(a.Connections, conn)
			}
		}
	}

	return nil
}

func (a *App) getDHCPLeases(ctx *fiber.Ctx) error {

	dhcpLeases, err := a.fetchDHCPLeases()
	if err != nil {
		return err
	}

	dhcpLeasesWithVendor, _ := AddVendorsToDHCPLeases(a.OuiHelper, dhcpLeases)
	allAPNetworks := []*APNetwork{}

	for _, conn := range a.Connections {
		nets, err := conn.WirelessDataService.GetApNetworks()
		if err != nil {
			return err
		}
		allAPNetworks = append(allAPNetworks, nets...)
	}

	dhcpLeasesWithWirelessDetails := []*DHCPLeaseWithWirelessDetails{}

	for _, lease := range dhcpLeasesWithVendor {
		leaseWithDetails := &DHCPLeaseWithWirelessDetails{
			DHCPLeaseWithVendor: lease,
		}
		for _, net := range allAPNetworks {
			for _, c := range net.Clients {
				if c.MACAddress == lease.MACAddress {
					leaseWithDetails.SSID = &net.SSID
					leaseWithDetails.APHostname = &net.APHostname
					leaseWithDetails.SignalStrength = &c.SignalStrength
					leaseWithDetails.WirelessNetworkType = &net.Type
				}
			}
		}
		dhcpLeasesWithWirelessDetails = append(dhcpLeasesWithWirelessDetails, leaseWithDetails)
	}
	return ctx.JSON(dhcpLeasesWithWirelessDetails)

}

func (a *App) getDevices(ctx *fiber.Ctx) error {
	var out []interface{}
	for _, c := range a.Connections {
		available := true
		if c.Hostname == "" {
			err := c.DiscoverHostname()

			if err != nil {
				log.Printf("failed to discover hostname of %v: %v", c.SSHAddress, err)
				available = false
			}
		}
		var uptime *time.Duration
		if available {
			ut, err := c.GetUptime()
			if err != nil {
				log.Print(err)
			}
			uptime = &ut
		}
		var uptimeSeconds *float64
		if uptime != nil {
			sec := uptime.Seconds()
			uptimeSeconds = &sec
		}

		var vendorP *string
		addrs, err := c.GetMacAddrs()
		if err != nil || addrs == nil || len(addrs) == 0 {
			log.Printf("failed to get mac address: %w", err)
		} else {
			vendor, err := a.OuiHelper.LookupVendor(addrs[0])

			if err != nil {
				log.Printf("vendor lookup failed: %v", err)
			}

			vendorP = &vendor

		}

		out = append(out, fiber.Map{
			"hasDHCP":       c.HasDHCP,
			"available":     available,
			"address":       c.SSHAddress,
			"hostname":      c.Hostname,
			"uptime":        uptime,
			"vendor":        vendorP,
			"uptimeSeconds": uptimeSeconds,
		})
	}

	return ctx.JSON(out)
}

func (a *App) fetchDHCPLeases() ([]*DHCPLease, error) {
	a.cachedDHCPLeasesMutex.Lock()
	defer a.cachedDHCPLeasesMutex.Unlock()
	if a.cachedDHCPLeases == nil || a.lastDHCPLeaseFetch.Before(time.Now().Add(-time.Second*10)) {
		fileContents, err := a.deviceWithDHCP().RunCommandAndGetString("cat /tmp/dhcp.leases")
		if err != nil {
			a.cachedDHCPLeases = nil
			return nil, fmt.Errorf("failed to read /tmp/dhcp.leases: %w", err)
		}

		a.cachedDHCPLeases, err = ParseDHCPLeases(fileContents)
		if err != nil {
			a.cachedDHCPLeases = nil
			return nil, fmt.Errorf("failed to parse /tmp/dhcp.leases: %w", err)
		}

	}

	return a.cachedDHCPLeases, nil
}

func (a *App) deviceWithDHCP() *OpenWrtConnection {
	var conn *OpenWrtConnection
	for _, c := range a.Connections {
		if c.HasDHCP {
			conn = c
		}
	}
	return conn
}
