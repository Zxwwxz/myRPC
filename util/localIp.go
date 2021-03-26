package util

import (
	"net"
)

//获取本地ip
func GetLocalIP() (ip string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, interfaces := range netInterfaces {
		if (interfaces.Flags & net.FlagUp) != 0 {
			addressAll, _ := interfaces.Addrs()
			for _, address := range addressAll {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ip = ipnet.IP.String()
						return
					}
				}
			}
		}
	}
	return
}

