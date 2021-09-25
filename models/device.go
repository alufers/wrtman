package models

import "time"

func init() {
	AllModels = append(AllModels, &Device{})
}

type Device struct {
	Model
	MACAddress      string    `json:"macAddress"`
	Hostname        string    `json:"hostname"`
	WirelessNetwork *string   `json:"wirelessNetwork"`
	WirelessAPName  *string   `json:"wirelessAPName"`
	LastSeen        time.Time `json:"lastSeen"`
	Note            string    `json:"note"`
}
