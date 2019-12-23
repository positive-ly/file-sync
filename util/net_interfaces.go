package util

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func LocalIpAddCIDR() {
	if !Scan_lan_server {
		return
	}
	localIps := localIpAddress()
	for _, v := range localIps {
		v_ips, err := hosts(v + "/24")
		if err != nil {
			exit("[ERROR] See Loacl Ip Address DICR Fail : " + err.Error())
		}
		for _, v_ip := range v_ips {
			if _interfaces_http_discovery(v_ip) {
				fmt.Println("[FIND] file-sync server : ", v_ip)
			}
		}
	}

}

func _interfaces_http_discovery(ip string) bool {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 200 * time.Millisecond,
		}).Dial,
		IdleConnTimeout:       200 * time.Millisecond,
		TLSHandshakeTimeout:   200 * time.Millisecond,
		ResponseHeaderTimeout: 200 * time.Millisecond,
	}
	client := &http.Client{
		Timeout:   3000 * time.Millisecond,
		Transport: transport,
	}
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/discovery", ip, _httpProd), nil)
	if err != nil {
		exit("[ERROR] New Request fail : " + err.Error())
	}
	resp, err := client.Do(request)
	// resp, err := http.Get(fmt.Sprintf("http://%s:%s/discovery", ip, _httpProd))
	if err != nil {
		return false
		// exit("[ERROR] request discovery fail : " + err.Error())
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		exit("[ERROR] read response body fail : " + err.Error())
	}
	if string(bytes) == _discovery_text {
		return true
	}
	return false
}

func localIpAddress() []string {
	//https://blog.csdn.net/ycf8788/article/details/83111080
	net_inter, err := net.Interfaces()
	if err != nil {
		exit("[ERROR] See Loacl Ip Address Fail : " + err.Error())
	}
	ips := make([]string, 0)
	for _, v := range net_inter {
		if (v.Flags & net.FlagUp) != 0 {
			addrs, err := v.Addrs()
			if err != nil {
				exit("[ERROR] See Loacl Ip Address Fail : " + err.Error())
			}
			for _, val := range addrs {
				if ipnet, ok := val.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ips = append(ips, ipnet.IP.String())
					}
				}
			}
		}
	}
	return ips
}

func hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
