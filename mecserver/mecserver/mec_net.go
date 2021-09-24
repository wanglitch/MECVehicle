package mecserver

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

//PingServer .
//是否可Ping通指定服务器
func pingServer(server string) bool {
	pinger, err := ping.NewPinger(server)
	CheckError(err)
	// pinger.SetPrivileged(false)
	// pinger.Interval = 5 * time.Microsecond
	pinger.Timeout = 5 * time.Second
	pinger.Count = 1
	pinger.Run() // blocks until finished
	stats := pinger.Statistics()
	// fmt.Printf("%+v\n", stats)
	if stats.PacketsRecv == pinger.Count {
		return true
	} else {
		return false
	}

}

//PingServerFromIP : 从指定的IP地址ping目标服务器
func pingServerFromIP(server string, ip string) bool {
	pinger, err := ping.NewPinger(server)
	CheckError(err)
	pinger.SetPrivileged(false)
	pinger.Interval = 1000 * time.Microsecond
	pinger.Timeout = 1 * time.Second
	pinger.Count = 1
	pinger.Source = ip
	err = pinger.Run() // blocks until finished
	CheckError(err)
	stats := pinger.Statistics()

	// fmt.Printf("%+v\n", stats)
	if stats.PacketsRecv == pinger.Count {
		return true
	} else {
		return false
	}
}

//GetConnMACAddress : 获取正在与zookeeper通信的接口的mac地址
//server : 目标服务器的地址
func GetConnMACAddress(server string) string {
	var macaddress string

	netInterfaces, err := net.Interfaces()
	CheckError(err)
	for _, netInterface := range netInterfaces {
		ift, _ := netInterface.Addrs()
		if ift != nil {
			// fmt.Println("Addrs:", ift[0])
			fmt.Printf("Cheaking %s ...\n", ift[0])
			sip := (ift[0]).String()
			sip = sip[:strings.Index(sip, "/")]
			if pingServerFromIP(server, sip) {
				macaddress = netInterface.HardwareAddr.String()
			}

		}

	}
	return macaddress
}
