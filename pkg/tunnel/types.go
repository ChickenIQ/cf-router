package tunnel

import (
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

type Tunnel struct {
	Dev   tun.Device
	Net   *netstack.Net
	WGDev *device.Device
}
