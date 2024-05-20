package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/vishvananda/netlink"
)

const (
	FAMILY_ALL = netlink.FAMILY_ALL
)

type IPAddress struct {
	Interface   string `json:"interface"`
	IP          string `json:"ip"`
	IPVersion   string `json:"ipVersion"`
	IsPublic    bool   `json:"isPublic"`
	IsInterface bool   `json:"isInterface"`
}

func Run(outfile string, out string) {
	ipList := IPs2()

	if out == "list" {
		for _, ip := range ipList {
			fmt.Println(ip.IP)
		}
	} else {
		data, err := json.MarshalIndent(ipList, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(data))

		if outfile != "" {
			err = os.WriteFile(outfile, data, 0o644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func IPs2() []IPAddress {
	links, err := netlink.LinkList()
	if err != nil {
		log.Fatal(err)
	}

	var ipList []IPAddress

	for _, link := range links {
		addrs, err := netlink.AddrList(link, FAMILY_ALL)
		if err != nil {
			fmt.Printf("Failed to retrieve IP addresses for %s: %v\n", link.Attrs().Name, err)
			continue
		}

		for _, addr := range addrs {
			ip := addr.IP
			ipVersion := "IPv4"
			if ip.To4() == nil {
				ipVersion = "IPv6"
			}
			ipList = append(ipList, IPAddress{
				Interface:   link.Attrs().Name,
				IP:          ip.String(),
				IPVersion:   ipVersion,
				IsPublic:    isPublicIP(ip.String()),
				IsInterface: true,
			})
		}
	}

	publicIP, err := getPublicIP()
	if err == nil {
		ipList = append(ipList, IPAddress{
			Interface:   "",
			IP:          publicIP,
			IPVersion:   "IPv4",
			IsPublic:    true,
			IsInterface: false,
		})
	}

	return ipList
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

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		IP string `json:"ip"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.IP, nil
}
