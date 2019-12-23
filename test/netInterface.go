package test

import (
	"fmt"
	"net"
	"os"
)

func LocalIpAddress() {
	//https://blog.csdn.net/ycf8788/article/details/83111080
	net_inter, err := net.Interfaces()
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
	for _, v := range net_inter {
		if (v.Flags & net.FlagUp) != 0 {
			addrs, err := v.Addrs()
			if err != nil {
				fmt.Print(err)
				os.Exit(2)
			}
			for _, val := range addrs {
				if ipnet, ok := val.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						fmt.Println(ipnet.IP.String())
					}
				}
			}
		}
	}
}
