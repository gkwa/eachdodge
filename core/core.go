package core

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"syscall"

	"github.com/vishvananda/netlink"
)

type IPAddress struct {
	IP       string `json:"ip"`
	IsPublic bool   `json:"isPublic"`
}

func GetIPs() {
	addrs, err := netlink.AddrList(nil, syscall.AF_UNSPEC)
	if err != nil {
		log.Fatal(err)
	}

	var ipList []IPAddress
	for _, addr := range addrs {
		if addr.IP.IsLoopback() {
			continue
		}
		ip := addr.IP.String()
		isPublic := isPublicIP(ip)
		ipList = append(ipList, IPAddress{IP: ip, IsPublic: isPublic})
	}

	data, err := json.MarshalIndent(ipList, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("ips.json", data, 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

func isPublicIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	private := false
	if parsedIP.IsPrivate() || parsedIP.IsLoopback() || parsedIP.IsLinkLocalUnicast() || parsedIP.IsLinkLocalMulticast() {
		private = true
	}

	return !private
}
