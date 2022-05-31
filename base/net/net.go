package net

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

// 获取局域网IP, copy from sonyflake
func PrivateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

// 寻找本机下一个可以被监听的端口
func NextAvaliablePort(port int) (int, error) {
	laddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return 0, err
	}
	for {
		if laddr.Port >= 65535 {
			return 0, fmt.Errorf("exceed max port number")
		}
		listener, err := net.ListenTCP("tcp4", laddr)
		if err != nil {
			laddr.Port += 1
		} else {
			if err := listener.Close(); err != nil {
				return 0, err
			}
			return laddr.Port, nil
		}
	}
}

func Addr2IP(addr string) string {
	if addr == "" {
		return ""
	}
	i := strings.LastIndex(addr, ":")
	if i < 0 {
		return addr
	}
	return addr[:i]
}

func GetXForwardFor(req *http.Request) string {
	x := req.Header.Get("X-Forwarded-For")
	if x == "" {
		return ""
	}
	xs := strings.Split(x, ",")
	if len(xs) > 0 {
		return strings.TrimSpace(xs[0])
	}
	return x
}

// 获取请求的IP
func GetRemoteAddr(req *http.Request) string {
	xip := GetXForwardFor(req)
	if xip != "" {
		return xip
	}
	return req.RemoteAddr
}

func GetIP(req *http.Request) string {
	addr := GetRemoteAddr(req)
	if addr != "" {
		return Addr2IP(addr)
	}
	return ""
}
