package wireguard

import (
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

type Config struct {
	Interface Interface
	Peers     []Peer
}

type Interface struct {
	PrivateKey string
	Address    string
	DNS        string
	MTU        int
}

type Peer struct {
	PublicKey  string
	AllowedIPs string
	Endpoint   string
}

type Tun struct {
	Dev tun.Device
	Net *netstack.Net
}
