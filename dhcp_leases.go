package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"
)

type DHCPLease struct {
	MACAddress string    `json:"macAddress"`
	IPAddress  string    `json:"ipAddress"`
	Hostname   string    `json:"hostname"`
	ExpiryTime time.Time `json:"expiryTime"`
}

type DHCPLeaseWithVendor struct {
	*DHCPLease
	Vendor string `json:"vendor"`
}

type DHCPLeaseWithWirelessDetails struct {
	*DHCPLeaseWithVendor
	SSID                *string  `json:"ssid"`
	SignalStrength      *float64 `json:"signalStrength"`
	WirelessNetworkType *string  `json:"wirelessNetworkType"`
	APHostname          *string  `json:"apHostname"`
}

func ParseDHCPLeases(fileContents string) (leases []*DHCPLease, err error) {
	leases = []*DHCPLease{}
	scanner := bufio.NewScanner(strings.NewReader(fileContents))
	for scanner.Scan() {
		l := &DHCPLease{}
		var timestampRaw int64
		fmt.Sscanf(scanner.Text(), "%v %v %v %v", &timestampRaw, &l.MACAddress, &l.IPAddress, &l.Hostname)
		if l.MACAddress == "" {
			continue
		}
		l.ExpiryTime = time.Unix(timestampRaw, 0)
		leases = append(leases, l)
	}
	return
}

func AddVendorsToDHCPLeases(oh *OuiHelper, leases []*DHCPLease) (out []*DHCPLeaseWithVendor, err error) {
	out = []*DHCPLeaseWithVendor{}
	for _, l := range leases {
		vendor, err := oh.LookupVendor(l.MACAddress)
		if err != nil {
			log.Printf("vendor lookup err: %v", err)
		}
		out = append(out, &DHCPLeaseWithVendor{

			DHCPLease: l,
			Vendor:    vendor,
		})
	}
	return
}
