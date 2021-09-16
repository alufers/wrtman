package main

//  ubus call network.wireless status
// ubus call hostapd.wlan0 get_clients
type RadioConfig struct {
	Hwmode      string `json:"hwmode"`
	Path        string `json:"path"`
	Htmode      string `json:"htmode"`
	Channel     string `json:"channel"`
	CellDensity int    `json:"cell_density"`
}
type UbusInterfaceConfig struct {
	Mode       string   `json:"mode"`
	Ssid       string   `json:"ssid"`
	Encryption string   `json:"encryption"`
	Key        string   `json:"key"`
	Network    []string `json:"network"`
}
type UbusInterfaces struct {
	Section  string              `json:"section"`
	Ifname   string              `json:"ifname"`
	Config   UbusInterfaceConfig `json:"config"`
	Vlans    []interface{}       `json:"vlans"`
	Stations []interface{}       `json:"stations"`
}
type UbusRadio struct {
	Up               bool             `json:"up"`
	Pending          bool             `json:"pending"`
	Autostart        bool             `json:"autostart"`
	Disabled         bool             `json:"disabled"`
	RetrySetupFailed bool             `json:"retry_setup_failed"`
	Config           RadioConfig      `json:"config"`
	Interfaces       []UbusInterfaces `json:"interfaces"`
}

type UbusGetClientsResponse struct {
	Freq    float64               `json:"freq"`
	Clients map[string]UbusClient `json:"clients"`
}

type UbusClient struct {
	Auth       bool          `json:"auth"`
	Assoc      bool          `json:"assoc"`
	Authorized bool          `json:"authorized"`
	Preauth    bool          `json:"preauth"`
	Wds        bool          `json:"wds"`
	Wmm        bool          `json:"wmm"`
	Ht         bool          `json:"ht"`
	Vht        bool          `json:"vht"`
	Wps        bool          `json:"wps"`
	Mfp        bool          `json:"mfp"`
	Rrm        []int         `json:"rrm"`
	Aid        int           `json:"aid"`
	Bytes      UbusRxTxStats `json:"bytes"`
	Airtime    UbusRxTxStats `json:"airtime"`
	Packets    UbusRxTxStats `json:"packets"`
	Rate       UbusRxTxStats `json:"rate"`
	Signal     float64       `json:"signal"`
}
type UbusRxTxStats struct {
	Rx float64 `json:"rx"`
	Tx float64 `json:"tx"`
}
