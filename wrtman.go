package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"golang.org/x/crypto/ssh"
)

//go:generate bash -c "cd wrtman-frontend && yarn install && yarn run build"

func main() {

	viper.SetConfigName("wrtman-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/wrtman/")
	viper.AddConfigPath("$HOME/.config/alufers/wrtman")
	viper.AddConfigPath(".")
	viper.SetDefault("http.addr", ":8080")
	viper.SetDefault("devices.main.addr", "192.168.1.1")
	viper.SetDefault("devices.main.user", "root")
	viper.SetDefault("extra.oui_db", "oui.txt")
	viper.SetDefault("ssh.key_path", "$HOME/.ssh/id_rsa")
	viper.SetDefault("db.type", "sqlite")
	viper.SetDefault("db.filename", "wrtman.db")
	viper.SetDefault("db.dsn", "SET FOR POSTGRES lol")

	// viper.SafeWriteConfigAs("wrtman-config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	log.Printf("Loaded config from: %v", viper.ConfigFileUsed())
	key, err := LoadKeyFile(os.ExpandEnv(viper.GetString("ssh.key_path")))
	if err != nil {
		panic(err)
	}
	mainConn := NewOpenWrtConnection(viper.GetString("devices.main.addr"), &ssh.ClientConfig{
		User: viper.GetString("devices.main.user"),
		Auth: []ssh.AuthMethod{
			key,
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	mainConn.HasDHCP = true
	application := NewApp([]*OpenWrtConnection{
		mainConn,
	})
	if err := application.ConnectToDB(); err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	app := fiber.New()

	app.Use(cors.New())
	application.MountEndpoints(app)
	application.MountHooks()
	err = application.AutodiscoverDHCPDevices()
	if err != nil {
		log.Fatal(err)
	}

	MountFrontend(app)
	err = app.Listen(viper.GetString("http.addr"))

	if err != nil {
		log.Fatal(err)
	}

}
