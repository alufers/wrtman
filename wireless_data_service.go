package main

import (
	"fmt"
	"sync"
	"time"
)

type APNetwork struct {
	APHostname    string      `json:"apHostname"`
	SSID          string      `json:"ssid"`
	Encryption    string      `json:"encryption"`
	Key           string      `json:"-"`
	Channel       string      `json:"channel"`
	Type          string      `json:"type"`
	InterfaceName string      `json:"interfaceName"`
	Clients       []*APClient `json:"clients"`
}

type APClient struct {
	MACAddress     string  `json:"macAddress"`
	SignalStrength float64 `json:"signalStrength"`
}

// WirelessDataService queries OpenWRT APs to
type WirelessDataService struct {
	conn           *OpenWrtConnection
	CacheDuration  time.Duration
	lastRequest    time.Time
	cachedNetworks []*APNetwork
	mutex          sync.Mutex
}

func NewWirelessDataService(conn *OpenWrtConnection) *WirelessDataService {
	return &WirelessDataService{
		conn:          conn,
		CacheDuration: 10 * time.Second,
		lastRequest:   time.Unix(0, 0),
	}
}

func (d *WirelessDataService) GetApNetworks() ([]*APNetwork, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.cachedNetworks == nil || time.Since(d.lastRequest) > d.CacheDuration {
		var err error
		d.cachedNetworks, err = d.getAPNetworksUncached()
		if err != nil {
			d.cachedNetworks = nil
			return nil, fmt.Errorf("failed to fetch wireless status for %v: %w", d.conn.Hostname, err)
		}
		d.lastRequest = time.Now()
	}
	return d.cachedNetworks, nil
}

func (d *WirelessDataService) getAPNetworksUncached() ([]*APNetwork, error) {
	status, err := d.getWirelessStatus()
	if err != nil {
		return nil, err
	}
	nets := []*APNetwork{}

	for k, v := range status {
		if len(v.Interfaces) <= 0 {
			return nil, fmt.Errorf("radio %v has no interfaces", k)
		}
		net := &APNetwork{
			APHostname:    d.conn.Hostname,
			SSID:          v.Interfaces[0].Config.Ssid,
			Key:           v.Interfaces[0].Config.Key,
			Encryption:    v.Interfaces[0].Config.Encryption,
			Channel:       v.Config.Channel,
			Type:          v.Config.Hwmode,
			InterfaceName: v.Interfaces[0].Ifname,
			Clients:       []*APClient{},
		}
		clientsResp, err := d.getInterfaceClients(net.InterfaceName)
		if err != nil {
			return nil, fmt.Errorf("failed to get clients for %v: %w", net.InterfaceName, err)
		}
		for macAddress, c := range clientsResp.Clients {
			net.Clients = append(net.Clients, &APClient{
				MACAddress:     macAddress,
				SignalStrength: c.Signal,
			})
		}
		nets = append(nets, net)
	}

	return nets, nil
}

func (d *WirelessDataService) getWirelessStatus() (out map[string]*UbusRadio, err error) {
	out = make(map[string]*UbusRadio)
	err = d.conn.UbusCall("network.wireless", "status", &out)
	return
}

func (d *WirelessDataService) getInterfaceClients(interfaceName string) (out *UbusGetClientsResponse, err error) {
	out = &UbusGetClientsResponse{}
	err = d.conn.UbusCall(fmt.Sprintf("hostapd.%v", interfaceName), "get_clients", out)
	return
}
