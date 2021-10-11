package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alufers/wrtman/models"
	"github.com/go-ping/ping"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	DB              *gorm.DB
	Connections     []*OpenWrtConnection
	OuiHelper       *OuiHelper
	MainDHCPService *DHCPLeasesService
	deviceWithDHCP  *OpenWrtConnection
}

func NewApp(Connections []*OpenWrtConnection) *App {
	var conn *OpenWrtConnection
	for _, c := range Connections {
		if c.HasDHCP {
			conn = c
		}
	}
	return &App{
		Connections:     Connections,
		OuiHelper:       NewOuiHelper(),
		MainDHCPService: NewDHCPLeasesService(conn),
		deviceWithDHCP:  conn,
	}
}

func (a *App) ConnectToDB() error {
	var err error
	switch viper.GetString("db.type") {
	case "sqlite":
		a.DB, err = gorm.Open(sqlite.Open(viper.GetString("db.filename")), &gorm.Config{})
	case "postgres":
		a.DB, err = gorm.Open(postgres.Open(viper.GetString("db.dsn")), &gorm.Config{})
	default:
		return fmt.Errorf("unknown database type '%v'", viper.GetString("db.type"))
	}
	if err != nil {
		return err
	}
	err = a.DB.AutoMigrate(models.AllModels...)
	return err
}

func (a *App) MountEndpoints(fiberApp *fiber.App) {
	fiberApp.Get("/api/devices", a.getDevices)
	fiberApp.Get("/api/dhcp-leases", a.getDHCPLeases)
	fiberApp.Get("/api/all-devices", a.getAllDevices)
	fiberApp.Post("/api/ping", a.postPingDevices)
}

func (a *App) MountHooks() {
	a.MainDHCPService.AddHook(a.dhcpLeasesFetchedHook)
	go func() {
		for {
			log.Printf("Fetching dhcp leases because of schedule...")
			a.MainDHCPService.GetDHCPLeases()
			time.Sleep(time.Minute * 10)
		}
	}()
}

func (a *App) getAllDevices(c *fiber.Ctx) error {
	devices := []*models.Device{}
	if err := a.DB.Find(&devices).Error; err != nil {
		return err
	}
	devicesWithVendor := []interface{}{}
	for _, dev := range devices {
		v, _ := a.OuiHelper.LookupVendor(dev.MACAddress)
		devicesWithVendor = append(devicesWithVendor, struct {
			*models.Device
			Vendor string `json:"vendor"`
		}{
			Device: dev,
			Vendor: v,
		})
	}
	return c.JSON(devicesWithVendor)
}

func (a *App) dhcpLeasesFetchedHook(leases []*DHCPLease) {
	log.Printf("dhcpLeasesFetchedHook %#v", leases)
	go func() {
		allAPNetworks, err := a.getAllApNetworks()
		if err != nil {
			log.Printf("Error fetching AP networks: %v", err)
			return
		}
		for _, l := range leases {
			dev := &models.Device{}
			if err := a.DB.Where("mac_address = ?", l.MACAddress).First(dev).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("Error fetching device %v: %v", l.MACAddress, err)
				continue
			}
			dev.MACAddress = l.MACAddress
			dev.Hostname = l.Hostname
			dev.LastSeen = time.Now()
			for _, net := range allAPNetworks {
				for _, client := range net.Clients {
					if client.MACAddress == l.MACAddress {
						dev.WirelessAPName = &net.APHostname
						dev.WirelessNetwork = &net.SSID
					}
				}
			}
			if err := a.DB.Save(dev).Error; err != nil {
				log.Printf("Error saving device to DB %v: %v", l.MACAddress, err)
			}
		}
	}()
}

// getAllApNetworks returns all AP networks from all connections
func (a *App) getAllApNetworks() ([]*APNetwork, error) {
	networks := []*APNetwork{}
	for _, conn := range a.Connections {
		apNetworks, err := conn.WirelessDataService.GetApNetworks()
		if err != nil {
			return nil, err
		}
		networks = append(networks, apNetworks...)
	}
	return networks, nil
}

func (a *App) AutodiscoverDHCPDevices() error {
	leases, err := a.MainDHCPService.GetDHCPLeases()
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
				conn := NewOpenWrtConnection(l.IPAddress+":22", a.deviceWithDHCP.SSHClientConfig)
				conn.Hostname = l.Hostname
				a.Connections = append(a.Connections, conn)
			}
		}
	}

	return nil
}

func (a *App) getDHCPLeases(ctx *fiber.Ctx) error {

	dhcpLeases, err := a.MainDHCPService.GetDHCPLeases()
	if err != nil {
		return err
	}

	dhcpLeasesWithVendor, _ := AddVendorsToDHCPLeases(a.OuiHelper, dhcpLeases)
	allAPNetworks, err := a.getAllApNetworks()
	if err != nil {
		return err
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
			log.Printf("failed to get mac address: %v", err)
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

type postPingBody struct {
	Addresses []string `json:"addresses"`
	Timeout   float64  `json:"timeout"`
}

func (a *App) postPingDevices(ctx *fiber.Ctx) error {
	dhcpLeases, err := a.MainDHCPService.GetDHCPLeases()
	if err != nil {
		return err
	}

	var body postPingBody
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}

	if body.Timeout == 0 {
		body.Timeout = 5
	}

	// check if addresses are in dhcp leases
	for _, addr := range body.Addresses {
		found := false
		for _, l := range dhcpLeases {
			if l.IPAddress == addr {
				found = true
			}
		}
		if !found {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("address %v not found in dhcp leases", addr),
			})
		}
	}

	resps := make(chan interface{})

	for _, addr := range body.Addresses {
		go func(addr string) {
			pinger, err := ping.NewPinger(addr)
			if err != nil {
				resps <- fiber.Map{
					"address": addr,
					"error":   err.Error(),
				}
				return
			}

			pinger.OnRecv = func(pkt *ping.Packet) {
				pinger.Stop()
				resps <- fiber.Map{
					"address": addr,
					"time":    float64(pkt.Rtt) / float64(time.Millisecond),
				}
			}
			go func() {
				time.Sleep(time.Duration(body.Timeout * float64(time.Second)))
				pinger.Stop()
				resps <- fiber.Map{
					"address": addr,
					"error":   "timeout",
				}
			}()
			err = pinger.Run()
			if err != nil {
				pinger.Stop()
				resps <- fiber.Map{
					"address": addr,
					"error":   err.Error(),
				}
			}

		}(addr)
	}

	ctx.Status(fiber.StatusOK)
	i := 0
	for resp := range resps {
		data, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		ctx.Write(data)
		ctx.Write([]byte("\n"))
		i++
		if i == len(body.Addresses) {
			break
		}
	}

	return nil

}
