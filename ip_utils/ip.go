package ip_utils

import (
	"net"
	"log"
	"strings"
)

var IP string=""


func init() {
	IP=getLocalhostIp()
}


func getLocalhostIp()string {
	ifaces, err := net.Interfaces()
	if err!=nil{
		log.Printf("error,net.Interfaces error,%v\n",err)
		return ""
	}

	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err!=nil{
			log.Printf("error,interfaces.Addrs error,%v\n",err)
			return ""
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip==nil{
				continue
			}
			if ip.IsLoopback()==false{
				var ipStr=ip.String()
				if strings.Contains(ipStr,"."){
					return ipStr
				}
			}
		}
	}
	return ""
}



